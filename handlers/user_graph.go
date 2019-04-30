package handlers

import (
	"crypto/md5"
	"fmt"
	"github.com/gocolly/colly"
	"ins/client"
)

const instagramUrl = `https://www.instagram.com/%s/`

type userGraph struct {
	gis              string
	variables        string
	instagramAccount string
}

func NewUserGraph(gis, variables, instagramAccount string) *userGraph {
	return &userGraph{
		gis:              gis,
		variables:        variables,
		instagramAccount: instagramAccount,
	}
}

func (u *userGraph) genGis() (gisHash string) {
	gis := fmt.Sprintf("%s:%s", u.gis, u.variables)
	h := md5.New()
	h.Write([]byte(gis))
	gisHash = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func (u *userGraph) Read(client *client.Client, url string) {

	client.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("Referrer", "https://www.instagram.com/"+u.instagramAccount)
		if u.variables != "" && u.gis != "" {
			r.Headers.Set("X-Instagram-GIS", u.genGis())
		}
	})

	client.OnResponse(func(response *colly.Response) {
		//if strings.Contains(response.Headers.Get("Content-Type"), "application/json") {
		//	u.Write(response.Body)
		//
		//	data := &types.GraphqlResponse{}
		//	err := json.Unmarshal(response.Body, data)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//
		//	for _,node:=range data.Data.User.EdgeChaining.Edges{
		//
		//		uri:=fmt.Sprintf(instagramUrl,node.Node.UserName)
		//
		//		task:= task.NewTask(client, uri, NewUserPage())
		//
		//		go task.Run()
		//	}
		//
		//}
	})

	client.Visit(url)
}

func (u *userGraph) Write(data interface{}) {
	fmt.Printf("response data %s", data)
}
