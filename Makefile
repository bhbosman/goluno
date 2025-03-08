DIR := ${CURDIR}

folders = \
    goerrors \
    goprotoextra \
    gomessageblock \
    goCommsDefinitions \
    gocommon \
    goConnectionManager \
	goUi \
    goCommonMarketData \
    gocomms \
    goCommsMultiDialer \
    goCommsNetDialer \
    goCommsNetListener \
    goCommsSshListener \
    goCommsStacks \
    goConn \
    goFxApp \
    goFxAppManager \
    gokraken \
    goMarketData \
    goMessages \
    gomock \
    goSocks5 \
    goTrader \
    pubsub \
    sshApplication \
    tcell \
    tview \
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

