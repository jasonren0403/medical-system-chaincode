# medical-system-chaincode

“基于区块链的智慧医疗系统”毕设项目的链码部分仓库

## API笔记

* `peer.Response`结构
  * `response.status`：200/400/500
  * `response.payload`：`base64`编码后的返回
  * `response.message`：错误信息，由`shim.Error`指定
* 布置链码
  > 教程参考：[Fabric官方教程 - 英文版](https://hyperledger-fabric.readthedocs.io/en/latest/test_network.html "Fabric官方教程")
  1. 开启
  ```shell
  ./network.sh up
  ```
    2. 建立频道
  ```shell
  ./network.sh createChannel -c <channelName>
  ./network.sh createChannel -c fa-jason
  ```
    3. 安装链码
  ```shell
  ./network.sh deployCC -ccn <chaincode-name> -ccp <chaincode-filepath> -ccl <chaincode-language>
  ./network.sh deployCC -ccn MedicalSystem -ccp ../chaincode/MedicalSystem -ccl go -c fa-jason
  ```
    4. 调用函数
  ```shell
  peer chaincode instantiate -C fa-jason
  peer chaincode invoke -C fa-jason -c '{"function":"initLedger","Args":[]}'
  peer chaincode query -C fa-jason -c '...'
  ```
    5. 使用结束
  ```shell
  ./network.sh down
  ```

### 手动安装合约

> 需要从go代码目录执行下列命令

1. 打包：

```shell
GO111MODULE=on go mod vendor
cd ../../test-network
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
peer lifecycle chaincode package ms.tar.gz --path ../chaincode/MedicalSystem --lang golang --label ms_1.0
```

2. 安装：
    1. 节点1名义
    ```shell
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051
    
    peer lifecycle chaincode install ms.tar.gz
    ```
    2. 节点2名义
    ```shell
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051
    
    peer lifecycle chaincode install ms.tar.gz
    ```
3. 接受链码定义：

```shell
peer lifecycle chaincode queryinstalled

export CC_PACKAGE_ID=<pkgid from last order>

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID fa-jason --name MedicalSystem --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
```

> 需切换至节点1，继续运行```approveformyorg```命令

```shell
  export CORE_PEER_LOCALMSPID="Org1MSP"
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
  export CORE_PEER_ADDRESS=localhost:7051
      
  peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID fa-jason --name MedicalSystem --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
```

4. 提交链码定义至频道

```shell
peer lifecycle chaincode checkcommitreadiness --channelID fa-jason --name MedicalSystem --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --output json
```

> 输出应均为```true```

```shell
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID fa-jason --name MedicalSystem --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
```

5. 检验提交状态

```shell
peer lifecycle chaincode querycommitted --channelID fa-jason --name MedicalSystem --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
```
