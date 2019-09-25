# 概要
Bookshelfチュートリアルを見ながらGCPサービスを一通り触ってみる  

# ローカル環境構築
## 事前準備
CloduSQLを利用する場合以下コマンドでプロキシー起動
```zsh
./scripts/cloud_sql_proxy
```

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
