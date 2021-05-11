package main

import (
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/ManinderSinghRajput/monitor-market/pkg/config"
)

func main() {

	file, err := ioutil.ReadFile("conf/config.json")
	if err != nil{
		log.Fatal("Unable to read file. Err:" + err.Error())
	}

	ci, err := config.UnmarshalCurrencyInfo(file)
	if err != nil{
		log.Fatal("Unable to Unmarshal. Err: " + err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(len(ci.CurrencyInfo))
	for _, v := range ci.CurrencyInfo{
		//go v.MonitorFromApi(ci.APIKey, &wg)
		go v.MonitorFromWeb(&wg)
		time.Sleep(1*time.Second)
	}
	wg.Wait()
	//ci.CurrencyInfo[0].MonitorFromWeb("https://www.coindesk.com/price/bitcoin")
}
