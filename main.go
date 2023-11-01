package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/brad4au57/nhl-schedule-scraper/database"
	"github.com/brad4au57/nhl-schedule-scraper/models"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	// Connect to Postgres DB
	err := database.InitPostgresDb(os.Getenv("DB_CONNECT_LOCAL"))
	if err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	// Run web scrapper
	ScheduleScraper()

	// Query database for the schedule saved in the DB
	var existingSchedule []models.AnnualSchedule
	existingSchedule = models.GetFullSchedule()

	// Get scraped schedule data from newly created JSON file
	jsonData, err := os.ReadFile("nhlSchedule.json")
	if err != nil {
		log.Fatal(err)
	}
	var scheduleData []models.AnnualSchedule
	if err := json.Unmarshal(jsonData, &scheduleData); err != nil {
		log.Fatal(err)
	}

	// Now, you can compare the scraped schedule data with the data retrieved from the database
	for _, scraped := range scheduleData {
		// Set default for whether the scraped data is found in database data
		found := false

		for _, dbData := range existingSchedule {
			if scraped.ID == dbData.ID {
				found = true
				// Check for differences and update if needed
				if scraped.Date != dbData.Date || scraped.Visitor != dbData.Visitor || scraped.HomeTeam != dbData.HomeTeam {
					// Update the corresponding row in the database
					_, err := database.DB.Exec("UPDATE schedules SET date = $1, visitor_team = $2, home_team = $3 WHERE id = $4", scraped.Date, scraped.Visitor, scraped.HomeTeam, scraped.ID)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Data for ID: %d updated in the database.\n", scraped.ID)
				} else {
					fmt.Printf("Data matched for ID: %d \n", scraped.ID)
				}
				break
			}
		}
		if !found {
			// Insert a new row into the database
			_, err := database.DB.Exec("INSERT INTO schedules (date, visitor_team, home_team) VALUES ($1, $2, $3)", scraped.Date, scraped.Visitor, scraped.HomeTeam)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("New data inserted into the database.")
		}
	}
}
