package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/brad4au57/nhl-schedule-scraper/models"
	"github.com/gocolly/colly"
)

func ScheduleScraper() {
	// Instantiate default collector and scope with restricted domains
	c := colly.NewCollector(colly.AllowedDomains("www.hockey-reference.com", "hockey-reference.com"))

	// Create a slice to store the scraped schedule
	var schedule []models.AnnualSchedule
	// Define a counter for generating IDs
	idCounter := 1

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Requesting from site...")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("The following error occured:", err)
	})

	c.OnHTML("div#all_games tbody tr:not(.thead)", func(h *colly.HTMLElement) {
		// Define variable for my TableRowData struct type
		var row models.AnnualSchedule

		childNodes := h.DOM.Children()

		childNodes.Each(func(index int, node *goquery.Selection) {
			childText := node.Text()
			switch index {
			case 0:
				row.Date = childText
			case 1:
				row.Visitor = node.Find("a").Text()
			case 3:
				row.HomeTeam = node.Find("a").Text()
			}
		})

		// Assign ID from counter and increment the counter
		row.ID = int64(idCounter)
		idCounter++

		schedule = append(schedule, row)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Your scrape was successful")
	})

	// Start scraping on the specific page with the 2024 schedule
	// Could variablize the year in the future
	scheduleURL := "https://www.hockey-reference.com/leagues/NHL_2024_games.html"
	c.Visit(scheduleURL)

	writeJSON(schedule)
}

func writeJSON(data []models.AnnualSchedule) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("Unable to create JSON file")
		return
	}

	_ = os.WriteFile("nhlSchedule.json", file, 0644)
}
