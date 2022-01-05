package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Loukis-13/currency-API-GO/internal/currencies"

	"github.com/gorilla/mux"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currencies.Currencies)
}

func Moedas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currencies.Currencies.Rates)
}

func Moeda(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	json.NewEncoder(w).Encode(currencies.Currencies.Rates[strings.ToUpper(params["moeda"])])
}

func Troca(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()

	moeda1 := "BRL"
	if urlParams["moeda1"] != nil {
		moeda1 = strings.ToUpper(urlParams["moeda1"][0])
	}

	var moeda2 string
	if urlParams["moeda2"] != nil {
		moeda2 = strings.ToUpper(urlParams["moeda2"][0])
	}

	var qnt float64 = 1
	var err error
	if urlParams["qnt"] != nil {
		if qnt, err = strconv.ParseFloat(urlParams["qnt"][0], 64); err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}
	}

	conversao, err := currencies.Trocar(moeda1, moeda2)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	response := fmt.Sprintf("%f %s igual a %f %s", qnt, moeda1, conversao*qnt, moeda2)

	json.NewEncoder(w).Encode(response)
}
