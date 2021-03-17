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

	ctx, cc := chromedp.NewExecAllocator(ctx, chromedp.Headless)
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

const jsCensorUser = `const sheet = new CSSStyleSheet();
		sheet.replaceSync(".mw-userlink{color: #000;background: #000;}");
document.adoptedStyleSheets = [sheet];`

// see https://github.com/chromedp/examples/blob/master/screenshot/main.go
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		// If the viewport height is too small, the lower part of the page is cut off
		chromedp.EmulateViewport(1800, 1080*4),
		chromedp.Navigate(urlstr),

		// Cannot pass nil, but we won't use the returned value
		chromedp.Evaluate(jsCensorUser, &[]byte{}),
		// This next is basically a copy of the source code of chromedp.Screenshot, except that the scale
		// is set to 2 so the text resolution is higher
		chromedp.QueryAfter(sel, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
			if len(nodes) < 1 {
				return fmt.Errorf("selector %q did not return any nodes", sel)
			}

			// get box model
			box, err := dom.GetBoxModel().WithBackendNodeID(nodes[0].BackendNodeID).Do(ctx)
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
					Width:  float64(box.Width),
					Height: float64(box.Height),
					Scale:  2,
				}).Do(ctx)
			if err != nil {
				return err
			}

			*res = buf
			return nil
		}),
	}
}
