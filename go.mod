module ccode

go 1.15

require (
	github.com/Nik-U/pbc v0.0.0-20181205041846-3e516ca0c5d6
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.5.0 // indirect
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20200424173110-d7076418f212
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/hyperledger/fabric-protos-go v0.0.0-20200424173316-dd554ba3746e
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/stretchr/testify v1.6.1
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

//run 'go mod tidy' to add lacked package and remove unused package
//run 'go mod vendor' to copy the requirements to /vendor
//GO111MODULE=on go mod vendor
