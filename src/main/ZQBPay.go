package main

import (
	"SqlModel"
	"fmt"
	"io/ioutil"
	"math"
	log "misc/seelog"
	"net/http"
	"strconv"
)

const ZQBApp_Key = "00107f79a204773999c7ce3171b55572"

func ZQBPayHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	plat_id := r.FormValue("plat_id")
	game_order := r.FormValue("game_order")
	plat_order := r.FormValue("plat_order")
	amount := r.FormValue("amount")
	server_id := r.FormValue("server_id")
	role_id := r.FormValue("role_id")
	//ext := r.FormValue("ext")
	sign := r.FormValue("sign")

	//验证签名
	//amount game_order plat_id plat_order role_id server_id
	curSingn := Tomd5("amount=" + amount + "&game_order=" + game_order + "&plat_id=" +
		plat_id + "&plat_order=" + plat_order + "&role_id=" + role_id + "&server_id=" + server_id + ZQBApp_Key)
	if curSingn != sign {
		log.Info("签名验证失败", sign, "   ", curSingn, " ", amount+" "+" "+game_order+" "+plat_id+" "+plat_order+" "+role_id+" "+server_id+" "+ZQBApp_Key)
		fmt.Fprint(w, "500")
		return
	}
	//换算金钱
	floatAmount, err := strconv.ParseFloat(amount, 32)
	if err != nil {
		log.Info("amount is not float number", amount)
		return
	}
	RMB := int(math.Ceil(float64(floatAmount) / 100.0))
	intServerId, err := strconv.Atoi(server_id)
	if err != nil {
		log.Info("serverid is not int number", server_id)
		return
	}

	serverData := SqlModel.ServerData_ById(intServerId)
	reqStr := "http://" + serverData.IP + ":" + serverData.Port + "/api/user?msg=1021&id=" +
		role_id + "&yb=" + strconv.Itoa(RMB) + "&orderId=" + game_order
	resp, err := http.Get(reqStr)
	//请求没有响应 插入错误记录并返回
	if err != nil {
		faildstr := "请求内部接口失败没有响应" + reqStr
		log.Info(faildstr)
		SqlModel.InsertPayHistory(&SqlModel.PayHistory{
			OrderId:    game_order,
			ServerId:   intServerId,
			RoleId:     role_id,
			PayItemId:  "0",
			State:      "F",
			FailedDesc: faildstr,
			Amount:     strconv.Itoa(RMB),
			UCJsonStr:  "",
		}, "")
		log.Info(faildstr)
		fmt.Fprint(w, "500")
		return
	}
	gameRet, _ := ioutil.ReadAll(resp.Body)
	//内部请求返回错误码 插入错误记录并返回
	if string(gameRet) != "0" {
		faildstr := "服务器内部请求错误 Code:" + string(gameRet) + reqStr
		log.Info(faildstr)
		SqlModel.InsertPayHistory(&SqlModel.PayHistory{
			OrderId:    game_order,
			ServerId:   intServerId,
			RoleId:     role_id,
			PayItemId:  "0",
			State:      "F",
			FailedDesc: faildstr,
			Amount:     strconv.Itoa(RMB),
			UCJsonStr:  "",
		}, "")
		log.Info(faildstr)
		return
	}

	SqlModel.InsertPayHistory(&SqlModel.PayHistory{
		OrderId:    game_order,
		ServerId:   intServerId,
		RoleId:     role_id,
		PayItemId:  "0",
		State:      "S",
		FailedDesc: "",
		Amount:     strconv.Itoa(RMB),
		UCJsonStr:  "",
	}, "")
	log.Info("充值成功", role_id, game_order, RMB)
	fmt.Fprint(w, "200")
}
