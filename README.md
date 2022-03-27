# stamp-dl

## Description

LINE STORE から LINE スタンプの画像(.png)をダウンロードするためのコマンドラインツール

## Download

ターミナルで ↓ を実行 (どこからでも`stamp-dl`コマンドが使えるようになるはず...)  
Homebrew でのインストールか、手動でのインストールか、どちらか選んでね。

### Homebrew

```console
brew tap makitune/stamp-dl
brew install stamp-dl
```

### Manual Install

搭載されている CPU アーキテクチャによって、少し異なるので注意！  
どちらにすべきかわからない場合は、ターミナルで ↓ を実行してみてね

```console
arch
```

#### Intel processors (amd64)

```console
curl -L 'https://github.com/makitune/stamp-dl/releases/download/v0.1.0/stamp-dl-0.1.0.darwin_amd64.bottle.tar.gz' | tar xz; chmod 755 stamp-dl; sudo mv stamp-dl /usr/local/bin/
```

実行するとパスワードを聞かれるので、Mac にログインするときに使うパスワードを入力する。  
成功すると`/usr/local/bin/`にバイナリが配置される。(削除する時はココを探して！)

#### Apple silicon (arm64)

```console
curl -L 'https://github.com/makitune/stamp-dl/releases/download/v0.1.0/stamp-dl-0.1.0.darwin_arm64.bottle.tar.gz' | tar xz; chmod 755 stamp-dl; sudo mv stamp-dl /opt/homebrew/bin/
```

実行するとパスワードを聞かれるので、Mac にログインするときに使うパスワードを入力する。  
成功すると`/opt/homebrew/bin/`にバイナリが配置される。(削除する時はココを探して！)

## Usage

引数に LINE STORE の LINE スタンプページの URL を指定する。

```console
stamp-dl "ここにLINEスタンプページのURLを入れて！"
```

### 例

- 通常  
  実行したディレクトリに保存される。

```console
stamp-dl "https://store.line.me/stickershop/product/5647174/ja"
```

![例](./images/sample.png)

- デスクトップに保存する場合

```console
stamp-dl -o ~/Desktop "https://store.line.me/stickershop/product/5647174/ja"
```

- 完了時に Finder で表示しない場合

```console
stamp-dl -q "https://store.line.me/stickershop/product/5647174/ja"
```

- いっぱい

```console
stamp-dl "https://store.line.me/stickershop/product/5647174/ja" "https://store.line.me/stickershop/product/1625544/ja" "https://store.line.me/stickershop/product/1435028/ja"
```
