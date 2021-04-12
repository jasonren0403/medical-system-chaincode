module ccode

go 1.15

require (
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20200424173110-d7076418f212
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/hyperledger/fabric-protos-go v0.0.0-20200424173316-dd554ba3746e
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/stretchr/testify v1.5.1
)

//run 'go mod tidy' to add lacked package and remove unused package
//run 'go mod vendor' to copy the requirements to /vendor
//GO111MODULE=on go mod vendor
