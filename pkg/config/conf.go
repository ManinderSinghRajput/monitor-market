package config

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func UnmarshalCurrencyInfo(data []byte) (CurrencyInfo, error) {
	var r CurrencyInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CurrencyInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CurrencyInfo struct {
	APIKey       string                `json:"apiKey"`
	Currencies   string                `json:"currencies"`
	ToCurrency   string                `json:"toCurrency"`
	CurrencyInfo []CurrencyInfoElement `json:"currencyInfo"`
}

type CurrencyInfoElement struct {
	FromCurrency string `json:"fromCurrency"`
	UpperLimit   string `json:"upperLimit"`
	LowerLimit   string `json:"lowerLimit"`
}

type APIResponse []APIResponseElement

func UnmarshalAPIResponse(data []byte) (APIResponse, error) {
	var r APIResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *APIResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *APIResponse) MarshalIndent() ([]byte, error) {
	return json.MarshalIndent(r, "", " ")
}

type APIResponseElement struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Timestamp string `json:"price_timestamp"`
}

func UnmarshalResponse(data []byte) (Response, error) {
	var r Response
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Response struct {
	Name                string              `json:"name"`
	Currency            string              `json:"currency"`
	CurrentExchangeRate CurrentExchangeRate `json:"currentExchangeRate"`
}

type CurrentExchangeRate struct {
	Price         float64 `json:"price"`
	PriceCurrency string  `json:"priceCurrency"`
}

func (r *CurrencyInfo) MonitorFromApi(wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("https://api.nomics.com/v1/currencies/ticker?key=%s&ids=%s&interval=1d&convert=%s", r.APIKey, strings.ToUpper(r.Currencies), r.ToCurrency)

	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err.Error())
			break
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err.Error())
			break
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err.Error())
			break
		}
		err = res.Body.Close()
		if err != nil {
			log.Println(err.Error())
			break
		}
		response,  err := UnmarshalAPIResponse(body)
		if err != nil {
			log.Println(err.Error())
			break
		}
		b, err := response.MarshalIndent()
		if err != nil {
			log.Println(err.Error())
			break
		}
		fmt.Println(string(b))
		//fmt.Println(string(body))
		time.Sleep(1*time.Minute)
	}
	return
}

func (r *CurrencyInfoElement) MonitorFromWeb(wg *sync.WaitGroup) {
	url := fmt.Sprintf("https://www.coindesk.com/price/%s", r.FromCurrency)
	defer wg.Done()
	for{
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err.Error())
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err.Error())
			return
		}

		z := html.NewTokenizer(res.Body)

		var response Response

		for {
			_ = z.Next()
			t := z.Token()

			if strings.Contains(t.Data, "ExchangeRateSpecification"){
				response, err = UnmarshalResponse([]byte(t.Data))
				if err != nil{
					log.Println("Response found but unable to unmarshal")
				}
				break
			}
		}
		fmt.Println(time.Now(), ": ",response)
		time.Sleep(10 * time.Second)
	}
}
