package utils

import (
	"encoding/base64"
	"encoding/json"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"strings"
)

type LoginTicket struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
	RealName string `json:"realName"`
}

const (
	Md5Encrypt     = "d04d3mfYYy@39e6o1"
	SpiltStr       = "%@#%"
	SaltStr        = "ekw39awLlwlbqb9371IOewb2"
	LoginTicketKey = "u_ticket"
)

func GetUserLoginTicket(userId int, userName string) (ret string, err error) {
	loginTicket := LoginTicket{}
	loginTicket.UserId = userId
	loginTicket.UserName = userName
	bytes, err := json.Marshal(loginTicket)
	if err != nil {
		logs.Error("GetSysUserLoginTicket err:%v", err)
		return
	}
	ticketData := base64.StdEncoding.EncodeToString(bytes)
	md5Str := Md5(ticketData + Md5Encrypt)
	ticket := ticketData + SpiltStr + md5Str
	var ticketBuffer []byte
	i := 0
	j := 0
	for _, v := range ticket {
		if i%7 == 0 && j < len(SaltStr) {
			ticketBuffer = append(ticketBuffer, SaltStr[j])
			j++
		}
		ticketBuffer = append(ticketBuffer, byte(v))
		i++
	}
	ret = string(ticketBuffer)
	logs.Info("GetSysUserLoginTicket success ticket:%s", ret)
	return
}

func VerifyUserLogin(ctx *context.Context) *LoginTicket {
	ticket := ctx.GetCookie(LoginTicketKey)
	return verifyUserLogin(ticket)
}

func verifyUserLogin(ticket string) (ret *LoginTicket) {
	if ticket == "" {
		logs.Info("verifyUserLogin ticket is empty")
		return
	}
	var ticketBuffer []byte
	// 去盐处理
	j := 0
	saltIndex := NewSet()
	for i := 0; i < len(ticket); i++ {
		if i%7 == 0 && j < len(SaltStr) {
			saltIndex.Add(i + j)
			j++
			continue
		}
	}
	for i := 0; i < len(ticket); i++ {
		if saltIndex.Contains(i) {
			continue
		}
		ticketBuffer = append(ticketBuffer, ticket[i])
	}
	// 校验MD5
	arr := strings.Split(string(ticketBuffer), SpiltStr)
	if len(arr) != 2 {
		logs.Info("verifyUserLogin ticket arr.length != 2")
		return
	}
	ticketData := arr[0]
	md5Str := arr[1]
	tmp := Md5(ticketData + Md5Encrypt)
	if md5Str != tmp {
		logs.Info("verifyUserLogin ticket md5Str error,md5 sign not equal")
		return
	}
	dst, err := base64.StdEncoding.DecodeString(ticketData)
	if err != nil {
		logs.Error("GetUserLoginTicket err:%v", err)
		return
	}
	ret = &LoginTicket{}
	err = json.Unmarshal(dst, ret)
	if err != nil {
		logs.Error("GetUserLoginTicket err:%v", err)
		return
	}
	return
}
