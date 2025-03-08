DIR := ${CURDIR}

folders = \
	goPolygon-io \
    goerrors \
    goprotoextra \
    gomessageblock \
    goCommsDefinitions \
    gocommon \
    goConnectionManager \
    pubsub \
    tcell \
    tview \
	goUi \
    goCommonMarketData \
    gocomms \
    goCommsStacks \
    goCommsNetDialer \
    goCommsMultiDialer \
    goCommsNetListener \
    goCommsSshListener \
    goConn \
    goFxApp \
    goFxAppManager \
    goMessages \
    gomock \
    goSocks5 \
    goTrader \
    sshApplication \
    goMarketData \
    gokraken \
    goLuno

$(folders):
	make --ignore-errors -f ${CURDIR}/Makefile -C ${CURDIR}/../$@ update-go-mod

update-go-mod:
	@set GOROOT=/opt/homebrew/opt/go/libexec
	@set GOPATH=/Users/ronelspijkerman
	/opt/homebrew/opt/go/libexec/bin/go get -d -v -u all
	git add go.mod
	git add go.sum
	git commit -m "fix go mod issues"
	git push

ww:
	make gocommon
	make goCommonMarketData
	make gocomms
	make goCommsDefinitions
	make goCommsMultiDialer
	make goCommsNetDialer
	make goCommsNetListener
	make goCommsSshListener
	make goCommsStacks
	make goConn
	make goConnectionManager
	make goerrors
	make goFxApp
	make goFxAppManager
	make gokraken
	make goMarketData
	make gomessageblock
	make goMessages
	make gomock
	make goprotoextra
	make goSocks5
	make goUi
	make pubsub
	make sshApplication
	make tcell
	make tview
	make goLuno
	make goPolygon-io
	#make goTrader