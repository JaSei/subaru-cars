package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const hyperinzerceURL = "https://autobazar.hyperinzerce.cz/subaru/"

func HyperinzerceScrape(db DB) {
	url := hyperinzerceURL

	for {
		// Request the HTML page.
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find(".inzerat__text").Find("h3").Find("a").Each(func(i int, s *goquery.Selection) {
			href := s.AttrOr("href", "")
			title := s.Text()
			log.Printf("Review %d: %s - %s\n", i, href, title)

			if !strings.Contains(strings.ToLower(title), "subaru") {
				log.Printf("Title doesn't looks like subaru ad")
				return
			}

			carUrl := href
			if db.touch(carUrl) {
				log.Printf("We track this car")
			} else {
				hyperCar(db, carUrl, title)
				log.Printf("New car was added")
			}
		})

		url = ""
		doc.Find("a[rel=\"next\"]").Each(func(i int, s *goquery.Selection) {
			url = s.AttrOr("href", "")
			log.Printf("next url: %s\n", url)
		})

		if url == "" {
			break
		}

		time.Sleep(3 * time.Second)
	}
}

func hyperCar(db DB, url, title string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	car := Car{URL: url, Title: title}

	doc.Find(".inz_description").Find("p").Each(func(i int, s *goquery.Selection) {
		log.Printf("Popisdetail %d: %s\n", i, s.Text())
		car.Description = s.Text()
	})

	InfoExtractor(&car)

	priceA := strings.Split(doc.Find(".price_tag").First().Text(), ":")
	car.Price = strings.TrimSpace(priceA[1])

	info := make(map[string]string)
	doc.Find(".inz_detail__table").Find(".row").Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		k := strings.TrimSpace(children.First().Text())
		v := strings.TrimSpace(children.Next().Text())

		info[k] = v
	})

	car.VIN = info["VIN kód"]
	car.Milage = info["Stav tachometru"]

	yearOfManufactory, err := strconv.ParseUint(info["Rok výroby"], 10, 32)
	if err != nil {
		log.Printf("Error %s", err)
	}
	car.YearOfManufactory = uint(yearOfManufactory)

	switch info["Palivo"] {
	case "benzin":
		car.Fuel = gas
	case "diesel":
		car.Fuel = diesel
	case "plyn (LPG, CNG atd.)":
		car.Fuel = lpg
	}

	db.add(car)
}
