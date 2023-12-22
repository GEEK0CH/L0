package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	stand "github.com/nats-io/nats-streaming-server/server"
	"github.com/nats-io/stan.go"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

var connection *sql.DB

const connectionString = "host=localhost port=5432 dbname=orders user=Orders_role password=1234 sslmode=disable"

type TableView struct {
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

type tableId struct {
	idx int `json:"id"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	var u TableView
	cache := make(map[int]string)
	fContent, _ := ioutil.ReadFile("cache.txt")
	json.Unmarshal(fContent, &cache)
	id := r.FormValue("id")
	id2, _ := strconv.Atoi(id)
	json.Unmarshal([]byte(cache[id2]), &u)
	values := reflect.ValueOf(u)
	types := values.Type()
	tmpl, _ := template.New("").Parse("<h1>Order table</h1><a>id:{{.}}</a><br>")
	_ = tmpl.Execute(w, id)
	for i := 0; i < values.NumField(); i++ {
		name := types.Field(i).Name
		if name == "Name" {
			tmpl, _ := template.New("id").Parse("<h1>Order delivery</h1><br>")
			_ = tmpl.Execute(w, id)
		} else if name == "Transaction" {
			tmpl, _ := template.New("id").Parse("<h1>Order payment</h1><br>")
			_ = tmpl.Execute(w, id)
		}
		str := "<a>" + name + ":{{.}}</a><br>"
		tmpl1, _ := template.New(name).Parse(str)
		_ = tmpl1.Execute(w, values.Field(i))
	}
	tmpl30, _ := template.New("button").Parse("<a href=\"/\"><input type=\"submit\" class=\"submit\" value=\"Назад\"></a>")
	_ = tmpl30.Execute(w, id)
}

func hh(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "form.html")
}

func main() {
	var e error
	connection, e = sql.Open("postgres", connectionString)
	if e != nil {
		fmt.Println(e)
		return
	}

	stand.RunServer("mystreamingserver")
	sc, err := stan.Connect("mystreamingserver", "mystreamingserver")
	if err != nil {
		fmt.Print(err)
		return
	}

	var u TableView
	var idx tableId
	cache := make(map[int]string)
	rowtable, _ := connection.Query(`SELECT * FROM "order"`)
	rowtable2, _ := connection.Query(`SELECT * FROM "delivery"`)
	rowtable3, _ := connection.Query(`SELECT * FROM "payment"`)
	defer rowtable.Close()
	defer rowtable2.Close()
	defer rowtable3.Close()
	for rowtable.Next() {
		rowtable2.Next()
		rowtable3.Next()
		rowtable.Scan(&idx.idx, &u.OrderUid, &u.TrackNumber, &u.Entry, &u.Items, &u.Locale, &u.InternalSignature, &u.CustomerId, &u.DeliveryService, &u.Shardkey, &u.SmId, &u.DateCreate, &u.OffShard)
		rowtable3.Scan(&idx.idx, &u.Transaction, &u.RequestId, &u.Currency, &u.Provider, &u.Amount, &u.PaymentDt, &u.Bank, &u.DeliveryCost, &u.GoodsTotal, &u.CustomFee)
		rowtable2.Scan(&idx.idx, &u.Name, &u.Phone, &u.Zip, &u.City, &u.Address, &u.Region, &u.Email)
		bytes, _ := json.Marshal(u)
		cache[idx.idx] = string(bytes)
	}
	bytes, _ := json.Marshal(cache)
	str := string(bytes)
	file, _ := os.Create("cache.txt")
	defer file.Close()
	file.WriteString(str)

	time.Sleep(10 * time.Second)
	sc.Subscribe("order",
		func(m *stan.Msg) {
			var idx tableId
			json.Unmarshal(m.Data, &u)
			row := connection.QueryRow(`SELECT MAX("id") FROM "order"`)
			row.Scan(&idx.idx)
			connection.QueryRow(`INSERT INTO "order"
    ("order_uid", "track_number", "entry", "items", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard")
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
				u.OrderUid, u.TrackNumber, u.Entry, u.Items, u.Locale, u.InternalSignature, u.CustomerId, u.DeliveryService, u.Shardkey, u.SmId, u.DateCreate, u.OffShard)
			connection.QueryRow(`INSERT INTO "delivery"
    ("name", "phone", "zip", "city", "address", "region", "email")
    VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				u.Name, u.Phone, u.Zip, u.City, u.Address, u.Region, u.Email)
			connection.QueryRow(`INSERT INTO public."payment"
    ("transaction", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee")
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
				u.Transaction, u.RequestId, u.Currency, u.Provider, u.Amount, u.PaymentDt, u.Bank, u.DeliveryCost, u.GoodsTotal, u.CustomFee)

			cache := make(map[int]string)
			bytes, _ := json.Marshal(u)
			idx.idx++
			cache[idx.idx] = string(bytes)
			bytesT, _ := json.Marshal(cache)
			strT := string(bytesT)
			fContent, _ := ioutil.ReadFile("cache.txt")
			file, _ := os.Create("cache.txt")
			defer file.Close()
			str := string(fContent)[:len(string(fContent))-1] + ","
			file.WriteString(str + strT[1:])

		},
		stan.StartWithLastReceived())

	http.HandleFunc("/", hh)
	http.HandleFunc("/db.html", hello)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
