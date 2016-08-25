package main

import (
	"SqlModel"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	log "misc/seelog"
	"net/http"
	"strconv"
)

func LJPayHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	orderId := r.FormValue("orderId")
	price := r.FormValue("price")
	freePrice := r.FormValue("freePrice")
	//channelCode := r.FormValue("channelCode")
	callbackInfo := r.FormValue("callbackInfo")
	sign := r.FormValue("sign")
	mysign := Tomd5(orderId + price + callbackInfo + "f9b4f36b13ae4d88b3a7d941f43dd126")
	if mysign != sign {
		log.Info("签名验证失败", sign, mysign)
		fmt.Fprint(w, "fail")
		return
	}
	fmt.Fprint(w, "success")
	callbackStr, _ := base64Decode(callbackInfo)
	callBackMap := SplitCallBackInfo(string(callbackStr))
	if callBackMap == nil {
		log.Info("自定义回调数据格式错误" + string(callbackStr))
		return
	}
	log.Info(callbackStr)
	ServerId, _ := strconv.Atoi(callBackMap["ServerId"])
	serverData := SqlModel.ServerData_ById(ServerId)
	YB, _ := strconv.Atoi(price)
	FreeYB, _ := strconv.Atoi(freePrice)
	YB /= 100
	FreeYB /= 100
	YB += FreeYB
	if YB < 1 {
		YB = 1
	}
	reqStr := "http://" + serverData.IP + ":" + serverData.Port + "/api/user?msg=1021&id=" + callBackMap["RoleId"] + "&yb=" + strconv.Itoa(YB)
	resp, err := http.Get(reqStr)
	//内部请求没有响应 插入错误记录并返回
	if err != nil {
		faildstr := "请求内部接口失败没有响应" + reqStr
		log.Info(faildstr)
		SqlModel.InsertPayHistory(&SqlModel.PayHistory{
			OrderId:    orderId,
			ServerId:   ServerId,
			RoleId:     callBackMap["RoleId"],
			PayItemId:  "0",
			State:      "F",
			FailedDesc: faildstr,
			Amount:     strconv.Itoa(YB),
			UCJsonStr:  "",
		}, "")
		return
	}

	gameRet, _ := ioutil.ReadAll(resp.Body)
	//内部请求返回错误码 插入错误记录并返回
	if string(gameRet) != "0" {
		faildstr := "服务器内部请求错误 Code:" + string(gameRet) + reqStr
		log.Info(faildstr)
		SqlModel.InsertPayHistory(&SqlModel.PayHistory{
			OrderId:    orderId,
			ServerId:   ServerId,
			RoleId:     callBackMap["RoleId"],
			PayItemId:  "0",
			State:      "F",
			FailedDesc: faildstr,
			Amount:     strconv.Itoa(YB),
			UCJsonStr:  "",
		}, "")
		return
	}
	SqlModel.InsertPayHistory(&SqlModel.PayHistory{
		OrderId:    orderId,
		ServerId:   ServerId,
		RoleId:     callBackMap["RoleId"],
		PayItemId:  "0",
		State:      "S",
		FailedDesc: "",
		Amount:     strconv.Itoa(YB),
		UCJsonStr:  "",
	}, "")
	log.Info("充值成功")
}

func base64Decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
