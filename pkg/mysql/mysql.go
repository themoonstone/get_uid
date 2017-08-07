package mysql

import (
	"database/sql"
	"fmt"
	//"fmt"
	"time"

	"github.com/caarlos0/env"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbhostip    = "127.0.0.1:3306"
	dbusername  = "root"
	dbpassword  = "123456"
	dbtablename = "guid_mutex"
	dbname      = "guid"
)

type mysqlConf struct {
	User     string `env:"MYSQL_USER" envDefault:"root"`
	Password string `env:"MYSQL_PASSWORD,required" envDefault:"123456"`
	Host     string `env:"MYSQL_HOST" envDefault:"localhost:3306"`
	DB       string `env:"MYSQL_DB" envDefault:"guid"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var db *sql.DB

var dbConf mysqlConf

func init() {
	db = new(sql.DB)
	env.Parse(&dbConf)
}

func MysqlConn() (*sql.DB, error) {
	var err error
	connSql := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", dbConf.User, dbConf.Password, dbConf.Host, dbConf.DB)
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
