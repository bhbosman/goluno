module github.com/bhbosman/goLuno

go 1.18

require (
	github.com/bhbosman/goCommonMarketData v0.0.0-20220802122727-698b9feba01e
	github.com/bhbosman/goCommsDefinitions v0.0.0-20220801175552-c5aa68065af3
	github.com/bhbosman/goCommsMultiDialer v0.0.0-20220728174928-21aff6c1232e
	github.com/bhbosman/goCommsNetDialer v0.0.0-20220726130315-bec9f09e45e7
	github.com/bhbosman/goCommsNetListener v0.0.0-20220802132650-9c07365eb427
	github.com/bhbosman/goCommsStacks v0.0.0-20220715172028-e8be81d8a4d9
	github.com/bhbosman/goFxApp v0.0.0-20220715185456-22d132c8b983
	github.com/bhbosman/goFxAppManager v0.0.0-20220717213634-dd153fc4beda
	github.com/bhbosman/goMessages v0.0.0-20220719163819-d38fc7e6d38c
	github.com/bhbosman/gocommon v0.0.0-20220725200742-9cdc334065f3
	github.com/bhbosman/gocomms v0.0.0-20220707044838-2892132b1367
	github.com/bhbosman/gomessageblock v0.0.0-20220617132215-32f430d7de62
	github.com/bhbosman/goprotoextra v0.0.2-0.20210817141206-117becbef7c7
	github.com/cskr/pubsub v1.0.2
	github.com/golang/protobuf v1.5.2
	github.com/reactivex/rxgo/v2 v2.5.0
	github.com/rivo/tview v0.0.0-20220709181631-73bf2902b59a
	go.uber.org/fx v1.17.1
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-20220513224357-95641704303c
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/bhbosman/goConnectionManager v0.0.0-20220721070628-0f4b3c036d93 // indirect
	github.com/bhbosman/goUi v0.0.0-20220725200743-ddc6ed05f1d6 // indirect
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf // indirect
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.5.1 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/icza/gox v0.0.0-20220321141217-e2d488ab2fbc // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	github.com/teivah/onecontext v0.0.0-20200513185103-40f981bfd775 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/dig v1.14.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20220617134815-f277ff266f47

replace github.com/rivo/tview => ../tview

//replace github.com/rivo/tview => ../tview latest

replace github.com/bhbosman/gocomms => ../gocomms

replace github.com/bhbosman/goFxAppManager => ../goFxAppManager

replace github.com/bhbosman/gocommon => ../gocommon

replace github.com/bhbosman/goCommsStacks => ../goCommsStacks

replace github.com/bhbosman/goCommsNetDialer => ../goCommsNetDialer

replace github.com/bhbosman/goCommsNetListener => ../goCommsNetListener

replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions

replace github.com/bhbosman/goFxApp => ../goFxApp

replace github.com/bhbosman/goUi => ../goUi

replace github.com/bhbosman/goerrors => ../goerrors

replace github.com/bhbosman/goConnectionManager => ../goConnectionManager

replace github.com/bhbosman/goprotoextra => ../goprotoextra

replace github.com/bhbosman/goMessages => ../goMessages

replace github.com/bhbosman/goCommonMarketData => ../goCommonMarketData

replace github.com/cskr/pubsub => ../pubsub

replace github.com/bhbosman/goCommsMultiDialer => ../goCommsMultiDialer
