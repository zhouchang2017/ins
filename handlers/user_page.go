package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"ins/client"
	"ins/types"
	"log"
	url1 "net/url"
	"regexp"
	"strings"
)

var queryIdPattern = regexp.MustCompile(`="([0-9a-zA-Z]{32})"`)

const graphqlEndpoint string = `https://www.instagram.com/graphql/query/?query_hash=%s&variables=%s`
const userProfilePayload string = `{"user_id":"%s","include_chaining":true}`

type userPage struct {
	query_hash string
}

func NewUserPage() *userPage {
	return &userPage{}
}

func (u *userPage) Read(client *client.Client, url string) {



	client.OnHTML("html", func(element *colly.HTMLElement) {

		f := client.Clone()
		f.OnResponse(func(r *colly.Response) {
			//submatch := queryIdPattern.FindAllSubmatch(r.Body, -1)
			submatch := queryIdPattern.FindSubmatch(r.Body)
			u.query_hash = string(submatch[1])
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
			log.Print("未找到ProfilePageContainer.js URL")
		}

		requestIDURL := element.Request.AbsoluteURL(ProfilePageContainer)

		f.Visit(requestIDURL)

		u.graphHandle(client, element)

		//go func() {
		//	dat := element.ChildText("body > script:first-of-type")
		//	jsonData := dat[strings.Index(dat, "{") : len(dat)-1]
		//	data := &types.MainPageData{}
		//	err := json.Unmarshal([]byte(jsonData), data)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//
		//	currenUser := data.EntryData.ProfilePage[0].Graphql.User
		//
		//	u.Write(currenUser)
		//
		//	currenUserId := currenUser.Id
		//
		//	userProfileVars := fmt.Sprintf(userProfilePayload, currenUserId)
		//
		//	uri := fmt.Sprintf(
		//		graphqlEndpoint,
		//		u.query_hash,
		//		url1.QueryEscape(userProfileVars),
		//	)
		//
		//	graphHandler := NewUserGraph(data.Rhxgis, userProfileVars, currenUser.Username)
		//
		//	task := task.NewTask(client, uri, graphHandler)
		//	go task.Run()
		//}()

	})

	client.Visit(url)

}

var count = 0

func (u *userPage) Write(data *types.MainPageData) {
	count++
	currenUser := data.EntryData.ProfilePage[0].Graphql.User
	log.Printf("已采集 %d 	,username:%s	,ins_id:%s  ,followers %d\n", count, currenUser.Username, currenUser.Id, currenUser.EdgeFollowedBy.Count)
}

func (u *userPage) graphHandle(client *client.Client, element *colly.HTMLElement) {
	dat := element.ChildText("body > script:first-of-type")
	jsonData := dat[strings.Index(dat, "{") : len(dat)-1]
	data := &types.MainPageData{}
	err := json.Unmarshal([]byte(jsonData), data)
	if err != nil {
		log.Fatal(err)
	}

	u.Write(data)

	currenUser := data.EntryData.ProfilePage[0].Graphql.User

	currenUserId := currenUser.Id

	userProfileVars := fmt.Sprintf(userProfilePayload, currenUserId)

	uri := fmt.Sprintf(
		graphqlEndpoint,
		u.query_hash,
		url1.QueryEscape(userProfileVars),
	)

	gqClient := client.Clone()

	gqClient.OnRequest(func(request *colly.Request) {
		request.Headers.Set("X-Requested-With", "XMLHttpRequest")
		request.Headers.Set("Referrer", "https://www.instagram.com/"+currenUser.Username)
		if userProfileVars != "" && data.Rhxgis != "" {
			request.Headers.Set("X-Instagram-GIS", u.genGis(data.Rhxgis, userProfileVars))
		}
	})

	gqClient.OnResponse(func(response *colly.Response) {
		if strings.Contains(response.Headers.Get("Content-Type"), "application/json") {
			// u.Write(response.Body)

			data := &types.GraphqlResponse{}
			err := json.Unmarshal(response.Body, data)
			if err != nil {
				log.Fatal(err)
			}

			for _, node := range data.Data.User.EdgeChaining.Edges {
				uri := fmt.Sprintf("https://www.instagram.com/%s/", node.Node.Username)
				client.Visit(uri)
			}

		}
	})

	gqClient.Visit(uri)
}

func (u *userPage) genGis(gis, variables string) (gisHash string) {
	giskey := fmt.Sprintf("%s:%s", gis, variables)
	h := md5.New()
	h.Write([]byte(giskey))
	gisHash = fmt.Sprintf("%x", h.Sum(nil))
	return
}
