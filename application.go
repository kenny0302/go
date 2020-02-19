package main

import (
	"net/http"
	"fmt"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"bufio"
	"GO/BTC/api/Gecok"
)

type BTC struct {
	ID bson.ObjectId  `bson:"_id"`
	SourceName string `bson:"sourcename"`
	Price string `bson:"price"`
	Time time.Time 
}
const = (
	URL ="127.0.0.1:27017"
)

func main() {
	session, err :=mgo.Dial(URL)
	if err != nil {
		panic(err)
	}
defer session.Close()
db := session.DB("mydb")
collection := db.C("btc")

var iuputReeader *bufio.Reeader
var choice string
var err error
inputReader = bufio.NewRader(os.Stdin)

fmt.Println("請輸入所要使用的api")
fmt.Println("1.Gecko")
fmt.Println("2.CoinMarketCap")
fmt.Println("3.coinbase")
choice, err = inputReader.ReadString('/n')
if err == nil {
	fmt.Printf("The input was: %s\n", input)
}


switch choice{
case 1:
	
	fmt.Println("1 BTC :"+price+"USD")
case 2:
	fmt.Println("1 BTC :"+price+"USD")
case 3:
	fmt.Println("1 BTC :"+price+"USD")
}
}
