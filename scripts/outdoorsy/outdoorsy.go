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
	outFile.WriteString(fmt.Sprintf("There are %d number of rentals \n", len(list.Rentals)))

	csvWriter.Write([]string{
		"id", "daily_rate_cents", "weekly_rate_cents", "monthly_rate_cents", "min_days",
	})

	for _, rental := range list.Rentals {
		bookings, err := lib.GetBookings(rental.ID, now, lib.EndOfYear(now))
		checkError(err, "failed to get bookings")
		csvWriter.Write([]string{
			fmt.Sprint(rental.ID),
			fmt.Sprint(rental.DailyPriceCents),
			fmt.Sprint(rental.WeeklyPriceCents),
			fmt.Sprint(rental.MonthlyPriceCents),
			fmt.Sprint(rental.MinDays),
		})

		fmt.Println("Bookings", len(bookings))
	}
}
