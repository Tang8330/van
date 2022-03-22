package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ListResponse struct {
	Rentals []Rental `json:"data"`
}

type Location struct {
	City string `json:"city"`
}

type Tier struct {
	PriceCents uint `json:"price"`
}

type MileageOption struct {
	Free  uint   `json:"included"`
	Tiers []Tier `json:"tiers"`
}

type Booking struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Rental struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"filtered_description"`
	FavoriteCount uint   `json:"favorite_count"`

	// Where's the car parked
	Location Location `json:"location"`

	VehicleMake       string        `string:"vehicle_make"`
	VehicleYear       uint          `json:"vehicle_year"`
	MinDays           uint          `json:"minimum_days"`
	DailyPriceCents   uint          `json:"price_per_day"`
	WeeklyPriceCents  uint          `json:"price_per_week"`
	MonthlyPriceCents uint          `json:"price_per_month"`
	MileageOption     MileageOption `json:"mileage_usage_item"`
}

func List(start, end time.Time) (*ListResponse, error) {
	// I only care about Sprinter Vans.
	urlString := fmt.Sprintf("https://search.outdoorsy.com/rentals?raw_json=true&seo_links=true&education=true&average_daily_pricing=true&address=San Francisco, California, United States&bounds[ne]=37.933896600579175,-121.92748823529058&bounds[sw]=37.77581594083472,-122.55651176470212&currency=USD&date[from]=%s&date[to]=%s&filter[exclude_type]=utility-trailer,tow-vehicle,other&filter[keywords]=sprinter&filter[type]=camper-van&locale=en-us&page[limit]=500&page[offset]=0&suggested=true", ToYYYMMDD(start), ToYYYMMDD(end))
	urlEncodedString, err := EncodeURL(urlString)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Get(urlEncodedString)
	if err != nil {
		return nil, err
	}

	listResponseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var listResp ListResponse
	err = json.Unmarshal(listResponseBytes, &listResp)
	return &listResp, err
}

func GetBookings(rentalID uint, start, end time.Time) ([]Booking, error) {
	var bookings []Booking
	url := fmt.Sprintf("https://api.outdoorsy.com/v0/availability?from=%s&to=%s&rental_id=%d",
		ToYYYMMDD(start), ToYYYMMDD(end), rentalID)

	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBytes, &bookings)
	return bookings, err
}
