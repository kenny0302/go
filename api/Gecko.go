package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"log"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/tidwall/gjson"
	gecko "github.com/superoo7/go-gecko/v3"
	"bufio"
	//"btc/go"
)


const (
	URL = "127.0.0.1:27017"
)

type BTC struct {
	ID bson.ObjectId  `bson:"_id"`
	SourceName string `bson:"sourcename"`
	Price string `bson:"price"`
	Time time.Time 
}

func main() {
	var inputReader *bufio.Reader
	var input string
	var err error
	inputReader = bufio.NewReader(os.Stdin)
	fmt.Println("Please enter some input: ")
	input, err = inputReader.ReadString('\n')
	if err == nil {
		fmt.Printf("The input was: %s\n", inputReader)
	}
	switch input{
	case "1":
		Gecko()
	case "2":
		Coinmarketcap()
	case "3":
		Coinapi()
	}
	
}

	//連接資料庫
func Init()*mgo.Database {
	session, err :=mgo.Dial(URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("mydb")
	return db
}

//插入資料
func Insert() {
	db := Init()
	c := db.C("btc")
	i := bson.NewObjectId()
	err := c.Insert(&BTC{
		ID: i,
		SourceName: sourcename,
		Price: values,
		Time:time.Now(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(err)
}

// 資料來源1(ok)

func Gecko()(string,string) {
	cg := gecko.NewClient(nil)

	price, err := cg.SimpleSinglePrice("bitcoin", "usd")
	if err != nil {
		log.Fatal(err)
	}
	Usd := price.MarketPrice
	usd := fmt.Sprintf("%f",Usd)
	values := usd
	sourcename := "Gecko"
	return sourcename,values
	//fmt.Println(Usd)
}

//資料來源2(ok)

func Coinmarketcap()(string,string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "1")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "b2268019-06ab-4a96-807d-8917ea500db5")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
   
	value := gjson.Get(string(body), "data.0.quote.USD.price")
	//println(value.String())
	values :=value.String()
	//return(value.String())
	sourcename := "Coinmarketcap"
	return sourcename,values
}

//資料來源3(ok)
func Coinapi()(string,string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://rest.coinapi.io/v1/symbols?filter_symbol_id=CHAOEX_SPOT_BTC_USDT", nil)
	if err != nil {
	 log.Print(err)
	 os.Exit(1)
	}
	req.Header.Add("X-CoinAPI-Key", "6C13E532-15B4-4CD1-A932-3D2BB74BAE89")
   
	resp, err := client.Do(req)
	if err != nil {
	 fmt.Println("Error sending request to server")
	 os.Exit(1)
	}
   
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	value := gjson.Get(string(body), "0.price")
	values := value.String()
	//println("USD : " + value.String())
	//return (value.String())
	sourcename := "Coinapi"
	return sourcename,values
   }
   
