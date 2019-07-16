module github.com/blocktree/tron-adapter

go 1.12

require (
	github.com/asdine/storm v2.1.2+incompatible
	github.com/astaxie/beego v1.11.1
	github.com/blocktree/go-owcdrivers v1.0.12
	github.com/blocktree/go-owcrypt v1.0.1
	github.com/blocktree/openwallet v1.4.6
	github.com/bndr/gotabulate v1.1.2
	github.com/btcsuite/btcutil v0.0.0-20190316010144-3ac1210f4b38
	github.com/golang/protobuf v1.3.1
	github.com/graarh/golang-socketio v0.0.0-20170510162725-2c44953b9b5f
	github.com/grpc-ecosystem/grpc-gateway v1.8.5
	github.com/imroc/req v0.2.3
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/tidwall/gjson v1.2.1
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107
	google.golang.org/grpc v1.19.1
)

//replace github.com/blocktree/openwallet => ../../openwallet
