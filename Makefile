folders = \
    gocommon
#    goCommonMarketData \
#    gocomms \
#    goCommsDefinitions \
#    goCommsMultiDialer \
#    goCommsNetDialer \
#    goCommsNetListener \
#    goCommsSshListener \
#    goCommsStacks \
#    goConn \
#    goConnectionManager \
#    goerrors \
#    goFxApp \
#    goFxAppManager \
#    gokraken \
#    goMarketData \
#    gomessageblock \
#    goMessages \
#    gomock \
#    goPolygon-io \
#    goprotoextra \
#    goSocks5 \
#    goTrader \
#    goUi \
#    pubsub \
#    sshApplication \
#    tcell \
#    tview



all_mods:

	go get -v -u all #gosetup


$(folders):
	echo $@
	cd ../$@
	go get -v -u all



ww:
	make $< gocommon