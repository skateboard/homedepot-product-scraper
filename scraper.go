package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/skateboard/ajson"

	"github.com/data-harvesters/goapify"
	goapifytask "github.com/data-harvesters/goapify-task"
	goapifytls "github.com/data-harvesters/goapify-tls"
)

type scraper struct {
	goapifytask.Base

	input *input

	client *goapifytls.TlsClient
}

func newScraper(input *input, actor *goapify.Actor) (*scraper, error) {
	tlsClient, err := goapifytls.NewTlsClient(actor, goapifytls.DefaultOptions())
	if err != nil {
		return nil, err
	}

	return &scraper{
		Base:   *goapifytask.New(actor),
		input:  input,
		client: tlsClient,
	}, nil
}

const (
	ScrapeProducts goapifytask.State = iota
	FinishedScraping
)

func (s *scraper) Next(state goapifytask.State) (goapifytask.State, error) {
	switch state {
	case goapifytask.Initialize:
		fmt.Println("initialzing homedepot product scraper")
		return ScrapeProducts, nil
	case ScrapeProducts:
		fmt.Println("scraping homedepot products")
		for _, productID := range s.input.ProductIds {

			for _, storeID := range s.input.StoreIds {
				productInfoSet := false
				productData := map[string]any{
					"sku":      productID,
					"store_id": storeID,
				}
				var stockLevels []map[string]any

				for _, zipCode := range s.input.ZipCodes {
					proxiedClient := s.client.ProxiedClient()

					fmt.Printf("%s: %s: %s: checking stock\n", productID, storeID, zipCode)
					payload := strings.NewReader(fmt.Sprintf(`{"operationName":"productClientOnlyProduct","variables":{"itemId":"%v","quantity":1,"storeId":"%v","zipCode":"%v"},"query":"query productClientOnlyProduct($storeId: String, $zipCode: String, $quantity: Int, $itemId: String!, $dataSource: String, $loyaltyMembershipInput: LoyaltyMembershipInput, $skipSubscribeAndSave: Boolean = false, $skipSpecificationGroup: Boolean = false, $skipInstallServices: Boolean = true, $skipKPF: Boolean = false) {\n  product(itemId: $itemId, dataSource: $dataSource, loyaltyMembershipInput: $loyaltyMembershipInput) {\n    fulfillment(storeId: $storeId, zipCode: $zipCode, quantity: $quantity) {\n      backordered\n      fulfillmentOptions {\n        type\n        services {\n          type\n          locations {\n            isAnchor\n            locationId\n            inventory {\n              isOutOfStock\n              quantity\n              isInStock\n              isLimitedQuantity\n              isUnavailable\n              maxAllowedBopisQty\n              minAllowedBopisQty\n              __typename\n            }\n            curbsidePickupFlag\n            isBuyInStoreCheckNearBy\n            distance\n            storeName\n            state\n            type\n            storePhone\n            __typename\n          }\n          hasFreeShipping\n          freeDeliveryThreshold\n          optimalFulfillment\n          deliveryTimeline\n          deliveryDates {\n            startDate\n            endDate\n            __typename\n          }\n          deliveryCharge\n          dynamicEta {\n            hours\n            minutes\n            __typename\n          }\n          totalCharge\n          __typename\n        }\n        fulfillable\n        __typename\n      }\n      backorderedShipDate\n      bossExcludedShipStates\n      excludedShipStates\n      seasonStatusEligible\n      anchorStoreStatus\n      anchorStoreStatusType\n      sthExcludedShipState\n      bossExcludedShipState\n      onlineStoreStatus\n      onlineStoreStatusType\n      inStoreAssemblyEligible\n      __typename\n    }\n    info {\n      dotComColorEligible\n      productSubType {\n        name\n        link\n        __typename\n      }\n      forProfessionalUseOnly\n      hidePrice\n      ecoRebate\n      quantityLimit\n      sskMin\n      sskMax\n      unitOfMeasureCoverage\n      wasMaxPriceRange\n      wasMinPriceRange\n      label\n      globalCustomConfigurator {\n        customButtonText\n        customDescription\n        customExperience\n        customExperienceUrl\n        customTitle\n        __typename\n      }\n      movingCalculatorEligible\n      prop65Warning\n      returnable\n      hasSubscription\n      isBuryProduct\n      isSponsored\n      isGenericProduct\n      isLiveGoodsProduct\n      sponsoredBeacon {\n        onClickBeacon\n        onViewBeacon\n        __typename\n      }\n      sponsoredMetadata {\n        campaignId\n        placementId\n        slotId\n        __typename\n      }\n      categoryHierarchy\n      samplesAvailable\n      customerSignal {\n        previouslyPurchased\n        __typename\n      }\n      productDepartmentId\n      productDepartment\n      augmentedReality\n      swatches {\n        isSelected\n        itemId\n        label\n        swatchImgUrl\n        url\n        value\n        __typename\n      }\n      totalNumberOfOptions\n      paintBrand\n      fiscalYear\n      classNumber\n      recommendationFlags {\n        visualNavigation\n        __typename\n      }\n      minimumOrderQuantity\n      projectCalculatorEligible\n      subClassNumber\n      calculatorType\n      pipCalculator {\n        coverageUnits\n        display\n        publisher\n        toggle\n        __typename\n      }\n      protectionPlanSku\n      hasServiceAddOns\n      consultationType\n      __typename\n    }\n    dataSources\n    identifiers {\n      productType\n      storeSkuNumber\n      brandName\n      productLabel\n      itemId\n      canonicalUrl\n      modelNumber\n      specialOrderSku\n      toolRentalSkuNumber\n      rentalCategory\n      rentalSubCategory\n      upc\n      upcGtin13\n      parentId\n      isSuperSku\n      roomVOEnabled\n      sampleId\n      __typename\n    }\n    availabilityType {\n      type\n      discontinued\n      status\n      buyable\n      __typename\n    }\n    media {\n      augmentedRealityLink {\n        usdz\n        image\n        __typename\n      }\n      images {\n        url\n        sizes\n        type\n        subType\n        __typename\n      }\n      video {\n        url\n        videoStill\n        link {\n          text\n          url\n          __typename\n        }\n        title\n        type\n        videoId\n        thumbnail\n        longDescription\n        shortDescription\n        __typename\n      }\n      threeSixty {\n        id\n        url\n        __typename\n      }\n      richContent {\n        content\n        displayMode\n        richContentSource\n        salsifyRichContent\n        __typename\n      }\n      __typename\n    }\n    itemId\n    taxonomy {\n      breadCrumbs {\n        browseUrl\n        creativeIconUrl\n        deselectUrl\n        dimensionName\n        label\n        refinementKey\n        url\n        __typename\n      }\n      brandLinkUrl\n      __typename\n    }\n    pricing(storeId: $storeId) {\n      value\n      alternatePriceDisplay\n      alternate {\n        bulk {\n          pricePerUnit\n          thresholdQuantity\n          value\n          __typename\n        }\n        unit {\n          caseUnitOfMeasure\n          unitsOriginalPrice\n          unitsPerCase\n          value\n          __typename\n        }\n        __typename\n      }\n      original\n      mapAboveOriginalPrice\n      message\n      preferredPriceFlag\n      promotion {\n        type\n        description {\n          shortDesc\n          longDesc\n          __typename\n        }\n        dollarOff\n        percentageOff\n        savingsCenter\n        savingsCenterPromos\n        specialBuySavings\n        specialBuyDollarOff\n        specialBuyPercentageOff\n        dates {\n          start\n          end\n          __typename\n        }\n        experienceTag\n        subExperienceTag\n        itemList\n        reward {\n          tiers {\n            minPurchaseAmount\n            minPurchaseQuantity\n            rewardPercent\n            rewardAmountPerOrder\n            rewardAmountPerItem\n            rewardFixedPrice\n            __typename\n          }\n          __typename\n        }\n        nvalues\n        brandRefinementId\n        __typename\n      }\n      specialBuy\n      unitOfMeasure\n      conditionalPromotions {\n        dates {\n          start\n          end\n          __typename\n        }\n        description {\n          shortDesc\n          longDesc\n          __typename\n        }\n        experienceTag\n        subExperienceTag\n        eligibilityCriteria {\n          itemGroup\n          minPurchaseAmount\n          minPurchaseQuantity\n          relatedSkusCount\n          omsSkus\n          __typename\n        }\n        reward {\n          tiers {\n            minPurchaseAmount\n            minPurchaseQuantity\n            rewardPercent\n            rewardAmountPerOrder\n            rewardAmountPerItem\n            rewardFixedPrice\n            __typename\n          }\n          __typename\n        }\n        nvalues\n        brandRefinementId\n        __typename\n      }\n      __typename\n    }\n    subscription @skip(if: $skipSubscribeAndSave) {\n      defaultfrequency @skip(if: $skipSubscribeAndSave)\n      discountPercentage @skip(if: $skipSubscribeAndSave)\n      subscriptionEnabled @skip(if: $skipSubscribeAndSave)\n      __typename\n    }\n    favoriteDetail {\n      count\n      __typename\n    }\n    badges(storeId: $storeId) {\n      label\n      name\n      color\n      creativeImageUrl\n      endDate\n      message\n      timerDuration\n      timer {\n        timeBombThreshold\n        daysLeftThreshold\n        dateDisplayThreshold\n        message\n        __typename\n      }\n      __typename\n    }\n    reviews {\n      ratingsReviews {\n        totalReviews\n        averageRating\n        __typename\n      }\n      __typename\n    }\n    details {\n      collection {\n        url\n        collectionId\n        name\n        __typename\n      }\n      description\n      descriptiveAttributes {\n        name\n        value\n        bulleted\n        sequence\n        __typename\n      }\n      highlights\n      additionalResources {\n        infoAndGuides {\n          name\n          url\n          __typename\n        }\n        installationAndRentals {\n          contentType\n          name\n          url\n          __typename\n        }\n        diyProjects {\n          contentType\n          name\n          url\n          __typename\n        }\n        __typename\n      }\n      installation {\n        leadGenUrl\n        __typename\n      }\n      __typename\n    }\n    seo {\n      seoKeywords\n      seoDescription\n      __typename\n    }\n    specificationGroup @skip(if: $skipSpecificationGroup) {\n      specifications {\n        specName\n        specValue\n        __typename\n      }\n      specTitle @skip(if: $skipSpecificationGroup)\n      __typename\n    }\n    dataSource\n    installServices(storeId: $storeId, zipCode: $zipCode) @skip(if: $skipInstallServices) {\n      scheduleAMeasure @skip(if: $skipInstallServices)\n      gccCarpetDesignAndOrderEligible @skip(if: $skipInstallServices)\n      __typename\n    }\n    keyProductFeatures @skip(if: $skipKPF) {\n      keyProductFeaturesItems {\n        features {\n          name\n          refinementId\n          refinementUrl\n          value\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    sizeAndFitDetail {\n      attributeGroups {\n        attributes {\n          attributeName\n          dimensions\n          __typename\n        }\n        dimensionLabel\n        productType\n        __typename\n      }\n      __typename\n    }\n    seoDescription\n    __typename\n  }\n}\n"}`,
						productID, storeID, zipCode))

					req, err := http.NewRequest("POST", "https://www.homedepot.com/federation-gateway/graphql?opname=productClientOnlyProduct", payload)
					if err != nil {
						continue
					}

					req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0")
					req.Header.Add("Accept", "*/*")
					req.Header.Add("Accept-Language", "en-US,en;q=0.5")
					req.Header.Add("Referer", "https://www.homedepot.com/p/DEWALT-20V-MAX-Lithium-Ion-Cordless-Cable-Stapler-with-2-0Ah-Battery-Charger-and-Bag-DCN701D1/308047450?store=4115")
					req.Header.Add("content-type", "application/json")
					req.Header.Add("X-Experience-Name", "b2b")
					req.Header.Add("apollographql-client-name", "b2b")
					req.Header.Add("apollographql-client-version", "0.0.0")
					req.Header.Add("X-current-url", "/p/DEWALT-20V-MAX-Lithium-Ion-Cordless-Cable-Stapler-with-2-0Ah-Battery-Charger-and-Bag-DCN701D1/308047450")
					req.Header.Add("x-hd-dc", "origin")
					req.Header.Add("x-customer-type", "B2B")
					req.Header.Add("x-customer-role", "ADMIN")
					req.Header.Add("x-segment-id", "Contractors")
					req.Header.Add("Origin", "https://www.homedepot.com")
					req.Header.Add("Connection", "keep-alive")
					req.Header.Add("Sec-Fetch-Dest", "empty")
					req.Header.Add("Sec-Fetch-Mode", "cors")
					req.Header.Add("Sec-Fetch-Site", "same-origin")
					req.Header.Add("Pragma", "no-cache")
					req.Header.Add("Cache-Control", "no-cache")
					req.Header.Add("TE", "trailers")

					res, err := proxiedClient.Do(req)
					if err != nil {
						fmt.Printf("%s: %s: %s: error making request: %v\n", productID, storeID, zipCode, err)
						continue
					}

					body, err := io.ReadAll(res.Body)
					if err != nil {
						fmt.Printf("%s: %s: %s: error reading body: %v\n", productID, storeID, zipCode, err)
						continue
					}
					j := ajson.Parse(string(body))

					if j.Get("errors").Exists() {
						for _, err := range j.Get("errors").Array() {
							fmt.Printf("%s: %s: %s: error reading body: %v\n", productID, storeID, zipCode, err.Get("message").String())
						}

						continue
					}

					product := j.Get("data").Get("product")

					if !productInfoSet {
						productInfoSet = true
						// image := product.Get("media").Get("images").Array()[0].Get("url").String()
						// image = strings.Replace(image, "<SIZE>", "400", 1)

						// productData["image"] = image
						productData["name"] = product.Get("identifiers").Get("productLabel").String()
						productData["brand"] = product.Get("identifiers").Get("brandName").String()
						productData["price"] = float32(product.Get("pricing").Get("value").Float())
					}

					inv := map[string]any{
						"zip_code":        zipCode,
						"stock_level":     int64(0),
						"delivery":        "N/A",
						"delivery_charge": "",
					}

					var fulfillment *ajson.Result
					for _, f := range product.Get("fulfillment").Get("fulfillmentOptions").Array() {
						if f.Get("type").Exists() {
							if f.Get("type").String() == "delivery" {
								fulfillment = &f
								break
							}
						}
					}

					if fulfillment == nil {
						fmt.Printf("%s: %s: %s: no fulfillment found\n", productID, storeID, zipCode)
						continue
					}

					var service *ajson.Result
					for _, serv := range fulfillment.Get("services").Array() {
						if serv.Get("type").String() == "express delivery" {
							service = &serv
							break
						}
					}

					if service == nil {
						fmt.Printf("%s: %s: %s: no service found\n", productID, storeID, zipCode)
						continue
					}

					if service.Get("locations").Exists() {
						for _, l := range service.Get("locations").Array() {
							if l.Get("inventory.isOutOfStock").Exists() {
								if l.Get("inventory.isOutOfStock").Bool() {
									inv["stock_level"] = int64(0)
									break
								}
							}
							stockLevel := inv["stock_level"].(int64)
							stockLevel += l.Get("inventory.quantity").Int()
							inv["stock_level"] = stockLevel
							break
						}
					}

					if inv["stock_level"] == 0 {
						continue
					}

					deliveryDate := service.Get("deliveryDates").Get("startDate").String()
					deliveryTimeline := service.Get("deliveryTimeline").String()
					deliveryCharge := service.Get("deliveryCharge").String()

					inv["delivery"] = deliveryTimeline + "-" + deliveryDate
					inv["delivery_charge"] = deliveryCharge

					stockLevels = append(stockLevels, inv)
				}

				if len(stockLevels) == 0 {
					continue
				}

				if s.input.ScrapeReviews {
					reviews, err := s.scrapeReviews(productID, s.client.ProxiedClient())
					if err == nil {
						productData["reviews"] = reviews
					} else {
						fmt.Printf("%s: %s: failed to scrape reviews: %v\n", productID, storeID, err)
					}
				}

				fmt.Printf("%s: %s: sending stock: %d\n", productID, storeID, len(stockLevels))

				productData["stock_levels"] = stockLevels

				err := s.Actor().Output(productData)
				if err != nil {
					fmt.Printf("%s: %s: failed send output: %d\n", productID, storeID, err)
					continue
				}
			}
		}

		return FinishedScraping, nil
	case FinishedScraping:
		fmt.Println("succesfully scrapped all product data!")
		s.Stop()
		return 0, nil
	}

	s.Stop() // should never get here, theoretically, unless wanted
	return 0, nil
}

