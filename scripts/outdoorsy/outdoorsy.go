package main

import (
	"fmt"
	"github.com/tang8330/van/lib"
	"log"
	"time"
)

func main() {
	// URL https://search.outdoorsy.com/rentals?raw_json=true&seo_links=true&education=true&average_daily_pricing=true&address=San Francisco, California, United States&bounds[ne]=37.933896600579175,-121.92748823529058&bounds[sw]=37.77581594083472,-122.55651176470212&currency=USD&date[from]=2022-03-22&date[to]=2022-03-31&filter[exclude_type]=utility-trailer,tow-vehicle,other&filter[keywords]=sprinter&filter[type]=camper-van&locale=en-us&page[limit]=24&page[offset]=0&suggested=true

	list, err := lib.List(time.Now(), time.Now().Add(7*24*time.Hour))
	if err != nil {
		log.Fatalf("Something broke. Error: %v", err)
	}

	fmt.Println(fmt.Sprintf("There are %d number of rentals", len(list.Rentals)))

	for _, rental := range list.Rentals {
		fmt.Println(lib.GetBookings(rental.ID, time.Now(), time.Now().Add(4*7*24*time.Hour)))
	}
}
