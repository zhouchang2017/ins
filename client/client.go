package client

import (
	"encoding/json"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/proxy"
	"ins/types"
	"log"
	"net/http"
	"strings"
)

var client *Client

type Client struct {
	*colly.Collector
	debug     bool
	proxy     string
	logged    bool
	csrftoken string
	username  string
	password  string
}

type ClientConfig map[string]interface{}

func NewClient(debug bool, proxy string) *Client {

	if client != nil {
		return client
	}

	c := &Client{}

	c.debug = debug

	if proxy != "" {
		c.proxy = proxy
	}

	c.init()

	client = c

	return c
}

// 开启Debug模式
func (c *Client) Debug() {
	c.debug = true
}

// 设置Proxy
func (c *Client) SetProxy(url string) {
	c.proxy = url
}

func (c *Client) init() {

	c.Collector = colly.NewCollector(
		//colly.CacheDir("./_instagram_cache/"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)



	if c.debug {
		c.Collector.SetDebugger(&debug.LogDebugger{})
	}

	if c.proxy != "" {

		rp, err := proxy.RoundRobinProxySwitcher(c.proxy)
		if err != nil {
			log.Fatal(err)
		}
		c.SetProxyFunc(rp)
	}

	c.OnError(func(response *colly.Response, e error) {
		log.Printf("response err, url->%s , msg->%s \n", response.Request.URL.String(), e.Error())
	})

}

func (c *Client) getCsrftoken() {
	// 获取csrftoken
	c.OnHTML("html", func(e *colly.HTMLElement) {
		dat := e.ChildText("body > script:first-of-type")

		jsonData := dat[strings.Index(dat, "{") : len(dat)-1]

		data := &types.MainPageData{}
		err := json.Unmarshal([]byte(jsonData), data)
		if err != nil {
			log.Fatal(err)
		}

		c.csrftoken = data.Config.CsrfToken
	})
	c.Visit("https://www.instagram.com/accounts/login/")
}

func (c *Client) doLogin() {
	c.getCsrftoken()

	d := c.Clone()



	d.OnRequest(func(request *colly.Request) {
		if c.csrftoken != "" {
			request.Headers.Set("X-Requested-With", "XMLHttpRequest")
			request.Headers.Set("X-Instagram-AJAX", "1")
			if request.Method == http.MethodPost {
				request.Headers.Set("X-CSRFToken", c.csrftoken)
			}
		}
	})

	d.Post(
		"https://www.instagram.com/accounts/login/ajax/",
		map[string]string{
			"username":      c.username,
			"password":      c.password,
			"queryParams":   `{"source":"auth_switcher"}`,
			"optIntoOneTap": "false"})

	c.Collector = d
}

func (c *Client) Login(username, password string) {
	c.username = username
	c.password = password
	c.doLogin()
}
