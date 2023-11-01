package models

import (
	"log"
	"time"

	"github.com/brad4au57/nhl-schedule-scraper/database"
)

// AnnualSchedule is a struct model containing an ID, Date, Visitor, and HomeTeam
type AnnualSchedule struct {
	ID       int64
	Date     string
	Visitor  string
	HomeTeam string
}

func GetFullSchedule() []AnnualSchedule {
	// Query the database to retrieve data
	query := "SELECT id, date, visitor_team, home_team FROM schedules"
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Create a slice to store the retrieved data
	var existingSchedule []AnnualSchedule

	// Iterate over the rows and scan the data into the struct
	for rows.Next() {
		var data AnnualSchedule
		if err := rows.Scan(&data.ID, &data.Date, &data.Visitor, &data.HomeTeam); err != nil {
			log.Fatal(err)
		}

		// Convert the timestamp to the desired format
		parsedTime, err := time.Parse(time.RFC3339, data.Date)
		if err != nil {
			log.Fatal(err)
		}
		data.Date = parsedTime.Format("2006-01-02")

		existingSchedule = append(existingSchedule, data)
	}

	return existingSchedule
}
