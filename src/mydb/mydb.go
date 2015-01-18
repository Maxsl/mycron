package mydb

import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
"time"
)

type Job struct {
	ID              int
	Name, Time, Cmd string
	STime ,ETime time.Duration
	Status uint8
}

func GetCronList() (jobss []Job, e error) {
	db, err := sql.Open("mysql", "root:@/mycron")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM cron")
	if err != nil {
		panic(err.Error())
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))
	jobs := make([]Job, len(values))
	i := 0
	// Fetch rows
	for rows.Next() {
		err = rows.Scan(&jobs[i].ID, &jobs[i].Name, &jobs[i].Time, &jobs[i].Cmd,&jobs[i].STime,&jobs[i].ETime,&jobs[i].Status)
		if err != nil {
			panic(err.Error())
		}
		i++
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
	return jobs, nil
}
