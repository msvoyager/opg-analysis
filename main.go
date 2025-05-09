//CREATING THE FILE LOADER
//THIS WILL MAKE A SLICE WITH THE DATA EXTRACTED FROM THE CSV FILE

package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
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

		sel := Selection{
			Ticker: stock.Ticker,
			Position: position,
		}

		selections = append(selections, sel)
	}



}