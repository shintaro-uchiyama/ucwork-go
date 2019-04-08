# 概要
とりあえずgoの最小環境

# 手順
## プログラム作成
```zsh
$ go install github.com/ucwork-go/hello
$ hello
```

## ライブラリ作成
```zsh
$ go build github.com/ucwork-go/stringutil
$ go install github.com/ucwork-go/hello
$ hello
```

## テスト
```zsh
$ go test github.com/ucwork-go/stringutil
```