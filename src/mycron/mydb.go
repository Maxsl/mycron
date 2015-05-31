package mycron

import (
    "database/sql"
    "fmt"
    "reflect"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
)

//插入
func insert(db *sql.DB, sqlstr string, args ...interface{}) (int64, error) {
    stmtIns, err := db.Prepare(sqlstr)
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close()

    result, err := stmtIns.Exec(args...)
    if err != nil {
        panic(err.Error())
    }
    return result.LastInsertId()
}

//修改和删除
func exec(db *sql.DB, sqlstr string, args ...interface{}) (int64, error) {
    stmtIns, err := db.Prepare(sqlstr)
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close()

    result, err := stmtIns.Exec(args...)
    if err != nil {
        panic(err.Error())
    }
    return result.RowsAffected()
}

//取一行数据，注意这类取出来的结果都是string
func fetchRow(db *sql.DB, sqlstr string, args ...interface{}) (*map[string]string, error) {
    stmtOut, err := db.Prepare(sqlstr)
    if err != nil {
        panic(err.Error())
    }
    defer stmtOut.Close()

    rows, err := stmtOut.Query(args...)
    if err != nil {
        panic(err.Error())
    }

    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error())
    }

    values := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{}, len(values))
    ret := make(map[string]string, len(scanArgs))

    for i := range values {
        scanArgs[i] = &values[i]
    }
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error())
        }
        var value string

        for i, col := range values {
            if col == nil {
                value = "NULL"
            } else {
                value = string(col)
            }
            ret[columns[i]] = value
        }
        break //get the first row only
    }
    return &ret, nil
}

//取多行，注意这类取出来的结果都是string
func fetchRows(db *sql.DB, sqlstr string, args ...interface{}) (*[]map[string]string, error) {
    stmtOut, err := db.Prepare(sqlstr)
    if err != nil {
        panic(err.Error())
    }
    defer stmtOut.Close()

    rows, err := stmtOut.Query(args...)
    if err != nil {
        panic(err.Error())
    }

    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error())
    }

    values := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{}, len(values))

    ret := make([]map[string]string, 0)
    for i := range values {
        scanArgs[i] = &values[i]
    }

    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error())
        }
        var value string
        vmap := make(map[string]string, len(scanArgs))
        for i, col := range values {
            if col == nil {
                value = "NULL"
            } else {
                value = string(col)
            }
            vmap[columns[i]] = value
        }
        ret = append(ret, vmap)
    }
    return &ret, nil
}

// set field value to row container
func setFieldValue(ind reflect.Value, value interface{}) {
    switch ind.Kind() {
        case reflect.Bool:
        if value == nil {
            ind.SetBool(false)
        } else if v, ok := value.(bool); ok {
            ind.SetBool(v)
        } else {
            v, _ := StrTo(ToStr(value)).Bool()
            ind.SetBool(v)
        }

        case reflect.String:
        if value == nil {
            ind.SetString("")
        } else {
            ind.SetString(ToStr(value))
        }

        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        if value == nil {
            ind.SetInt(0)
        } else {
            val := reflect.ValueOf(value)
            switch val.Kind() {
                case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                ind.SetInt(val.Int())
                case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
                ind.SetInt(int64(val.Uint()))
                default:
                v, _ := StrTo(ToStr(value)).Int64()
                ind.SetInt(v)
            }
        }
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        if value == nil {
            ind.SetUint(0)
        } else {
            val := reflect.ValueOf(value)
            switch val.Kind() {
                case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
                ind.SetUint(uint64(val.Int()))
                case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
                ind.SetUint(val.Uint())
                default:
                v, _ := StrTo(ToStr(value)).Uint64()
                ind.SetUint(v)
            }
        }
        case reflect.Float64, reflect.Float32:
        if value == nil {
            ind.SetFloat(0)
        } else {
            val := reflect.ValueOf(value)
            switch val.Kind() {
                case reflect.Float64:
                ind.SetFloat(val.Float())
                default:
                v, _ := StrTo(ToStr(value)).Float64()
                ind.SetFloat(v)
            }
        }

        case reflect.Struct:
        if value == nil {
            ind.Set(reflect.Zero(ind.Type()))
        }
    }
}
