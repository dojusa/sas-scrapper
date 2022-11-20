package main

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly"
)

type Store struct {
	Name   string
	URL    string
	ImgURL string
	Points string
	PerKr  string
}

func (s Store) AsCSVRow() []string {
	return []string{s.Name, s.Points, s.PerKr, s.URL, s.ImgURL}
}

func RunOnlineStores() {
	fmt.Println("Started scraping online stores")
	defer fmt.Println("Finished scraping online stores")

	const onlineStoresUrl = "https://onlineshopping.flysas.com/se/alle-butiker/"
	numbers := regexp.MustCompile("[0-9]+")
	stores := []Store{}

	c := colly.NewCollector()
	c.OnHTML("div .ShopGrid > div", func(e *colly.HTMLElement) {
		s := Store{
			ImgURL: e.ChildAttr("div .ShopImage img", "src"),
			Name:   e.ChildText("div .ShopInfo.BorderBottom"),
			Points: numbers.FindString(e.ChildText(".ShopInfo span")),
			PerKr:  numbers.FindString(getPer(e)),
			URL:    e.ChildAttr("div a", "href"),
		}
		stores = append(stores, s)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(onlineStoresUrl)

	fmt.Printf("Found %d stores\n", len(stores))

	head := []string{"Name", "Points", "PerKr", "URL", "ImgURL"}
	WriteCSV("output/online_stores.csv", head, stores)
}

func getPer(e *colly.HTMLElement) string {
	var per string
	e.ForEachWithBreak(".ShopInfo string", func(i int, h *colly.HTMLElement) bool {
		if i == 2 {
			per = h.Text
			return false
		}
		return true
	})
	return per
}
