module github.com/bhbosman/goLuno

go 1.15

require (
	github.com/bhbosman/goMessages v0.0.0-20210819131032-dfe3cad9135f
	github.com/bhbosman/gocommon v0.0.0-20220120133819-681d385f0463
	github.com/bhbosman/gocomms v0.0.0-20211124080017-8ffe6f0e804c
	github.com/bhbosman/gomessageblock v0.0.0-20211029070951-75b9d5ae1fe6
	github.com/bhbosman/goprotoextra v0.0.2-0.20210817141206-117becbef7c7
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/cskr/pubsub v1.0.2
	github.com/emirpasic/gods v1.12.0
	github.com/gdamore/tcell/v2 v2.5.1
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/kardianos/service v1.1.0
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/tview v0.0.0-20220307222120-9994674d60a8
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/fx v1.14.2
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.21.0
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f
	golang.org/x/tools v0.1.9 // indirect
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect

)

replace github.com/bhbosman/gocomms => ../gocomms

replace github.com/golang/mock => ../gomock

replace github.com/bhbosman/gocommon => ../gocommon
