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
#	make goerrors
#	make goprotoextra
#	make gomessageblock
#	make goCommsDefinitions
#	make gocommon
#	make goConnectionManager
#	make pubsub
#	make tcell
#	make tview
	make goUi
#	make goCommonMarketData
#	make gocomms
#	make goCommsStacks
#	make goCommsNetDialer
#	make goCommsMultiDialer
#	make goCommsNetListener
#	make goCommsSshListener
#	make goConn
#	make goFxApp
#	make goFxAppManager
#	make goMessages
#	make gomock
#	make goSocks5
#	make goTrader
#	make sshApplication
#	make goMarketData
#	make gokraken
#	make goLuno
#	make goPolygon-io
