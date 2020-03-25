package main

import (
	"context"
	//"io/ioutil"
	"log"
	//"time"
	"github.com/chromedp/chromedp"
	"strconv"
)

func main() {

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	//ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	//defer cancel()

	// navigate to a page, wait for an element, click
	var example string
	
	for i :=108000001 ; i<108000112 ; i++{
		j := strconv.Itoa(i)
		err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.taiwanlottery.com.tw/lotto/Lotto649/history.aspx`),
		chromedp.WaitVisible(`#Lotto649Control_history_dlQuery_L649_CategA5_0`),
		chromedp.SendKeys(`#Lotto649Control_history_txtNO`,j,chromedp.ByID),
		chromedp.Click(`#Lotto649Control_history_btnSubmit`, chromedp.NodeVisible),
		chromedp.Text(`#Lotto649Control_history_dlQuery > tbody > tr > td > .table_org > tbody > tr:nth-child(5)`, &example),
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("example: %s", example)
	}
}
