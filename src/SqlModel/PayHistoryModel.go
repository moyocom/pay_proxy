package SqlModel

import (
	"fmt"
	_ "misc/mysql"
	"time"
)

type PayHistory struct {
	Id          int
	OrderId     string
	ServerId    int
	RoleId      string
	PayItemId   string
	State       string
	FailedDesc  string
	Amount      string
	UCAccountId string
	UCJsonStr   string
}

func InsertPayHistory(payHistory *PayHistory, ucJsonStr string) {
	stmt, err := SQLDB.Prepare("insert into go_payhistory(OrderId,ServerId,RoleId,PayItemId,State,FailedDesc,Amount,UCAccountId,UCJsonStr,Time)values(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(payHistory.OrderId, payHistory.ServerId, payHistory.RoleId,
		payHistory.PayItemId, payHistory.State, payHistory.FailedDesc,
		payHistory.Amount, payHistory.UCAccountId, ucJsonStr, int(time.Now().Unix()))
	if err != nil {
		fmt.Println(err)
	}
}
