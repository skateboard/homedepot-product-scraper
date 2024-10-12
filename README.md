![enter image description here](https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fwww.pngmart.com%2Ffiles%2F16%2FHome-Depot-Logo-PNG-Image.png&f=1&nofb=1&ipt=def5757aced3b327c03b99e5ef3280e8f0c80bc271e4c8682ef7f1c5bd1f976f&ipo=images)

# Homedepot Product Data Scraper

## About This Actor

This Actor is a powerful, user-fiendly tool made to scrape products from specific [Homedepot.com](https://homedepot.com) Products. This tool will save you time and provide you with reliable data on products from Homedepot.

Made with Golang 1.22.1

## Tutorial

Basic Usage

```json
{
    "productIds": ["328984390"],
    "zipCodes": ["60607"],
    "storeIds": ["1950"],
    "scrapeReviews": true,
    "reviewsLimit": 100
}
```

| parameter | type | argument | description |
| --------- | ----- | ------------------------- | ---------------------------- |
| productIds | array | _["1223", "12312312", ...]_ | An array of Homedepot Product IDs |
| zipCodes | array | _["1223", "12312312", ...]_ | An array of USA Zip-codes |
| storeIds | array | _["1223", "12312312", ...]_ | An array of Homedepot Store IDs |
| scrapeReviews | bool | _default=true_ | A boolean if you want to scrape product reviews |
| reviewsLimit | int | _default=100_ | limit the number of reviews you want to scrape |

### Output Sample

```json
[
  {
    "brand": "Milwaukee",
    "name": "M12 12V Lithium-Ion Compact 2.0 Ah Battery Pack (2-Pack) Starter Kit with Charger",
    "price": 129,
    "sku": "328984390",
    "stock_levels": [
      {
        "delivery": "today-2024-10-11",
        "delivery_charge": "0.0",
        "stock_level": 1767,
        "zip_code": "60607"
      }
    ],
    "reviews": [
    {
      "rating": 4,
      "text": "I really love it!",
      "title": "Super awesome really like it",
      "username": "SuperAwesomeReviewer"
    }
  ],
    "store_id": "1950"
  }
]
```
