package main

import (
	"log"
	"os"
	"encoding/csv"
	"net/url"
	"fmt"
	"strconv"
	"github.com/ChimeraCoder/anaconda"
)

type savedtweet struct {
	id int
	tweetString string
	tweetedDate string
}

const tweetsFileName string = "tweets.csv"

const consumerKey = "5jAS8RyLoCFH1eZOZ6VZYA"
const consumerSecret = "6UJAQl7WbBexqJXCvPaIOQxf1eMapr4CFqENSWpiSoA"
const accessKey = "64090321-yFCT66IFZOGhm0qVE34CdhBpSBwo8EO4ftvK1Dfzz"
const accessSecret = "LRrzFJXFvzJGvGSXXr3cPVpmeg78tcv5XueAVwnBpEs"

var oldestTweetId int64 = -1
var params = url.Values{}
var api = anaconda.NewTwitterApi(accessKey, accessSecret)

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)

	if _, err := os.Stat(tweetsFileName); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", tweetsFileName)
		newFile, _ := os.Create(tweetsFileName)
		newFile.Close()
	}

	params.Set("screen_name", "ShintaroTay")
	params.Add("count", strconv.Itoa(200)) // maximum of 200 tweets per request

	initialCrawl()
}

func initialCrawl() {
	log.Println("Starting crawl...")
	file, err := os.Open(tweetsFileName)
	if err != nil {
		log.Println("Error: ", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading all lines: %v", err)
	}

	if len(lines) > 0 {
		lastTweetIdString := lines[len(lines) - 1][0]

		if s, err := strconv.ParseInt(lastTweetIdString, 10, 64); err == nil {
			oldestTweetId = s
		}
		log.Printf("Starting crawl from %d\n", oldestTweetId)
	} else {
		log.Println("Fresh Crawl")
	}

	crawlFrom(oldestTweetId)
}

func crawlFrom(tweetId int64) {
	newTweets, err := api.GetUserTimeline(params)

	if err != nil {
		log.Println(err)
	}

	for len(newTweets) > 0 {
		fmt.Printf("Getting tweets before %d\n", oldestTweetId)
		if oldestTweetId != -1 {
			// all subsequent requests use the max_id param to prevent duplicates
			params.Set("max_id", strconv.FormatInt(oldestTweetId, 10))
		}

		newTweets, err = api.GetUserTimeline(params)
		if err != nil {
			log.Println(err)
		}

		file, err := os.Open(tweetsFileName)
		if err != nil {
			log.Println("Error: ", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		lines, err := reader.ReadAll()
		if err != nil {
			log.Println("Error reading all lines: %v", err)
		}
		newTweetsStringArray := parseToString2dArray(newTweets)
		lines = append(lines, newTweetsStringArray...)

		writer := csv.NewWriter(file)
		writer.WriteAll(lines)
		if err := writer.Error(); err != nil {
			log.Println("Error Writing CSV: ", err)
		}

		if s, err := strconv.ParseInt(lines[len(lines) - 1][0], 10, 64); err == nil {
			oldestTweetId = s
		}

		log.Printf("%d tweets downloaded so far\n", len(lines))
	}
}

func parseToString2dArray(tweets []anaconda.Tweet) (tweetStringsArray [][]string) {
	newTweetStringsArray := make([][]string, len(tweets))
	for i, _ := range tweets {
		tweetStringsArray[i] = make([]string, 3)
		tweetStringsArray[i][0] = tweet.IdStr
		tweetStringsArray[i][1] = tweet.Text
		tweetStringsArray[i][2] = tweet.CreatedAt
	}

	return newTweetStringsArray
}