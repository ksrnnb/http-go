# http-go

## 立ち上げ
下記コマンドで実行。
```bash
docker-compose up -d
docker-compose exec app ash
go run main.go
```

socketの場合は別プロセスでcurlなどを実行する。
```bash
curl localhost:3000/index.html
```

ファイルをinputとする場合は、`test.txt`を読み込んで、レスポンスを標準出力に書き込む。

## socket通信する場合
`app/config.ini`のSERVER_ENVをsocketに設定する
```bash
SERVER_ENV=socket
```

## ファイル（test.txt）を使用する場合
socket以外であればなんでもOK
```bash
SERVER_ENV=***
```

## 脆弱性
ディレクトリトラバーサルなどの脆弱性を含んでいるので注意。

## 参考書籍
[ふつうのLinuxプログラミング 第2版](https://www.sbcr.jp/product/4797386479/)