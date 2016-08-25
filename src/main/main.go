package main

import (
	"SqlModel"
	"crypto/md5"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	cfg "misc/config"
	log "misc/seelog"
	"net/http"
	"strconv"
)

func main() {
	initlog()
	cfg.Load()
	cfgDic := cfg.Get()
	SqlModel.Init(cfgDic["sqlstr"])

	//http.HandleFunc("/api/pay", PayHandler)
	//http.HandleFunc("/api/pay/lj", LJPayHandler)
	http.HandleFunc("/api/pay/zqb", ZQBPayHandler)
	log.Info("server start success port:" + cfgDic["port"] + " sqlstr:" + cfgDic["sqlstr"])
	http.ListenAndServe(":"+cfgDic["port"], nil)
}

func initlog() {
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")
	if err != nil {
		log.Critical("err parsing config log file", err)
		return
	}
	log.ReplaceLogger(logger)
}

func PayHandler(w http.ResponseWriter, r *http.Request) {
	bydata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		return
	}

	ucJson := &UCPayJson{}
	json.Unmarshal(bydata, ucJson)
	log.Info(string(bydata))
	//验证签名
	if !checkSign(ucJson) {
		log.Info("签名验证失败", ucJson)
		fmt.Fprint(w, "FAILURE")
		return
	}

	fmt.Fprint(w, "SUCCESS")
	//UC充值失败 返回并存储本次记录
	if ucJson.Data.OrderStatus == "F" {
		log.Info("充值失败", ucJson)
		InsertPayHistory2Sql(ucJson, string(bydata))
		return
	}
	//UC充值成功
	callBackMap := SplitCallBackInfo(ucJson.Data.CallBackInfo)
	if callBackMap == nil {
		log.Info("自定义回调数据格式错误" + ucJson.Data.CallBackInfo)
		return
	}
	ServerId, _ := strconv.Atoi(callBackMap["ServerId"])
	serverData := SqlModel.ServerData_ById(ServerId)
	YB, _ := strconv.ParseFloat(ucJson.Data.Amount, 32)
	if YB < 1 {
		YB = 1
	}
	YBInt := int(YB)
	reqStr := "http://" + serverData.IP + ":" + serverData.Port + "/api/user?msg=1021&id=" + callBackMap["RoleId"] + "&yb=" + strconv.Itoa(YBInt)
	resp, err := http.Get(reqStr)
	//内部请求没有响应 插入错误记录并返回
	if err != nil {
		ucJson.Data.FailedDesc = "请求内部接口失败没有响应" + reqStr
		log.Info(ucJson.Data.FailedDesc, ucJson)
		InsertPayHistory2Sql(ucJson, string(bydata))
		return
	}
	gameRet, _ := ioutil.ReadAll(resp.Body)
	//内部请求返回错误码 插入错误记录并返回
	if string(gameRet) != "0" {
		ucJson.Data.FailedDesc = "服务器内部请求错误 Code:" + string(gameRet) + reqStr
		log.Info(ucJson.Data.FailedDesc, ucJson)
		InsertPayHistory2Sql(ucJson, string(bydata))
		return
	}
	InsertPayHistory2Sql(ucJson, string(bydata))
	log.Info("充值成功", ucJson)
}

func checkSign(jsonData *UCPayJson) bool {
	checkStr := "accountId=" + jsonData.Data.AccountId + "amount=" +
		jsonData.Data.Amount + "callbackInfo=" + jsonData.Data.CallBackInfo +
		"creator=" + jsonData.Data.Creator + "failedDesc=" + jsonData.Data.FailedDesc +
		"gameId=" + jsonData.Data.GameId + "orderId=" + jsonData.Data.OrderId +
		"orderStatus=" + jsonData.Data.OrderStatus + "payWay=" + jsonData.Data.PayWay + UCAPIKEY
	md5Str := Tomd5(checkStr)
	if md5Str == jsonData.Sign {
		return true
	}

	return false
}

func Tomd5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func InsertPayHistory2Sql(ucJson *UCPayJson, jsonStr string) {
	callBackMap := SplitCallBackInfo(ucJson.Data.CallBackInfo)
	ServerId, _ := strconv.Atoi(callBackMap["ServerId"])

	SqlModel.InsertPayHistory(&SqlModel.PayHistory{
		OrderId:     ucJson.Data.OrderId,
		ServerId:    ServerId,
		RoleId:      callBackMap["RoleId"],
		PayItemId:   "0",
		State:       ucJson.Data.OrderStatus,
		FailedDesc:  ucJson.Data.FailedDesc,
		Amount:      ucJson.Data.Amount,
		UCAccountId: ucJson.Data.AccountId,
		UCJsonStr:   jsonStr,
	}, jsonStr)
}
