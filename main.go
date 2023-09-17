package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Create a new Collector
	c := colly.NewCollector()
	// Set up the URL to scrape
	url := "https://www.bcel.com.la/bcel/exchange-rate.html?lang=en"
	// form data for query
	exDate := "2023-09-13"
	round := "1"

	// Get current date
	// today := time.Now().Format("2006-01-02")

	// Attach the form data to the POST request
	c.OnRequest(func(r *colly.Request) {
		r.Method = "POST"
		r.Body = strings.NewReader(fmt.Sprintf("exDate=%s&round=%s", exDate, round))
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
	})

	var tableData [][]string
	c.OnHTML("#fxRateAll tbody tr", func(e *colly.HTMLElement) {
		// Extract and concatenate the data from each table cell
		var rowDataArr []string
		e.ForEach("td", func(_ int, el *colly.HTMLElement) {
			rowDataArr = append(rowDataArr, el.Text)
		})

		// Append the rowData array to tableData
		tableData = append(tableData, rowDataArr)

		// fmt.Fprintf(file, "%s\n", strings.Join(rowDataArr, "\t"))
	})

	// Start the scraping process
	scrErr := c.Visit(url)
	if scrErr != nil {
		log.Println("Error scraping:", scrErr)
		log.Fatal(scrErr)
	}

	// Format the scraped data and write to a json file
	jsonFile, err := os.Create("exchanges-rate.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	if len(tableData) > 0 {
		for _, row := range tableData {
			// Remove index 1 of column in the row which is the contry logo
			row = append(row[:1], row[2:]...)
			// Convert the row data to json
			jsonRow, err := json.Marshal(row)
			if err != nil {
				log.Println("Error converting row to json:", err)
				log.Fatal(err)
			}
			// Write the row data to the json file

			fmt.Fprintf(jsonFile, "%s\n", string(jsonRow))
			// fmt.Println(string(jsonRow))
		}
	}
	fmt.Println("Success, entire page content saved to exchanage-rate.json")
}
