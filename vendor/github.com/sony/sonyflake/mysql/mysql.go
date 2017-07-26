package mysql

import (
	"database/sql"
	//"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbhostip    = "127.0.0.1:3306"
	dbusername  = "root"
	dbpassword  = "123456"
	dbtablename = "guid_mutex"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var db *sql.DB

func init() {
	db = new(sql.DB)
}

func MysqlConn() (*sql.DB, error) {
	var err error
	db, err = sql.Open("mysql", "root:123456@tcp(192.168.1.175:3306)/guid?charset=utf8")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MysqlSelect(db *sql.DB) (int64, error) {
	Tx, _ := db.Begin()
	stmt, err := db.Prepare("INSERT guid_mutex SET user=?,time=?")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec("troy", MysqlTime())
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	//fmt.Println(id)
	err = Tx.Commit()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func MysqlTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
