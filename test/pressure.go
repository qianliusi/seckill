package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"seckill/model"
	"strings"
)

func PressureKill(url, cookie, referer, productId string) (res model.SecResult, err error) {
	var r http.Request
	err = r.ParseForm()
	if err != nil {
		return
	}
	r.Form.Add("productId", productId)
	request, err := http.NewRequest("POST", url, strings.NewReader(strings.TrimSpace(r.Form.Encode())))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("Cookie", cookie)
	request.Header.Set("Referer", referer)
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}
	return
}
