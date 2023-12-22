package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"math/rand"
)

type TableView1 struct {
	OrderUid          string `json:"order_uid"`
	TrackNumber       string `json:"track_number"`
	Entry             string `json:"entry"`
	Items             string `json:"items"`
	Locale            string `json:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerId        string `json:"customer_id"`
	DeliveryService   string `json:"delivery_service"`
	Shardkey          int    `json:"shardkey"`
	SmId              int    `json:"sm_id"`
	DateCreate        string `json:"date_created"`
	OffShard          int    `json:"oof_shard"`

	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     int    `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`

	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func tableComp() []byte {
	var u TableView1
	u.OrderUid = RandStringRunes(12)
	u.TrackNumber = RandStringRunes(8)
	u.Entry = RandStringRunes(4)
	u.Items = `[{"chrt_id": 9934930,"track_number": "WBILMTESTTRACK","price": 453,"rid": "ab4219087a764ae0btest","name": "Mascaras","sale": 30,"size": "0","total_price": 317,"nm_id": 2389212,"brand": "Vivienne Sabo","status": 202}]`
	u.Locale = "ru"
	u.InternalSignature = "2"
	u.CustomerId = "test2"
	u.DeliveryService = RandStringRunes(3)
	u.Shardkey = rand.Intn(30-1) + 1
	u.SmId = rand.Intn(220-1) + 1
	u.DateCreate = "2021-11-26T06:22:19Z"
	u.OffShard = rand.Intn(12-1) + 1

	u.Name = RandStringRunes(5)
	u.Phone = "+79928763542"
	u.Zip = rand.Intn(400000-300000) + 300000
	u.City = RandStringRunes(10)
	u.Address = RandStringRunes(13)
	u.Region = "region2"
	u.Email = RandStringRunes(5)

	u.Transaction = u.OrderUid
	u.RequestId = "12345"
	u.Currency = "RUB"
	u.Provider = "WBPAY"
	u.Amount = rand.Intn(2000-1000) + 1000
	u.PaymentDt = rand.Intn(900000-100000) + 100000
	u.Bank = "WB"
	u.DeliveryCost = rand.Intn(300-200) + 200
	u.GoodsTotal = rand.Intn(10-1) + 1
	u.CustomFee = rand.Intn(2000-1000) + 1000
	bytes, _ := json.Marshal(u)
	return bytes
}

func main() {
	bytes := tableComp()
	sc, err := stan.Connect("mystreamingserver", "Natan")
	if err != nil {
		fmt.Print(err)
		return
	}

	err2 := sc.Publish("order", bytes)
	if err2 != nil {
		fmt.Print(err2)
		return
	}
}
