package main

import (
	"data-app/functions"
	"data-app/utils"
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
	var openApiUrl = os.Getenv("API_HOST")
	var apiKey = os.Getenv("API_KEY")
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
	// var today = time.Now().AddDate(0, 0, -1).Format("20060102150405")
	// today = today[:8]

	// fetchAPIAndStoreToS3(today)

	lambda.Start(HandleLambdaEvent)
}
