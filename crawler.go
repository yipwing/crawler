package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PuerkitoBio/goquery"

	"github.com/parnurzeal/gorequest"
)

// var 变量定义
var (
	StartDay string
)

// const 常量定义
const (
	BaseURL      string = "https://mp.weixin.qq.com/mp/profile_ext"
	BaseDaySec   int64  = 86400
	ReadCountURL string = "https://mp.weixin.qq.com/mp/getappmsgext"
	UserAgent    string = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 MicroMessenger/6.5.2.501 NetType/WIFI WindowsWechat QBCore/3.43.901.400 QQBrowser/9.0.2524.400"
	DJGZ         string = "动静贵州"
	GZXWLB       string = "贵州新闻联播"
)

// PostData 提交给微信服务端的数据结构
type PostData struct {
	Biz        string `json:"__biz"`
	Uin        string `json:"uin"`
	Offset     int    `json:"offset"`
	Count      int    `json:"count"`
	Key        string `json:"key"`
	PassTicket string `json:"pass_ticket"`
	Format     string `json:"f"`
	IsOk       int    `json:"is_ok"`
	X5         int    `json:"x5"`
	Devicetype string `json:"devicetype"`
	// AppmsgToken string `json:appmsg_token`  // 这个值估计会有什么作用，暂时用不上，固注释掉
}

// ReadingGet 用于提交时的参数
type ReadingGet struct {
	UIN         string
	KEY         string
	PassTicket  string
	AppMsgToken string
}

//	ReadingPost 提交给微信端的参数
type ReadingPost struct {
	Biz        string `json:"__biz"`
	MID        string `json:"mid"`
	IDx        string `json:"idx"`
	SN         string `json:"sn"`
	PassTicket string `json:"pass_ticket"`
	ISOnlyRead int    `json:"is_only_read"`
}

// ReceiveData Json数据接收内部有个list的json
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

// GeneralMsgList list数据
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

// ExcelData 用于存储要用来保存的数据
type ExcelData struct {
	PublicName     string
	Date           string
	ReadCound      int
	LikeCount      int
	CommentCount   int
	ArticleTitle   string
	ArticleContent string
}

// ExcelTitle excel头文本
type ExcelTitle []struct {
	Axis string
	Name string
}

// ReadArticleAndTitle 获取文章内容与文章标题，还有公众号名称
func ReadArticleAndTitle(url string) (name string, article string, articletitle string) {
	resp, body, respErr := gorequest.New().Get(url).End()
	if respErr != nil || resp.StatusCode != 200 {
		panic(respErr)
	}
	defer resp.Body.Close()
	doc, reqErr := goquery.NewDocumentFromReader(strings.NewReader(body))
	if reqErr != nil {
		panic(reqErr)
	}
	PublicName := doc.Find("#meta_content").Find("#profileBt").Find("#js_name").Text() // 公众号名称
	re, _ := regexp.Compile("\\s+|\n")
	pn := re.ReplaceAllString(PublicName, "")
	ArticleString, _ := doc.Find(".rich_media_content").Html()
	atitle := doc.Find(".rich_media_title").Text()
	title := re.ReplaceAllString(atitle, "")
	fmt.Print(title)
	return pn, ArticleString, title
}

// 初始化全局变量
var (
	excelTitle = ExcelTitle{
		{
			"A1",
			"公众号名称",
		},
		{
			"B1",
			"日期",
		},
		{
			"C1",
			"阅读数",
		},
		{
			"D1",
			"点赞数",
		},
		{
			"E1",
			"评论数",
		},
		{
			"F1",
			"文章标题",
		},
		{
			"G1",
			"文章内容",
		},
	}
)

// 微信服务端json数据结构结束

