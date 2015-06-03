package db

import (
    "testing"
    "database/sql"
    "fmt"
)
func TestMyDb(t *testing.T) {
    db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/mycron?charset=utf8")
    if err != nil {
        t.Error(err.Error())
    }
    defer db.Close()
    d := Item{}
    err = fetchRow(&d,db, "SELECT * FROM cron where id=?",2)
    if err != nil{
        t.Error(err.Error())
    }
    fmt.Println(d)
    fmt.Println(d["uid"] == 1)

    s := & []Item{}

    err = fetchRows(s,db,"SELECT * FROM cron")
    if err != nil{
        t.Error(err.Error())
    }
    fmt.Println(s)
}

