# 概要
Bookshelfチュートリアルを見ながらGCPサービスを一通り触ってみる  

# ローカル環境構築
## ビルド&実行
```zsh
go build -o build/ucwork cmd/ucwork/main.go
./build/ucwork
```

## URLへアクセス
```zsh
curl http://localhost:8080/members                                                                                            ![add_datastore_connection#6]
[{"Name":"Name1"},{"Name":"Name2"}]%
```
