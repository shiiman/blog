---
id: 1594
title: 【mac】シェルをカスタマイズしてみた |【bash】【zsh】【fish】
slug: shell
status: publish
excerpt: こんばんは、しーまんです。 皆さんはmacで作業する時、コマンドラインを使用しますか？エンジニアなら必ず使用しますよね。 しかしエンジニアであっても意外と標準のコマンドラインをそのまま使っている方が多いような気がします。 \[…\]
categories:
    - 106
    - 18
tags:
    - 110
    - 111
    - 112
featured_media: 1596
date: 2022-02-25T19:30:00
modified: 2022-02-17T15:46:42
---

こんばんは、しーまんです。

皆さんはmacで作業する時、コマンドラインを使用しますか？

エンジニアなら必ず使用しますよね。

しかしエンジニアであっても意外と標準のコマンドラインをそのまま使っている方が多いような気がします。そんな方をみると「あぁー、勿体ないなぁ〜」と思ってしまいます。

シェルには色々な機能を追加したり、自分好みにカスタマイズする機能が備わっています。コマンドラインを多く使う人ほど、標準の機能では物足りなく感じるはずです。

今回はシェルをカスタマイズして、**作業効率をグンッとアップ**させる方法を紹介したいと思います！

## シェルの歴史

まずはシェルといっても種類が豊富にあるので、自分に合ったシェルを見つけるところから始めます。

そのために軽くシェルの歴史に関して触れておきましょう。

シェルには主に「 **Bシェル系**」と「 **Cシェル系**」と呼ばれる分類があります。

macではデフォルトで数種類のシェルが入っていて、両系統のシェルが入っていることが確認出来ます。

```
$ cat /etc/shells

/bin/bash
/bin/csh
/bin/dash
/bin/ksh
/bin/sh
/bin/tcsh
/bin/zsh
```

[![シェルの歴史](https://shiimanblog.com/wp-content/uploads/2022/02/history.png)](https://shiimanblog.com/wp-content/uploads/2022/02/history.png)

それぞれのシェルの簡単な特徴を下記に上げておきます。

sh様々な亜種が誕生したため、現在はそれらと区別するためにBourneシェル（略してBシェル）と呼ばれている。現在ではログインシェルとして使われることはほとんどないが、シェルスクリプトを書くためによく使われている。kshKornシェル(ksh)はBourneシェルの上位互換シェルであり、Kシェルとも呼ばれる。Bourneシェルの構文を引き継ぎつつ、Bourneシェルにはない便利な機能を備えている。bashBourne Againシェル(bash)はBourneシェルの上位互換シェルである。Bourneシェルの構文を引き継ぎつつ、Bourneシェルにはない便利な機能を備えている。zshZシェル(zsh)はシェルの中でもとくに多機能なシェルである。cshCシェル(csh)はC言語に似た構文を持つシェルである。Bourneシェルにはない便利な機能を備えている。Bourneシェルとは互換性がなく、同じ機能でもその構文が異なる。tcshTCシェル(tcsh)はCシェルの上位互換シェルである。

今回はその中でも長年macのデフォルトシェルであった「 **bash**」とmacOS Catalina からデフォルトシェルになった「 **zsh**」、そして導入/設定が簡単だと好評な「 **fish**」をみていきます。

因みにmacでは下記のコマンドでデフォルトシェルを変更出来ますので、自分でお気に入りのシェルに変更してみましょう。

```
chsh -s [シェルのパス]

# 例.
chsh -s /bin/zsh
```

## bash

まずは **bash** をみていきます。bashは多くの **Linux OSデフォルトシェル** として設定されています。

シェルといえば多くの方はこの **bash** のことだと思っていると思います。それだけ多く使用されているシェルですね。

bashはデフォルトだととてもシンプルな設定になっています。

しかし、 **設定ファイルに拡張条件を記入する** ことで多くの恩恵を受けることができます。

設定ファイルですが、シェルの起動には複数の設定ファイルが存在し、読み込まれるタイミングの違いにより、どのファイルに記載した方がよいのかが変わってきます。

下記がその起動時に読み込まれるファイル順です。

[![bash起動処理](https://shiimanblog.com/wp-content/uploads/2022/02/d9da7ab28fab49a1a05f541abb2788fe-1-800x429.png)](https://shiimanblog.com/wp-content/uploads/2022/02/d9da7ab28fab49a1a05f541abb2788fe-1.png)

とはいえ、基本的には「 **~/.bashrc**」に記載しておけば問題ありませんので、こちらに設定を記入していきましょう。

### bashのカスタマイズ例

シェルのカスタマイズは何をやりたいのかによって、人それぞれの設定がありますので、正解とかはありません。今回はどんなことができるのかを少し紹介するのと、私自身の設定ファイルを紹介します。

▼ aliasの設定 ======================

ex) alias la=’ls -a’

説明) laとコマンドで打つとls -aを打つのと同じ挙動になる

▼ プロンプト設定 ======================

ex) PS1=”\[\\\[\\e\[0;32m\\\]\\u\\\[\\e\[0m\\\]@\\\[\\e\[0;36m\\\]\\h\\\[\\e\[0m\\\]\] % “

説明) プロンプトの端にログイン名などを表示させる

▼ 環境変数の設定 ======================

ex) export CLICOLOR=true

説明) lsした時にファイルの種類ごとに色をつけてくれる

▼ 関数定義 ======================

ex) function cdlspwd() {cd $1; la; pwd;}

alias cd=cdlspwd

説明) cdの後にlsとpwdを実行する関数

