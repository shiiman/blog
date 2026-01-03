---
id: 1626
title: 【mac】Macの初期設定をワンライナーで完結させてみた
slug: mac-setting
status: publish
excerpt: こんばんは、しーまんです。 今回はmacの初期設定についての記事の完結編になります。こちらの記事をご覧になれば、毎回手動で設定していた面倒くさいPCの初期設定が、コマンド一つ叩くだけで完了致します。 私自身PCをWind \[…\]
categories:
    - 106
    - 18
tags:
    - 113
    - 114
featured_media: 1628
date: 2022-03-11T19:30:00
modified: 2022-03-09T15:02:12
---

こんばんは、しーまんです。

今回はmacの初期設定についての記事の完結編になります。

こちらの記事をご覧になれば、毎回手動で設定していた面倒くさいPCの初期設定が、 **コマンド一つ** 叩くだけで完了致します。

私自身PCをWindowsからMacに変えてから、既に4、5回macの初期設定を行っておりますが、こちらの設定をしておいたおかげで、かなりスムーズに初期設定を終えることが出来ています。

ですので、今回は**macを使用しているエンジニア**は誰でも設定しておく事をおすすめする、 **ワンライナーで初期設定が完結する方法** を解説致します。

## Macの初期設定とは

まずはmacの初期設定にあたり、何を毎回していて時間がかかっているのかを確認しましょう。

※ エンジニアがやっていることで時間の掛かる設定についてみていきます。

1つ目は **コンソールの設定**

2つ目に **アプリケーションのインストール**

3つ目が **アプリケーションの設定**

主にはこの3点に時間がかかっています。

ではそれぞれを効率化する方法をみていきましょう。

### コンソールの設定

こちらは既に別記事で紹介しておりますが、自分好みのコンソールにするために、bashやzshなどを導入してカスタマイズすることを指します。

下記記事でおすすめの設定方法を解説しておりますので、参考にして下さい。

