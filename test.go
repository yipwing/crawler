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

type ReceiveData struct {
	Ret            int    `json:"ret"`
	Errmsg         string `json:"errmsg"`
	MsgCount       int    `json:"msg_count"`
	CanMsgContinue int    `json:"can_msg_continue"`
	GeneralMsgList string `json:"general_msg_list"`
	NextOffset     int    `json:"next_offset"`
	VideoCount     int    `json:"video_count"`
	UseVideoTab    int    `json:"use_video_tab"`
	RealType       int    `json:"real_type"`
}
type GeneralMsgList struct {
	List []struct {
		AppMsgExtInfo struct {
			AudioFileid            int    `json:"audio_fileid"`
			Author                 string `json:"author"`
			Content                string `json:"content"`
			ContentURL             string `json:"content_url"`
			CopyrightStat          int    `json:"copyright_stat"`
			Cover                  string `json:"cover"`
			DelFlag                int    `json:"del_flag"`
			Digest                 string `json:"digest"`
			Duration               int    `json:"duration"`
			Fileid                 int    `json:"fileid"`
			IsMulti                int    `json:"is_multi"`
			ItemShowType           int    `json:"item_show_type"`
			MaliciousContentType   int    `json:"malicious_content_type"`
			MaliciousTitleReasonID int    `json:"malicious_title_reason_id"`
			MultiAppMsgItemList    []struct {
				AudioFileid            int    `json:"audio_fileid"`
				Author                 string `json:"author"`
				Content                string `json:"content"`
				ContentURL             string `json:"content_url"`
				CopyrightStat          int    `json:"copyright_stat"`
				Cover                  string `json:"cover"`
				DelFlag                int    `json:"del_flag"`
				Digest                 string `json:"digest"`
				Duration               int    `json:"duration"`
				Fileid                 int    `json:"fileid"`
				ItemShowType           int    `json:"item_show_type"`
				MaliciousContentType   int    `json:"malicious_content_type"`
				MaliciousTitleReasonID int    `json:"malicious_title_reason_id"`
				PlayURL                string `json:"play_url"`
				SourceURL              string `json:"source_url"`
				Title                  string `json:"title"`
			} `json:"multi_app_msg_item_list"`
			PlayURL   string `json:"play_url"`
			SourceURL string `json:"source_url"`
			Subtype   int    `json:"subtype"`
			Title     string `json:"title"`
		} `json:"app_msg_ext_info"`
		CommMsgInfo struct {
			Content  string `json:"content"`
			Datetime int64  `json:"datetime"`
			Fakeid   string `json:"fakeid"`
			ID       int64  `json:"id"`
			Status   int    `json:"status"`
			Type     int    `json:"type"`
		} `json:"comm_msg_info"`
	} `json:"list"`
}

// func ReadArticleAndTitle(url string) (name string, article string, articletitle string) {
// 	resp, body, respErr := gorequest.New().Get(url).End()
// 	if respErr != nil || resp.StatusCode != 200 {
// 		panic(respErr)
// 	}
// 	defer resp.Body.Close()
// 	doc, reqErr := goquery.NewDocumentFromReader(strings.NewReader(body))
// 	if reqErr != nil {
// 		panic(reqErr)
// 	}
// 	PublicName := doc.Find("#meta_content").Find("#profileBt").Find("#js_name").Text() // 公众号名称
// 	re, _ := regexp.Compile("\\s+|\n")
// 	pn := re.ReplaceAllString(PublicName, "")
// 	ArticleString, _ := doc.Find(".rich_media_content").Html()
// 	atitle := doc.Find(".rich_media_title").Text()
// 	title := re.ReplaceAllString(atitle, "")
// 	fmt.Print(title)
// 	return pn, ArticleString, title
// }

// const DATEFORMAT = "2006-01-02"

// JSONStruct 配置文件
type JSONStruct struct {
	Biz     string
	Uin     string
	PubName string
}

func main() {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	startDay := thisMonth.AddDate(0, -1, 0)
	h := time.Now()
	fmt.Println(startDay.Month() < h.Month())
	// for _, row := range rows {
	// 	fmt.Println(len(row))
	// 	for _, col := range row {
	// 		fmt.Println(len(col))
	// 	}
	// }
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
	// gml := &GeneralMsgList{}

	// level1Err := json.Unmarshal([]byte(rec.GeneralMsgList), &gml)
	// if level1Err != nil {
	// 	panic(level1Err)
	// }
	// for _, First := range gml.List {
	// 	contentURL, _ := url.Parse(First.AppMsgExtInfo.ContentURL)
	// 	contentQuery := &url.Values{}
	// 	for key, value := range contentURL.Query() {
	// 		switch key {
	// 		case "action":
	// 			break
	// 		case "lang":
	// 			break
	// 		case "winzoom":
	// 			break
	// 		case "a8scene":
	// 			break
	// 		case "version":
	// 			break
	// 		case "scene":
	// 			break
	// 		case "devicetype":
	// 			break
	// 		case "chksm":
	// 			break
	// 		case "amp":
	// 			break
	// 		default:
	// 			contentQuery.Add(key, value[0])
	// 		}
	// 	}
	// 	fmt.Println(contentQuery)
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
	// f, _ := time.ParseInLocation("2006年01月02日", "2018年08月29日", time.Local)

}
