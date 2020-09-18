build:
	go build

install: build
	go install
	#luno-stream-service -service stop
	luno-stream-service -service uninstall
	luno-stream-service -service install
	luno-stream-service -service start



clean:
	rm luno-stream-service
