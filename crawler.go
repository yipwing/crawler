package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/gommon/log"
	"github.com/parnurzeal/gorequest"
)

// var 变量定义
var (
	StartDay   string
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

// const 常量定义
const (
	BaseURL      string = "https://mp.weixin.qq.com/mp/profile_ext"
	BaseDaySec   int64  = 86400
	ReadCountURL string = "https://mp.weixin.qq.com/mp/getappmsgext"
	UserAgent    string = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 MicroMessenger/6.5.2.501 NetType/WIFI WindowsWechat QBCore/3.43.901.400 QQBrowser/9.0.2524.400"
	DJGZ         string = "动静贵州"
	GZXWLB       string = "贵州新闻联播"
	PostType     string = "application/x-www-form-urlencoded"
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

// MsgStatus 文章内容评论数、点赞数与阅读数
type MsgStatus struct {
	AdvertisementInfo []interface{} `json:"advertisement_info"`
	AdvertisementNum  int           `json:"advertisement_num"`
	Appmsgstat        struct {
		IsLogin     bool `json:"is_login"`
		LikeNum     int  `json:"like_num"`
		Liked       bool `json:"liked"`
		ReadNum     int  `json:"read_num"`
		RealReadNum int  `json:"real_read_num"`
		Ret         int  `json:"ret"`
		Show        bool `json:"show"`
	} `json:"appmsgstat"`
	BaseResp struct {
		Wxtoken int `json:"wxtoken"`
	} `json:"base_resp"`
	CommentCount         int           `json:"comment_count"`
	CommentEnabled       int           `json:"comment_enabled"`
	FriendCommentEnabled int           `json:"friend_comment_enabled"`
	IsFans               int           `json:"is_fans"`
	LogoURL              string        `json:"logo_url"`
	NickName             string        `json:"nick_name"`
	OnlyFansCanComment   bool          `json:"only_fans_can_comment"`
	RewardHeadImgs       []interface{} `json:"reward_head_imgs"`
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

// var log = logging.MustGetLogger("crawler")

// var format = logging.MustStringFormatter(
// 	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
// )

// 微信服务端json数据结构结束

func main() {
	// todo 获取微信主要数据以访问
	logFile, fErr := os.OpenFile("logger.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if fErr != nil {
		fmt.Println("无法打开或创建文件\"logging.log\"")
		panic(fErr)
	}
	defer logFile.Close()
	fmt.Println("注意：由于程序会获取当前系统时区作为判断。请确认你的时区是正确的")
	fmt.Println("注意：程序会直接获取上个月的微信公众号内容作为任务目标，为期一个月")
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
	// CurrentTime := time.Now()
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	startDay := thisMonth.AddDate(0, -1, 0)
	endDay := thisMonth.AddDate(0, 0, -1)
	MonthOfDay := (endDay.Unix()-startDay.Unix())/86400 + 1

	// 初始化excel表头标题
	xlsx := excelize.NewFile()
	timeLayout := "2006-01-02"
	currentTime := time.Now().Format(timeLayout)
	for _, v := range excelTitle {
		xlsx.SetCellValue(pubname, v.Axis, v.Name)
	}
	ed := []ExcelData{}
	index := 0
	timer := 0
	for index <= p.Offset {
		client := &http.Client{}
		URL, _ := url.Parse(BaseURL)
		param := &url.Values{}
		param.Add("action", Action)
		param.Add("__biz", p.Biz)
		param.Add("uin", p.Uin)
		param.Add("key", p.Key)
		param.Add("pass_ticket", p.PassTicket)
		// param.Add("is_ok", string(p.IsOk))
		param.Add("f", p.Format)
		// param.Add("x5", string(p.X5))
		param.Add("offset", strconv.Itoa(p.Offset))
		param.Add("count", strconv.Itoa(p.Count))
		URL.RawQuery = param.Encode()
		req, reqErr := http.NewRequest("GET", URL.String(), nil)
		req.Header.Set("User-Agent", UserAgent)
		if reqErr != nil {
			panic(reqErr)
		}
		resp, respErr := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		if respErr != nil {
			log.Debugf("Get Request : %s", body)
			logFile.Write([]byte(body))
			panic(respErr)
		}
		// req, body, reqErr := request.Get(BaseURL).Query(param).Set("User-Agent", UserAgent).End()
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			panic("status code error")
		}
		receive := &ReceiveData{}
		orgErr := json.Unmarshal([]byte(body), &receive)
		if orgErr != nil {
			// buffer := []byte(strings.Join([]string{"json transform error", orgErr.Error(), "\n"}, ""))
			log.Debugf("json transform error: %s", orgErr.Error())
			logFile.Write([]byte(body))
			panic(orgErr)
		}
		gml := &GeneralMsgList{}
		level1Err := json.Unmarshal([]byte(receive.GeneralMsgList), &gml)
		if level1Err != nil {
			// buffer := []byte(strings.Join([]string{"json sub list transform error", orgErr.Error(), "\n"}, ""))
			log.Debugf("json sub list transform error : %s", level1Err.Error())
			logFile.Write([]byte(body))
			panic(level1Err)
		}
		articleTime, _ := time.ParseInLocation(timeLayout, time.Unix(gml.List[index].CommMsgInfo.Datetime, 0).Format("2006-01-02"), time.Local)
		// todo 这里需要判断时间 并不清楚具体情况
		if endDay.Unix()-articleTime.Unix() > BaseDaySec && timer == 0 {
			fmt.Printf("上月最后一天: %s\n", endDay.String())
			fmt.Printf("文章最后一天: %s\n", articleTime.String())
			fmt.Printf("timer : %s\n", timer)
			panic("文章时间与本地获取时间不符。")
		}
		if startDay.Unix()-articleTime.Unix() < BaseDaySec && timer == int(MonthOfDay) {
			break
		} else {
			p.Offset = receive.NextOffset
		}
		fmt.Printf("loop: %d times\n", timer+1)
		sublink := &http.Client{}
		for _, First := range gml.List {
			pubDate := time.Unix(First.CommMsgInfo.Datetime, 0).Format(timeLayout)
			contentURL, _ := url.Parse(First.AppMsgExtInfo.ContentURL)
			pn, article, title := ReadArticleAndTitle(First.AppMsgExtInfo.ContentURL)
			postForm := &url.Values{}
			for key, value := range contentURL.Query() {
				switch key {
				case "action":
					break
				case "lang":
					break
				case "winzoom":
					break
				case "a8scene":
					break
				case "version":
					break
				case "scene":
					break
				case "devicetype":
					break
				case "chksm":
					break
				case "amp":
					break
				default:
					postForm.Add(key, value[0])
				}
			}
			postForm.Add("pass_ticket", p.PassTicket)
			postForm.Add("is_only_read", "1")
			postForm.Add("msg_daily_idx", "1")
			req, postErr := http.NewRequest("POST", ReadCountURL, strings.NewReader(postForm.Encode()))
			if postErr != nil {
				panic(postErr)
			}
			getForm := req.URL.Query()
			getForm.Add("wxtoken", "777")
			getForm.Add("f", "json")
			getForm.Add("uin", p.Uin)
			getForm.Add("key", p.Key)
			getForm.Add("pass_ticket", p.PassTicket)
			req.URL.RawQuery = getForm.Encode()
			req.Header.Set("User-Agent", UserAgent)
			req.Header.Set("Content-Type", PostType)
			reading, readErr := sublink.Do(req)
			if readErr != nil {
				fmt.Printf("get parameter %s\npost parameter %s", getForm, postForm)
				panic("发送数据时发生异常")
			}
			defer reading.Body.Close()
			if reading.StatusCode != 200 {
				panic("http状态码返回错误")
			}
			body, _ := ioutil.ReadAll(reading.Body)
			msgStat := &MsgStatus{}
			msgErr := json.Unmarshal([]byte(body), msgStat)
			if msgErr != nil {
				panic("无法获取文章扩展数据")
			}
			fmt.Printf("loop out %s\n", msgStat)
			ed = append(ed, ExcelData{pn, pubDate, msgStat.Appmsgstat.ReadNum, msgStat.Appmsgstat.LikeNum, msgStat.CommentCount, title, article}) // 保存已获取的数据到数据结构中。完成所有循环后，写入到文件中
			if First.AppMsgExtInfo.IsMulti == 1 {
				for _, Sec := range First.AppMsgExtInfo.MultiAppMsgItemList {
					contentURL, _ := url.Parse(Sec.ContentURL)
					pn, article, title := ReadArticleAndTitle(Sec.ContentURL)
					postForm := &url.Values{}
					for key, value := range contentURL.Query() {
						switch key {
						case "action":
							break
						case "lang":
							break
						case "winzoom":
							break
						case "a8scene":
							break
						case "version":
							break
						case "scene":
							break
						case "devicetype":
							break
						case "chksm":
							break
						case "amp":
							break
						default:
							postForm.Add(key, value[0])
						}
					}
					postForm.Add("pass_ticket", p.PassTicket)
					postForm.Add("is_only_read", "1")
					postForm.Add("msg_daily_idx", "1")
					req, postErr := http.NewRequest("POST", ReadCountURL, strings.NewReader(postForm.Encode()))
					if postErr != nil {
						panic(postErr)
					}
					getForm := req.URL.Query()
					getForm.Add("wxtoken", "777")
					getForm.Add("f", "json")
					getForm.Add("uin", p.Uin)
					getForm.Add("key", p.Key)
					getForm.Add("pass_ticket", p.PassTicket)
					req.URL.RawQuery = getForm.Encode()
					req.Header.Set("User-Agent", UserAgent)
					req.Header.Set("Content-Type", PostType)
					reading, readErr := sublink.Do(req)
					if readErr != nil {
						fmt.Printf("get parameter %s\npost parameter %s", getForm, postForm)
						panic("发送数据时发生异常")
					}
					defer reading.Body.Close()
					if reading.StatusCode != 200 {
						panic("http状态码返回错误")
					}
					body, _ := ioutil.ReadAll(reading.Body)
					msgStat := &MsgStatus{}
					msgErr := json.Unmarshal([]byte(body), msgStat)
					if msgErr != nil {
						panic("无法获取文章扩展数据")
					}
					fmt.Printf("loop in %s\n", msgStat)
					ed = append(ed, ExcelData{pn, pubDate, msgStat.Appmsgstat.ReadNum, msgStat.Appmsgstat.LikeNum, msgStat.CommentCount, title, article}) // 保存已获取的数据到数据结构中。完成所有循环后，写入到文件中
					time.Sleep(5 * time.Second)
				}
			}
			time.Sleep(5 * time.Second)
		}
		excelLoop := 2
		for i := 0; i < len(ed); i++ {
			excelLoop += i
			axis := excelize.ToAlphaString(i) + strconv.Itoa(excelLoop)
			for _, excel := range ed {
				xlsx.SetCellValue(pubname, axis, excel.PublicName)
			}
		}
		timer++
	}
	if xlsxErr := xlsx.SaveAs(currentTime + ".xlsx"); xlsxErr != nil {
		panic("保存excel失败")
	}
	time.Sleep(60 * time.Second)
}
