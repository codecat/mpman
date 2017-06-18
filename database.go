package main

import "strconv"
import "database/sql"
import "reflect"

import _ "github.com/go-sql-driver/mysql"
import "github.com/codecat/go-libs/log"

type Row map[string]interface{}

var db *sql.DB

func dbMakeDsn() string {
	ret := Config.Database.Username

	if Config.Database.Password != "" {
		ret += ":" + Config.Database.Password
	}

	ret += "@tcp(" + Config.Database.Hostname + ":" + strconv.Itoa(Config.Database.Port) + ")"
	ret += "/" + Config.Database.Database

	return ret
}

func dbOpen() bool {
	var err error
	db, err = sql.Open("mysql", dbMakeDsn())
	if err != nil {
		log.Fatal("Failed to connect to database: %s", err.Error())
		return false
	}

	if !dbExec("SELECT 1") {
		return false
	}

	return true
}

func dbExec(query string, args ...interface{}) bool {
	_, err := db.Exec(query, args...)
	if err != nil {
		log.Error("Query error: %s", err.Error())
		log.Debug("Query was: \"%s\"", query)
		return false
	}
	return true
}

func dbQuery(query string, args ...interface{}) []Row {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Error("Query error: %s", err.Error())
		log.Debug("Query was: \"%s\"", query)
		return nil
	}
	defer rows.Close()

	cols, err := rows.Columns()
	colTypes, err := rows.ColumnTypes()

	row := make([]interface{}, len(cols), len(cols))
	for i, t := range colTypes {
		switch vt := t.ScanType(); vt {
			case reflect.TypeOf(int32(0)): row[i] = new(int)
			case reflect.TypeOf(string("")): row[i] = new(string)
			default: log.Warn("Unhandled set-up type %s", vt.Name())
		}
	}

	ret := []Row{}

	for rows.Next() {
		err = rows.Scan(row...)
		if err != nil {
			log.Error("Couldn't scan row: %s", err.Error())
			continue
		}

		newRow := Row{}
		for i, t := range colTypes {
			switch vt := t.ScanType(); vt {
				case reflect.TypeOf(int32(0)): newRow[cols[i]] = *row[i].(*int)
				case reflect.TypeOf(string("")): newRow[cols[i]] = *row[i].(*string)
				default: log.Warn("Unhandled unmarshal type %s", vt.Name())
			}
		}
		ret = append(ret, newRow)
 	}

	return ret
}
