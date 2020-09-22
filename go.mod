module github.com/bhbosman/goLuno

go 1.15

require (
	github.com/bhbosman/goMessages v0.0.0-20200918071950-29c6c3c09ea4
	github.com/bhbosman/gocommon v0.0.0-20200921215456-bfddd9bb050e
	github.com/bhbosman/gocomms v0.0.0-20200921215103-85dd2b219cf0
	github.com/bhbosman/gologging v0.0.0-20200921180328-d29fc55c00bc
	github.com/bhbosman/gomessageblock v0.0.0-20200921180725-7cd29a998aa3
	github.com/bhbosman/goprotoextra v0.0.1
	github.com/cskr/pubsub v1.0.2
	github.com/emirpasic/gods v1.12.0
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/kardianos/service v1.1.0
	github.com/reactivex/rxgo/v2 v2.1.0
	go.uber.org/fx v1.13.1
	google.golang.org/protobuf v1.25.0
)

replace github.com/reactivex/rxgo/v2 v2.1.0 => github.com/bhbosman/rxgo/v2 v2.1.1-0.20200922152528-6aef42e76e00

//for DEV
replace github.com/bhbosman/gocommon => /Users/brendanbosman/src/github.com/bhbosman/gocommon

replace github.com/bhbosman/gocomms => /Users/brendanbosman/src/github.com/bhbosman/gocomms

replace github.com/bhbosman/goMessages => /Users/brendanbosman/src/github.com/bhbosman/goMessages

//replace github.com/reactivex/rxgo/v2 v2.1.0 => /Users/brendanbosman/src/github.com/ReactiveX/RxGo
