package httpmgr

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

//向一个服务器发送消息(发送标准格式)(默认超时 10秒)
func Post(url string, data *bytes.Buffer) ([]byte, error) {
	to, err := time.ParseDuration("10s")// ---设置消息发送超时时间
	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: to} // ---超时时间
	resp, err := c.Post(url, "application/json;charset=utf-8", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

//Get(目前用来向api发送消息返回数据(这个函数只负责返回从Api获得的数据,不管数据正确性,由于数据正确性的判断很复杂,所以放在外部判断) (默认超时 10秒)；
func Get(url string) ([]byte, error) {
	to, err := time.ParseDuration("10s")
	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: to}
	//发送消息获取记录
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取响应的数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
