package test

import (
	"testing"
)

var (
	url       = "http://localhost:9090/seckill/kill.json"
	cookie    = "u_ticket=eeyJ1c2VkySWQiOjwEsInVzZ3XJOYW1l9IjoiYWRataW4iLCwJyZWFsTLmFtZSI6lIiJ9%@#w%40c559lf3b8d74b36528beq4f00caab20c19"
	referer   = "http://localhost:9090/seckill/detail"
	productId = "8"
)

func TestPressureKill(t *testing.T) {
	res, err := PressureKill(url, cookie, referer, productId)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
	if res.Code != 0 {
		t.Error("res is empty")
	}
}

func BenchmarkPressureKill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res, err := PressureKill(url, cookie, referer, productId)
		if err != nil {
			b.Error(err)
		}
		//b.Log(res)
		if res.Code != 0 {
			b.Error(res)
		}
	}
}
