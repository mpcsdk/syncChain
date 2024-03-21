module syncChain

go 1.15

require (
	github.com/ethereum/go-ethereum v1.13.14
	github.com/glacjay/goini v0.0.0-20161120062552-fd3024d87ee2
	github.com/gogf/gf/contrib/drivers/pgsql/v2 v2.6.4
	github.com/gogf/gf/v2 v2.6.4
	github.com/holiman/uint256 v1.2.4
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/mpcsdk/mpcCommon v0.0.0
	github.com/nats-io/nats.go v1.34.0
	golang.org/x/crypto v0.18.0
)

replace github.com/mpcsdk/mpcCommon v0.0.0 => ./mpcCommon