上記で紹介した拡張はほんの一部です。他にも **入力補完** の設定などはよく行いますね。

さまざまな拡張が出来ますので、ぜひ自分に合った設定を試してみて下さい。

因みに私の設定ファイルは下記になります。参考にしてみて下さい。

 [![](https://opengraph.githubassets.com/3a0cdc0b8f8500a5baa19560b979fdae003e32fdc6dd23f73693da549d6e98c8/shiiman/dotfiles)\
\
dotfiles/.bashrc at master · shiiman/dotfiles\
\
Contribute to shiiman/dotfiles development by creating an account on GitHub.\
\
![](https://www.google.com/s2/favicons?domain=https://github.com/shiiman/dotfiles/blob/master/.bashrc)\
\
github.com](https://github.com/shiiman/dotfiles/blob/master/.bashrc "dotfiles/.bashrc at master · shiiman/dotfiles")

## zsh

次に紹介するのは私もメインに使用している「 **zsh**」です。最新のmacのデフォルトがこちらになったため、使用している方も増えたのではないでしょうか。

zshの特徴として、bashよりも圧倒的にカスタマイズができる点とプラグインの豊富さが上げられます。

カスタマイズにはbashと同じように「 **~/.zshrc**」というファイルに記述していきます。

プラグインには今まで「 **zplug**」というプラグインマネージャを使用していたのですが、今回の執筆にあたり「 **zinit(旧 zplugin)**」に変更しました。変更理由としてはより起動が高速になりそうだったからです。

変更後にしばらく使用してみましたが、確かに体感で起動が高速になったように感じます。zinitはoss作者のリポジトリが急に削除され現在は有志の方が新しいリポジトリを作成して運用されていたりと、少し不安な部分もありますが、使用感はとてもいいのでしばらくは使用していこうと思っています。

ということで、私のzshの設定も共有しておきます。

参考にしていただけると助かります。

 [![](https://opengraph.githubassets.com/3a0cdc0b8f8500a5baa19560b979fdae003e32fdc6dd23f73693da549d6e98c8/shiiman/dotfiles)\
\
dotfiles/.zshrc at master · shiiman/dotfiles\
\
Contribute to shiiman/dotfiles development by creating an account on GitHub.\
\
![](https://www.google.com/s2/favicons?domain=https://github.com/shiiman/dotfiles/blob/master/.zshrc)\
\
github.com](https://github.com/shiiman/dotfiles/blob/master/.zshrc "dotfiles/.zshrc at master · shiiman/dotfiles")

## fish

最後に紹介するのが、今回はじめて出てきた「 **fish**」というシェルです。

設定ファイルの作成が面倒くさい！最初から色々設定しておいてよ！！という方にはこちらのシェルがおすすめです。

 [![](https://s.wordpress.com/mshots/v1/https%3A%2F%2Ffishshell.com%2F?w=320&h=180)\
\
fish shell\
\
A smart and user-friendly command line shell\
\
![](https://www.google.com/s2/favicons?domain=https://fishshell.com/)\
\
fishshell.com](https://fishshell.com/ "fish shell")

fishはzshのように **高機能** でありながら、インストールしただけで高性能な補完やシンタックスハイライトなど様々な便利機能がデフォルト使用できます。

その特徴は下記になります。

▼ 特徴 =======================

- デフォルトの設定が貧弱→fishならば無設定で便利機能満載
- シェルスクリプトが苦手→fishスクリプトならばシンプルでわかりやすい
- 無駄に多機能で覚えきれない→fishは必要な機能を厳選している
- マニュアルがわかりづらい→fishのマニュアルは具体例満載でわかりやすい
- 設定ファイルを書くのが面倒→fishならばWebブラウザで設定できる！
- 補完設定を書くのが面倒→fishならばmanページを解析して自動で補完設定をしてくれる
- シングルクォート中にシングルクォートを書けない→fishならば \\’と書ける
- 関数や環境変数を保存するのに設定ファイルの書換えが面倒→fishならばその場で永続化できる
- 和訳マニュアルがない→fishの最新版公式文書を全訳済

▼ デミリット =======================

- まだ新しいシェルなので使っている人が少ない
- 現状のシェルスクリプトが動かない可能性がある

### macでの導入方法

fishはmacにデフォルトで備わっていないので、インストールする必要があります。

とはいえインストールは簡単です。

```
brew install fish

# デフォルトシェルに設定.
echo $(which fish) | sudo tee -a /etc/shells
chsh -s $(which fish)
```

上記3行で終わりです。

簡単ですよね。

大本の初期設定は下記コマンドで行います。

コマンドを入力すると各種設定や履歴などが確認できるページがブラウザで表示されます。

```
fish_config
```

こちらで大体の設定は出来ます。

さらにzshのように独自でカスタマイズをされたい方は「 **fisher**」というパッケージ管理を導入してみましょう。(設定が面倒な方はここまでの設定でも十分だと思います。)

 [![](https://repository-images.githubusercontent.com/43098836/48e9a909-99be-41ae-852a-06469c412e45)\
\
GitHub - jorgebucaran/fisher: A plugin manager for Fish\
\
A plugin manager for Fish. Contribute to jorgebucaran/fisher development by creating an account on GitHub.\
\
![](https://www.google.com/s2/favicons?domain=https://github.com/jorgebucaran/fisher)\
\
github.com](https://github.com/jorgebucaran/fisher "GitHub - jorgebucaran/fisher: A plugin manager for Fish")

私は自身で細かく設定したい派なのでzshを普段使いしておりますが、 **設定が面倒だけどbashよりはもっと便利なシェルを使いたい** という方はぜひ一度 **fish** を試してみて下さい。

## まとめ

今回はコマンドラインをカスタマイズするために「 **bash**」「 **zsh**」「 **fish**」のそれぞれの特徴や設定を紹介しました。

私自身は **zsh** をカスタマイズしてメインで使用していますが皆さんはどのシェルがお好みでしょうか？

デフォルトのまま使用しているよ！という方はカスタマイズがほとんど不要だけど高性能な **fish** を一度試していただきたいです。作業効率が大幅に上がること間違いなしです。

今回の記事がエンジニアの効率アップに役立ちましたら幸いです。