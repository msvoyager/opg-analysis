//CREATING THE FILE LOADER
//THIS WILL MAKE A SLICE WITH THE DATA EXTRACTED FROM THE CSV FILE

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"
)

//struct to group the reading values

type stock struct{
	Ticker string
	Gap float64
	OpeningPrice float64
}

func Load(path string) ([]stock, error) {
	f, err := os.Open(path)
	
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err //exit the main function
	}
	defer f.Close() //close the file as soon as the file open(successfully or with error) and make sure it close before the function end 
	r := csv.NewReader(f)
	rows, err := r.ReadAll()

	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err //exit the main function
	}

	rows = slices.Delete(rows, 0, 1)

	var stocks []stock

	for _, row := range rows {

		ticker := row[0]
		//convert string value of row[1] to float 64 and pass to variable
		gap, err := strconv.ParseFloat(row[1], 64)
		//handling conversion error
		if err != nil {
			continue // skip the row
		} 

		
		openingPrice, err := strconv.ParseFloat(row[2], 64)

		if err != nil {
			continue
		}

		stocks = append(stocks, stock{Ticker: ticker, Gap: gap, OpeningPrice: openingPrice})

	}

	return stocks, nil
}

// How much money in the trading account
var accountBalance = 10000.0

// What percentage of that balance I can tolerate losing
var lossTolerance = .02 // 2%

// Maximum amount I can tolerate losing
var maxLossPerTrade = accountBalance * lossTolerance

// Percentage of the gap I want to take as profit
var profitPercent = .8 // 80%

type Position struct {
    // The price at which to buy or sell
    EntryPrice float64
    // How many shares to buy or sell
    Shares int
    // The price at which to exit and take my profit
    TakeProfitPrice float64
    // The price at which to stop my loss if the stock doesn’t go our way
    StopLossPrice float64
    // Expected final Profit
    Profit float64
}

func Calculate(gapPercent, openingPrice float64) Position {
    closingPrice := openingPrice / (1 + gapPercent)
    gapValue := closingPrice - openingPrice
    profitFromGap := profitPercent * gapValue

    stopLoss := openingPrice - profitFromGap
    takeProfit := openingPrice + profitFromGap

    shares := int(maxLossPerTrade / math.Abs(stopLoss - openingPrice))

    profit := math.Abs(openingPrice - takeProfit) * float64(shares)
    profit = math.Round(profit * 100) / 100

    return Position{
        EntryPrice:      math.Round(openingPrice * 100) / 100,
        Shares:          shares,
        TakeProfitPrice: math.Round(takeProfit * 100) / 100,
        StopLossPrice:   math.Round(stopLoss * 100) / 100,
        Profit:          math.Round(profit * 100) / 100,
    }
}

type Selection struct {
	Ticker string
	Position
	Articles []Article
}

//Fetch the new using API

const (
	url  = "https://seeking-alpha.p.rapidapi.com/news/v2/list-by-symbol?size=5&id="
	apiKeyHeader = "X-Rapidapi-Key"
	apiKey = "3ed75a19b8mshcdda51e3421b503p1df4d5jsn3027836278da"
	

)

type attributes struct {
	PublishOn time.Time `json:"publishOn"`
	Title 	  string `json:"title"`

	//above properties names match the json properties that is enough for go to automatically map the data  but for good practice we use json struct TAG

}

type seekingAlphaNews struct {
	Attributes attributes `json:"attributes"`
}

type seekingAlphaResponse struct {
	Data []seekingAlphaNews `json:"data"` 
}

//we cant pass api response in seekingalpharesponse type around the programme so we make a new type

type Article struct {
	PublishOn time.Time
	Headline  string
}

func FetchNews(ticker string) ([]Article, error) {
	req, err := http.NewRequest(http.MethodGet, url+ticker, nil )

	if err != nil {
		return nil, err
	}


	req.Header.Add(apiKeyHeader, apiKey)
	
	//make the request

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	//http response status code {successful responses 200-299}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("unsuccessful status code recieved %d", resp.StatusCode)
	}

	res := &seekingAlphaResponse{}

	//response id ioReadCloser interface not a string
	//newdecoder return a pointer to json.decoder it has a decode method
	json.NewDecoder(resp.Body).Decode(res)

	var articles []Article
	

	for _, item := range res.Data {
		art := Article{
			PublishOn: item.Attributes.PublishOn,
			Headline: item.Attributes.Title,
		}
		articles = append(articles, art)
	}

	return articles, nil

}



func main() {
	
	stocks,err := Load("./opg.csv")

	if err != nil {
		fmt.Print(err)
		return
	}

	//WE ONLY CONSIDER THE VALUES WITH THE GAP >= 10%(0.1)


	stocks = slices.DeleteFunc(stocks, func(s stock) bool {
		return math.Abs(s.Gap) < .1
	})

	fmt.Println(stocks)

	var selections []Selection

	for _, stock := range stocks{
		position := Calculate(stock.Gap, stock.OpeningPrice)

		articles, err := FetchNews(stock.Ticker)
		if err != nil {
			log.Printf("errpr loading news on %s, %v", stock.Ticker, err)

			continue
		} else {
			log.Printf("Found %d news on %s", len(articles), stock.Ticker)
		}
		sel := Selection{
			Ticker: stock.Ticker,
			Position: position,
			Articles: articles,
		}

		selections = append(selections, sel)
	}

	fmt.Println(selections)

}