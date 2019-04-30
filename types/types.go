package types

type PageInfo struct {
	EndCursor string `json:"end_cursor"`
	NextPage  bool   `json:"has_next_page"`
}


type MainPageData struct {
	Config struct {
		CsrfToken string `json:"csrf_token"`
	} `json:"config"`
	Rhxgis    string `json:"rhx_gis"`
	EntryData struct {
		ProfilePage []struct {
			Graphql struct {
				User struct {
					Biography      string `json:"biography"`
					Id             string `json:"id"` // id
					EdgeFollowedBy struct {
						Count int `json:"count"` // 粉丝数
					} `json:"edge_followed_by"`
					EdgeFollow struct {
						Count int `json:"count"`
					} `json:"edge_follow"`
					ProfilePicUrlHd      string `json:"profile_pic_url_hd"` // 头像
					FullName             string `json:"full_name"`
					Username             string `json:"username"` //用户名
					BusinessCategoryName string `json:"business_category_name"`
					Media                struct {
						Edges []struct {
							Node struct {
								ImageURL     string `json:"display_url"`
								ThumbnailURL string `json:"thumbnail_src"`
								IsVideo      bool   `json:"is_video"`
								Date         int    `json:"date"`
								Dimensions   struct {
									Width  int `json:"width"`
									Height int `json:"height"`
								} `json:"dimensions"`
							} `json::node"`
						} `json:"edges"`
						PageInfo PageInfo `json:"page_info"`
					} `json:"edge_owner_to_timeline_media"`
				} `json:"user"`
			} `json:"graphql"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
}

type NextPageData struct {
	Data struct {
		User struct {
			Container struct {
				PageInfo PageInfo `json:"page_info"`
				Edges    []struct {
					Node struct {
						ImageURL     string `json:"display_url"`
						ThumbnailURL string `json:"thumbnail_src"`
						IsVideo      bool   `json:"is_video"`
						Date         int    `json:"taken_at_timestamp"`
						Dimensions   struct {
							Width  int `json:"width"`
							Height int `json:"height"`
						}
					}
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
		}
	} `json:"data"`
}

type GraphqlResponse struct {
	Data struct{
		User struct{
			EdgeChaining struct{
				Edges []struct{
					Node struct{
						Id string `json:"id"`
						FullName string `json:"full_name"`
						Username string `json:"username"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_chaining"`
		} `json:"user"`
	} `json:"data"`
}