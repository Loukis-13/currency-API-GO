package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//

type Response struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

var currencies Response

func init() {
	fmt.Println("Calling API...")

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

	json.Unmarshal(bodyBytes, &currencies)
	fmt.Printf("API Response as struct %+v\n", currencies)
}

//

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currencies)
}

func moedas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currencies.Rates)
}

func moeda(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(currencies.Rates[strings.ToUpper(params["moeda"])])
}

func troca(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()

	moeda1 := "BRL"
	if urlParams["moeda1"] != nil {
		if moeda1 = strings.ToUpper(urlParams["moeda1"][0]); currencies.Rates[moeda1] == 0 {
			moeda1 = "BRL"
		}
	}

	var moeda2 string
	if urlParams["moeda2"] != nil {
		if moeda2 = strings.ToUpper(urlParams["moeda2"][0]); currencies.Rates[moeda2] == 0 {
			json.NewEncoder(w).Encode("'moeda2' inválida")
			return
		}
	} else {
		json.NewEncoder(w).Encode("'moeda2' inválida")
		return
	}

	var qnt float64 = 1
	var errQnt error
	if urlParams["qnt"] != nil {
		if qnt, errQnt = strconv.ParseFloat(urlParams["qnt"][0], 64); errQnt != nil {
			qnt = 1
		}
	}

	conversao := currencies.Rates[moeda2] / currencies.Rates[moeda1] * qnt

	response := fmt.Sprintf("%f %s igual a %f %s", qnt, moeda1, conversao, moeda2)

	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/moedas", moedas).Methods("GET")
	router.HandleFunc("/moedas/{moeda}", moeda).Methods("GET")
	router.HandleFunc("/troca", troca).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
