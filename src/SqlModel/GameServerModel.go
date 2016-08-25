package SqlModel

import (
	"database/sql"
	"fmt"
	. "misc/lisp_core"
)

type ServerData struct {
	Id      int
	Name    string
	Desc    string
	IP      string
	Port    string
	DBUser  string
	DBPwd   string
	State   int
	AddTime int64
}

func ServerData_Table() []*ServerData {
	query, _ := SQLDB.Query("select * from go_server_list")
	retData := make([]*ServerData, 0)
	for query.Next() {
		serverData := goQueryServerData2Struct(query)
		retData = append(retData, serverData)
	}
	return retData
}

func ServerData_ById(id int) *ServerData {
	sqlStr := "select * from go_server_list where id = " + Str(id)
	fmt.Println(sqlStr)
	query := SQLDB.QueryRow(sqlStr)
	retServerData := &ServerData{}
	err := query.Scan(&retServerData.Id, &retServerData.Name, &retServerData.Desc, &retServerData.IP, &retServerData.Port, &retServerData.DBUser,
		&retServerData.DBPwd, &retServerData.State, &retServerData.AddTime)
	if err != nil {
		fmt.Println(err, "query err")
	}
	return retServerData
}

func goQueryServerData2Struct(row *sql.Rows) *ServerData {
	retServerData := &ServerData{}
	row.Scan(&retServerData.Id, &retServerData.Name, &retServerData.Desc, &retServerData.IP, &retServerData.Port, &retServerData.DBUser,
		&retServerData.DBPwd, &retServerData.State, &retServerData.AddTime)
	return retServerData
}
