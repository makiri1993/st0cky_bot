package api

import (
	"bing-news-api/models"
	"bing-news-api/setup"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetBingNews(searchTerm string) *models.BingNews {
	req, err := http.NewRequest("GET", setup.Endpoint, nil)

	if err != nil {
		panic(err)
	}

	param := req.URL.Query()
	param.Add("q", searchTerm)
	req.URL.RawQuery = param.Encode()

	req.Header.Add("Ocp-Apim-Subscription-Key", setup.Token)

	client := new(http.Client)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	ans := new(models.BingNews)
	err = json.Unmarshal(body, &ans)

	if err != nil {
		fmt.Print(err)
	}
	return ans
}
