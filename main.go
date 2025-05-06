package main

import (
	"encoding/csv"
	"fmt"
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
	defer f.Close() //close the files before the functions call
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

func main() {
	
	stocks,err := Load("./opg.csv")

	if err != nil {
		return
	}

	fmt.Println(stocks)

}