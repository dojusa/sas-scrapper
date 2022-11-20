package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Page struct {
	Pages int    `json:"pages"`
	Shops []Shop `json:"shops"`
}

type Shop struct {
	Name        string `json:"name_common"`
	Description string `json:"merchant_sas.description"`
	RewardRate  string `json:"purchase_reward_rate"`

	// Address
	Address string `json:"address"`
	City    string `json:"city"`

	// Cards
	Amex       string `json:"merchant_amex_sop.status"`
	MasterCard string `json:"merchant_mc_cls.status"`
	Visa       string `json:"merchant_visa_vlps.status"`
}

func (s Shop) AsCSVRow() []string {
	return []string{s.Name, s.RewardRate, s.Amex, s.MasterCard, s.Visa, s.Description}
}

func RunPhysicalStores() {
	fmt.Println("Started scraping physical stores")
	defer fmt.Println("Finished scraping physical stores")

	baseUrl := "https://eb-member-portal-api.loyaltfacts.com/stores?autoComplete=0&country=2&hideComingSoon=0&webShops=0&specialCampaign=0&sortBy=purchase_reward_rate&sortDirection=desc&offset=%d"
	shops := []Shop{}
	offset := 0

	for {
		url := fmt.Sprintf(baseUrl, offset)

		fmt.Printf("Visiting %s\n", url)
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		page := Page{}
		if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
			panic(err)
		}

		if len(page.Shops) == 0 {
			break
		}

		shops = append(shops, page.Shops...)
		offset++
	}

	fmt.Printf("Found %d shops\n", len(shops))

	head := []string{"Name", "RewardRate", "Amex", "MasterCard", "Visa", "Description"}
	WriteCSV("output/physical_stores.csv", head, shops)
}