func main() {
	// todo 获取微信主要数据以访问
	fmt.Println("注意，由于程序会获取当前系统时区作为判断。请确认你的时区是正确的")
	var pubname string
	p := &PostData{}
	fmt.Print("输入__biz参数值：")
	fmt.Scanln(&p.Biz)
	fmt.Print("输入uin参数值(先执行下decode)：")
	fmt.Scanln(&p.Uin)
	fmt.Print("输入key参数值：")
	fmt.Scanln(&p.Key)
	fmt.Print("输入pass_ticket参数值(先执行下decode)：")
	fmt.Scanln(&p.PassTicket)
	fmt.Print("输入要抓取的公众号名称：")
	fmt.Scanln(&pubname)
	p.IsOk = 1
	p.Offset = 0
	p.Count = 10
	p.X5 = 0
	p.Format = "json"
	Action := "getmsg"
	// 计算循环次数 todo
	CurrentTime := time.Now()
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	startDay := thisMonth.AddDate(0, -1, 0)
	endDay := thisMonth.AddDate(0, 0, -1)
	MonthOfDay := (endDay.Unix()-startDay.Unix())/86400 + 1

	// 初始化excel表头标题
	xlsx := excelize.NewFile()
	timeLayout := "2006-01-02"
	today := time.Now().Unix()
	for _, v := range excelTitle {
		xlsx.SetCellValue(pubname, v.Axis, v.Name)
	}
	ed := []ExcelData{}
	for index := 0; index <= int(MonthOfDay); index++ {
		request := gorequest.New()
		param := &url.Values{}
		param.Add("action", Action)
		param.Add("__biz", p.Biz)
		param.Add("uin", p.Uin)
		param.Add("key", p.Key)
		param.Add("pass_ticket", p.PassTicket)
		param.Add("is_ok", string(p.IsOk))
		param.Add("f", p.Format)
		param.Add("x5", string(p.X5))
		param.Add("offset", string(p.Offset))
		param.Add("count", string(p.Count))
		req, body, reqErr := request.Get(BaseURL).Query(param).Set("User-Agent", UserAgent).End()
		if reqErr != nil {
			panic(reqErr)
		}
		defer req.Body.Close()
		if req.StatusCode != 200 {
			panic("status code error")
		}
		// startday := time.Now().Format(StartDay)
		receive := &ReceiveData{}
		orgErr := json.Unmarshal([]byte(body), &receive)
		if orgErr != nil {
			panic(orgErr)
		}
		p.Offset = receive.NextOffset
		gml := &GeneralMsgList{}
		level1Err := json.Unmarshal([]byte(receive.GeneralMsgList), &gml)
		if level1Err != nil {
			panic(level1Err)
		}
		for _, First := range gml.List {
			pubDate := time.Unix(First.CommMsgInfo.Datetime, 0).Format(timeLayout)

			at := (today - First.CommMsgInfo.Datetime) / BaseDaySec

		}

		// url.QueryUnescape() // urldecode
		// todo 这里需要循环内部的list，list中还有list需要循环。
		// for i := 0; i < receive.GeneralMsgList.Array.Len(); i++ {
		// 	cmi = &CommMsgInfo{}
		// 	amei = &AppMsgExtInfo{}
		// 	mamil = &MultiAppMsgItemList{}
		// 	cmi = receive.GeneralMsgList.Array[i].comm_msg_info
		// 	// for j :=0; j <
		// }
	}

	// 微信端拉取数据这里要注意下一页方法是action=getmsg
	/*
		mp.weixin.qq.com/mp/profile_ext?action=getmsg&__biz=MzA5ODMyODYzMQ==&f=json&offset=10&count=10&is_ok=1&scene=124&uin=MjY0MjQ3MDM0MA==&key=d33469543f977294b762012380263097f662eeb79fec9b6e0fac1cf2420b209fb6a67247cb90c62fa759d5357b4e59bc53875bef9ebe2a8ee82dec17aad20b478bea67c6aac317a5d1c3ecbf3e48fb21&pass_ticket=yk%2FXz65ktUUZp%2BkksTLzJz2BXrGWndRSGzc6ViN%2BFBEWvS3XTsuJgL5s02BPDOb3&wxtoken=&appmsg_token=970_aAXg8wPyGe%252F46ggtzaMsyEdlI6gix6nHpfu-_w~~&x5=0&f=json

		mp.weixin.qq.com/mp/profile_ext?action=home&__biz=MzA3MTA1MDU3OA==&uin=MjY0MjQ3MDM0MA%3D%3D&key=b8264eebf6d6109a6b2b704a25a42fae44cf1117e611cf34d0db580e1428d56252c019368099da01a43d13069c41806f85ced315163a5483cdf9fc5bfc88cd2b08c711598d9dde1ffbfeed8e367b060a&devicetype=Windows+10&version=62060426&lang=zh_CN&a8scene=7&pass_ticket=yk%2FXz65ktUUZp%2BkksTLzJz2BXrGWndRSGzc6ViN%2BFBEWvS3XTsuJgL5s02BPDOb3&winzoom=1
	*/
	// 微信错误返回
	/*
		{
			"base_resp": {
					"ret": -3,
					"errmsg": "no session",
					"cookie_count": 0,
					"csp_nonce": 1916841621
			},
				"ret": -3,
				"errmsg": "no session",
				"cookie_count": 0
		}
	*/
	// doc.Find("div.list a.item").Each(func(i int, s *goquery.Selection) {
	// https://movie.douban.com/explore#!type=movie&sort=recommend&page_limit=20&page_start=0
	// 豆瓣这个地方的属于js请求范畴，故用此方法获取不到内容。 下面是js请求内容
	// https://movie.douban.com/j/search_subjects?type=movie&tag=%E7%83%AD%E9%97%A8&sort=recommend&page_limit=20&page_start=0
	// 	href, _ := s.Attr("href")
	// 	name := s.Find("p").Text()
	// 	fmt.Printf("number of %d\nname: %s\nhref: %s", i, name, href)
	// })
}
