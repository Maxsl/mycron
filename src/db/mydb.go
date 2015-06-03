package db

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "reflect"
    "errors"
)

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

func fetchRow(ptr interface{},db *sql.DB, sqlstr string, args ...interface{})  error {
    rows,columns,err := rows(db,sqlstr, args)
    defer rows.Close()
    columnsLen := len(columns)
    kind, _, scan, err := scanVariables(ptr, columnsLen, false)
    if err != nil {
        return err
    }

    // Return data
    val := reflect.ValueOf(ptr).Elem()
    defer rows.Close()
    for rows.Next() {
        err = rows.Scan(scan...)
        if err != nil {
            return err
        }

        switch kind {
            /*   case reflect.Struct: // struct
               val.Set(reflect.ValueOf(ptrRow).Elem())*/

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

func fetchRows(ptr interface{},db *sql.DB, sqlstr string, args ...interface{}) ( error) {
    rows, columns, err := rows(db,sqlstr, args)
    if err != nil {
        panic(err.Error())
        return err
    }

    defer rows.Close()
    columnsLen := len(columns)

    kind, _, scan, err := scanVariables(ptr, columnsLen, true)
    if err != nil {
        panic(err.Error())
        return err
    }

    //return data
    val := reflect.ValueOf(ptr).Elem()

    for rows.Next() {
        if err := rows.Scan(scan...); err != nil {
            panic(err.Error())
            return err
        }

        switch kind {
            /*
                        case reflect.Struct: // struct
                        val.Set(reflect.Append(val, reflect.ValueOf(ptrRow).Elem()))
            */
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
    }

    if err = rows.Err(); err != nil {
        panic(err.Error())
        return err
    }

    return nil
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


// Item
type Item map[string]interface{}

// Where
type Where map[string]interface{}

// Get scan variables
func scanVariables(ptr interface{}, columnsLen int, isRows bool) (reflect.Kind, interface{}, []interface{}, error) {
    typ := reflect.ValueOf(ptr).Type()

    if typ.Kind() != reflect.Ptr {
        return 0, nil, nil, errors.New("ptr is not a pointer")
    }

    //log.Printf("%s\n", dataType.Elem().Kind())
    elemTyp := typ.Elem()

    if isRows { // Rows
        if elemTyp.Kind() != reflect.Slice {
            return 0, nil, nil, errors.New("ptr is not point a slice")
        }

        elemTyp = elemTyp.Elem()
    }

    elemKind := elemTyp.Kind()

    // element(value) is point to row
    scan := make([]interface{}, columnsLen)

    //log.Printf("%s\n", elemKind)

    if elemKind == reflect.Struct {
        if columnsLen != elemTyp.NumField() {
            return 0, nil, nil, errors.New("columnsLen is not equal elemTyp.NumField()")
        }

        row := reflect.New(elemTyp) // Data
        for i := 0; i < columnsLen; i++ {
            f := elemTyp.Field(i)
            if !f.Anonymous { // && f.Tag.Get("json") != ""
                scan[i] = row.Elem().FieldByIndex([]int{i}).Addr().Interface()
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