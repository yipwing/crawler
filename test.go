package main

import (
	"fmt"
	"time"
)

// const var
const (
	ReqURL       string = "https://mp.weixin.qq.com/mp/profile_ext?action=home&__biz=MzA5ODMyODYzMQ==&scene=124&uin=MjY0MjQ3MDM0MA%3D%3D&key=da65dc85b3e2c9b9bd0a758c5330cb66b68a380fb0a71aeb1cd4139ee73e237dd35e3e4443cc2da85bcee0de352b81ab9dbe85850b13469af8896ec04192c7d58b929aae1e1a400e69af5130a949c934&devicetype=Windows+10&version=62060426&lang=zh_CN&a8scene=7&pass_ticket=sPvpyYM58mb3GFrrwD6D1P88sW6bbrM1BhWCe5jLaYeCZdVbwgm66g8cPWifWBI8&winzoom=1"
	UserAgent    string = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 MicroMessenger/6.5.2.501 NetType/WIFI WindowsWechat QBCore/3.43.901.400 QQBrowser/9.0.2524.400"
	ReadCountURL string = "https://mp.weixin.qq.com/mp/getappmsgext"
)

// type RequestParameter struct {
// 	Biz        string `json:"__biz"`
// 	MID        string `json:"mid"`
// 	IDx        string `json:"idx"`
// 	SN         string `json:"sn"`
// 	KEY        string `json:"key"`
// 	UIN        string `json:"uin"`
// 	PassTicket string `json:"pass_ticket"`
// 	ISOnlyRead int    `json:"is_only_read"`
// }
const DATEFORMAT = "2006-01-02"

func main() {
	// fileObj, _ := ioutil.ReadFile("origin.json")
	// rec := &ReceiveData{}

	// level0Err := json.Unmarshal(fileObj, &rec)
	// if level0Err != nil {
	// 	panic(level0Err)
	// }
	// dstFile, err := os.Create("output.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// defer dstFile.Close()
	// dstFile.WriteString(rec.GeneralMsgList)
	// rep := strings.Replace(rec.GeneralMsgList, "\\", "", -1)
	// gml := &GeneralMsgList{}

	// level1Err := json.Unmarshal([]byte(rec.GeneralMsgList), &gml)
	// if level1Err != nil {
	// 	panic(level1Err)
	// }

	// for _, v := range gml.List {
	// 	pubDate := time.Unix(v.CommMsgInfo.Datetime, 0).Format(timeLayout)
	// 	at := (today - v.CommMsgInfo.Datetime) / BaseDaySec
	// 	if at > 0 {
	// 		contentURL, _ := url.Parse(v.AppMsgExtInfo.ContentURL)
	// 		pn, article, title := ReadArticleAndTitle(v.AppMsgExtInfo.ContentURL)
	// 		count := 0
	// 		postURL := ReadCountURL
	// for key, value := range contentURL.Query() {
	// 	switch key {
	// 	case "action":
	// 		break
	// 	case "lang":
	// 		break
	// 	case "winzoom":
	// 		break
	// 	case "a8scene":
	// 		break
	// 	case "version":
	// 		break
	// 	case "scene":
	// 		break
	// 	case "devicetype":
	// 		break
	// 	default:
	// 		if count > 0 {
	// 			postURL += "&" + key + "=" + value[0]
	// 		} else {
	// 			postURL += "?" + key + "=" + value[0]
	// 			count++
	// 		}
	// 	}
	// }
	// 		reading := gorequest.New().Post(ReadCountURL).Set("User-Agent", UserAgent).Type("urlencoded").
	// 			Send()
	// 		ed = append(ed, ExcelData{pn, pubDate, 1, 1, 1, title, article})

	// 		if v.AppMsgExtInfo.IsMulti == 1 {
	// 			for _, mul := range v.AppMsgExtInfo.MultiAppMsgItemList {
	// 				_, subarticle, subtitle := ReadArticleAndTitle(mul.ContentURL)
	// 				strings.NewReader(subarticle)
	// 				fmt.Println(subtitle)
	// 			}
	// 		}
	// 		continue
	// 	}
	// }

	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start := thisMonth.AddDate(0, -1, 0)
	end := thisMonth.AddDate(0, 0, -1)
	MonthOfDay := (end.Unix()-start.Unix())/86400 + 1
	fmt.Println(MonthOfDay)
	fmt.Println(start)
	fmt.Println(end)
}
