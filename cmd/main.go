package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type viaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type viaCDN struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func main() {
	canal1 := make(chan viaCDN)
	canal2 := make(chan viaCEP)
	var cep string
	fmt.Println("Digite o cep ex: 49160-000")
	fmt.Scan(&cep)
	go getEnderecoViaCEP(cep, canal2)
	go getEnderecoViaCDN(cep, canal1)
	select {
	case msg := <-canal1:
		fmt.Printf("Via CDN, %v \n", msg)
	case msg := <-canal2:
		fmt.Printf("Via CEP, %v \n", msg)
	case <-time.After(time.Second):
		fmt.Println("Timeout :/")
	}
}

func getEnderecoViaCDN(cep string, canal chan viaCDN) {
	client := http.Client{}
	url := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	res, err := client.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	var viaCDN viaCDN
	payload, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err := json.Unmarshal(payload, &viaCDN); err != nil {
		fmt.Println(err.Error())
	}
	canal <- viaCDN
}

func getEnderecoViaCEP(cep string, canal chan<- viaCEP) {
	client := http.Client{}
	url := "http://viacep.com.br/ws/" + cep + "/json/"
	res, err := client.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	var viaCEP viaCEP
	payload, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	if err := json.Unmarshal(payload, &viaCEP); err != nil {
		fmt.Println(err.Error())
	}
	canal <- viaCEP

}
