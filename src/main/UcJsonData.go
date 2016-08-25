package main

import (
	"strings"
)

const (
	UCAPIKEY = "b045edd7173a5058330d09da5fbc1e71"
)

type UCPayJson struct {
	Ver  string        //版本号
	Sign string        //签名参数
	Data UCPayJsonData //Data
}

type UCPayJsonData struct {
	OrderId      string
	AccountId    string
	Creator      string
	PayWay       string `json:"payWay"`
	Amount       string
	CallBackInfo string
	OrderStatus  string
	FailedDesc   string
	CpOrderId    string
	GameId       string `json:"gameId"`
}

func SplitCallBackInfo(callBackStr string) map[string]string {
	strAttrs := strings.Split(callBackStr, "#")
	if !(len(strAttrs) == 2) {
		return nil
	}
	retMap := make(map[string]string)
	for i := 0; i < len(strAttrs); i++ {
		strkv := strings.Split(strAttrs[i], "=")
		if len(strkv) == 2 {
			retMap[strkv[0]] = strkv[1]
		} else {
			return nil
		}
	}
	return retMap
}
