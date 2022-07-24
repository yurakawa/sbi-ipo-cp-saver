package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/yurakawa/sbi-ipo-cp-miner/config"
)

var (
	userName         string
	password         string
	torihikiPassword string
)

// memo: actionを別のrunに分けるとcontext cancelが発生しやすく鳴る
func main() {
	e, err := config.LoadEnvVariables()
	if err != nil {
		log.Fatal(err)
	}

	userName = e.SbiUserName
	password = e.SbiPassword
	torihikiPassword = e.SbiTorihikiPassword

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		// chromedp.Flag("headless", false),
	)...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx,
		func() chromedp.ContextOption {
			co := chromedp.WithLogf(log.Printf)
			if e.LogLevel == "DEBUG" {
				co = chromedp.WithDebugf(log.Printf)
			}
			return co
		}(),
	)
	defer cancel()

	err = orderBookBuilding(taskCtx)
	if err != nil {
		log.Fatal(err)
	}

}

func orderBookBuilding(ctx context.Context) error {
	log.Println("login start")
	err := chromedp.Run(ctx, login())
	if err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}
	log.Println("login end")

	if err = chromedp.Run(ctx, chromedp.Navigate("https://m.sbisec.co.jp/oeliw011?type=21")); err != nil {
		return err
	}
	if err = chromedp.Run(ctx, chromedp.WaitVisible(`//h2[contains(text(),"新規上場株式ブックビルディング")]`)); err != nil {
		return err
	}

	var company []*cdp.Node
	if err = chromedp.Run(ctx, chromedp.Nodes(`//img[@alt="申込"]`, &company, chromedp.AtLeast(0))); err != nil {
		return err
	}
	if len(company) == 0 {
		log.Println("unapplied for does not exist.")
		return nil
	}
	var companyName string
	for _, _ = range company {
		log.Println("apply start")
		err = chromedp.Run(ctx, apply(&companyName))
		log.Println("apply", strings.TrimSpace(companyName))
		if err != nil {
			return fmt.Errorf("failed to apply: %v", err)
		}
		log.Println("apply end")
	}
	return nil
}

func login() chromedp.Tasks {
	urlStr := `https://site2.sbisec.co.jp/ETGate/`
	usernameSel := `//input[@name="user_id"]`
	passwordSel := `//input[@name="user_password"]`
	loginSel := `//input[@name="ACT_login"]`

	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(usernameSel),
		chromedp.WaitVisible(passwordSel),
		chromedp.SendKeys(usernameSel, userName),
		chromedp.SendKeys(passwordSel, password),
		chromedp.Click(loginSel),
		chromedp.WaitNotPresent(`//div[@id="new_login"]`),
		chromedp.Sleep(time.Second),
	}
}

func apply(companyName *string) chromedp.Tasks {
	applySel := `//img[@alt="申込"]`
	suryoSel := `//input[@name="suryo"]`
	strikePriceSel := `//label[@for="strPriceRadio"]`
	torihikiPasswordSel := `//input[@name="tr_pass"]`
	submitOrderSel := `//input[@name="order_kakunin"]`
	submitOrderConfirmSel := `//input[@name="order_btn"]`
	return chromedp.Tasks{
		chromedp.WaitVisible(applySel),
		chromedp.Click(applySel),
		chromedp.WaitVisible(suryoSel),
		chromedp.Text(`.lbody`, companyName),
		chromedp.SendKeys(suryoSel, "10000"),
		chromedp.SendKeys(torihikiPasswordSel, torihikiPassword),
		chromedp.Click(strikePriceSel),
		chromedp.Click(submitOrderSel),
		chromedp.WaitVisible(submitOrderConfirmSel),
		chromedp.Click(submitOrderConfirmSel),
		chromedp.Sleep(time.Second),
	}

}
