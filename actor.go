package main

import (
	"fmt"
	"os"

	"github.com/data-harvesters/goapify"
	goapifyscraper "github.com/data-harvesters/goapify-scraper"
)

type input struct {
	*goapify.ProxyConfigurationOptions `json:"proxyConfiguration"`

	ProductIds    []string `json:"productIds"`
	ZipCodes      []string `json:"zipCodes"`
	StoreIds      []string `json:"storeIds"`
	ScrapeReviews bool     `json:"scrapeReviews"`
}

func main() {
	a := goapify.NewActor(
		os.Getenv("APIFY_DEFAULT_KEY_VALUE_STORE_ID"),
		os.Getenv("APIFY_TOKEN"),
		os.Getenv("APIFY_DEFAULT_DATASET_ID"),
	)

	i := new(input)

	err := a.Input(i)
	if err != nil {
		fmt.Printf("failed to decode input: %v\\n", err)
		panic(err)
	}

	if i.ProxyConfigurationOptions != nil {
		err = a.CreateProxyConfiguration(i.ProxyConfigurationOptions)
		if err != nil {
			panic(err)
		}
	}

	s, err := newScraper(i, a)
	if err != nil {
		fmt.Printf("failed to create scraper: %v\\n", err)
		panic(err)
	}

	goapifyscraper.Run(s)
}
