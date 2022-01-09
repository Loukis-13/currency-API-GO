package currencies

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

var Currencies Response

func init() {
	fmt.Println("Requisitando API...")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.currency-api.com/rates?base=BRL", nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err.Error())
	}

	json.Unmarshal(bodyBytes, &Currencies)
}

func Trocar(moeda1, moeda2 string) (float64, error) {
	if Currencies.Rates[moeda1] == 0 {
		return 0, fmt.Errorf("'moeda1' Inválida")
	}
	if Currencies.Rates[moeda2] == 0 {
		return 0, fmt.Errorf("'moeda2' Inválida")
	}

	return Currencies.Rates[moeda2] / Currencies.Rates[moeda1], nil
}
