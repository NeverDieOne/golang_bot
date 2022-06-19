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

func getReviews(c *http.Client, t string, u string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", u, nil)
	req.Header.Set("Authorization", "Token "+t)
	if err != nil {
		log.Fatal(err)
	}

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
	var timestamp interface{}

	for {
		reviews, err := getReviews(c, token, url)
		if err != nil {
			log.Println(err)
			continue
		}

		status := reviews["status"]
		switch status {
		case "timeout":
			timestamp = reviews["timestamp_to_request"]
		case "found":
			timestamp = reviews["last_attempt_timestamp"]
			attempts, _ := reviews["new_attempts"].([]interface{})
			log.Println(attempts)
		default:
			log.Println("Unexpected status: " + status.(string))
		}

		log.Println(timestamp)
	}
}
