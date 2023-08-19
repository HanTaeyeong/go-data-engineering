package functions

import (
	"data-app/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func FetchAPI(url string) utils.MarketDataType {

	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)

	var responseData utils.MarketDataType
	json.Unmarshal(b, &responseData)

	return responseData
}

func FetchAPIAsync(url string, resultChannel chan utils.MargetDataGrid, errors chan error) {

	response, err := http.Get(url)
	if err != nil {
		errors <- err
		os.Exit(1)
	}
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)

	var responseData utils.MarketDataType
	json.Unmarshal(b, &responseData)
	resultChannel <- responseData.Grid_20161221000000000429_1

}
