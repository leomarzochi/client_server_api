package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/leomarzochi/client-server-api/server/db"
	"github.com/leomarzochi/client-server-api/server/models"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", QuotationHandler)

	panic(http.ListenAndServe(":8080", mux))
}

func QuotationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	quotation, err := getQuotation()
	if err != nil {
		log.Default().Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = saveQuotationToDB(quotation)
	if err != nil {
		log.Default().Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(quotation.USDBRL)
	if err != nil {
		log.Default().Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func saveQuotationToDB(q *models.Quotation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		log.Printf("Context done: %v", ctx.Err())
		return ctx.Err()
	default:
		// do nothing
	}

	db, err := db.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO cotacao (name, code, codein, bid) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(q.USDBRL.Name, q.USDBRL.Code, q.USDBRL.CodeIN, q.USDBRL.Bid)
	if err != nil {
		return err
	}

	return nil
}

func getQuotation() (*models.Quotation, error) {
	var quotation models.Quotation

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil,
	)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error doing request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	err = json.Unmarshal(b, &quotation)

	return &quotation, err
}
