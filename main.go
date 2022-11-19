package main

import (
	"encoding/csv"
	"fmt"
	"os"
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

var numbers = regexp.MustCompile("[0-9]+")

func main() {
	c := colly.NewCollector()

	stores := []Store{}

	// Find and visit all links
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

	c.Visit("https://onlineshopping.flysas.com/se/alle-butiker/")

	fmt.Printf("Found %d stores\n", len(stores))

	writeCSV(stores)
}

func writeCSV(stores []Store) error {
	if len(stores) == 0 {
		fmt.Println("Skipping csv write as no stores where found")
		return nil
	}

	csvFile, err := os.Create("stores.csv")
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	writer.Write([]string{"Name", "Points", "PerKr", "URL", "ImgURL"})
	for _, s := range stores {
		writer.Write([]string{s.Name, s.Points, s.PerKr, s.URL, s.ImgURL})
	}

	return nil
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
