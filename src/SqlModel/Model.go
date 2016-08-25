package SqlModel

import (
	"database/sql"
	"fmt"
	_ "misc/mysql"
)

var SQLDB *sql.DB

func Init(sqlstr string) {
	var err error
	SQLDB, err = sql.Open("mysql", sqlstr)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = SQLDB.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("mysql open success")
}