func (s *scraper) scrapeReviews(productID string, tlsClient tls_client.HttpClient) ([]map[string]any, error) {
	payload := strings.NewReader(fmt.Sprintf(`{
	"operationName": "reviews",
	"variables": {
		"filters": {
			"isVerifiedPurchase": false,
			"prosCons": null,
			"starRatings": null
		},
		"itemId": "%v",
		"pagesize": "10",
		"recfirstpage": "10",
		"searchTerm": null,
		"sortBy": "photoreview",
		"startIndex": 0
	},
	"query": "query reviews($itemId: String!, $searchTerm: String, $sortBy: String, $startIndex: Int, $recfirstpage: String, $pagesize: String, $filters: ReviewsFilterInput) {\n  reviews(\n    itemId: $itemId\n    searchTerm: $searchTerm\n    sortBy: $sortBy\n    startIndex: $startIndex\n    recfirstpage: $recfirstpage\n    pagesize: $pagesize\n    filters: $filters\n  ) {\n    Results {\n      AuthorId\n      Badges {\n        DIY {\n          BadgeType\n          __typename\n        }\n        top250Contributor {\n          BadgeType\n          __typename\n        }\n        IncentivizedReview {\n          BadgeType\n          __typename\n        }\n        EarlyReviewerIncentive {\n          BadgeType\n          __typename\n        }\n        top1000Contributor {\n          BadgeType\n          __typename\n        }\n        VerifiedPurchaser {\n          BadgeType\n          __typename\n        }\n        __typename\n      }\n      BadgesOrder\n      CampaignId\n      ContextDataValues {\n        Age {\n          Value\n          __typename\n        }\n        VerifiedPurchaser {\n          Value\n          __typename\n        }\n        __typename\n      }\n      ContextDataValuesOrder\n      Id\n      IsRecommended\n      IsSyndicated\n      Photos {\n        Id\n        Sizes {\n          normal {\n            Url\n            __typename\n          }\n          thumbnail {\n            Url\n            __typename\n          }\n          __typename\n        }\n        __typename\n      }\n      ProductId\n      SubmissionTime\n      TagDimensions {\n        Pro {\n          Values\n          __typename\n        }\n        Con {\n          Values\n          __typename\n        }\n        __typename\n      }\n      Title\n      TotalNegativeFeedbackCount\n      TotalPositiveFeedbackCount\n      ClientResponses {\n        Response\n        Date\n        Department\n        __typename\n      }\n      Rating\n      RatingRange\n      ReviewText\n      SecondaryRatings {\n        Quality {\n          Label\n          Value\n          __typename\n        }\n        Value {\n          Label\n          Value\n          __typename\n        }\n        EnergyEfficiency {\n          Label\n          Value\n          __typename\n        }\n        Features {\n          Label\n          Value\n          __typename\n        }\n        Appearance {\n          Label\n          Value\n          __typename\n        }\n        EaseOfInstallation {\n          Label\n          Value\n          __typename\n        }\n        EaseOfUse {\n          Label\n          Value\n          __typename\n        }\n        __typename\n      }\n      SecondaryRatingsOrder\n      SyndicationSource {\n        LogoImageUrl\n        Name\n        __typename\n      }\n      UserNickname\n      UserLocation\n      Videos {\n        VideoId\n        VideoThumbnailUrl\n        VideoUrl\n        __typename\n      }\n      __typename\n    }\n    Includes {\n      Products {\n        store {\n          Id\n          FilteredReviewStatistics {\n            AverageOverallRating\n            TotalReviewCount\n            TotalRecommendedCount\n            RecommendedCount\n            NotRecommendedCount\n            SecondaryRatingsAveragesOrder\n            RatingDistribution {\n              RatingValue\n              Count\n              __typename\n            }\n            ContextDataDistribution {\n              Age {\n                Values {\n                  Value\n                  Count\n                  __typename\n                }\n                __typename\n              }\n              Gender {\n                Values {\n                  Value\n                  Count\n                  __typename\n                }\n                __typename\n              }\n              Expertise {\n                Values {\n                  Value\n                  __typename\n                }\n                __typename\n              }\n              HomeGoodsProfile {\n                Values {\n                  Value\n                  Count\n                  __typename\n                }\n                __typename\n              }\n              VerifiedPurchaser {\n                Values {\n                  Value\n                  Count\n                  __typename\n                }\n                __typename\n              }\n              __typename\n            }\n            __typename\n          }\n          __typename\n        }\n        items {\n          Id\n          FilteredReviewStatistics {\n            AverageOverallRating\n            TotalReviewCount\n            TotalRecommendedCount\n            RecommendedCount\n            NotRecommendedCount\n            SecondaryRatingsAveragesOrder\n            RatingDistribution {\n              RatingValue\n              Count\n              __typename\n            }\n            ContextDataDistribution {\n              Age {\n                Values {\n                  Value\n                  Count\n                  __typename\n                }\n                __typename\n              }\n              Gender {\n                Values {\n                  Value\n                  Count\n                  __typename\n                }\n                __typename\n              }\n              Expertise {\n                Values {\n                  Value\n                  __typename\n                }\n                __typename\n              }\n              HomeGoodsProfile {\n                Values {\n                  Value\n                  Count\n                  __typename\n                }\n                __typename\n              }\n              VerifiedPurchaser {\n                Values {\n                  Value\n                  Count\n                  __typename\n                }\n                __typename\n              }\n              __typename\n            }\n            __typename\n          }\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    FilterSelected {\n      StarRatings {\n        is5Star\n        is4Star\n        is3Star\n        is2Star\n        is1Star\n        __typename\n      }\n      VerifiedPurchaser\n      SearchText\n      __typename\n    }\n    pagination {\n      previousPage {\n        label\n        isNextPage\n        isPreviousPage\n        isSelectedPage\n        __typename\n      }\n      pages {\n        label\n        isNextPage\n        isPreviousPage\n        isSelectedPage\n        __typename\n      }\n      nextPage {\n        label\n        isNextPage\n        isPreviousPage\n        isSelectedPage\n        __typename\n      }\n      __typename\n    }\n    SortBy {\n      mosthelpfull\n      newest\n      oldest\n      highestrating\n      lowestrating\n      photoreview\n      __typename\n    }\n    TotalResults\n    __typename\n  }\n}"
}`, productID))

	req, err := http.NewRequest("POST", "https://www.homedepot.com/federation-gateway/graphql?opname=reviews", payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Referer", "https://www.homedepot.com/p/DEWALT-20V-MAX-Lithium-Ion-Cordless-Cable-Stapler-with-2-0Ah-Battery-Charger-and-Bag-DCN701D1/308047450?store=4115")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-Experience-Name", "b2b")
	req.Header.Add("apollographql-client-name", "b2b")
	req.Header.Add("apollographql-client-version", "0.0.0")
	req.Header.Add("X-current-url", "/p/DEWALT-20V-MAX-Lithium-Ion-Cordless-Cable-Stapler-with-2-0Ah-Battery-Charger-and-Bag-DCN701D1/308047450")
	req.Header.Add("x-hd-dc", "origin")
	req.Header.Add("x-customer-type", "B2B")
	req.Header.Add("x-customer-role", "ADMIN")
	req.Header.Add("x-segment-id", "Contractors")
	req.Header.Add("Origin", "https://www.homedepot.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("TE", "trailers")

	res, err := tlsClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	j := ajson.Parse(string(body))

	reviews := j.Get("data").Get("reviews").Get("Results").Array()

	if len(reviews) == 0 {
		return nil, errors.New("no reviews found")
	}

	var revs []map[string]any
	for _, review := range reviews {
		revs = append(revs, map[string]any{
			"title":    review.Get("Title").String(),
			"rating":   review.Get("Rating").Float(),
			"text":     review.Get("ReviewText").String(),
			"username": review.Get("UserNickname").String(),
		})
	}

	return revs, nil
}
