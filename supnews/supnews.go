package supnews

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aiomonitors/godiscord"
	proxymanager "github.com/aiomonitors/goproxymanager"
	"github.com/fatih/color"
	colors "github.com/fatih/color"
)

type Config struct {
	Webhook   string `json:"webhook"`
	Hexcode   string `json:"hexcode"`
	Icon      string `json:"icon"`
	Groupname string `json:"groupName"`
}

type Monitor struct {
	LatestText string
	Proxies    proxymanager.ProxyManager
	Config     Config
}

type ResponseData struct {
	Title       string
	Description string
	Time        string
	Images      []string
}

type Response http.Response

//GetPage gets the latest news from https://supremenewyork.com/news.
//Can pass in an empty string for proxy if a proxy is not needed in the request.
//Returns struct ResponseData
func GetPage(proxy string) (ResponseData, error) {
	var client *http.Client
	resData := ResponseData{}

	if len(proxy) > 0 {
		proxyURL, parseError := url.Parse(proxy)
		if parseError != nil {
			return resData, parseError
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	} else {
		client = &http.Client{}
	}

	headers := map[string]string{
		"authority":                 "www.supremenewyork.com",
		"cache-control":             "max-age=0",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
		"sec-fetch-dest":            "document",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"sec-fetch-site":            "none",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-user":            "?1",
		"accept-language":           "en-US,en;q=0.9",
		"cookie":                    "__utmz=74692624.1577776302.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utmc=74692624; mp_c5c3c493b693d7f413d219e72ab974b2_mixpanel=%7B%22distinct_id%22%3A%20%2216f5acbc96c38e-0b60a3e08f974d-1d316b5b-1fa400-16f5acbc96db3a%22%2C%22%24device_id%22%3A%20%2216f5acbc96c38e-0b60a3e08f974d-1d316b5b-1fa400-16f5acbc96db3a%22%2C%22Store%20Location%22%3A%20%22US%20Web%22%2C%22Platform%22%3A%20%22Web%22%2C%22%24initial_referrer%22%3A%20%22https%3A%2F%2Fwww.supremenewyork.com%2Fshop%22%2C%22%24initial_referring_domain%22%3A%20%22www.supremenewyork.com%22%7D; __utma=74692624.1437660561.1577776307.1585746821.1585788160.133; __utmb=74692624.1.10.1585788160; _ticket=b001a8c01ff1a9a77a40ba39603b3622bf97f883aa6ac6ea25def88b21a32a1d9da202eb201456103f44f6cac27da50eddbedb48692e5cc085a4cc73ad50261a1585789067",
	}
	//HTTP Request
	req, reqErr := http.NewRequest("GET", "https://www.supremenewyork.com/news", nil)
	if reqErr != nil {
		return resData, reqErr
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, resErr := client.Do(req)
	if resErr != nil {
		return resData, resErr
	}
	defer res.Body.Close()

	doc, docErr := goquery.NewDocumentFromReader(res.Body) //Creates a goquery document
	if docErr != nil {
		return resData, docErr
	}

	newsElem := doc.Find("div.news_container.page-1") //locates the news container
	if class, _ := newsElem.Attr("class"); newsElem != nil && class != "news_container page-1 archived" {
		//finds figure and parses the array structured as "["string"]"
		fig := newsElem.Find("figure")
		links, linksExists := fig.Attr("data-image-urls")
		if linksExists {
			var images []string
			jsonErr := json.Unmarshal([]byte(links), &images)
			if jsonErr != nil {
				return resData, jsonErr
			}
			resData.Images = images
		} else {
			return resData, errors.New("Error parsing news page, links not found")
		}

		time := doc.Find("#news_scroll_container > div.news_page_container.page-1 > div:nth-child(1) > article > time").Text()
		title := newsElem.Find("#news_scroll_container > div.news_page_container.page-1 > div:nth-child(1) > article > h2").Text()
		description := newsElem.Find("#news_scroll_container > div.news_page_container.page-1 > div:nth-child(1) > article > div.blurb").Text()

		resData.Title = title
		resData.Description = description
		resData.Time = time
	}
	return resData, nil
}

//NewMonitor returns a new Monitor object. A path to a txt file containing proxies can be passed in, or an empty string instead.
//If a proxy file is specified, it initializes a proxymanager instance. Proxies can be accessed from Monitor.Proxies
//To access next / random proxy for example: Monitor.Proxies.NextProxy() / Monitor.Proxies.RandomProxy()
func NewMonitor(proxyFile string, c Config) (*Monitor, error) {
	monitor := &Monitor{}
	if len(proxyFile) > 1 {
		manager, managerErr := proxymanager.NewManager(proxyFile)
		if managerErr != nil {
			return nil, managerErr
		}
		monitor.Proxies = *manager
	}

	proxy := ""
	var grabErr error
	if monitor.Proxies.Proxies != nil && len(monitor.Proxies.Proxies) > 0 {
		proxy, grabErr = monitor.Proxies.NextProxy()
		if grabErr != nil {
			return nil, grabErr
		}
	}

	res, resErr := GetPage(proxy)
	if resErr != nil {
		return nil, resErr
	}
	monitor.LatestText = res.Title
	monitor.Config = c
	return monitor, nil
}

func getTime() string {
	return time.Now().Format("15:02:06.111")
}

func printErr(err error) {
	colors.Red("[ %s ] Error in monitor %v", getTime(), err)
}

func greenMessage(msg string) {
	colors.Green("[ %s ] %s", getTime(), msg)
}

func yellowMessage(msg string) {
	colors.Yellow("[ %s ] %s", getTime(), msg)
}

func SendToWebhook(res ResponseData, c Config) {
	e := godiscord.NewEmbed(res.Title, res.Description, "https://www.supremenewyork.com/news")
	e.SetAuthor("Supreme News", "https://www.shihab.dev", "")
	e.SetColor(c.Hexcode)
	e.SetFooter(fmt.Sprintf("%s c/o @aiomonitors", c.Groupname), c.Icon)
	e.SetImage(fmt.Sprintf("https:%s", res.Images[0]))
	e.SendToWebhook(c.Webhook)

	color.Green("[ %s ] Sent image 1", time.Now().Format("2006-01-02 15:04:05"))
	for i := 1; i < len(res.Images)-1; i++ {
		e := godiscord.NewEmbed(res.Title, fmt.Sprintf("Image %v/%v", i+1, len(res.Images)), "https://www.supremenewyork.com/news")
		e.SetImage(fmt.Sprintf("https:%s", res.Images[i]))
		e.SetAuthor("Supreme News", "https://www.shihab.dev", "")
		e.SetColor(c.Hexcode)
		e.SetFooter(fmt.Sprintf("%s c/o @aiomonitors", c.Groupname), c.Icon)
		e.SendToWebhook(c.Webhook)

		color.Green("[ %s ] Delivered picture %v/%v", time.Now().Format("2006-01-02 15:04:05"), i+1, len(res.Images))
		time.Sleep(time.Millisecond * 500)
	}
}

func (m *Monitor) StartMonitor() {
	yellowMessage("Starting Supreme News Monitor")
	yellowMessage("Created by Shihab Chowdhury, @aiomonitors")
	i := true
	for i == true {
		time.Sleep(time.Second * 10)
		proxy := ""
		var grabErr error
		if m.Proxies.Proxies != nil && len(m.Proxies.Proxies) > 0 {
			proxy, grabErr = m.Proxies.NextProxy()
			if grabErr != nil {
				printErr(grabErr)
				continue
			}
		}
		res, resErr := GetPage(proxy)
		if resErr != nil {
			printErr(resErr)
			continue
		}
		if m.LatestText != res.Title {
			greenMessage(fmt.Sprintf("Found new page %s", res.Title))
			m.LatestText = res.Title
			SendToWebhook(res, m.Config)
			continue
		} else {
			yellowMessage("No changes")
			continue
		}
	}
}
