package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const PLAY_BUTTON string = "#play_without_captchas_button > span"
const TIMER string = "#time_remaining > span > span:nth-child(1) > span"
const ROLL_BUTTON string = "#free_play_form_button"
const MODAL_200 string = "#push_notification_modal"
const MODAL_200_BUTTON string = "#push_notification_modal > div.push_notification_big > div:nth-child(2) > div > div.pushpad_deny_button"
const MODAL = "#myModal4"
const MODAL_BUTTON = "#myModal4"
const RECAPTCHA = "#free_play_recaptcha > form > div"
const BONUS_GROUP = "#reward_points_bonuses_main_div"

const (
	MODAL_22                string = "#myModal22"
	MODAL_22_CLOSE          string = "#myModal22 > a"
	MODAL_PUSH_NOTIFICATION string = "#push_notification_modal"
	MODAL_PUSH_
	NOTIFICATION_CLOSE string = "#push_notification_modal > div.push_notification_small > div:nth-child(2) > div > div.pushpad_deny_button"
)

const (
	RECAPTCHA_SELECTOR string = ""
)

type modal struct {
	node          *cdp.Node
	selector      string
	closeSelector string
}

func (m *modal) isVisible() bool {
	visible := false
	attrs := m.node.Attributes

	for _, attr := range attrs {
		//fmt.Println(attr)
		if strings.Contains(attr, "open") {
			visible = true
		}
	}

	return visible
}

type captcha struct {
	exists    bool
	completed bool
	sitekey   string
}

type timer struct {
	exists  bool
	minutes time.Duration
	seconds time.Duration
}

func main() {
	modals := []*modal{{selector: MODAL_22, closeSelector: MODAL_22_CLOSE}, {selector: MODAL_PUSH_NOTIFICATION, closeSelector: MODAL_PUSH_NOTIFICATION_CLOSE}}
	var ws string = "ws://127.0.0.1:9222/devtools/browser/6c2f3c1c-02d1-45da-85e0-e1604a7e4f32"
	//fmt.Printf("Entre com o websocket: ")
	//fmt.Scanln(&ws)

	remoteAlloc, cancel := chromedp.NewRemoteAllocator(context.Background(), ws)
	defer cancel()

	ctx, cancel := chromedp.NewContext(remoteAlloc)
	defer cancel()

	err := chromedp.Run(ctx, chromedp.Navigate("https://freebitco.in/?op=home"))
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println("Trying..")

		ok, sleepTime := checkTimer(ctx)

		if ok {
			updateModals(ctx, modals)
			ok, openModals := checkModals(ctx, modals)

			if !ok {
				fmt.Println("Found visible modals..")
				closeModals(ctx, openModals)
			} else {

			}
		} else {
			fmt.Printf("Sleeping for %v", sleepTime)
			time.Sleep(sleepTime)
		}

		time.Sleep(5 * time.Second)
	}

	/*
		fmt.Println("Checking TIMER")
		var timer []*cdp.Node
		err = chromedp.Run(ctx, chromedp.Nodes(TIMER, &timer, chromedp.AtLeast(0)))
		if err != nil {
			panic(err)
		}
		if len(timer) > 0 {
			fmt.Println("Timer exists!")
			fmt.Println(timer[0].Children[0].NodeValue)
			sleepTime, err := strconv.Atoi(timer[0].Children[0].NodeValue)
			if err != nil {
				return
			}
			if sleepTime > 0 {
				fmt.Printf("sleeping for: %v", sleepTime)
				time.Sleep(time.Duration(sleepTime) * time.Minute)
			}
		} else {
			fmt.Println("Checking for captcha")
			var captcha []*cdp.Node
			err = chromedp.Run(ctx, chromedp.Nodes(RECAPTCHA, &captcha, chromedp.AtLeast(0)))
			if err != nil {
				panic(err)
			}
			if len(captcha) > 0 {
				fmt.Println("Captcha exists!")
				fmt.Println(captcha[0])

				fmt.Println("Play button doesnt exist, checking roll button")
				var rollButton []*cdp.Node
				err = chromedp.Run(ctx, chromedp.Nodes(ROLL_BUTTON, &rollButton, chromedp.BySearch))
				if err != nil {
					panic(err)
				}
				if len(playButton) > 0 {
					fmt.Println("Roll button exists!")
					err = chromedp.Run(ctx,
						chromedp.ScrollIntoView(ROLL_BUTTON, chromedp.BySearch),
						chromedp.Click(ROLL_BUTTON, chromedp.BySearch),
					)
					if err != nil {
						panic(err)
					}
				} else {
					fmt.Println("Roll Button doest exist..")
				}
			}
		}
	*/
}

func checkTimer(ctx context.Context) (ok bool, sleepTime time.Duration) {
	ok = true
	fmt.Println("Checking TIMER")
	var timer []*cdp.Node
	err := chromedp.Run(ctx, chromedp.Nodes(TIMER, &timer, chromedp.AtLeast(0)))
	if err != nil {
		panic(err)
	}
	if len(timer) > 0 {
		ok = false
		fmt.Println("Timer exists!")
		st, err := strconv.Atoi(timer[0].Children[0].NodeValue)
		if err != nil {
			return
		}
		if st > 0 {
			sleepTime = time.Duration(st) * time.Minute
		} else {
			sleepTime = 1 * time.Minute
		}
	}

	return ok, sleepTime
}

func updateModals(ctx context.Context, modals []*modal) {
	for _, mod := range modals {
		var n []*cdp.Node
		err := chromedp.Run(ctx, chromedp.Nodes(mod.selector, &n, chromedp.AtLeast(0)))
		if err != nil {
			panic(err)
		}

		if len(n) > 0 {
			mod.node = n[0]
		}
	}
}

func checkModals(ctx context.Context, modals []*modal) (ok bool, openModals []*modal) {
	ok = true
	for _, mod := range modals {
		if mod.isVisible() {
			ok = false
			openModals = append(openModals, mod)
		}
	}

	time.Sleep(5 * time.Second)
	return ok, openModals
}

func closeModals(ctx context.Context, openModals []*modal) {
	for _, op := range openModals {
		err := chromedp.Run(ctx,
			chromedp.Click(op.closeSelector, chromedp.AtLeast(0)),
		)
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
}

/*
func seila(ctx context.Context, modals []*modal) {
	for _, mod := range modals {
		var n []*cdp.Node
		err := chromedp.Run(ctx, chromedp.Nodes(mod.selector, &n, chromedp.AtLeast(0)))
		if err != nil {
			panic(err)
		}

		if len(n) > 0 {
			fmt.Println("Modal exists!")
			mod.node = n[0]
			fmt.Println("printando aquii")
			fmt.Println(n[0])
			fmt.Println("------------------")
			if mod.isVisible() {
				canProceed = true
				fmt.Println("Modal is visible!")
				err = chromedp.Run(ctx,
					chromedp.Click(MODAL_22_CLOSE, chromedp.BySearch),
				)
				if err != nil {
					panic(err)
				}
			} else {
				fmt.Println("Modal is not visible!")
			}
		}
	}
	/*
		if len(modals) > 0 {
			fmt.Println("Modal exists!")
			if strings.Contains(modals[0].Attributes[7], "visible") {
				fmt.Println("MODAL 200 Visible!")
				err := chromedp.Run(ctx, chromedp.Click(MODAL_200_BUTTON))
				if err != nil {
					panic(err)
				}
			}
		}
}
*/
