package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	token := os.Getenv("DVMN_TOKEN")
	c := &http.Client{Timeout: 95 * time.Second}
	url := "https://dvmn.org/api/long_polling/"

	for {
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Token "+token)
		if err != nil {
			log.Fatal(err)
		}

		response, err := c.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(body))
	}
}
