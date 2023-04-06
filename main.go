package main

import (
	"booking-app/functions"
	"booking-app/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyResponse struct {
	Message string `json:"Answer:"`
}

func createApiUrl(date string, startRow int, endRow int) string {
	var openApiUrl = "http://211.237.50.150:7080/openapi"
	var apiKey = os.Getenv("apiKey")
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

	for i := 0; i <= totalCnt/BatchSize+1; i++ {
		var url = createApiUrl(today, i*BatchSize+1, i*BatchSize+BatchSize)

		var data = functions.FetchAPI(url).Grid_20161221000000000429_1
		results = append(results, data.Row...)
	}
	b, err := json.Marshal(results)
	if err != nil {
		log.Fatal("Data Marshaling failed!")
	}
	fmt.Println(totalCnt, len(results))
	functions.UploadToS3("data-"+today, string(b))
}

func HandleLambdaEvent(event any) (MyResponse, error) {

	var today = time.Now().AddDate(0, 0, -1).Format("20060102150405")
	today = today[:8]

	fetchAPIAndStoreToS3(today)
	return MyResponse{Message: fmt.Sprintf("upload success! %s ", today)}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)

	//goroutine.GoRoutinePractice()
}
