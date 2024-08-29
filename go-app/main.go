package main

import (
	"encoding/json"
	"net/http"
)

type (
	InputData struct {
		Amount   float64 `json:"amount`
		Currency string  `json:"currency"`
	}

	OutputData struct {
		Total float64 `json:"total"`
	}
)

func main() {
	http.HandleFunc("/convert", func(writer http.ResponseWriter, request *http.Request) {
		data := new(InputData)
		err := json.NewDecoder(request.Body).Decode(data)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		url := "https://open.er-api.com/v6/latest/USD"
		response, err := http.Get(url)
		if err != nil || response.StatusCode != 200 {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		exchangeRates := make(map[string]interface{})
		err = json.NewDecoder(response.Body).Decode(&exchangeRates)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		rates := exchangeRates["rates"].(map[string]interface{})
		if rate, ok := rates[data.Currency].(float64); ok {
			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(OutputData{Total: data.Amount * rate})
		} else {
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}
	})
	http.ListenAndServe(":8080", nil)
}
