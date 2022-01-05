package routes

import (
	"time"

	"github.com/Loukis-13/currency-API-GO/internal/currencies"
)

type pedidoTroca struct {
	Moeda1      string  `json:"moeda1"`
	Moeda2      string  `json:"moeda2"`
	Quantidade1 float64 `json:"quantidade1"`
}

type respostaTroca struct {
	Moeda1      string  `json:"moeda1"`
	Moeda2      string  `json:"moeda2"`
	Quantidade1 float64 `json:"quantidade1"`
	Quantidade2 float64 `json:"quantidade2"`
	ValorMoeda2 float64 `json:"valorMoeda2"`
	FeitaEm     string  `json:"feitaEm"`
}

func NewRespostaTroca(body pedidoTroca) (respostaTroca, error) {
	valor, err := currencies.Trocar(body.Moeda1, body.Moeda2)
	if err != nil {
		return respostaTroca{}, err
	}

	r := respostaTroca{
		Moeda1:      body.Moeda1,
		Moeda2:      body.Moeda2,
		Quantidade1: body.Quantidade1,
		Quantidade2: valor * body.Quantidade1,
		ValorMoeda2: valor,
		FeitaEm:     time.Now().Format("2006/01/02 15:04:05"),
	}

	return r, nil
}
