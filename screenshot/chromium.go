package screenshot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

var ErrNotInteresting = errors.New("this change is not interesting")

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

// This snippet counts the number of "interesting" changes on a wiki diff page, e.g.
// it filters out very small changes and changes to metadata (e.g. link lists)
// TODO: Maybe filter better; only check/count diff-addedline and deletedline if there's no diffchange-inline within it?
const jsCountInteresting = `
var changes = [...document.querySelectorAll(".diff-addedline"), ...document.querySelectorAll(".diff-deletedline"), ...document.querySelectorAll(".diffchange-inline")]

changes.map(x => {
    var change = x.innerText.trim(); 
    return (change.startsWith("[[") || change.startsWith("*[") || change.startsWith("* [") || change.length <= 10) ? 0 : 1;
}).reduce((a, b) => a + b);
`

// This JS snippet creates a css style that censors the user name text.
const jsCensorUser = `const sheet = new CSSStyleSheet();
sheet.replaceSync(".censored{color: #000 !important;background: #000 !important;}");
document.adoptedStyleSheets = [sheet];
[...document.querySelectorAll(".mw-userlink"),...document.querySelector(".diff").querySelectorAll("a[href^=\\/wiki]")]
.filter(x => !x.parentElement.classList.contains("autocomment")).filter(x => ["Visuelle Bearbeitung", "Markierung", "Markierungen", "Diskussion", "Beiträge"].indexOf(x.innerText) === -1)
.forEach(x => { x.innerText="censored"; x.className = "censored"; });`

// see https://github.com/chromedp/examples/blob/master/screenshot/main.go
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	var interestingCount []byte

	return chromedp.Tasks{
		// If the viewport height is too small, the lower part of the page is cut off
		// So now we just take the maximum image height twitter allows
		chromedp.EmulateViewport(1800, 8192),
		chromedp.Navigate(urlstr),

		// Now count the number of interesting changes. Changes are interesting, if they are **not** changes
		// to metadata. Changes to metadata typically start with "[[", as that's part of the Wiki syntax
		// We return an error if no objects are interesting
		chromedp.Evaluate(jsCountInteresting, &interestingCount),
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			c, err := strconv.Atoi(string(interestingCount))
			if err != nil {
				log.Printf("[Screenshot] error: invalid number format while parsing interestingCount (%q): %s\n", interestingCount, err.Error())
			}

			if c == 0 {
				return ErrNotInteresting
			}

			return nil
		}),

		// Cannot pass nil, but we won't use the returned value
		chromedp.Evaluate(jsCensorUser, &[]byte{}),

		// This next is basically a copy of the source code of chromedp.Screenshot, except that the scale
		// is set to 1.25 so the text resolution is higher
		chromedp.QueryAfter(sel, func(ctx context.Context, execCtx runtime.ExecutionContextID, nodes ...*cdp.Node) error {
			if len(nodes) < 1 {
				return fmt.Errorf("selector %q did not return any nodes", sel)
			}

			// get box model / ViewPort
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
					Scale:  1.25,
				}).Do(ctx)
			if err != nil {
				return err
			}

			*res = buf
			return nil
		}),
	}
}
