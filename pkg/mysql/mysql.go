package mysql

import (
	"database/sql"
	"fmt"
	//"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbhostip    = "192.168.1.175:3306"
	dbusername  = "root"
	dbpassword  = "123456"
	dbtablename = "guid_mutex"
	dbname      = "guid"
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
	connSql := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbusername, dbpassword, dbhostip, dbname)
	db, err = sql.Open("mysql", connSql)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MysqlSelect(db *sql.DB) (int64, error) {
	Tx, _ := db.Begin()
	insertSql := fmt.Sprintf("INSERT %s SET time=?", dbtablename)
	stmt, err := db.Prepare(insertSql)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(MysqlTime())
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
