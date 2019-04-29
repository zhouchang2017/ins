package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/proxy"
	"ins/types"
	"log"
	"net/url"
	"regexp"
	"strings"
)

var csrftoken string

var csrfTokenReg = regexp.MustCompile(`csrf_token":"([a-zA-Z0-9]+)"`)

var queryIdPattern = regexp.MustCompile(`queryId:".{32}"`)
var requestID string
var requestIds [][]byte


const nextPageURL string = `https://www.instagram.com/graphql/query/?query_hash=%s&variables=%s`
const userProfilePayload string = `{"user_id":"%s","include_chaining":true,"include_reel":true,"include_suggested_users":false,"include_logged_out_extras":false,"include_highlight_reels":false}`



func main() {

	// create a new collector
	c := NewColly()


	c.OnHTML("html", func(e *colly.HTMLElement) {
		dat := e.ChildText("body > script:first-of-type")

		jsonData := dat[strings.Index(dat, "{") : len(dat)-1]

		data := &types.MainPageData{}
		err := json.Unmarshal([]byte(jsonData), data)
		if err != nil {
			log.Fatal(err)
		}

		submatch := csrfTokenReg.FindStringSubmatch(dat)

		if len(submatch) < 1 {
			log.Printf("reg csrftoken notfond!, submatch: %v", submatch)
		}
		csrftoken = data.Config.CsrfToken

		e.Request.Ctx.Put("variables", data.Rhxgis)

	})

	// start scraping
	c.Visit("https://www.instagram.com/accounts/login/")

	d := c.Clone()

	d.OnResponse(func(response *colly.Response) {
		cookies := d.Cookies("https://www.instagram.com")
		fmt.Printf("post login OnRequest Cookies :%v\n", cookies)
		fmt.Printf("post login OnRequest Headers :%s\n", response.Headers)

	})

	d.OnRequest(func(request *colly.Request) {
		if csrftoken != "" {
			request.Headers.Set("X-Requested-With", "XMLHttpRequest")
			request.Headers.Set("X-Instagram-AJAX", "1")
			request.Headers.Set("X-CSRFToken", csrftoken)
		}
		if request.Ctx.Get("gis") != "" {
			gis := fmt.Sprintf("%s:%s", request.Ctx.Get("gis"), request.Ctx.Get("variables"))
			h := md5.New()
			h.Write([]byte(gis))
			gisHash := fmt.Sprintf("%x", h.Sum(nil))
			request.Headers.Set("X-Instagram-GIS", gisHash)
		}

	})

	d.Post("https://www.instagram.com/accounts/login/ajax/", map[string]string{"username": "avvlover", "password": "qazwsx123", "queryParams": `{"source":"auth_switcher"}`, "optIntoOneTap": "false"})

	e := d.Clone()


	e.OnHTML("html", func(element *colly.HTMLElement) {
		f:=e.Clone()
		f.OnResponse(func(r *colly.Response) {

			requestIds = queryIdPattern.FindAll(r.Body, -1)
			requestID = string(requestIds[len(requestIds)-1][9:41])
			fmt.Printf("%s\n",requestID)
		})

		scripts := element.ChildAttrs(`link[as="script"]`, "href")

		var ProfilePageContainer string

		for _, v := range scripts {
			if strings.Contains(v, "ProfilePageContainer") {
				ProfilePageContainer = v
				break
			}
		}

		if ProfilePageContainer == "" {
			panic("未找到ProfilePageContainer.js URL")
		}

		requestIDURL := element.Request.AbsoluteURL(ProfilePageContainer)

		f.Visit(requestIDURL)

		dat := element.ChildText("body > script:first-of-type")
		jsonData := dat[strings.Index(dat, "{") : len(dat)-1]
		data := &types.MainPageData{}
		err := json.Unmarshal([]byte(jsonData), data)
		if err != nil {
			log.Fatal(err)
		}


		currenUserId:=data.EntryData.ProfilePage[0].Graphql.User.Id
		userProfileVars := fmt.Sprintf(userProfilePayload, currenUserId)
		element.Request.Ctx.Put("variables", userProfileVars)
		element.Request.Ctx.Put("gis", data.Rhxgis)
		u := fmt.Sprintf(
			nextPageURL,
			requestID,
			url.QueryEscape(userProfileVars),
		)
		element.Request.Visit(u)

	})
	
	e.OnRequest(func(r *colly.Request) {
		log.Printf("request url=%s",r.URL.String())
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		//r.Headers.Set("Referrer", "https://www.instagram.com/"+instagramAccount)
		if r.Ctx.Get("gis") != "" {
			gis := fmt.Sprintf("%s:%s", r.Ctx.Get("gis"), r.Ctx.Get("variables"))
			h := md5.New()
			h.Write([]byte(gis))
			gisHash := fmt.Sprintf("%x", h.Sum(nil))
			log.Printf("set gisHash =%s",gisHash)
			r.Headers.Set("X-Instagram-GIS", gisHash)
		}
	})

	e.Visit("https://www.instagram.com/aj_duan/")
}

func NewColly() *colly.Collector {
	c := colly.NewCollector(
		//colly.CacheDir("./_instagram_cache/"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.Debugger(&debug.LogDebugger{}),
	)

	// Rotate two socks5 proxies
	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1086")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)
	return c
}

