package parser

import (
	_ "fmt"
	"log"
	_ "math/rand"
	_ "net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const DEFAULT_URL string = "https://steamcommunity.com/market/search"

func get_num_items() int {
	res_f := 0
	f := func(e *colly.HTMLElement) {
		resultCount := e.Text
		//fmt.Println(resultCount)
		res := strings.Join(strings.Split(resultCount, ","), "")
		//fmt.Println(res)

		res_local, err := strconv.Atoi(res)
		if err != nil {
			log.Fatal("Issue with converting into int in func() in get_pages()")
		}
		res_f = res_local
	}
	c := Collect{colly.NewCollector(), "", nil}
	c.set_conf("span#searchResults_total", f)
	//TODO colly.OnResponse c.c.
	c.c.Visit(DEFAULT_URL+"?appid=730")
	return res_f
}

func get_items(total_pg uint8, count uint8) {
	params := url.Values{}
	//params.Add("start", string(pg * 10))
	if count != 0 {
		params.Add("count", string(count))
	}
	params.Add("sort_dir", "desc")
	params.Add("sort_column", "popular")
	params.Add("appid", string(730))
	params.Add("search_descriptions", string(0))
	params.Add("norender", string(1))
	params.Add("query", "")

	c := Collect{colly.NewCollector(), "", nil}
	f := func(e *colly.HTMLElement) {

	}
	c.set_conf("span#searchResults_total", f)
}

func get_items_page(pg uint8, count uint8) int {
	file, err := os.ReadFile("proxies.txt")
	if err != nil {
		log.Fatal("Файл proxies.txt вызывает ошибку в get_items()")
	}
	proxies := strings.Split(string(file), "\n")
	processed_proxies := make([]string, 30)
	for _, proxie := range proxies {
		processed_proxies = append(processed_proxies, process_proxy(proxie))
	}

	return 1
}

func process_proxy(s string) string { //return like "http://username:password@proxyserver:port"
	user, pass := "DnK0dp", "yfsbdp"
	new_s := strings.Split(s, "://")
	res_string := new_s[0] + "://" + user + ":" + pass + "@" + new_s[1]
	return res_string
}

type HTTPConf interface {
	set_conf(s string, f func(e *colly.HTMLElement))
	set_proxy(s string)
}

type Collect struct {
	c   *colly.Collector
	css string
	fun func(e *colly.HTMLElement)
}

func (p Collect) set_proxy(s string) {
	err := p.c.SetProxy(s)
	if err != nil {
		log.Fatal("Ошибка при соединении с прокси-сервером в set_proxy()")
	}
}

func (p Collect) set_conf(s string, f func(e *colly.HTMLElement)) {
	p.c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36 OPR/72.0.3815.465 (Edition Yx GX)"
	p.c.OnHTML(s, f)
}

func Parse() {
	if len(os.Args) > 1 {
		nums, _ := strconv.Atoi(os.Args[1])
		runtime.GOMAXPROCS(nums)
	}
	//print(get_num_items())
	get_items(0, 10)
}
