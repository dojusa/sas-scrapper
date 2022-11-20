package main

import (
	"fmt"
)

func main() {
	fmt.Println("Started Scraping")
	defer fmt.Println("Finished Scraping")

	RunPhysicalStores()
	RunOnlineStores()
}
