package lib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type StatsResponse struct {
	DailyRate uint
	WeeklyRate uint
	MonthlyRate uint
	MinimumStay uint
	RefundableSecurityDeposit uint
	Mileage uint
	MileageOverageFee uint
	Location string
	// 60 days out, how many nights are booked?
	TwoMonthsUtilization uint
}

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
	Free uint `json:"included"`
	Tiers []Tier `json:"tiers"`
}

type Rental struct {
	Name string `json:"name"`
	Description string `json:"filtered_description"`
	FavoriteCount uint `json:"favorite_count"`

	// Where's the car parked
	Location Location `json:"location"`

	VehicleYear uint `json:"vehicle_year"`
	MinDays uint `json:"minimum_days"`
	DailyPriceCents uint `json:"price_per_day"`
	WeeklyPriceCents uint `json:"price_per_week"`
	MonthlyPriceCents uint `json:"price_per_month"`
	MileageOption MileageOption `json:"mileage_usage_item"`
}

type OutdoorsyAPIResponse struct {
	DealerMinDays uint `json:"dealer_minimum_days"`

}

func List(start, end time.Time) (*ListResponse, error) {
	start.Format("YYY")

	urlString := "https://search.outdoorsy.com/rentals?raw_json=true&seo_links=true&education=true&average_daily_pricing=true&address=San%20Francisco%2C%20California%2C%20United%20States&bounds[ne]=37.933896600579175%2C-121.92748823529058&bounds[sw]=37.77581594083472%2C-122.55651176470212&currency=USD&date[from]=2022-03-22&date[to]=2022-03-31&filter[exclude_type]=utility-trailer%2Ctow-vehicle%2Cother&filter[keywords]=sprinter&filter[type]=camper-van&locale=en-us&page[limit]=24&page[offset]=0&suggested=true"
	resp, err := http.DefaultClient.Get(urlString)
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

// https://api.outdoorsy.com/v0/quotes/401bede0-aa0a-11ec-95f0-46468eaab59c
//func GetStats(car string) StatsResponse {
//}