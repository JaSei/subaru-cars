package main

import (
	"errors"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

type Fuel string
type SubaruModel string
type Transmission string

const (
	diesel Fuel = "diesel"
	gas    Fuel = "benzin"
	lpg    Fuel = "lpg"

	legacy   SubaruModel = "legacy"
	forester SubaruModel = "forester"
	outback  SubaruModel = "outback"
	impreza  SubaruModel = "impreza"
	levorg   SubaruModel = "levorg"
	wrxsti   SubaruModel = "wrx sti"
	tribeca  SubaruModel = "tribeca"
	justy    SubaruModel = "justy"
	xv       SubaruModel = "xv"

	automatic Transmission = "automatic"
	cvt       Transmission = "cvt"
	manual    Transmission = "manual"

	Unknown uint = iota
	OutbackV36R
)

type ModelInfo struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	YearFrom     uint
	YearTo       uint
	Volume       uint
	EnginePower  uint
	Transmission Transmission
}

type Car struct {
	gorm.Model
	URL               string
	Title             string
	Description       string
	Price             string
	Location          string
	VIN               string
	Milage            string
	Fuel              Fuel
	SubaruModel       SubaruModel
	ServiceBook       bool
	YearOfManufactory uint
}

func newDB() DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Car{})

	return DB{db}
}

func (db DB) touch(url string) bool {
	var car Car
	err := db.db.First(&car, "url = ?", url).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	car.UpdatedAt = time.Now()
	db.db.Updates(&car)

	return true
}

const batch = 16

func reextractAll(db DB) {
	var cars []Car

	log.Println("Start reextarct")

	i := 0
	for {
		err := db.db.Limit(batch).Offset(i).Find(&cars).Error
		if err != nil {
			log.Fatal(err)
		}

		if len(cars) == 0 {
			break
		}

		for _, car := range cars {
			InfoExtractor(&car)
			db.db.Updates(&car)
		}

		i += batch
	}

	log.Println("Reextarct done")

}

func (db DB) add(car Car) {
	db.db.Create(&car)
}

func setupInitialModeInfo(db DB) {
	db.insertOrUpdateModelInfo(ModelInfo{
		ID:           OutbackV36R,
		Name:         "Outback 3.6R",
		YearFrom:     2009,
		YearTo:       2010,
		Volume:       3600,
		EnginePower:  191,
		Transmission: automatic,
	})

	//db.insertOrUpdateModelInfo(ModelInfo{
	//	ID:           OutbackV25,
	//	Name:         "Outback 2.5",
	//	YearFrom:     2009,
	//	YearTo:       2014,
	//	Volume:       2457,
	//	EnginePower:  ,
	//	Transmission: automatic,
	//})

	//db.insertOrUpdateModelInfo(ModelInfo{
	//	ID:           LegacyVKombi25,
	//	Name:         "Legacy 2.5",
	//	YearFrom:     2013,
	//	YearTo:       2015,
	//	Volume:       2498,
	//	EnginePower:  127,
	//	Transmission: cvt,
	//})

	//db.insertOrUpdateModelInfo(ModelInfo{
	//	ID:           LegacyVKombiGT,
	//	Name:         "Legacy GT",
	//	YearFrom:     2009,
	//	YearTo:       2010,
	//	Volume:       2498,
	//	EnginePower:  195,
	//	Transmission: automatic,
	//})

}

func (db DB) insertOrUpdateModelInfo(model ModelInfo) {
	if err := db.db.Model(&model).Where("id = ?", model.ID).Updates(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			db.db.Create(&model)
		}
	}
}
