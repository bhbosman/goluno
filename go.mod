module github.com/bhbosman/goLuno

go 1.15

require (
	github.com/bhbosman/goMessages v0.0.0-20210414134625-4d7166d206a6
	github.com/bhbosman/gocommon v0.0.0-20211028204315-fbe0347a4b3f
	github.com/bhbosman/gocomms v0.0.0-20210901083025-a45dff1c542b
	github.com/bhbosman/gomessageblock v0.0.0-20211029070951-75b9d5ae1fe6
	github.com/bhbosman/goprotoextra v0.0.2-0.20210817141206-117becbef7c7
	github.com/cskr/pubsub v1.0.2
	github.com/emirpasic/gods v1.12.0
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/kardianos/service v1.1.0
	go.uber.org/fx v1.14.2
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.16.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/bhbosman/goMessages => ../goMessages

replace github.com/bhbosman/gocomms => ../gocomms

replace github.com/bhbosman/gocommon => ../gocommon
