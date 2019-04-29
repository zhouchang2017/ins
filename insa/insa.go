package insa

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
)

const (
	urlMsisdnHeader = "accounts/read_msisdn_header/"
	urlLogin        = "accounts/login/ajax/"
	urlLogout       = "accounts/logout/"
)

type Instagram struct {
	user string
	pass string
	// device id: android-1923fjnma8123
	dID string
	// uuid: 8493-1233-4312312-5123
	uuid string
	// rankToken
	rankToken string
	// token
	token string
	// phone id
	pid string
	// ads id
	adid string

	c *http.Client
}

// New creates Instagram structure
func New(username, password string) *Instagram {
	// this call never returns error
	jar, _ := cookiejar.New(nil)
	inst := &Instagram{
		user: username,
		pass: password,
		dID: generateDeviceID(
			generateMD5Hash(username + password),
		),
		uuid: generateUUID(), // both uuid must be differents
		pid:  generateUUID(),
		c: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			Jar: jar,
		},
	}
	inst.init()

	return inst
}

func (inst *Instagram) init() {

}

// SetProxy sets proxy for connection.
func (inst *Instagram) SetProxy(url string, insecure bool) error {
	uri, err := neturl.Parse(url)
	if err == nil {
		inst.c.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		}
	}
	return err
}

// UnsetProxy unsets proxy for connection.
func (inst *Instagram) UnsetProxy() {
	inst.c.Transport = nil
}

func (inst *Instagram) readMsisdnHeader() error {
	data, err := json.Marshal(
		map[string]string{
			"device_id": inst.uuid,
		},
	)
	if err != nil {
		return err
	}
	_, err = inst.sendRequest(
		&reqOptions{
			Endpoint:   urlMsisdnHeader,
			IsPost:     true,
			Connection: "keep-alive",
			Query:      generateSignature(b2s(data)),
		},
	)
	return err
}

// Login performs instagram login.
//
// Password will be deleted after login
func (inst *Instagram) Login() error {
	err := inst.readMsisdnHeader()
	//if err != nil {
	//	return err
	//}

	//err = inst.syncFeatures()
	//if err != nil {
	//	return err
	//}
	//
	//err = inst.zrToken()
	//if err != nil {
	//	return err
	//}
	//
	//err = inst.sendAdID()
	//if err != nil {
	//	return err
	//}
	//
	//err = inst.contactPrefill()
	//if err != nil {
	//	return err
	//}

	result, err := json.Marshal(
		map[string]interface{}{
			"guid":                inst.uuid,
			"login_attempt_count": 0,
			"_csrftoken":          inst.token,
			"device_id":           inst.dID,
			"adid":                inst.adid,
			"phone_id":            inst.pid,
			"username":            inst.user,
			"password":            inst.pass,
			"google_tokens":       "[]",
		},
	)
	if err != nil {
		return err
	}


	body, err := inst.sendRequest(
		&reqOptions{
			Endpoint: urlLogin,
			Query:    generateSignature(b2s(result)),
			IsPost:   true,
			Login:    true,
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s", body)
	inst.pass = ""

	// getting account data
	//res := accountResp{}
	//err = json.Unmarshal(body, &res)
	//if err != nil {
	//	return err
	//}
	//
	//inst.Account = &res.Account
	//inst.Account.inst = inst
	//inst.rankToken = strconv.FormatInt(inst.Account.ID, 10) + "_" + inst.uuid
	//inst.zrToken()

	return err
}

// Logout closes current session
func (inst *Instagram) Logout() error {
	_, err := inst.sendSimpleRequest(urlLogout)
	inst.c.Jar = nil
	inst.c = nil
	return err
}
