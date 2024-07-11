package main

import (
	"encoding/json"
	"fmt"
	"log"
	_ "math/rand"
	_ "net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	_ "steamtrade.shop/parser/internal/repository"
	"steamtrade.shop/parser/internal/repository/database"
	"steamtrade.shop/parser/pkg"
)

const DEFAULT_URL string = "https://steamcommunity.com/market/search"

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
	if s != "" {
		p.c.OnHTML(s, f)
	}
}

func process_proxy(s string) string { //return like "http://username:password@proxyserver:port"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	user := os.Getenv("PROXY_USER")
	pass := os.Getenv("PROXY_PASS")
	new_s := strings.Split(s, "://")
	res_string := new_s[0] + "://" + user + ":" + pass + "@" + new_s[1]
	return res_string
}

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
	collector := Collect{colly.NewCollector(), "", nil}
	collector.set_conf("span#searchResults_total", f)
	//TODO colly.OnResponse c.c.
	collector.c.Visit(DEFAULT_URL + "?appid=730")
	return res_f
}

func get_items(total_pg uint8, count uint8, db *gorm.DB) {
	params := url.Values{}
	//params.Add("start", string(pg * 10))
	if count != 0 {
		params.Add("count", strconv.Itoa(int(count)))
	}
	params.Add("sort_dir", "desc")
	params.Add("sort_column", "popular")
	params.Add("appid", "730")
	params.Add("search_descriptions", "0")
	params.Add("norender", "1")
	params.Add("query", "")

	collector := Collect{colly.NewCollector(), "", nil}

	collector.c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36 OPR/72.0.3815.465 (Edition Yx GX)"
	collector.c.OnResponse(func(r *colly.Response) {

		var result map[string]interface{}
		json.Unmarshal(r.Body, &result) //result - JSON-like ответ на запрос

		if results, ok := result["results"].([]interface{}); ok {
			new_items := make([]model.Product, 0) //для добавления целого батча
			//fmt.Println(results)
			for _, item := range results { //item = interface{}

				if listing, ok := item.(map[string]interface{}); ok {

					name := listing["name"].(string)
					//fmt.Println(name)
					price := listing["sell_price"].(float64) / 100

					if contain, _, this_item := database.Contains(name, db); contain { //проверка, что в базе уже есть элемент с таким именем
						this_item.Price = price
						this_item.ID = 0
						new_items = append(new_items, this_item)

					} else {
						name20 := strings.Replace(name, " ", "%20", -1)
						ref := "https://steamcommunity.com/market/listings/730/" + name20
						var img_small, img_big string
						if img, ok := listing["asset_description"].(map[string]interface{}); ok {
							img_small = "https://community.akamai.steamstatic.com/economy/image/" + img["icon_url"].(string) + "/62fx62f"
							img_big = "https://community.akamai.steamstatic.com/economy/image/" + img["icon_url"].(string) + "/360fx360f"
						} else {
							img_small = ""
							img_big = ""
						}
						new_items = append(new_items, model.Product{Name: name, ImageBig: img_big, ImageSmall: img_small, Ref: ref, Price: price})
					}

				}
			}
			fmt.Println(new_items)
			//добавляем все цены
			db.CreateInBatches(new_items, int(count))
		}
	})
	fmt.Println(DEFAULT_URL + "/render/?" + params.Encode())
	collector.c.Visit(DEFAULT_URL + "/render/?" + params.Encode())

}

func get_page_items(pg uint8, count uint8) int {
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

func main() {
	db := database.Conn()
	//db.Migrator().DropTable("products")
	//db.AutoMigrate(&model.Product{})
	if len(os.Args) > 1 {
		nums, _ := strconv.Atoi(os.Args[1])
		runtime.GOMAXPROCS(nums)
	}
	//print(get_num_items())
	get_items(0, 10, db)
}
