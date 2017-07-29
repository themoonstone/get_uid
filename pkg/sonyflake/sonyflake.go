// Package sonyflake implements Sonyflake, a distributed unique ID generator inspired by Twitter's Snowflake.
//
// A Sonyflake ID is composed of
//     39 bits for time in units of 10 msec
//      8 bits for a sequence number
//     16 bits for a machine id
package sonyflake

import (
	"database/sql"
	"errors"

	"fmt"
	"net"
	"sync"
	"time"

	"get_uid/pkg/mysql"
)

// These constants are the bit lengths of Sonyflake ID parts.
const (
	BitLenTime      = 39                               // bit length of time
	BitLenSequence  = 8                                // bit length of sequence number
	BitLenMachineID = 63 - BitLenTime - BitLenSequence // bit length of machine id
)

// Settings configures Sonyflake:
//
// StartTime is the time since which the Sonyflake time is defined as the elapsed time.
// If StartTime is 0, the start time of the Sonyflake is set to "2014-09-01 00:00:00 +0000 UTC".
// If StartTime is ahead of the current time, Sonyflake is not created.
//
// MachineID returns the unique ID of the Sonyflake instance.
// If MachineID returns an error, Sonyflake is not created.
// If MachineID is nil, default MachineID is used.
// Default MachineID returns the lower 16 bits of the private IP address.
//
// CheckMachineID validates the uniqueness of the machine ID.
// If CheckMachineID returns false, Sonyflake is not created.
// If CheckMachineID is nil, no validation is done.
type Settings struct {
	StartTime      time.Time
	MachineID      func() (uint16, error)
	CheckMachineID func(uint16) bool
}

// Sonyflake is a distributed unique ID generator.
type Sonyflake struct {
	mutex       *sync.Mutex
	startTime   int64
	elapsedTime int64
	sequence    uint16
	machineID   uint16
}

// GlobalVal is a goroutine pool
type GlobalVal struct {
	ch       chan uint64
	Poolsize int
	Slice    []*Sonyflake
}

// a mysql db
var db *sql.DB

// NewGlobal returns a new Global Sonyflake configured with the given Settings.
// NewGlobal Initialize poolsize with passed parameters
// Set the length of the channel to 100
func (gl *GlobalVal) NewGlobal(size int) *GlobalVal {

	gbl := &GlobalVal{
		ch:       make(chan uint64, 100),
		Poolsize: size,
	}
	var err error
	db, err = mysql.MysqlConn()
	if err != nil {
		panic(err)
	}
	for i := 0; i < gbl.Poolsize; i++ {
		var st Settings
		st.MachineID = genMachineID

		st.StartTime = time.Now()
		sr := NewSonyflake(st)
		if sr == nil {
			panic("sonyflake not created")
		}
		gbl.Slice = append(gbl.Slice, sr)

	}

	return gbl
}

// Generates a globally unique id in a concurrency way
func (gl *GlobalVal) GenId() {
	for i := 0; i < len(gl.Slice); i++ {

		go func(a int) {

			for {
				gl.Slice[a].ChanNextID(&gl.ch)
			}
		}(i)

	}
}

//GetId return a uint64 id which client requested by a channel
//the default timeout is 10sec,If more than 10 seconds has not yet taken the data
//return the time out error to the client
func (gl *GlobalVal) GetId() (uint64, error) {
	select {
	case id := <-gl.ch:
		return id, nil
	case <-time.After(10 * time.Second):
		return 0, errors.New("time out")
	}

}

//getMysqlID return a mysql auto_increment id.
func getMysqlID() (uint16, error) {
	id, err := mysql.MysqlSelect(db)
	if err != nil {
		return 0, err
	}

	return uint16(id), nil

}

//genMachineID return a machineID Consists of lower16BitPrivateIP and mysql auto_increment id.
func genMachineID() (uint16, error) {
	/*addrs, err := net.InterfaceAddrs()

	if err != nil {
		return 0, err
	}
		var ip net.IP
		for _, address := range addrs {

			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = net.ParseIP(ipnet.IP.String()).To4()

					if ip == nil {
						return 0, err
					}
					continue
				}

			}
		}*/
	ip, err := lower16BitPrivateIP()
	if err != nil {
		return 0, err
	}
	mysql_id, err := getMysqlID()
	if err != nil {
		return 0, err
	}
	return ip + mysql_id, nil
}

