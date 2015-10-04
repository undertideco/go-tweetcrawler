package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	TweetsFileName string `json:"tweetsFileName"`
	ConsumerKey    string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
	AccessKey      string `json:"accessKey"`
	AccessSecret   string `json:"accessSecret"`
	TargetUsername string `json:"targetUsername"`
}

var oldestTweetId int64 = -1
var params = url.Values{}
var api *anaconda.TwitterApi
var cfg Config

func main() {
	cfg_file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("Error reading config.json")
		return
	}

	json.Unmarshal(cfg_file, &cfg)
	api = anaconda.NewTwitterApi(cfg.AccessKey, cfg.AccessSecret)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	anaconda.SetConsumerKey(cfg.ConsumerKey)
	anaconda.SetConsumerSecret(cfg.ConsumerSecret)

	if _, err := os.Stat(cfg.TweetsFileName); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", cfg.TweetsFileName)
		newFile, _ := os.Create(cfg.TweetsFileName)
		newFile.Close()
	}

	params.Set("screen_name", cfg.TargetUsername)
	params.Add("count", strconv.Itoa(200)) // maximum of 200 tweets per request

	initialCrawl()
}

func initialCrawl() {
	log.Println("Starting crawl...")
	file, err := os.Open(cfg.TweetsFileName)
	if err != nil {
		log.Println("Error: ", err)
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading all lines: %v", err)
	}

	file.Close()

	if len(lines) > 0 {
		lastTweetIdString := lines[len(lines)-1][0]

		if s, err := strconv.ParseInt(lastTweetIdString, 10, 64); err == nil {
			oldestTweetId = s
		}
		oldestTweetId, _ = strconv.ParseInt(lines[0][0], 10, 64)
		log.Printf("Starting crawl from %d\n", oldestTweetId)

		afterLastTimeTweets, err := multiCrawl(-1, oldestTweetId)

		if err != nil {
			log.Println(err)
		}

		saveCrawls(afterLastTimeTweets)

		log.Printf("[SINCE LAST TIME CRAWL] %d tweets in total were retrieved\n", len(afterLastTimeTweets))

	} else {
		log.Println("Fresh Crawl")

		beforeNowTweets, err := multiCrawl(-1, -1)
		if err != nil {
			log.Println(err)
		}

		saveCrawls(beforeNowTweets)

		log.Printf("[INITIAL CRAWL] %d tweets in total were retrieved\n", len(beforeNowTweets))
	}

	//Interval crawl
	t := time.NewTicker(time.Duration(5) * time.Minute)
	for {
		intervalTweets, err := multiCrawl(-1, oldestTweetId)
		if err != nil {
			log.Println(err)
		}
		saveCrawls(intervalTweets)
		log.Printf("[INTERVAL CRAWL] %d tweets in total were retrieved\n", len(intervalTweets))
		<-t.C
	}
}

//Crawl tweets with max_id or since_id until finished. Crawls 200 at a time and returns everything.
func multiCrawl(max_id, since_id int64) ([][]string, error) {
	allTweets := make([][]string, 0)
	tweets, err := crawl(max_id, since_id)
	if err != nil {
		log.Println(err)
		return [][]string{}, err
	}

	for len(tweets) > 1 {
		allTweets = append(allTweets, tweets...)
		if since_id != -1 {
			oldestTweetId, _ = strconv.ParseInt(allTweets[0][0], 10, 64)
			tweets, err = crawl(-1, oldestTweetId)
		} else {
			oldestTweetId, _ = strconv.ParseInt(allTweets[len(allTweets)-1][0], 10, 64)
			tweets, err = crawl(oldestTweetId, -1)
		}

		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("[MULTI-CRAWLER] Progress: %d tweets retrieved\n", len(allTweets))
	}

	return allTweets, nil
}

//Crawls tweets with max_id or since_id. Max 200, according to Twitter's API.
func crawl(max_id, since_id int64) ([][]string, error) {
	if max_id != -1 {
		params.Set("max_id", strconv.FormatInt(max_id, 10))
	}
	if since_id != -1 {
		params.Set("since_id", strconv.FormatInt(since_id, 10))
	}

	tweets, err := api.GetUserTimeline(params)

	if err != nil {
		log.Println(err)
		return [][]string{}, err
	}

	return parseToString2dArray(tweets), nil
}

//Saves new tweets. Appends to existing tweets.
func saveCrawls(tweets [][]string) {
	file, err := os.Open(cfg.TweetsFileName)
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading all lines: %v", err)
	}

	file.Close()

	lines = append(lines, tweets...)

	file, err = os.Create(cfg.TweetsFileName)
	if err != nil {
		log.Println("Error: ", err)
		return
	}

	// withHeader := make([][]string, len(lines) + 1)
	// withHeader = append(withHeader, []string{"id", "text", "createdAt"})
	// withHeader = append(withHeader, lines...)

	writer := csv.NewWriter(file)
	writer.WriteAll(lines)
	if err := writer.Error(); err != nil {
		log.Println("Error Writing CSV: ", err)
	}
	file.Close()
}

//Parses anaconda library tweets into a 2d array
func parseToString2dArray(tweets []anaconda.Tweet) [][]string {
	newTweetStringsArray := make([][]string, len(tweets))
	for i, tweet := range tweets {
		newTweetStringsArray[i] = []string{tweet.IdStr, tweet.Text, tweet.CreatedAt}
	}

	return newTweetStringsArray
}
