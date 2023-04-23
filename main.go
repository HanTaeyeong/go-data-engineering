package main

import (
	"booking-app/functions"
	"booking-app/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyResponse struct {
	Message string `json:"Answer:"`
}

func createApiUrl(date string, startRow int, endRow int) string {
	var openApiUrl = "http://211.237.50.150:7080/openapi"
	var apiKey = "d2c63e5d185d5f97f718aa884df70a8d34813ff4808f97a2c2808e585b5fe327"
	var dataType = "json"
	var dataGrid = fmt.Sprintf("Grid_20161221000000000429_1/%d/%d", startRow, endRow)

	var url = fmt.Sprintf("%s/%s/%s/%s?AUCNG_DE=%s", openApiUrl, apiKey, dataType, dataGrid, date)
	return url
}

func fetchAPIAndStoreToS3(today string) {
	utils.ReadAndSetEnvFile(".env")

	const BatchSize = 1000
	const RowLimit = 10e6

	var url = createApiUrl(today, 1, 1)

	var data = functions.FetchAPI(url).Grid_20161221000000000429_1
	var totalCnt = data.TotalCnt

	if totalCnt > RowLimit {
		log.Fatal("RowLimit exceeds!")
		return
	}

	var results []utils.MarketDataRowType
	resultChannel := make(chan utils.MargetDataGrid, totalCnt/BatchSize+1)
	errors := make(chan error, totalCnt/BatchSize+1)

	for i := 0; i < totalCnt/BatchSize+1; i++ {
		var url = createApiUrl(today, i*BatchSize+1, i*BatchSize+BatchSize)
		go functions.FetchAPIAsync(url, resultChannel, errors)
	}

	for i := 0; i < totalCnt/BatchSize+1; i++ {
		select {
		case result := <-resultChannel:
			b, err := json.Marshal(result)
			if err != nil {
				log.Fatal("Data Marshaling failed!")
			}
			functions.UploadToS3("datas/"+today+"/"+fmt.Sprint(i), string(b))
			results = append(results, result.Row...)
		case err := <-errors:
			log.Println(err)
		}
	}

	fmt.Println(totalCnt)

}

func HandleLambdaEvent(event any) (MyResponse, error) {

	var today = time.Now().AddDate(0, 0, -1).Format("20060102150405")
	today = today[:8]

	fetchAPIAndStoreToS3(today)
	return MyResponse{Message: fmt.Sprintf("upload success! %s ", today)}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
