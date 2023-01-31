package logic

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ipuppet/gtools/config"
	"github.com/ipuppet/gtools/utils"
)

var (
	conf *config.Config

	httpClient *http.Client
)

func init() {
	conf = &config.Config{
		Filename: "config.json",
	}
	conf.Init()
	conf.AddNotifyer(config.LoggerNotifyer())

	httpClient = InitHttpClient()
}

func getConfigForKey(key string) string {
	value, err := conf.Get(key)
	if err != nil {
		log.Fatal(err)
	}

	return value.(string)
}

func InitHttpClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   6 * time.Second,
		KeepAlive: 3 * time.Second,
	}
	transport := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			jackett_host := getConfigForKey("jackett_host")
			jackett_port := getConfigForKey("jackett_port")

			if addr == jackett_host+":"+jackett_port {
				_jackett_ip, err := conf.Get("jackett_ip")
				if err == nil {
					jackett_ip := strings.TrimSpace(_jackett_ip.(string))
					if jackett_ip != "" {
						addr = jackett_ip + ":" + jackett_port
					}
				}
			}

			return dialer.DialContext(ctx, network, addr)
		},
	}
	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: &transport,
	}
}

func ParseQuery(q string, season int, ep int) string {
	cache := map[string]interface{}{}
	err := utils.GetStorageJSON("", "tmdb.json", cache)
	if err == nil {
		name, ok := cache[q].(string)
		if ok {
			return name
		}
	}

	u, _ := url.Parse(getConfigForKey("tmdb_search_url"))

	params := url.Values{}
	tmdb_api_key := getConfigForKey("tmdb_api_key")
	params.Set("api_key", tmdb_api_key)
	params.Set("query", q)
	params.Set("language", "zh")

	u.RawQuery = params.Encode()
	resp, err := httpClient.Get(u.String())
	if err != nil {
		return q
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var tmdb map[string]interface{}
	err = json.Unmarshal(body, &tmdb)
	if err != nil {
		return q
	}

	results, ok := tmdb["results"].([]interface{})
	if ok && len(results) > 0 {
		name, ok := results[0].(map[string]interface{})["name"].(string)
		if ok {
			cache[q] = name
			s, err := json.Marshal(cache)
			if err == nil {
				utils.SetStorageContent("", "tmdb.json", string(s))
			}

			if season > -1 {
				var epStr string
				if ep < 10 {
					epStr = "0" + strconv.Itoa(ep)
				} else {
					epStr = strconv.Itoa(ep)
				}
				return name + " S" + strconv.Itoa(season) + " " + epStr
			}

			return name
		}
	}

	return q
}

func Proxy(request *http.Request) (*http.Response, error) {
	jackett_host := getConfigForKey("jackett_host")
	jackett_port := getConfigForKey("jackett_port")
	host := jackett_host + ":" + jackett_port

	jackett_scheme := getConfigForKey("jackett_scheme")
	targetUrl := jackett_scheme + "://" + host + request.URL.String()

	// 构造请求
	requ, _ := http.NewRequest(request.Method, targetUrl, request.Body)

	// 重写 header
	for key, values := range request.Header {
		if len(values) == 1 {
			requ.Header.Set(key, values[0])
		} else {
			requ.Header.Set(key, values[0])
			for _, value := range values[1:] {
				requ.Header.Add(key, value)
			}
		}
	}

	requQuery := requ.URL.Query()

	if q := requQuery.Get("q"); q != "" {
		season := -1
		ep := -1
		if requQuery.Has("season") {
			season, _ = strconv.Atoi(requQuery.Get("season"))
			ep, _ = strconv.Atoi(requQuery.Get("ep"))
			requQuery.Del("season")
			requQuery.Del("ep")
		}

		requQuery.Set("q", ParseQuery(q, season, ep))
		requ.URL.RawQuery = requQuery.Encode()

		// 日志打印修饰后的 URL
		log.Println(requ.URL.String())
	}
	requ.Host = host

	resp, err := httpClient.Do(requ)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
