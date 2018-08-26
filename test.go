package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
)

// const var
const (
	ReqURL       string = "https://mp.weixin.qq.com/s?__biz=MzA5ODMyODYzMQ==&mid=2650702544&idx=1&sn=19527bf20990e78686f8b9657dd52960&chksm=88996b88bfeee29eb4c5385f0c19c0291aa9d9c0ef49ab7f04d634cdef2a4e98acd4201d97b1&scene=38&key=d33469543f97729414a203c5808eea82a10294b23e6b4a5d60635646141ca55c0334a8f90930999d3a5be1330138e5884dc8d56aec8231b3f8ac6362143efcf05015eccfa01d656fb9e2c4e482a3450f&ascene=7&uin=MjY0MjQ3MDM0MA%3D%3D&devicetype=Windows+10&version=62060426&lang=zh_CN&pass_ticket=nVG1Vd7Q%2B5eDb9RAkmvA9Nx3WYqogZ7N3ixOwzhXOzpfoMDUPfjm5Zzz9cpAmeh0&winzoom=1"
	UserAgent    string = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 MicroMessenger/6.5.2.501 NetType/WIFI WindowsWechat QBCore/3.43.901.400 QQBrowser/9.0.2524.400"
	ReadCountURL string = "https://mp.weixin.qq.com/mp/getappmsgext"
)

// RequestParameter struct
type RequestParameter struct {
	Biz        string `json:"__biz"`
	MID        string `json:"mid"`
	IDx        string `json:"idx"`
	SN         string `json:"sn"`
	KEY        string `json:"key"`
	UIN        string `json:"uin"`
	PassTicket string `json:"pass_ticket"`
	ISOnlyRead int    `json:"is_only_read"`
}

func main() {

	// reading := gorequest.New()

	// URLParse, PraseErr := url.Parse(ReqURL)
	// if PraseErr != nil {
	// 	panic(PraseErr)
	// }
	// query := URLParse.Query()
	// param := &RequestParameter{}
	// param.Biz = query["__biz"][0]
	// param.IDx = query["idx"][0]
	// param.KEY = query["key"][0]
	// param.MID = query["mid"][0]
	// param.SN = query["sn"][0]
	// param.PassTicket = query["pass_ticket"][0]
	// param.UIN = query["uin"][0]
	// fmt.Println(URLParse.RawQuery)
	// return
	// // resp, respErr := http.Get(ReqURL)
	// // if respErr != nil {
	// // 	panic(respErr)
	// // }
	// // readHandle := &http.Client{}
	// // reading, readErr := http.NewRequest("POST", ReadCountURL, )
	// // reading.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// urlParam := &url.Values{}
	// urlParam.Add("__biz", query["__biz"][0])
	// urlParam.Add("idx", query["idx"][0])
	// urlParam.Add("key", query["key"][0])
	// urlParam.Add("mid", query["mid"][0])
	// urlParam.Add("sn", query["sn"][0])
	// urlParam.Add("pass_ticket", query["pass_ticket"][0])
	// urlParam.Add("uin", query["uin"][0])
	// fmt.Print(urlParam)
	// reading.Type("urlencoded")
	// // reading.Post(ReadCountURL).Type("urlencoded")
	// // reading.QueryData.Add("__biz", query["__biz"][0])
	// // reading.QueryData.Add("idx", query["idx"][0])
	// // reading.QueryData.Add("key", query["key"][0])
	// // reading.QueryData.Add("mid", query["mid"][0])
	// // reading.QueryData.Add("sn", query["sn"][0])
	// // reading.QueryData.Add("pass_ticket", query["pass_ticket"][0])
	// // reading.QueryData.Add("uin", query["uin"][0])
	// // fmt.Print(reading.QueryData)
	// return

	resp, body, respErr := gorequest.New().Get(ReqURL).End()
	if respErr != nil || resp.StatusCode != 200 {
		panic(respErr)
	}
	doc, docErr := goquery.NewDocumentFromReader(strings.NewReader(body))
	if docErr != nil {
		panic(docErr)
	}
	PublicName := doc.Find("#meta_content").Find("#profileBt").Find("#js_name").Text()
	re, _ := regexp.Compile("\\s+|\n")
	name := re.ReplaceAllString(PublicName, "")
	fmt.Println(name)
	ArticleString, _ := doc.Find(".rich_media_content").Html()
	fmt.Println(ArticleString)
}
