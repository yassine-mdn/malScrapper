package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/tidwall/pretty"
)

func main() {
	userName := flag.String("userName", "error", "Specify myanimelist user name")
	status := flag.String("status", "ptw", "can either be 'ptw' for plan to watch or 'w' for watching or 'd' for dropped or 'c' for completed\nplan to watch by default")
	mode := flag.String("mode", "all", "can either be 'all' or 'rand'\nall by default")

	flag.Parse()

	// Instantiate default collector
	c := colly.NewCollector(
	//colly.AllowedDomains("https://myanimelist.net"),
	) //only domain allowed is myanimelist ma bghatch tkhdem
	c.SetRequestTimeout(30 * time.Second)
	var animelist []string
	c.OnHTML("table.list-table", func(e *colly.HTMLElement) {
		rawList := e.Attr("data-items")
		tempSlices := strings.Split(Decluter(rawList), "\n\n")
		animelist = append(animelist, tempSlices...)
		if strings.Contains(*mode, "rand") {
			printRandom(animelist, len(animelist))
		} else {
			printAll(animelist)
		}

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL.String())
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("User name not found")
	})
	if *status == "ptw" {
		c.Visit("https://myanimelist.net/animelist/" + *userName + "?status=6")
	} else if *status == "w" {
		c.Visit("https://myanimelist.net/animelist/" + *userName + "?status=1")
	} else if *status == "c" {
		c.Visit("https://myanimelist.net/animelist/" + *userName + "?status=2")
	} else if *status == "d" {
		c.Visit("https://myanimelist.net/animelist/" + *userName + "?status=4")
	} else {
		fmt.Println("invalid argument")
	}

}

func printAll(animelist []string) {
	for _, s := range animelist {
		fmt.Println(s)
	}
}

func printRandom(animelist []string, size int) {
	fmt.Println("\nmaybe this could be a great choice:")
	fmt.Println(animelist[rand.Intn(size)])
}

func Decluter(s string) string {
	tempSlices := strings.Split(string(pretty.Pretty([]byte(s))), "\n")
	var newString string
	for _, subString := range tempSlices {
		if strings.Contains(subString, "\"anime_title\"") {
			newString += subString + "\n"
		} else if strings.Contains(subString, "\"anime_num_episodes\"") {
			newString += subString + "\n\n"
		}
	}
	return newString
}