// NewSonyflake returns a new Sonyflake configured with the given Settings.
// NewSonyflake returns nil in the following cases:
// - Settings.StartTime is ahead of the current time.
// - Settings.MachineID returns an error.
// - Settings.CheckMachineID returns false.
func NewSonyflake(st Settings) *Sonyflake {

	sf := new(Sonyflake)
	sf.mutex = new(sync.Mutex)
	sf.sequence = uint16(1<<BitLenSequence - 1)

	if st.StartTime.After(time.Now()) {
		fmt.Println("time error")
		return nil
	}
	if st.StartTime.IsZero() {
		sf.startTime = toSonyflakeTime(time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC))
	} else {
		sf.startTime = toSonyflakeTime(st.StartTime)
	}

	var err error
	if st.MachineID == nil {
		sf.machineID, err = lower16BitPrivateIP()
	} else {
		sf.machineID, err = st.MachineID()
	}
	if err != nil || (st.CheckMachineID != nil && !st.CheckMachineID(sf.machineID)) {
		fmt.Println(err.Error())
		return nil
	}
	return sf
}

// NextID generates a next unique ID.
// After the Sonyflake time overflows, NextID returns an error.
func (sf *Sonyflake) NextID() (uint64, error) {
	const maskSequence = uint16(1<<BitLenSequence - 1)

	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	current := currentElapsedTime(sf.startTime)
	if sf.elapsedTime < current {
		sf.elapsedTime = current
		sf.sequence = 0
	} else { // sf.elapsedTime >= current
		sf.sequence = (sf.sequence + 1) & maskSequence
		if sf.sequence == 0 {
			sf.elapsedTime++
			overtime := sf.elapsedTime - current
			time.Sleep(sleepTime((overtime)))
		}
	}

	return sf.toID()
}

// ChanNextID generates a next unique ID.
// Write id to the incoming parameter channel
func (sf *Sonyflake) ChanNextID(ch *chan uint64) {
	const maskSequence = uint16(1<<BitLenSequence - 1)

	current := currentElapsedTime(sf.startTime)
	if sf.elapsedTime < current {
		sf.elapsedTime = current
		sf.sequence = 0
	} else { // sf.elapsedTime >= current
		sf.sequence = (sf.sequence + 1) & maskSequence
		if sf.sequence == 0 {
			sf.elapsedTime++
			overtime := sf.elapsedTime - current
			time.Sleep(sleepTime((overtime)))
		}
	}

	UID, err := sf.toID()
	if err != nil {
		panic(err)
	}
	*ch <- UID
}

const sonyflakeTimeUnit = 1e7 // nsec, i.e. 10 msec

func toSonyflakeTime(t time.Time) int64 {
	return t.UTC().UnixNano() / sonyflakeTimeUnit
}

func currentElapsedTime(startTime int64) int64 {
	return toSonyflakeTime(time.Now()) - startTime
}

func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime)*10*time.Millisecond -
		time.Duration(time.Now().UTC().UnixNano()%sonyflakeTimeUnit)*time.Nanosecond
}

func (sf *Sonyflake) toID() (uint64, error) {
	if sf.elapsedTime >= 1<<BitLenTime {
		return 0, errors.New("over the time limit")
	}

	return uint64(sf.elapsedTime)<<(BitLenSequence+BitLenMachineID) |
		uint64(sf.sequence)<<BitLenMachineID |
		uint64(sf.machineID), nil
}

func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func lower16BitPrivateIP() (uint16, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

// Decompose returns a set of Sonyflake ID parts.
func Decompose(id uint64) map[string]uint64 {
	const maskSequence = uint64((1<<BitLenSequence - 1) << BitLenMachineID)
	const maskMachineID = uint64(1<<BitLenMachineID - 1)

	msb := id >> 63
	time := id >> (BitLenSequence + BitLenMachineID)
	sequence := id & maskSequence >> BitLenMachineID
	machineID := id & maskMachineID
	return map[string]uint64{
		"id":         id,
		"msb":        msb,
		"time":       time,
		"sequence":   sequence,
		"machine-id": machineID,
	}
}
