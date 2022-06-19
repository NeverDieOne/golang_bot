package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Attempt struct {
	Timestamp   float64 `json:"timestamp"`
	Title       string  `json:"lesson_title"`
	Url         string  `json:"lesson_url"`
	IsNegative  bool    `json:"is_negative"`
	SubmittedAt string  `json:"submitted_at"`
}

type Reviews struct {
	FoundTimestamp   float64    `json:"last_attempt_timestamp,omitempy"`
	TimeoutTimestamp float64    `json:"timestamp_to_request,omitempy"`
	Status           string     `json:"status"`
	RequestQuery     [][]string `json:"request_query"`
	Attempts         []Attempt  `json:"new_attempts,omitempty"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	dvmnToken := os.Getenv("DVMN_TOKEN")
	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	chatId := os.Getenv("TG_CHAT_ID")

	c := &http.Client{Timeout: 95 * time.Second}
	timestamp := ""

	for {
		reviews, err := getReviews(c, dvmnToken, timestamp)
		if err != nil {
			log.Println(err)
			continue
		}

		switch reviews.Status {
		case "timeout":
			timestamp = fmt.Sprint(reviews.TimeoutTimestamp)
		case "found":
			timestamp = fmt.Sprint(reviews.FoundTimestamp)
			for _, attempt := range reviews.Attempts {
				message := prepareMessage(attempt)
				if err := sendTelegramNotification(c, tgBotToken, message, chatId); err != nil {
					log.Println(err)
					continue
				}
			}
		default:
			log.Println("Unexpected status: " + reviews.Status)
		}
	}
}

func makeRequest(c *http.Client, method string, url string, headers map[string]string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for header, headerValue := range headers {
		req.Header.Set(header, headerValue)
	}

	q := req.URL.Query()
	for param, paramValue := range params {
		q.Add(param, paramValue)
	}
	req.URL.RawQuery = q.Encode()

	response, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getReviews(c *http.Client, token string, timestamp string) (Reviews, error) {
	url := "https://dvmn.org/api/long_polling/"

	headers := make(map[string]string)
	params := make(map[string]string)

	headers["Authorization"] = "Token " + token
	params["timestamp"] = timestamp

	body, err := makeRequest(c, "GET", url, headers, params)
	if err != nil {
		log.Fatal(err)
	}

	reviews := Reviews{}
	if err := json.Unmarshal(body, &reviews); err != nil {
		log.Fatal(err)
	}

	return reviews, nil
}

func prepareMessage(attempt Attempt) string {
	var result string
	if attempt.IsNegative {
		result = "К сожалению в работе нашлись ошибки."
	} else {
		result = "Отличная работа! Преподаватель её принял!"
	}

	return fmt.Sprintf("Вашу работу '%s' проверили.\n%s\n%s", attempt.Title, result, attempt.Url)
}

func sendTelegramNotification(c *http.Client, token string, text string, chatId string) error {
	url := "https://api.telegram.org/bot" + token + "/sendMessage"

	headers := make(map[string]string)
	params := make(map[string]string)

	params["chat_id"] = chatId
	params["text"] = text

	_, err := makeRequest(c, "GET", url, headers, params)
	return err
}
