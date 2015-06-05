package mydb

import (
    "testing"
    "fmt"
)
type Job struct {
    ID              int
    Name, Time, Cmd string
    STime, ETime    int
    Status          uint8
    Running         uint8
    IsModify        uint8
    Process         uint8
    Ip              string
    Singleton       uint8
}
func TestMyDb(t *testing.T) {
    db, err := Open("mysql", "root:@tcp(127.0.0.1:3306)/mycron?charset=utf8")
    if err != nil {
        t.Error(err.Error())
    }
    db.DB.SetMaxOpenConns(30)
    db.DB.SetMaxIdleConns(10)
    db.DB.Ping()

    defer db.Close()
/*    d := Item{}
    err = db.Raw("SELECT * FROM cron where id=?",2).FetchRow(&d)
    if (err != nil ){
        t.Error(err.Error())
    }
    fmt.Println(d)
    fmt.Println(d["uid"] == 1)

    s := &[]Item{}
    i,err:= db.Raw("SELECT * FROM cron").FetchRows(s)
    if err != nil{
        t.Error(err.Error())
    }
    fmt.Println(s,i)*/

    var data Job
    err = db.Raw("SELECT * FROM cron where id=?",2).FetchRow(&data)
    if (err != nil ){
        t.Error(err.Error())
    }
    fmt.Println(data)
}

