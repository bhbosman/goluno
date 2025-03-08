folders = \
    gocommon \
    goCommonMarketData \
    gocomms \
    goCommsDefinitions \
    goCommsMultiDialer \
    goCommsNetDialer \
    goCommsNetListener \
    goCommsSshListener \
    goCommsStacks \
    goConn \
    goConnectionManager \
    goerrors \
    goFxApp \
    goFxAppManager \
    gokraken \
    goMarketData \
    gomessageblock \
    goMessages \
    gomock \
    goPolygon-io \
    goprotoextra \
    goSocks5 \
    goTrader \
    goUi \
    pubsub \
    sshApplication \
    tcell \
    tview

$(folders):
	echo $@
	cd ../$@
	go get -v -u all



ww:
	make $< gocommon
	make $< goCommonMarketData
	make $< gocomms
	make $< goCommsDefinitions
	make $< goCommsMultiDialer
	make $< goCommsNetDialer
	make $< goCommsNetListener
	make $< goCommsSshListener
	make $< goCommsStacks
	make $< goConn
	make $< goConnectionManager
	make $< goerrors
	make $< goFxApp
	make $< goFxAppManager
	make $< gokraken
	make $< goMarketData
	make $< gomessageblock
	make $< goMessages
	make $< gomock
	make $< goPolygon
	make $< goprotoextra
	make $< goSocks5
	make $< goTrader
	make $< goUi
	make $< pubsub
	make $< sshApplication
	make $< tcell
	make $< tview

