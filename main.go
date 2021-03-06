package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/gocolly/colly"
)

/*
 Webscraper for https://www.factretriever.com/rhino-facts.
The structure of the info I try to get is rapresented by the Fact struct

*/

//Capitalized fields are available outside the package
type Fact struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

// https://www.factretriever.com/rhino-facts
func main() {
	allFacts := make([]Fact, 0) // empty Fact array of size 0

	// collector will contains all info I nees
	collector := colly.NewCollector(
		colly.AllowedDomains("www.factretriever.com", "factretriever.com"), // allowed domains, must not write http://www... or https://www... . Only domain
	)

	// go query selector, my info are inside a div of class=factList and each one is a li. For each element the callback function is executed
	collector.OnHTML(".factsList li", func(element *colly.HTMLElement) {
		factID, err := strconv.Atoi(element.Attr("id")) // for each li inside .factList, I want to save the ID and convert it from string to int
		if err != nil {
			log.Println("Could not get element ID")
		}

		factDesc := element.Text

		fact := Fact{
			Id:          factID,
			Description: factDesc,
		}

		allFacts = append(allFacts, fact)
	})

	// every time the collector make a request on the website, this callback function is executed
	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting: ", request.URL.String())
	})

	// I want to tell the scraper where to start
	collector.Visit("https://www.factretriever.com/rhino-facts")

	writeJSON(allFacts)

}

func writeJSON(data []Fact) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Cannot create file")
		return
	}

	// Undescore (blank identifier) must be used if I don't need the return value of a function

	_ = ioutil.WriteFile("rhinofacts.json", file, 0644) // the file must be created, 0644 is the permission code to use

}
