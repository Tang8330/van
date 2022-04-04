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
	City  string `json:"city"`
	State string `json:"state"`
}

type Tier struct {
	PriceCents uint `json:"price"`
}

type MileageOption struct {
	Free  uint   `json:"included"`
	Tiers []Tier `json:"tiers"`
}

type Booking struct {
	From TimeFormatYYYYMMDD `json:"from"`
	To   TimeFormatYYYYMMDD `json:"to"`
}

// This all came from
// https://search.outdoorsy.com/rentals?raw_json=true&seo_links=true&education=true&average_daily_pricing=true&address=San%20Francisco%2C%20California%2C%20United%20States&bounds[ne]=37.933896600579175%2C-121.92748823529058&bounds[sw]=37.77581594083472%2C-122.55651176470212&currency=USD&date[from]=2022-03-22&date[to]=2022-03-31&filter[exclude_type]=utility-trailer%2Ctow-vehicle%2Cother&filter[keywords]=sprinter&filter[type]=camper-van&locale=en-us&page[limit]=24&page[offset]=0&suggested=true

type ActiveOption struct {
	DailyPrice int `json:"day_price"`
	WeekPrice  int `json:"week_price"`
	MonthPrice int `json:"month_price"`
}

type Rental struct {
	ID            uint      `json:"id"`
	Name          string    `json:"filtered_name"`
	CreatedAt     time.Time `json:"created"`
	Description   string    `json:"filtered_description"`
	FavoriteCount uint      `json:"favorite_count"`

	// Where's the car parked
	Location Location `json:"location"`

	VehicleMake string `json:"vehicle_make"`
	VehicleYear uint   `json:"vehicle_year"`
	MinDays     uint   `json:"minimum_days"`
	// Don't use, use ActiveOption instead
	DeprecatedDailyPriceCents   int `json:"price_per_day"`
	DeprecatedWeeklyPriceCents  int `json:"price_per_week"`
	DeprecatedMonthlyPriceCents int `json:"price_per_month"`

	SecurityDepositCents uint          `json:"security_deposit"`
	MileageOption        MileageOption `json:"mileage_usage_item"`
	ActiveOption         ActiveOption  `json:"active_options"`
}

func List() (*ListResponse, error) {
	// I only care about Sprinter Vans.
	urlString := "https://search.outdoorsy.com/rentals?raw_json=true&seo_links=true&education=true&average_daily_pricing=true&address=San Francisco, California, United States&bounds[ne]=37.933896600579175,-121.92748823529058&bounds[sw]=37.77581594083472,-122.55651176470212&currency=USD&filter[exclude_type]=utility-trailer,tow-vehicle,other&filter[keywords]=sprinter&filter[type]=camper-van&locale=en-us&page[limit]=500&page[offset]=0&suggested=true"
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

func (b *Booking) NumberOfNights() int {
	return int(b.To.Time().Sub(b.From.Time()).Hours() / 24)
}

func (b *Booking) ApproximateRevenue(dailyCents, weeklyTotalCents, monthlyTotalCents int) int {
	// Using number of nights, let's now apply the monthly and weekly discount.
	nights := b.NumberOfNights()
	var numOfMonths, numOfWeeks, numOfDaysRemain int
	numOfMonths = nights / 30
	numOfWeeks = (nights - (numOfMonths * 30)) / 7
	numOfDaysRemain = nights - (numOfMonths * 30) - (numOfWeeks * 7)
	return numOfDaysRemain*dailyCents + numOfMonths*monthlyTotalCents + numOfWeeks*weeklyTotalCents
}

func (r *Rental) URL() string {
	// https://www.outdoorsy.com/rv-rental/mill-valley_ca/2021_mercedes-benz_sprinter_253662-listing?
	return fmt.Sprintf("https://www.outdoorsy.com/rv-rental/%s_%s/%s_%d-listing?", SpaceToDash(r.Location.City), SpaceToDash(r.Location.State), SpaceToDash(r.Name), r.ID)
}
