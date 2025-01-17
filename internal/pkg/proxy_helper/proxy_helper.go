package proxy_helper

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func ProxyTest(proxyAddr string) (Speed int, Status int, err error) {
	// 首先检测 proxyAddr 是否合法，必须是 http 的代理，不支持 https 代理
	re := regexp.MustCompile(`(http):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`)
	result := re.FindAllStringSubmatch(proxyAddr, -1)
	if result == nil {
		err = errors.New("proxy address illegal, only support http://xx:xx")
		return 0, 0, err
	}
	// 检测代理iP访问地址
	var testUrl string
	testUrl = "http://google.com"
	// 解析代理地址
	proxy, err := url.Parse(proxyAddr)
	// 设置网络传输
	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}
	// 创建连接客户端
	httpClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	begin := time.Now() //判断代理访问时间
	// 使用代理IP访问测试地址
	res, err := httpClient.Get(testUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	speed := int(time.Now().Sub(begin).Nanoseconds() / 1000 / 1000) //ms
	// 判断是否成功访问，如果成功访问StatusCode应该为200
	if res.StatusCode != http.StatusOK {
		return
	}
	return speed, res.StatusCode, nil
}
