# medical-system-chaincode

“基于区块链的智慧医疗系统”毕设项目的链码部分仓库

## API笔记

* `peer.Response`结构
    * `response.status`：200/400/500
    * `response.payload`：`base64`编码后的返回
    * `response.message`：错误信息，由`shim.Error`指定