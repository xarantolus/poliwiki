module github.com/xarantolus/poliwiki

go 1.16

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/chromedp/cdproto v0.0.0-20210313213058-f5c5a7a06834
	github.com/chromedp/chromedp v0.6.8
	github.com/dghubble/go-twitter v0.0.0-20201011215211-4b180d0cc78d
	github.com/dghubble/oauth1 v0.7.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/sys v0.0.0-20210317225723-c4fcb01b228e // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/dghubble/go-twitter => ./bot/go-twitter
