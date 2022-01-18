module github.com/bhbosman/goLuno

go 1.15

require (
	github.com/bhbosman/goMessages v0.0.0-20210819131032-dfe3cad9135f
	github.com/bhbosman/gocommon v0.0.0-20211112042828-c309f7957be1
	github.com/bhbosman/gocomms v0.0.0-20211124080017-8ffe6f0e804c
	github.com/bhbosman/gomessageblock v0.0.0-20211029070951-75b9d5ae1fe6
	github.com/bhbosman/goprotoextra v0.0.2-0.20210817141206-117becbef7c7
	github.com/cskr/pubsub v1.0.2
	github.com/emirpasic/gods v1.12.0
	github.com/golang/protobuf v1.4.2
	github.com/kardianos/service v1.1.0
	go.uber.org/fx v1.14.2
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.16.0
	google.golang.org/protobuf v1.25.0
)

//replace github.com/bhbosman/gocomms => ../gocomms

replace github.com/bhbosman/gocommon => ../gocommon
