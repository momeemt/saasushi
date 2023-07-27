package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type RawItem struct {
	SrcContent     string
	NameContent    string
	PriceContent   string
	CalorieContent string
	NoteContent    string
}

func parsePrice(rawItem RawItem) (uint64, error) {
	plate120yen := "/assets/images/pages/common/icon_price_yellow.png"
	plate180yen := "/assets/images/pages/common/icon_price_red.png"
	plate260yen := "/assets/images/pages/common/icon_price_black.png"

	if rawItem.PriceContent == plate120yen {
		return 120, nil
	} else if rawItem.PriceContent == plate180yen {
		return 180, nil
	} else if rawItem.PriceContent == plate260yen {
		return 260, nil
	} else {
		rawPrice := strings.Split(rawItem.CalorieContent, "\n")[0]
		rawPrice = strings.ReplaceAll(rawPrice, "円(税込)", "")
		rawPrice = strings.ReplaceAll(rawPrice, ",", "")
		rawPrice = strings.TrimSpace(rawPrice)
		price, err := strconv.ParseUint(rawPrice, 10, 0)
		if err != nil {
			return 0, err
		}
		return price, nil
	}
}

func parseCalorie(rawItem RawItem) (uint64, error) {
	rawCalorie := ""
	if rawItem.PriceContent == "" {
		rawCalorie = strings.Split(rawItem.CalorieContent, "\n")[1]
	} else {
		rawCalorie = rawItem.CalorieContent
	}
	// TODO: 特殊なカロリー計算の場合はそれがわかるようにする
	rawCalorie = strings.ReplaceAll(rawCalorie, "100gあたり", "")
	rawCalorie = strings.ReplaceAll(rawCalorie, "100mlあたり", "")
	rawCalorie = strings.ReplaceAll(rawCalorie, "砂糖なしの場合", "")

	rawCalorie = strings.ReplaceAll(rawCalorie, "kcal", "")
	rawCalorie = strings.ReplaceAll(rawCalorie, ",", "")
	rawCalorie = strings.TrimSpace(rawCalorie)
	calorie, err := strconv.ParseUint(rawCalorie, 10, 0)
	if err != nil {
		return 0, err
	}
	return calorie, nil
}

func getRawItemInfo(e *colly.HTMLElement) RawItem {
	return RawItem{
		SrcContent:     e.ChildAttr("a span.img img", "src"),
		NameContent:    e.ChildText("a span.txt-wrap span.ttl"),
		PriceContent:   e.ChildAttr("a span.txt-wrap span.plate img", "src"),
		CalorieContent: e.ChildText("a span.txt-wrap span.price"),
		NoteContent:    e.ChildText("a span.txt-wrap span.note"),
	}
}

func parseItem(e *colly.HTMLElement, kind ItemKind) (Item, error) {
	rawItem := getRawItemInfo(e)

	price, err := parsePrice(rawItem)
	if err != nil {
		return Item{}, err
	}

	calorie, err := parseCalorie(rawItem)
	if err != nil {
		return Item{}, err
	}

	return Item{
		ImageSource: rawItem.SrcContent,
		Kind:        kind,
		Name:        rawItem.NameContent,
		Price:       price,
		Calorie:     calorie,
	}, nil
}

type SushiroTsukubaGakuenNoMoriItemsResponse struct {
	Status int
	Items  []Item
}

func sushiroTsukubaGakuenNoMoriItemsHandler(w http.ResponseWriter, _ *http.Request) {
	c := colly.NewCollector()

	var items []Item
	sectionMap := make(map[string]ItemKind)
	sectionMap["#anchor-sec01"] = LimitedEdition
	sectionMap["#anchor-sec03"] = Nigiri
	sectionMap["#anchor-sec04"] = GunkanMakimono
	sectionMap["#anchor-sec05"] = SideMenu
	sectionMap["#anchor-sec06"] = Drink
	sectionMap["#anchor-sec07"] = Dessert

	for section, kind := range sectionMap {
		c.OnHTML(section, func(e *colly.HTMLElement) {
			e.ForEach("ul li", func(_ int, e *colly.HTMLElement) {
				item, err := parseItem(e, kind)
				if err != nil {
					log.Fatal(err)
				}
				items = append(items, item)
			})
		})
	}

	c.Visit("https://www.akindo-sushiro.co.jp/menu/menu_detail/?s_id=528")

	res, err := json.Marshal(items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(res)
}

func main() {
	http.HandleFunc("/sushiro/tsukuba-gakuen-no-mori/items", sushiroTsukubaGakuenNoMoriItemsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
