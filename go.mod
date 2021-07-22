module github.com/xarantolus/poliwiki

go 1.16

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/chromedp/cdproto v0.0.0-20210721224921-12abf3292481
	github.com/chromedp/chromedp v0.7.4
	github.com/dghubble/go-twitter v0.0.0-20210609183100-2fdbf421508e
	github.com/dghubble/oauth1 v0.7.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/dghubble/go-twitter => ./bot/go-twitter
