package main

import (
	model "github.com/Pedro-Henrique-Sampaio/Desafio-Multithreading/cmd/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type EnderecoApresentacao struct {
	cep        string "cep"
	logradouro string "logradouro"
	bairro     string "bairro"
	localidade string "localidade"
	uf         string "uf"
}

const cep string = "01153000"

func main() {

	fmt.Println("Iniciando o programa Busca de Cep atraves da ViaCEP ou BrasilAPI para ver qual o mais rapido !!!!")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	chviaCEP := make(chan io.ReadCloser)
	chbrasilAPI := make(chan io.ReadCloser)

	//Thread numero 2 - anonima
	go func() {
		// time.Sleep(2 * time.Second) // Simular o tempo de resposta
		chviaCEP <- chamarApi(ctx, montarChamadaViaCep(cep))
	}()

	//Thread numero 3 - anonima
	go func() {
		// time.Sleep(1 * time.Second) // Simular o tempo de resposta
		chbrasilAPI <- chamarApi(ctx, montarChamadaBrasilApi(cep))
	}()

	select {

	case msg1 := <-chviaCEP: //Recebe do viaCep
		fmt.Println("Received from ViaCEP")
		var viaCep model.ViaCEP
		err := json.NewDecoder(msg1).Decode(&viaCep)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v \n", err)
		}
		tratarRetorno(viaCep.Cep, viaCep.Logradouro, viaCep.Bairro, viaCep.Localidade, viaCep.Uf)
		defer msg1.Close()

	case msg2 := <-chbrasilAPI: //Recebe do brasilApi
		fmt.Println("Received from BrasilAPI")
		var brasilApi model.BrasilAPI
		err := json.NewDecoder(msg2).Decode(&brasilApi)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v \n", err)
		}
		tratarRetorno(brasilApi.Cep, brasilApi.Street, brasilApi.Neighborhood, brasilApi.City, brasilApi.State)
		defer msg2.Close()

	case <-ctx.Done(): //Timeout
		fmt.Println("Error: timeout")

	}

	fmt.Println("Programa encerrado")

}

func chamarApi(ctx context.Context, url string) io.ReadCloser {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer a requisição: %v \n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return res.Body
}

func montarChamadaViaCep(cep string) string {
	return "http://viacep.com.br/ws/" + cep + "/json/"
}

func montarChamadaBrasilApi(cep string) string {
	return "https://brasilapi.com.br/api/cep/v1/" + cep
}

func tratarRetorno(cepEnt string, logradouroEnt string, bairroEnt string, localidadeEnt string, ufEnt string) {

	endretorno := EnderecoApresentacao{
		cep:        cepEnt,
		logradouro: logradouroEnt,
		bairro:     bairroEnt,
		localidade: localidadeEnt,
		uf:         ufEnt,
	}

	fmt.Println("CEP: " + endretorno.cep + " Logradouro: " + endretorno.logradouro + " Bairro: " + endretorno.bairro + " Localidade: " + endretorno.localidade + " UF: " + endretorno.uf)
}
