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

func main() {
	e, err := config.LoadEnvVariables()
	if err != nil {
		log.Fatal(err)
	}

	userName = e.SbiUserName
	password = e.SbiPassword
	torihikiPassword = e.SbiTorihikiPassword

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:],
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
	ctx, cancel := context.WithTimeout(taskCtx, 30*time.Second)
	defer cancel()

	err = orderBookBuilding(ctx)
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
	for i := 0; i < len(company); i++ {
		var companyName string
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
	urlStr := `https://site1.sbisec.co.jp/ETGate/?_ControlID=WPLETlgR001Control&_PageID=WPLETlgR001Rlgn50&_DataStoreID=DSWPLETlgR001Control&_ActionID=login&getFlg=on`
	usernameSel := `//input[@name="user_id"]`
	passwordSel := `//input[@name="user_password"]`
	loginSel := `//input[@name="logon"]`

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
	torihikiPasswordSel := `//input[@name="tr_pass"]`
	submitOrderSel := `//input[@name="order_kakunin"]`
	// submitOrderConfirmSel := `//input[@type="submit"]`
	submitOrderConfirmSel := `//input[@name="order_btn"]`
	backSel := `//a[@href="/oeliw011?type=21"]`
	urlStr := "https://m.sbisec.co.jp/oeliw011?type=21"

	return chromedp.Tasks{
		chromedp.WaitVisible(applySel),
		chromedp.Click(applySel),
		chromedp.WaitVisible(submitOrderSel),
		chromedp.Text(`.lbody`, companyName),
		chromedp.SendKeys(suryoSel, "10000"),
		chromedp.Click(`//label[@for="strPriceRadio"]`),
		chromedp.SendKeys(torihikiPasswordSel, torihikiPassword),
		chromedp.WaitVisible(submitOrderSel),
		chromedp.Click(submitOrderSel),
		chromedp.Sleep(1 * time.Second),
		chromedp.Click(submitOrderConfirmSel),
		chromedp.WaitVisible(backSel),
		chromedp.Navigate(urlStr),

		chromedp.Sleep(time.Second),
	}
}

// $x('//input[@name="order_btn"]')のようにcdpのconsoleで確認する
