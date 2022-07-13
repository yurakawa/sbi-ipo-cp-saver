package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"strings"
	"time"
)

// memo: actionを別のrunに分けるとcontext cancelが発生しやすく鳴る
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		// chromedp.Flag("enable-automation", false),
		// chromedp.Flag("restore-on-startup", false),
		// chromedp.Flag("new-window", true),
	)...)
	defer cancel()
	taskCtx, cancel := chromedp.NewContext(allocCtx,
		chromedp.WithLogf(log.Printf),
		// chromedp.WithDebugf(log.Printf), // debug onにすると初期段階でcontext canceledが発生する。
		// chromedp.WithErrorf(log.Printf),
	)
	defer cancel()

	// list awesome go projects for the "Selenium and browser control tools."
	//res, err := listAwesomeGoProjects(ctx, "Selenium and browser control tools.")
	err := listSbiIpoCompany(taskCtx)
	if err != nil {
		log.Fatalf("could not list awesome go projects: %v", err)
	}
}
func listSbiIpoCompany(ctx context.Context) error {
	var u string
	fmt.Println("login start")
	err := chromedp.Run(ctx, login(&u))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("login end")
	fmt.Printf("URL => %s\n", u)

	// fmt.Println("move to ipo page start")
	// err = chromedp.Run(ctx, moveToIpoPage(&u))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("URL => %s\n", u)
	// fmt.Println("move to ipo page end")

	// fmt.Println("move to application page start")
	// err = chromedp.Run(ctx, moveToApplicationPage(&u))
	// if err != nil {
	// 	fmt.Println("エラー!!")
	// 	log.Fatal(err)
	// }
	// fmt.Printf("URL => %s\n", u)
	// fmt.Println("move to application page end")

	// fmt.Println("apply application start")
	// err = chromedp.Run(ctx, applyApplication(&u))
	// if err != nil {
	// 	fmt.Println("エラーだっ!!!")
	// 	log.Fatal(err)
	// }
	// fmt.Printf("URL => %s\n", u)

	// fmt.Println("apply application end")

	// fmt.Println("confirm start")
	// err = chromedp.Run(ctx, confirmPage(&u))
	// fmt.Printf("URL => %s\n", u)
	// fmt.Println("confirm end")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return nil
}

var (
	username         = os.Getenv("SBI_USERNAME")
	password         = os.Getenv("SBI_PASSWORD")
	torihikiPassword = os.Getenv("SBI_TORIHIKI_PASSWORD")
)

func login(url *string) chromedp.Tasks {
	urlStr := `https://site2.sbisec.co.jp/ETGate/`
	usernameSel := `//input[@name='user_id']`
	passwordSel := `//input[@name='user_password']`

	urlB := "https://m.sbisec.co.jp/oeliw011?type=21"

	// urlC := "https://m.sbisec.co.jp/switchnaviMain"
	applicationSel := `//img[@alt="申込"]`

	moushikomikabusuu := `//input[@name="suryo"]`
	strikePriceSel := `//label[@for="strPriceRadio"]`
	// submitSel := `//input[@name="order_kakunin"]`
	submitSel := `//input[@value="申込確認画面へ"]`

	// submitSelB := `//input[@name='order_btn']`

	// submitSelC := `//input[@type='submit']`
	submitSelC := `//input[@name='order_btn']`
	// submitSelC := `//input[@value="  申込  "]`
	// submitSelC := `//input[@value="&nbsp;&nbsp;申込&nbsp;&nbsp;"]`
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(usernameSel),
		chromedp.WaitVisible(passwordSel),
		chromedp.SendKeys(usernameSel, username),
		chromedp.SendKeys(passwordSel, password),
		chromedp.Submit(passwordSel),
		chromedp.WaitNotPresent(`//div[@id="new_login"]`),

		chromedp.Navigate(urlB),

		// chromedp.Navigate(urlC),
		chromedp.WaitVisible(applicationSel),
		chromedp.Click(applicationSel, chromedp.NodeVisible),

		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("a")
			return nil
		}),
		chromedp.WaitVisible(moushikomikabusuu),
		chromedp.SendKeys(moushikomikabusuu, "1000"), // ok
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("b")
			return nil
		}),
		chromedp.SendKeys(`//input[@name="tr_pass"]`, torihikiPassword), // ok
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("c")
			return nil
		}),
		chromedp.Click(strikePriceSel),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("d")
			return nil
		}),
		chromedp.Click(submitSel),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("e")
			return nil
		}),

		// chromedp.Navigate("https://m.sbisec.co.jp/oeapw021"),
		// chromedp.WaitNotPresent(`//input[@name="tr_pass"]`),
		chromedp.WaitVisible(submitSelC),
		chromedp.Click(submitSelC),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("f")
			return nil
		}),
		chromedp.Location(url),
	}
}

func moveToIpoPage(url *string) chromedp.Tasks {
	urlStr := "https://m.sbisec.co.jp/oeliw011?type=21"
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.Location(url),
	}
}

