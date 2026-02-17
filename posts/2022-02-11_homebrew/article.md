---
id: 1580
title: 【mac】HomeBrewでパッケージ管理|macの初期設定を毎回手動でやっている方必見！
slug: homebrew
status: publish
date: 2022-02-11T19:30:00
modified: 2022-02-11T12:27:06
excerpt: MacのパッケージマネージャーHomebrewとmasの導入方法を解説。PCの初期設定を効率化する第一歩を紹介します。
categories:
    - 106
    - 18
tags:
    - 107
    - 108
    - 109
featured_media: 1591
---

こんばんは、しーまんです！

皆さんはPCはWindowsとMacどちらを使用しているでしょうか。

一般の方はPCといえば割とWindowsをイメージされる方のほうが多いのかなと思います。

私の周りのエンジニアでもWindowsとMacは半々くらいのイメージですかね。私自身は2017年にインフラエンジニアとしてジョブチェンジをしたのですが、それ以前はWindowsを使用していました。しかし、ジョブチェンジ以降はもっぱらMacを使用しております。もうWindowsには戻れません(笑)

今回はそんなMacユーザに向けた記事になります。

開発を行う上でまず最初にするのはPCの設定ではないでしょうか？フリーのエンジニアですと職場が変更になるたびに社用PCに初期設定を行います。有期でない方も、社用PCはレンタルのものが多いです。そのレンタルの期間が大体2〜3年だったりするので、その期限毎にPCを変更するのが普通だったりします。

つまり何回も何回もPCの初期設定をすることになるんですよね。今回はそんな何度も行う初期設定を手動で行っている方は必見です。 **HomeBrew** を使ってパッケージ管理していきましょう！というお話になります。

## HomeBrewとは

ということで **HomeBrew** についてみていきましょう！

そもそも **HomeBrew** とは何でしょうか。

**HomeBrew** とはMacのアプリをパッケージ管理するものです。

といってもよく分からないですよね。

RedHat系の「 **yum**」とかDebian系の「 **apt-get**」みたいなものだと思ってください。

つまりは今まで手動でインストールしていたMac用のアプリケーションを「 **Brewfile**」というファイル1つで管理し、コマンドでインストールすることができるツールです。

## 導入方法

では早速HomeBrewをインストールしていきましょう。

### HomeBrewのインストール

macでコンソールを開き、以下のコマンドを打つことでインストール出来ます！

```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

インストールの途中でパスワードを聞かれるので、PCにログインする時のパスワードを入力しましょう

インストールが完了したら以下のコマンドを入力して動作確認を行います。

```
brew -v
```

バージョンが表示されたらインストール完了です。

### Brew Bundleのインストール

次に **BrewBundle** もインストールしておきましょう。

「 **brew**」コマンドがコンソールでアプリをインストールできるのに対して、「 **brew bundle**」コマンドは現状の使用しているアプリの情報を「Brewfile」に書き込んだり、逆に「Brewfile」からアプリをまとめてインストールすることができるコマンドです。

下記のコマンドでインストールしておきましょう。

```
brew tap Homebrew/bundle
```

こちらで「brew bundle」コマンドも使用できるようになりました。

### masのインストール

最後に「 **mas**」CLIもインストールしておきましょう。

管理するアプリの対象により管理するコマンドが違いますので、こちらもインストールしておくと便利です。以下が各コマンドと管理するアプリの大まかな枠組みです。

ツール管理対象HomebrewCUIアプリHomebrew CaskGUIアプリmas-cliMac App Storeアプリ※ Homebrew Caskとmas-cliは同アプリが含まれている場合があります。

ということで、masのインストールは以下になります。

```
brew install mas
```

ここまで設定出来れば準備OKです。

## よく使うコマンド

ではまず、「brew」コマンドでよく使用するものを確認しておきましょう。ファイルで管理する前に手動でもアプリをインストールできるようにしておきます。

・バージョン確認

**brew -v**

・アップデート(homebrew自体のアップデート)

**brew update**

・アップグレード(formulaeのアップデート)

**brew upgrade**

**mas upgrade**

・homebrewの情報表示

**brew –config**

・formulae一覧

**brew list**

**mas list**

・formulae検索

**brew search \[formulae\] (–cask)**

**mas search \[formulae\]**

・ヘルスチェック

**brew doctor**

・formulaeインストール

**brew install \[formulae\]**

**mas install \[formulae\]**

・formulaeの有効化/無効化

**brew link \[formulae\]**

**brew unlink \[formulae\]**

・formulaeの削除

**brew remove \[formulae\]**

**mas uninstall \[formulae\]**

**Homebrew** では管理するアプリ(パッケージ)のことをformulaeと言います。

## Brewfile

手動でアプリの管理ができるようになったところで、ここからが本番です。

### Brewfileの作成

まずは現在使用しているアプリを「 **Brewfile**」に吐き出してみましょう！

やり方は簡単で下記のコマンドを実行することで「 **Brewfile**」が出来上がります。

```
brew bundle dump
```

下記のようなファイルが出来上がっていると思います。

[![Brewfile](https://shiimanblog.com/wp-content/uploads/2022/02/brewfile.png)](https://shiimanblog.com/wp-content/uploads/2022/02/brewfile.png)

**cask** や **mas** については全てのアプリがリスト化されない場合があります。

その場合は **search** コマンドを使用し、必要なアプリを手動で追加しましょう！

こちらのファイルをgitなどで管理することをおすすめします。

### Brewfileからアプリのインストール

「 **Brewfile」** が出来上がったら、そのファイルを利用しアプリをインストールしてみましょう。

PCの初期設定時にまず、上記で行った「 **Homebrew** の導入設定」を行います。

そしてBrewfileを配置し、下記コマンドを入力しましょう。

```
brew bundle
```

するとアプリのインストールを自動で行ってくれます！

これで毎回行っていたPCの初期設定で時間の掛かっていたアプリのインストールは完了です。

とっても簡単ですよね。また設定時間も大幅に短縮できる他、同じアプリを直ぐに使用可能になるため環境作りも簡単です。

## しーまんのBrewfile

ここまでで **Homebrew** の魅力は十分伝わっているかと思います。

そこで私自身の「 **Brewfile**」をサンプルとして紹介しておきます。

管理アプリの参考や、書き方の参考にしてください。

 [![](https://opengraph.githubassets.com/3a0cdc0b8f8500a5baa19560b979fdae003e32fdc6dd23f73693da549d6e98c8/shiiman/dotfiles)\
\
dotfiles/Brewfile at master · shiiman/dotfiles\
\
Contribute to shiiman/dotfiles development by creating an account on GitHub.\
\
![](https://www.google.com/s2/favicons?domain=https://github.com/shiiman/dotfiles/blob/master/Brewfile)\
\
github.com](https://github.com/shiiman/dotfiles/blob/master/Brewfile "dotfiles/Brewfile at master · shiiman/dotfiles")

## まとめ

今回は **Mac** の初期設定を楽にする方法として「 **Homebrew**」を使用してアプリを管理する方法を紹介しました。

**Homebrew** を使用して「 **Brewfile**」を作成し、管理することでご自身で使用しているアプリケーションを直ぐにインストールすることが可能です。こちらのファイルを作っておくことでPCの設定に掛かる時間は大幅に減少されます。

また、今まで使っていた環境をそのまま再現出来ますので、環境構築する上ではとても重宝するのではないでしょうか。

ぜひまだ、使ったことないよーという方はこの機会に「 **Homebrew**」を使ってみてください。

今回の記事がMacユーザの参考になれば幸いです！