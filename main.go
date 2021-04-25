package main

import (
	"context"
	"fmt"
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

type modal struct {
	visible       bool
	selector      string
	closeSelector string
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
	modals := []*modal{{selector: MODAL, closeSelector: MODAL_BUTTON}, {selector: MODAL_200, closeSelector: MODAL_200_BUTTON}}
	var ws string
	fmt.Printf("Entre com o websocket: ")
	fmt.Scanln(&ws)

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

		updateModals(ctx, modals)

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

func updateModals(ctx context.Context, modals []*modal) {
	for _, mod := range modals {
		var n []*cdp.Node
		err := chromedp.Run(ctx, chromedp.Nodes(mod.selector, &n, chromedp.AtLeast(0)))
		if err != nil {
			panic(err)
		}

		if len(n) > 0 {
			fmt.Println("Modal exists!")
			attrs := n[0].Attributes

			for _, attr := range attrs {
				fmt.Println(attr)
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
	*/
}
