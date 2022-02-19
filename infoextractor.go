package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func InfoExtractor(car *Car) {
	re := regexp.MustCompile("[A-Z0-9]{17}")
	car.VIN = re.FindString(car.Description)

	re = regexp.MustCompile("(?i)(?:najeto|tachometru?):?[ ]*([0-9x \\.]{3,}[ ]*(?:tis\\.?)?)[ ]*(?:km)?")
	milage := re.FindStringSubmatch(car.Description)

	if len(milage) > 0 {
		car.Milage = milage[1]
	}

	current_year := time.Now().Year()
	years := make([]string, current_year-1990+1)
	for i := 1990; i <= current_year; i++ {
		years[i-1990] = strconv.Itoa(i)
	}
	yearRange := strings.Join(years, "|")

	//1990 - 2021
	yearOfManufacturTitleRe := regexp.MustCompile(fmt.Sprintf("(%s)", yearRange))
	ymt := yearOfManufacturTitleRe.FindStringSubmatch(car.Title)
	if len(ymt) > 1 {
		year, err := strconv.ParseUint(ymt[1], 10, 32)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(year)

		car.YearOfManufactory = uint(year)
	}

	yearOfManufacturDescriptionRe := regexp.MustCompile(fmt.Sprintf("(?i)(?:r[^a-z]?v[^a-z]*|rok v[ýy]r|(?:do )?provozu(?: od)?|rok|registrace)[^0-9a-z]*(%s)", yearRange))
	ymd := yearOfManufacturDescriptionRe.FindStringSubmatch(car.Description)
	if len(ymd) > 1 {
		year, err := strconv.ParseUint(ymd[1], 10, 32)
		if err != nil {
			log.Fatal(err)
		}

		car.YearOfManufactory = uint(year)
	}

	dieselRe := regexp.MustCompile("(?i)nafta|diesel")
	lpgRe := regexp.MustCompile("(?i)lpg")
	benzinRe := regexp.MustCompile("(?i)benz[ií]n")
	if dieselRe.MatchString(car.Description) {
		car.Fuel = diesel
	} else if lpgRe.MatchString(car.Title) || lpgRe.MatchString(car.Description) {
		car.Fuel = lpg
	} else if benzinRe.MatchString(car.Title) || benzinRe.MatchString(car.Description) {
		car.Fuel = gas
	}

	legacyRe := regexp.MustCompile("(?i)legacy")
	foresterRe := regexp.MustCompile("(?i)forester")
	outbackRe := regexp.MustCompile("(?i)outback")
	imprezaRe := regexp.MustCompile("(?i)impreza")
	levorgRe := regexp.MustCompile("(?i)levorg")
	tribecaRe := regexp.MustCompile("(?i)tribeca|b9")
	wrxstiRe := regexp.MustCompile("(?i)wrx|sti")
	xvRe := regexp.MustCompile("(?i)xv")
	justyRe := regexp.MustCompile("(?i)justy")

	if legacyRe.MatchString(car.Title) || legacyRe.MatchString(car.Description) {
		car.SubaruModel = legacy
	} else if foresterRe.MatchString(car.Title) || foresterRe.MatchString(car.Description) {
		car.SubaruModel = forester
	} else if outbackRe.MatchString(car.Title) || outbackRe.MatchString(car.Description) {
		car.SubaruModel = outback
	} else if imprezaRe.MatchString(car.Title) || imprezaRe.MatchString(car.Description) {
		car.SubaruModel = impreza
	} else if levorgRe.MatchString(car.Title) || levorgRe.MatchString(car.Description) {
		car.SubaruModel = levorg
	} else if wrxstiRe.MatchString(car.Title) || wrxstiRe.MatchString(car.Description) {
		car.SubaruModel = wrxsti
	} else if tribecaRe.MatchString(car.Title) || tribecaRe.MatchString(car.Description) {
		car.SubaruModel = tribeca
	} else if justyRe.MatchString(car.Title) || justyRe.MatchString(car.Description) {
		car.SubaruModel = justy
	} else if xvRe.MatchString(car.Title) || xvRe.MatchString(car.Description) {
		car.SubaruModel = xv
	}

	serviceBookRe := regexp.MustCompile("(?i)SERVISN[ÍI] KN[IÍ][ZŽ]K")
	if serviceBookRe.MatchString(car.Description) {
		car.ServiceBook = true
	}
}
