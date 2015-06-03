package db

import (
    "testing"
    "git.oschina.net/wida/mycron/src/mycron"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
)
var (
    db  *sql.DB
    err error
)
func init() {
    db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
    mycron.Mysql_user, mycron.Mysql_pwd, mycron.Mysql_host, mycron.Mysql_prot,mycron.Mysql_dbname))
    if err != nil {
        panic(err.Error())
    }
    db.SetMaxOpenConns(30)
    db.SetMaxIdleConns(10)
    db.Ping()
}

func TestQueryRow (t *testing.T) {
    rawset := rawSet{db:db}
    rawset.query ="SELECT * FROM cron where id =?"
    job := mycron.Job{}
    rawset.SetArgs(1).QueryRow(&job)
    fmt.Println(job)
}