[![](https://shiimanblog.com/wp-content/uploads/2022/02/eyecatch_shell-320x180.png)\
\
【mac】シェルをカスタマイズしてみた \|【bash】【zsh】【fish】\
\
今回は開発には欠かせないコマンドラインについて自分好みにカスタマイズする方法を紹介します。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2022.02.25](https://shiimanblog.com/engineering/shell/ "【mac】シェルをカスタマイズしてみた |【bash】【zsh】【fish】")

### アプリケーションのインストール

次に行うのが作業で使用するアプリケーションのインストールです。

ブラウザやテキストエディタなどを1つ1つインストールしていくのはとても大変ですよね。

そこで登場するのが **HomeBrew** です。

HomeBrewを使用することで、アプリケーションのインストールを自動化出来ます。

こちらも以前の記事で詳しくまとめておりますので、ぜひ御覧ください。

[![](https://shiimanblog.com/wp-content/uploads/2022/02/eyecatch_homebrew-1-320x180.png)\
\
【mac】HomeBrewでパッケージ管理\|macの初期設定を毎回手動でやっている方必見！\
\
HomeBrewを使用してmacのアプリ導入を管理する方法を紹介します。毎回PCの初期設定を手動で行っている方はぜひ参考にしてください。\
\
![](https://www.google.com/s2/favicons?domain=https://shiimanblog.com)\
\
shiimanblog.com\
\
2022.02.11](https://shiimanblog.com/engineering/homebrew/ "【mac】HomeBrewでパッケージ管理|macの初期設定を毎回手動でやっている方必見！")

### アプリケーションの設定

最後に行うのが、インストールしたアプリケーションの設定になります。

chromeなどのようにgoogleアカウントにデータを紐付けて管理できるものもありますが、毎回設定をやり直しているものもありますよね。

そんな時にアプリケーションの設定を自動化するのが **ドットファイルズ** と呼ばれるものの保存になります。

もちろん例外は多くありますが、macのアプリケーションの設定の多くはホームディレクトリ配下の隠しファイルとして存在しております。それらを総称してファイルの頭に「.」がつくことから「ドットファイルズ」と呼ばれます。

この「ドットファイルズ」をgitなどに保存することでアプリケーションの設定を毎回手動設定することなく、反映させることが可能になります。

シェルの設定も「.bashrc」とか「.zshrc」などのように「ドットファイルズ」でしたよね。

これらを保存しておくというイメージになります。

そして下記のようなシェルを用いて、シンボリックリンクとして配置することで、アプリケーションの設定を反映させることが出来ます。

```
#!/bin/bash

# 通常のドットファイルを定義
DOT_FILES=(.bashrc .zshrc .gitconfig .gitignore_global)

# ホームディレクトリ配下にシンボリックリンクをはる
for file in ${DOT_FILES[@]}
do
    ln -sf ~/dotfiles/$file ~/$file
done
```

同じような考えで、ファイルではなくディレクトリとしてアプリケーションの設定が保存されているものもありますので、そのあたりはご自身で使用しているアプリケーションの設定ファイルを保存して置くようにしましょう！

## ワンライナーで初期設定

さて、ここまでmacの初期設定でエンジニアがしていることを1つずつ自動化してきました。

最後はこれを1コマンドで実行できようにしたら完成です。

やり方は簡単です。

1. gitで必要なファイルを管理する(publicリポジトリ)
2. シェルで一連の操作を実行するようにする

この2です。

私の場合は下記のリポジトリに必要なファイルを配置しております。

ぜひ参考にしてみて下さい。

※ 私以外にも多くの方が **dotfiles** というリポジトリをgitに上げているので参考になると思います。

 [![](https://opengraph.githubassets.com/3a0cdc0b8f8500a5baa19560b979fdae003e32fdc6dd23f73693da549d6e98c8/shiiman/dotfiles)\
\
GitHub - shiiman/dotfiles\
\
Contribute to shiiman/dotfiles development by creating an account on GitHub.\
\
![](https://www.google.com/s2/favicons?domain=https://github.com/shiiman/dotfiles)\
\
github.com](https://github.com/shiiman/dotfiles "GitHub - shiiman/dotfiles")

その中にセットアップシェルを用意しておりますので、少しみていきましょう！

```
#!/bin/bash

# Homebrewのインストール.
if ! type brew >/dev/null 2>&1; then
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

# GitHubから設定ファイルをclone.
if ! type git >/dev/null 2>&1; then
    brew install git
fi
git clone https://github.com/shiiman/dotfiles.git ~/dotfiles

# 設定ファイルフォルダに移動.
cd ~/dotfiles
# ローカルリポジトリにユーザのメールアドレス登録.
git config user.email hsnonsense5@gmail.com

# bundleインストール.
brew tap Homebrew/bundle
# アプリインストール.
brew bundle

##################################################################

# dotfileの設定.
sh ~/dotfiles/dotfile_setup.sh

# zshをシェルリストに追加.
sudo sh -c "echo /usr/local/bin/zsh >> /etc/shells"
# デフォルトシェルをzshに変更.
chsh -s /usr/local/bin/zsh
# ディレクトリ権限変更.
chmod 755 /usr/local/share/zsh
chmod 755 /usr/local/share/zsh/site-functions

...
```

シェルの中でやっていることをリストにすると以下になります。

1. **HomeBrew** のインストール
2. **Git** のインストール
3. git cloneでご自身が用意したリポジトリをクローン
4. brewコマンドでアプリケーションのインストール
5. ドットファイルズの配置(アプリケーションの設定反映)
6. デフォルトシェルの変更

## 実行方法

上記のシェルを実行することで、全ての初期設定を一括で行うことが出来ます。

gitのpublicリポジトリで作成していれば下記のコマンドをmacのコンソールで入力するだけで、実行されます。

```
curl https://raw.githubusercontent.com/[gitユーザ名]/dotfiles/master/mac_setup.sh | bash
```

※ 上記コマンドの「gitユーザ名」の部分をご自身のgitユーザ名に変更して実行して下さい

## まとめ

今回はmacの初期設定をワンライナーで行う方法を紹介しました。

エンジニアは自身の開発したい環境にこだわりがある方が多いので、自身のパフォーマンスが一番発揮できる状態を「 **いつでも**」「 **すぐに**」作れる状態にしておく必要があると思っています。

もしまだ手作業で初期設定をしている方がいらっしゃいましたら、ぜひ今回の記事を参考に、いつでも同じ状態に構築できるように準備しておいてはいかがでしょうか。

今回の記事がエンジニアの効率アップに繋がれば幸いです。