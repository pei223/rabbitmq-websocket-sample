## Setup
RabbitMQが立ち上がっていること.

Go 1.18で書いていて、標準以外は以下のようなライブラリを使用.
- chi
- envprocess
- amqp
- zerolog

```
go mod tidy
go run ./cmd/apiserver/.
go run ./cmd/websocketserver/.
```

## 内容
### apiserver
APIを介してRabbitMQにメッセージを積んでいく。

### websocketserver
クライアントとのセッションを開始・保持する.

#### worker
goroutineでwebsocketserverと並列で動かす.

RabbitMQからメッセージを取り出して、ユーザー名がメッセージのToのセッションに対してメッセージ送信する,


## RabbitMQについてのメモ
ACKでメッセージは破棄される。
https://www.rabbitmq.com/confirms.html#acknowledgement-modes

NACK/Rejectで再度キューに突っ込むことができるらしい。
https://www.rabbitmq.com/confirms.html#consumer-nacks-requeue



## フォルダ構成の参考資料
https://github.com/golang-standards/project-layout/blob/master/README_ja.md


