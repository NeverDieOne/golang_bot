package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func getReviews(c *http.Client, token string, url string, timestamp string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Token "+token)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("timestamp", timestamp)
	req.URL.RawQuery = q.Encode()

	response, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	reviews := map[string]interface{}{}
	if err := json.Unmarshal(body, &reviews); err != nil {
		log.Fatal(err)
	}

	return reviews, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("DVMN_TOKEN")
	c := &http.Client{Timeout: 95 * time.Second}
	url := "https://dvmn.org/api/long_polling/"
	timestamp := ""

	for {
		reviews, err := getReviews(c, token, url, timestamp)
		if err != nil {
			log.Println(err)
			continue
		}

		status := reviews["status"]
		switch status {
		case "timeout":
			timestamp, _ = reviews["timestamp_to_request"].(string)
		case "found":
			timestamp, _ = reviews["timestamp_to_request"].(string)
			attempts, _ := reviews["new_attempts"].([]interface{})
			log.Println(attempts)
		default:
			log.Println("Unexpected status: " + status.(string))
		}
	}
}
