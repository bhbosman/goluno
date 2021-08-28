module github.com/bhbosman/goLuno

go 1.15

require (
	github.com/bhbosman/goMessages v0.0.0-20210414134625-4d7166d206a6
	github.com/bhbosman/gocommon v0.0.0-20210817162338-4486f68fa3f4
	github.com/bhbosman/gocomms v0.0.0-20210414144344-fb75f75793be
	github.com/bhbosman/gologging v0.0.0-20200921180328-d29fc55c00bc
	github.com/bhbosman/gomessageblock v0.0.0-20210414135653-cd754835d03b
	github.com/bhbosman/goprotoextra v0.0.2-0.20210817141206-117becbef7c7
	github.com/cskr/pubsub v1.0.2
	github.com/emirpasic/gods v1.12.0
	github.com/golang/protobuf v1.4.2
	github.com/kardianos/service v1.1.0
	github.com/reactivex/rxgo/v2 v2.5.0 // indirect
	go.uber.org/fx v1.13.1
	google.golang.org/protobuf v1.25.0
)

replace github.com/bhbosman/goMessages => ../goMessages

replace github.com/bhbosman/gocomms => ../gocomms
replace github.com/bhbosman/gomessageblock => ../gomessageblock