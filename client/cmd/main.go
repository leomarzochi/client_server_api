package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/leomarzochi/client-server-api/client/models"
)

func main() {
	price, err := getQuotation()
	if err != nil {
		panic(err)
	}

	writeFile(price)

}

func writeFile(price *models.Price) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}

	f.Write([]byte("Dólar: " + price.Bid))
}

func getQuotation() (*models.Price, error) {
	var price *models.Price

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&price)
	if err != nil {
		return nil, err
	}

	return price, nil
}
