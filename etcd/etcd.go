// Write 2017 The etcd guid
//
//Take advantage of the key-value and strong consistency of etcd to dynamically generate a globally unique id
//you can get the instruction of this package in the test file 'etcd_test.go'
package etcd

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

//generates a 32-bit md5 string
func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//genernte a guid
func uniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return getMd5String(base64.URLEncoding.EncodeToString(b))
}

type EtcdConfig struct {
	cfg     client.Config
	options client.SetOptions
	k       client.KeysAPI
}

//to initalize the config of etcd's connection and client
func (ecfg *EtcdConfig) NewEtcdConfig(hostSlice []string) {
	if len(hostSlice) == 0 {
		hostSlice = []string{"http://127.0.0.1:2379"}
	}
	ecfg.cfg = client.Config{
		Endpoints:               hostSlice,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	ecfg.options = client.SetOptions{
		PrevExist: client.PrevNoExist,
	}

	c, err := client.New(ecfg.cfg)
	if err != nil {
		log.Fatal(err)
	}
	ecfg.k = client.NewKeysAPI(c)

}

//func (ecfg *EtcdConfig) SetUid() (*client.Response, error) {
//	s := UniqueId()
//	resp, err := ecfg.k.Set(context.Background(), s, s, &ecfg.options)
//	if err != nil {
//		return nil, err
//	}
//	return resp, nil
//}

// obtain a GUID
func (ecfg *EtcdConfig) GetUid() (string, error) {
	s := uniqueId()
	resp, err := ecfg.k.Set(context.Background(), s, s, &ecfg.options)
	if err != nil {
		return "", err
	}
	resp, err = ecfg.k.Get(context.Background(), s, nil)
	if err != nil && resp == nil {
		return "", err
	}
	return resp.Node.Value, nil
}
