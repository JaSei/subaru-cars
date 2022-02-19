package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const bazosURL = "https://auto.bazos.cz"

func BazosScrape(db DB) {
	url := bazosURL + "/ostatni/?hledat=subaru&rubriky=auto&hlokalita=&humkreis=25&cenaod=&cenado=&Submit=Hledat&kitx=ano"

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

		doc.Find(".nadpis").Find("a").Each(func(i int, s *goquery.Selection) {
			href := s.AttrOr("href", "")
			title := s.Text()
			log.Printf("Review %d: %s - %s\n", i, href, title)

			if !strings.Contains(strings.ToLower(title), "subaru") {
				log.Printf("Title doesn't looks like subaru ad")
				return
			}

			carUrl := bazosURL + href
			if db.touch(carUrl) {
				log.Printf("We track this car")
			} else {
				car(db, carUrl, title)
				log.Printf("New car was added")
			}
		})

		url = ""
		doc.Find(".strankovani").Find("a").Each(func(i int, s *goquery.Selection) {
			if s.Text() == "Další" {
				url = bazosURL + s.AttrOr("href", "")
				log.Printf("next url: %s\n", url)
			}
		})

		if url == "" {
			break
		}

		time.Sleep(3 * time.Second)
	}
}

func car(db DB, url, title string) {
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

	doc.Find(".popisdetail").Each(func(i int, s *goquery.Selection) {
		log.Printf("Popisdetail %d: %s\n", i, s.Text())
		car.Description = s.Text()
	})

	InfoExtractor(&car)

	info := make(map[string]string)
	doc.Find("td.listadvlevo").Find("tr").Each(func(i int, s *goquery.Selection) {
		splited := strings.Split(s.Text(), ":")
		info[splited[0]] = splited[1]
	})

	car.Price = info["Cena"]
	car.Location = info["Lokalita"]

	db.add(car)
}
