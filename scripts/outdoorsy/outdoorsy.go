package main

import (
	"encoding/csv"
	"fmt"
	"github.com/tang8330/van/lib"
	"log"
	"os"
	"time"
)

func checkError(err error, context string) {
	if err != nil {
		log.Fatalf(fmt.Sprintf("%s, err: %v", context, err))
	}
}

func main() {
	now := time.Now()
	folderName := lib.ToYYYMMDD(now)

	list, err := lib.List(now, now.Add(7*24*time.Hour))
	checkError(err, "failed to list")

	err = os.Mkdir(folderName, 0755)
	checkError(err, "failed to make folder")

	outFile, err := os.Create(fmt.Sprintf("%s/log.txt", folderName))
	checkError(err, "failed to create out file")
	defer outFile.Close()

	csvFile, err := os.Create(fmt.Sprintf("%s/results.csv", folderName))
	checkError(err, "failed to create csv file")
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	outFile.WriteString(fmt.Sprintf("###### This is the %s run ###### \n", lib.ToYYYMMDD(now)))
	outFile.WriteString(fmt.Sprintf("There are %d number of camper vans available for rental \n", len(list.Rentals)))

	csvWriter.Write([]string{
		"id", "date_on_outdoorsy", "make", "year", "daily_rate_cents", "weekly_rate_cents", "monthly_rate_cents", "min_days",
		"security_deposit_cents", "num_miles_free_daily", "per_mile_cents", "location", "num_favorites",
		"this_year_bookings", "this_year_revenue", "last_year_bookings", "last_year_revenue", "url",
	})

	monthsToRentalCount := make(map[string]int)
	lastYearMonthsToRentalCount := make(map[string]int)

	for _, rental := range list.Rentals {
		var lastYearRevenue, currentYearRevenue int
		thisYearBookings, err := lib.GetBookings(rental.ID, lib.StartOfYear(now), lib.EndOfYear(now))
		checkError(err, "failed to get this year bookings")
		lastYear := now.Add(-1 * 365 * 24 * time.Hour)
		lastYearBookings, err := lib.GetBookings(rental.ID, lib.StartOfYear(lastYear), lib.EndOfYear(lastYear))
		checkError(err, "failed to get last year bookings")

		for _, booking := range thisYearBookings {
			currentYearRevenue += booking.ApproximateRevenue(rental.ActiveOption.DailyPrice, rental.ActiveOption.WeekPrice, rental.ActiveOption.MonthPrice)
			month := booking.From.Time().Month().String()
			_, isOk := monthsToRentalCount[month]
			if !isOk {
				monthsToRentalCount[month] = 0
			}

			monthsToRentalCount[month] = monthsToRentalCount[month] + 1
		}

		for _, booking := range lastYearBookings {
			lastYearRevenue += booking.ApproximateRevenue(rental.ActiveOption.DailyPrice, rental.ActiveOption.WeekPrice, rental.ActiveOption.MonthPrice)
			month := booking.From.Time().Month().String()
			_, isOk := lastYearMonthsToRentalCount[month]
			if !isOk {
				lastYearMonthsToRentalCount[month] = 0
			}

			lastYearMonthsToRentalCount[month] = lastYearMonthsToRentalCount[month] + 1
		}

		mileageOverageCents := "N/A"
		if len(rental.MileageOption.Tiers) > 0 {
			mileageOverageCents = fmt.Sprint(rental.MileageOption.Tiers[0].PriceCents)
		}

		csvWriter.Write([]string{
			fmt.Sprint(rental.ID),
			lib.ToYYYMMDD(rental.CreatedAt),
			rental.VehicleMake,
			fmt.Sprint(rental.VehicleYear),
			fmt.Sprint(rental.ActiveOption.DailyPrice),
			fmt.Sprint(rental.ActiveOption.WeekPrice),
			fmt.Sprint(rental.ActiveOption.MonthPrice),
			fmt.Sprint(rental.MinDays),
			fmt.Sprint(rental.SecurityDepositCents),
			fmt.Sprint(rental.MileageOption.Free),
			mileageOverageCents,
			rental.Location.City,
			fmt.Sprint(rental.FavoriteCount),
			fmt.Sprint(len(thisYearBookings)),
			fmt.Sprint(currentYearRevenue / 100),
			fmt.Sprint(len(lastYearBookings)),
			fmt.Sprint(lastYearRevenue / 100),
			rental.URL(),
		})
	}

	sortedMonths := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	outFile.WriteString("This is the raw rental numbers for this region for this year and last year: \n")
	for _, month := range sortedMonths {
		outFile.WriteString(fmt.Sprintf("2022 - %s, value: %v, 2021 - %s, value: %v \n", month, monthsToRentalCount[month], month, lastYearMonthsToRentalCount[month]))
	}
}
