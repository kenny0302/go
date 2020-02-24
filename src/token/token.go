package main
 
import (
    "net/http"
	"golang.org/x/time/rate"
	"os"
	"fmt"
	"io/ioutil"
	"net/url"
	"log"
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/tidwall/gjson"
	gecko "github.com/superoo7/go-gecko/v3"
)

const (
	URL = "localhost:27017"
)

type BTC struct {
	ID bson.ObjectId  `bson:"_id"`
	SourceName string `bson:"sourcename"`
	Price string `bson:"price"`
	Time  time.Time
}

func Token()([]byte) {
	var d string
	var i []byte
	Insert(Gecko())
	Insert(Coinmarketcap())
	Insert(Coinapi())
	a := string([]byte(QueryGecko()))
	b := string([]byte(QueryCoinmarketcap()))
	c := string([]byte(QueryCoinapi()))
	d = a+b+c
	i = []byte(d)
	return i
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Flow)
	http.ListenAndServe(":8080", limit(mux))
}

//使用Token Bucket實作流量管制
var limiter = rate.NewLimiter(2, 5)
func limit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if limiter.Allow() == false {
            http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}

//api
func Flow(w http.ResponseWriter, r *http.Request)  {
	Insert(Gecko())
	Insert(Coinmarketcap())
	Insert(Coinapi())
	w.Write([]byte(QueryGecko()))
	w.Write([]byte(QueryCoinmarketcap()))
	w.Write([]byte(QueryCoinapi()))
}

//連接資料庫
func Init()*mgo.Database {
	session, err :=mgo.Dial(URL)
	if err != nil {
		panic(err)
	}
	//defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("mydb")
	return db
}

//插入資料
func Insert(a string,b string) {
	db := Init()
	c := db.C("btc")
	i := bson.NewObjectId()
	//j,k := v()
	err := c.Insert(&BTC{
		ID: i,
		SourceName: a,
		Price: b,
		Time:time.Now(),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(err)
}

//查詢資料Gecko
func QueryGecko()([]byte) {
	db := Init()
	c:= db.C("btc")
	btc:=make([]BTC,0,100)
	err := c.Find(bson.M{"sourcename":"Gecko"}).Sort("timestamp").Limit(1).All(&btc)
	if err != nil {
		panic(err)
	}
	btcs := []byte(fmt.Sprintln(btc))
	return btcs
}

//查詢資料Coinmarketcap
func QueryCoinmarketcap()([]byte) {
	db := Init()
	c:= db.C("btc")
	btc:=make([]BTC,0,100)
	err := c.Find(bson.M{"sourcename":"Coinmarketcap"}).Sort("timestamp").Limit(1).All(&btc)
	if err != nil {
		panic(err)
	}
	btcs := []byte(fmt.Sprintln(btc))
	return btcs
}

//查詢資料Coinapi
func QueryCoinapi()([]byte) {
	db := Init()
	c:= db.C("btc")
	btc:=make([]BTC,0,100)
	err := c.Find(bson.M{"sourcename":"Coinapi"}).Sort("timestamp").Limit(1).All(&btc)
	if err != nil {
		panic(err)
	}
	btcs := []byte(fmt.Sprintln(btc))
	return btcs
}

//查詢3筆最新資料
func QueryAll() {
	db := Init()
	c:= db.C("btc")
	btc:=make([]BTC,0,100)
	err := c.Find(bson.M{}).Sort("timestamp").Limit(3).All(&btc)
	if err != nil {
		panic(err)
	}
	fmt.Println(btc)
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
   
