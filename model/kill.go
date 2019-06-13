package model

import (
	"sync"
	"time"
)

type SecResult struct {
	ProductId int    `json:"productId"`
	UserId    int    `json:"userId"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Token     string `json:"token"`
	TokenTime int64  `json:"tokenTime"`
}

func NewSecRequest() (secRequest *SecRequest) {
	secRequest = &SecRequest{
		ResultChan: make(chan *SecResult, 1),
	}
	return
}

type SecRequest struct {
	ProductId     int             `json:"productId" form:"productId"`
	UserId        int             `json:"userId"`
	AccessTime    time.Time       `json:"accessTime"`
	ClientAddr    string          `json:"clientAddr"`
	ClientReferer string          `json:"clientReferer"`
	CloseNotify   <-chan bool     `json:"-"`
	ResultChan    chan *SecResult `json:"-"`
}

type RateLimitMgr struct {
	UserLimitMap map[int]*SecMinCounter
	IpLimitMap   map[string]*SecMinCounter
	sync.Mutex
}
