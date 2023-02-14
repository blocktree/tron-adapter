module github.com/blocktree/tron-adapter

go 1.12

require (
	github.com/Sereal/Sereal v0.0.0-20191211210414-3a6c62eca003 // indirect
	github.com/asdine/storm v2.1.2+incompatible
	github.com/astaxie/beego v1.12.2
	github.com/blocktree/go-owcdrivers v1.2.0
	github.com/blocktree/go-owcrypt v1.1.1
	github.com/blocktree/openwallet/v2 v2.0.5
	github.com/bndr/gotabulate v1.1.2
	github.com/btcsuite/btcutil v0.0.0-20191219182022-e17c9730c422
	github.com/golang/protobuf v1.4.2
	github.com/graarh/golang-socketio v0.0.0-20170510162725-2c44953b9b5f
	github.com/grpc-ecosystem/grpc-gateway v1.8.5
	github.com/imroc/req v0.2.4
	github.com/shopspring/decimal v0.0.0-20200105231215-408a2507e114
	github.com/tidwall/gjson v1.6.5
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107
	google.golang.org/grpc v1.19.1
)

//replace github.com/blocktree/openwallet => ../../openwallet
