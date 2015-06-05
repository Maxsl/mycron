package mydb

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "reflect"
    "errors"
    "strings"
  //  "fmt"
)

type MyDB struct {
   DB *sql.DB
}

func Open(driverName, dataSourceName string) (MyDB,error) {
    db := MyDB{}
    d,err:= sql.Open(driverName, dataSourceName)
    db.DB =d
    return db,err
}

func (my *MyDB) Close() {
    my.DB.Close()
}

func (my *MyDB) Raw(query string, args ...interface{}) RawSeter {
    o := NewRawSet(my.DB,query,args)
    return o
}

// Item
type Item map[string]interface{}

type rawSet struct {
    db *sql.DB
    sql string
    args []interface{}
}

// raw query seter
type RawSeter interface {
    Insert() (int64, error)
    Exec() (int64, error)
    FetchRow(interface{}) error
    FetchRows(interface{}) (int64, error)
 //   SetArgs(...interface{}) RawSeter
}


func NewRawSet(db *sql.DB ,query string, args []interface{}) RawSeter {
    o := new(rawSet)
    o.sql = query
    o.args = args
    o.db = db
    return o
}

//插入
func (r *rawSet) Insert() (int64, error) {
    stmtIns, err := r.db.Prepare(r.sql)
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close()

    result, err := stmtIns.Exec(r.args...)
    if err != nil {
        panic(err.Error())
    }
    return result.LastInsertId()
}

//修改和删除
func (r *rawSet) Exec() (int64, error) {
    stmtIns, err := r.db.Prepare(r.sql)
    if err != nil {
        panic(err.Error())
    }
    defer stmtIns.Close()

    result, err := stmtIns.Exec(r.args...)
    if err != nil {
        panic(err.Error())
    }
    return result.RowsAffected()
}

func (r *rawSet) FetchRow(ptr interface{}) ( error) {
    rows,columns,err := rows(r.db,r.sql, r.args)
    defer rows.Close()
    columnsLen := len(columns)
    kind,ptrRow, scan, err := scanVariables(ptr, columns, false)
    if err != nil {
        return err
    }
    val := reflect.ValueOf(ptr).Elem()
    defer rows.Close()
    for rows.Next() {
        err = rows.Scan(scan...)
        if err != nil {
            return err
        }
        switch kind {
            case reflect.Struct: // struct
            val.Set(reflect.ValueOf(ptrRow).Elem())
            case reflect.Map: //map
            row := make(map[string]interface{}, columnsLen)
            for i := 0; i < columnsLen; i++ {
                row[columns[i]] = typeAssertion(*(scan[i].(*interface{})))
            }
            val.Set(reflect.ValueOf(row))

            case reflect.Slice: //slice
            row := make([]interface{}, columnsLen)
            for i := 0; i < columnsLen; i++ {
                row[i] = typeAssertion(*(scan[i].(*interface{})))
            }
            val.Set(reflect.ValueOf(row))
        }
    }
    if err = rows.Err(); err != nil {
        return err
    }
    return nil
}

func (r *rawSet) FetchRows(ptr interface{}) (int64, error) {
    rows,columns,err := rows(r.db,r.sql, r.args)
    if err != nil {
        panic(err.Error())
        return 0,err
    }

    defer rows.Close()
    columnsLen := len(columns)

    kind,ptrRow, scan, err := scanVariables(ptr, columns, true)
    if err != nil {
        panic(err.Error())
        return 0,err
    }

    //return data
    val := reflect.ValueOf(ptr).Elem()
    var rowNum int64
    for rows.Next() {
        if err := rows.Scan(scan...); err != nil {
            panic(err.Error())
            return 0,err
        }

        switch kind {
            case reflect.Struct: // struct
            val.Set(reflect.Append(val, reflect.ValueOf(ptrRow).Elem()))
            case reflect.Map: // map
            row := make(map[string]interface{}, columnsLen)
            for i := 0; i < columnsLen; i++ {
                row[columns[i]] = typeAssertion(*(scan[i].(*interface{})))
            }
            val.Set(reflect.Append(val, reflect.ValueOf(row)))

            case reflect.Slice: // slice
            row := make([]interface{}, columnsLen)
            for i := 0; i < columnsLen; i++ {
                row[i] = typeAssertion(*(scan[i].(*interface{})))
            }
            val.Set(reflect.Append(val, reflect.ValueOf(row)))
        }
        rowNum++
    }

    if err = rows.Err(); err != nil {
        panic(err.Error())
        return 0,err
    }

    return rowNum,nil
}

func rows(db *sql.DB, sqlstr string, args []interface{}) (*sql.Rows, []string, error) {
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
    return rows,columns,nil
}

// Get scan variables
func scanVariables(ptr interface{}, columns []string, isRows bool) (reflect.Kind, interface{}, []interface{}, error) {
    columnsLen := len(columns)
    typ := reflect.ValueOf(ptr).Type()

    if typ.Kind() != reflect.Ptr {
        return 0, nil, nil, errors.New("ptr is not a pointer")
    }
    elemTyp := typ.Elem()

    if isRows { // Rows
        if elemTyp.Kind() != reflect.Slice {
            return 0, nil, nil, errors.New("ptr is not point a slice")
        }

        elemTyp = elemTyp.Elem()
    }

    elemKind := elemTyp.Kind()
    scan := make([]interface{}, columnsLen)
    if elemKind == reflect.Struct {
        row2 := make([]interface{}, columnsLen)
        row := reflect.New(elemTyp) // Data
        for i := 0; i < columnsLen; i++ {
            f := elemTyp.Field(i)
            if !f.Anonymous { // && f.Tag.Get("json") != ""
                //fmt.Println(row.Elem().FieldByName(strings.Title(columns[i])).IsValid())
                if row.Elem().FieldByName(strings.Title(columns[i])).IsValid(){
                    scan[i] = row.Elem().FieldByName(strings.Title(columns[i])).Addr().Interface()
                }else{
                    scan[i] = &row2[i]
                }
            }
        }
        return elemKind, row.Interface(), scan, nil
    }

    if elemKind == reflect.Map || elemKind == reflect.Slice {
        row := make([]interface{}, columnsLen) // Data
        for i := 0; i < columnsLen; i++ {
            scan[i] = &row[i]
        }

        return elemKind, &row, scan, nil
    }

    return 0, nil, nil, errors.New("ptr is not a point struct, map or slice")
}

// Type assertions
func typeAssertion(v interface{}) interface{} {
    switch v.(type) {
        case bool:
        //log.Printf("bool\n")
        return v.(bool)
        case int64:
        //log.Printf("int64\n")
        return v.(int64)
        case float64:
        //log.Printf("float64\n")
        return v.(float64)
        case string:
        //log.Printf("string\n")
        return v.(string)
        case []byte:
        //log.Printf("[]byte\n")
        return string(v.([]byte))
        default:
        //log.Printf("Unexpected type %#v\n", v)
        return ""
    }
}