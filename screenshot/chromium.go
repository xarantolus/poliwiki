package screenshot

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func Take(webpage string) (pngData []byte, err error) {
	ctx, c := context.WithTimeout(context.Background(), 5*time.Minute)
	defer c()

	ctx, cc := chromedp.NewExecAllocator(ctx, chromedp.WindowSize(1920, 1080), chromedp.Headless)
	defer cc()

	// create context
	ctx, ccc := chromedp.NewContext(ctx)
	defer ccc()

	// capture screenshot of an element
	if err = chromedp.Run(ctx, elementScreenshot(webpage, "table.diff", &pngData)); err != nil {
		return
	}

	if len(pngData) == 0 {
		err = fmt.Errorf("couldn't take screenshot, no error but no data received")
	}

	return
}

// see https://github.com/chromedp/examples/blob/master/screenshot/main.go
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.QueryAfter(sel, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
			if len(nodes) < 1 {
				return fmt.Errorf("selector %q did not return any nodes", sel)
			}

			// get box model
			box, err := dom.GetBoxModel().WithNodeID(nodes[0].NodeID).Do(ctx)
			if err != nil {
				return err
			}
			if len(box.Margin) != 8 {
				return chromedp.ErrInvalidBoxModel
			}

			// take screenshot of the box
			buf, err := page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatPng).
				WithFromSurface(false).
				WithClip(&page.Viewport{
					// Round the dimensions, as otherwise we might
					// lose one pixel in either dimension.
					X:      math.Round(box.Margin[0]),
					Y:      math.Round(box.Margin[1]),
					Width:  math.Round(box.Margin[4] - box.Margin[0]),
					Height: math.Round(box.Margin[5] - box.Margin[1]),
					Scale:  4,
				}).Do(ctx)
			if err != nil {
				return err
			}

			*res = buf
			return nil
		}),
	}
}