func moveToApplicationPage(url *string) chromedp.Tasks {
	urlStr := "https://m.sbisec.co.jp/switchnaviMain"
	applicationSel := `//img[@alt="申込"]`

	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(applicationSel),
		chromedp.Click(applicationSel, chromedp.NodeVisible),
		chromedp.Location(url),
	}
}
func applyApplication(url *string) chromedp.Tasks {
	moushikomikabusuu := `//input[@name="suryo"]`
	strikePriceSel := `//label[@for="strPriceRadio"]`
	submitSel := `//input[@type="submit"][@name="order_kakunin"]`
	return chromedp.Tasks{
		chromedp.SendKeys(moushikomikabusuu, "1000", chromedp.NodeVisible),                    // ok
		chromedp.SendKeys(`//input[@name="tr_pass"]`, torihikiPassword, chromedp.NodeVisible), // ok
		chromedp.Click(strikePriceSel, chromedp.NodeVisible),
		chromedp.Sleep(time.Second),
		chromedp.Click(submitSel, chromedp.NodeVisible),
		chromedp.Location(url),
	}
}

var SetCookiesAction = func(cookies []*network.Cookie) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		cc := make([]*network.CookieParam, 0, len(cookies))
		for _, c := range cookies {
			cc = append(cc, &network.CookieParam{
				Name:   c.Name,
				Value:  c.Value,
				Domain: c.Domain,
			})
		}
		return network.SetCookies(cc).Do(ctx)
	})
}

var GetCookiesAction = chromedp.ActionFunc(func(ctx context.Context) error {
	cookies, err := network.GetAllCookies().Do(ctx)
	if err != nil {
		return err
	}
	log.Printf("cookies got: %v", cookies)
	// e.g. save to file
	return nil
})

func old_applyApplication(url *string) chromedp.Tasks {
	moushikomikabusuu := `//input[@name='suryo']`
	//strikePriceSel := `//label[@for='strPriceRadio']`
	// ipoChallengePoint := `//input[@name='useKdn']`
	// passwordSel := `//input[@name='tr_pass']`
	// submitSel := `//input[@type='submit']`
	// submitSel := `//input[@name='order_kakunin']`
	return chromedp.Tasks{
		chromedp.WaitVisible(moushikomikabusuu),
		// chromedp.SendKeys(moushikomikabusuu, "100", chromedp.BySearch),
		//chromedp.SendKeys(`//input[@type='password']`, os.Getenv("SBI_TORIHIKI_PASSWORD"), chromedp.ByQuery),
		//chromedp.Click(strikePriceSel, chromedp.NodeVisible),
		// chromedp.Click(ipoChallengePoint, chromedp.NodeVisible),
		// chromedp.WaitVisible(passwordSel),
		// chromedp.Submit(submitSel, chromedp.NodeVisible),
		//chromedp.Click(submitSel, chromedp.NodeVisible),
		// chromedp.Location(url),
	}
}
func confirmPage(url *string) chromedp.Tasks {
	urlStr := "https://m.sbisec.co.jp/oeapw021"
	submitSel := `//input[@name='order_btn']`
	// submitSel := `//input[@value='  申込  ']`
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(submitSel),
		chromedp.Submit(submitSel),
		chromedp.Location(url),
	}
}

type projectDesc struct {
	URL, Description string
}

func listAwesomeGoProjects(ctx context.Context, sect string) (map[string]projectDesc, error) {

	sel := fmt.Sprintf(`//p[text()[contains(., '%s')]]`, sect)

	// navigate
	if err := chromedp.Run(ctx, chromedp.Navigate(`https://github.com/avelino/awesome-go`)); err != nil {
		return nil, fmt.Errorf("could not navigate to github: %v", err)
	}

	// wait visible
	if err := chromedp.Run(ctx, chromedp.WaitVisible(sel)); err != nil {
		return nil, fmt.Errorf("could not get section: %v", err)
	}

	sib := sel + `/following-sibling::ul/li`

	// get project link text
	var projects []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+`/child::a/text()`, &projects)); err != nil {
		return nil, fmt.Errorf("could not get projects: %v", err)
	}

	// get links and description text
	var linksAndDescriptions []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+`/child::node()`, &linksAndDescriptions)); err != nil {
		return nil, fmt.Errorf("could not get links and descriptions: %v", err)
	}

	// check length
	if 2*len(projects) != len(linksAndDescriptions) {
		return nil, fmt.Errorf("projects and links and descriptions lengths do not match (2*%d != %d)", len(projects), len(linksAndDescriptions))
	}

	// process data
	res := make(map[string]projectDesc)
	for i := 0; i < len(projects); i++ {
		res[projects[i].NodeValue] = projectDesc{
			URL:         linksAndDescriptions[2*i].AttributeValue("href"),
			Description: strings.TrimPrefix(strings.TrimSpace(linksAndDescriptions[2*i+1].NodeValue), "- "),
		}
	}

	return res, nil

}
