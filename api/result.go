package api

import (
	"github.com/jinzhu/gorm"
)

type Cookie struct {
	gorm.Model
	ResultID uint
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Expires  int64  `json:"expires"`
}

type Script struct {
	gorm.Model
	ResultID uint
	Value    string
}

type Link struct {
	gorm.Model
	ResultID uint
	Value    string
}

type Header struct {
	gorm.Model
	ResultID uint
	Key      string
	Value    string
}

type HTMLSource struct {
	gorm.Model
	ResultID uint
	Value    string `gorm:"TEXT"`
}

type Result struct {
	gorm.Model
	WebsiteID   WebsiteID   `json:"websiteId"`
	UserID      UserID      `json:"userId"`
	Scripts     []*Script   `json:"scripts"`
	Links       []*Link     `json:"links"`
	Headers     []*Header   `json:"headers"`
	Cookies     []*Cookie   `json:"cookies"`
	HTML        *HTMLSource `json:"html"`
	Description string      `json:"string"`
	LoadedIn    string      `json:"loadedIn"`
	ProcessedIn string      `json:"processedIn"`
	TotalIn     string      `json:"totalIn"`
}

type ResultStorage interface {
	Save(*Result) error
}
