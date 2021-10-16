package models

import (
	"strconv"
	"time"
	internalError "yula/internal/error"
)

type Advert struct {
	Id          int64     `json:"id" valid:"-"`
	Name        string    `json:"name" valid:"type(string),stringlength(1|100)"`
	Description string    `json:"description" valid:"optional,stringlength(1|2000)"`
	Price       int       `json:"price" valid:"optional,type(int)"`
	Location    string    `json:"location" valid:"type(string)"`
	Latitude    float64   `json:"latitude" valid:"latitude"`
	Longitude   float64   `json:"longitude" valid:"longitude"`
	PublishedAt time.Time `json:"published_at" valid:"-"`
	DateClose   time.Time `json:"date_close" valid:"-"`
	IsActive    bool      `json:"is_active" valid:"-"`
	PublisherId int64     `json:"publisher_id" valid:"-"`
	Category    string    `json:"category" valid:"type(string)"`
	Images      []string  `json:"images" valid:"-"`
	Views       int64     `json:"views" valid:"-"`
}

type AdvertShort struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Location string `json:"location"`
	Image    string `json:"image"`
}

func (a *Advert) ToShort() *AdvertShort {
	var imageStr string
	if len(a.Images) == 0 {
		imageStr = ""
	} else {
		imageStr = a.Images[0]
	}
	return &AdvertShort{
		Id: a.Id, Name: a.Name, Price: a.Price, Location: a.Location, Image: imageStr,
	}
}

type Page struct {
	PageNum int64
	Count   int64
}

const (
	DefaultPageNum     int64 = 1
	DefaultCountAdvert int64 = 50
)

func NewPage(pageNumS string, countS string) (*Page, error) {
	pageNum := DefaultPageNum
	count := DefaultCountAdvert
	var err error
	if pageNumS != "" {
		pageNum, err = strconv.ParseInt(pageNumS, 10, 64)
		if err != nil {
			return nil, internalError.BadRequest
		}
	}

	if countS != "" {
		count, err = strconv.ParseInt(countS, 10, 64)
		if err != nil {
			return nil, internalError.BadRequest
		}
	}

	return &Page{PageNum: pageNum - 1, Count: count}, nil
}
