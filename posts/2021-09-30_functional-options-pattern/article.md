---
id: 858
title: 【Golang】関数でデフォルト引数を実現させる方法 | Functional Options Pattern
slug: functional-options-pattern
status: publish
excerpt: こんばんは、しーまんです。 プログラムで関数を作成する時に引数を設定すると思いますが、Golangではデフォルト引数(オプションパラメータ)の設定をすることができません。関数で引数を設定した場合には呼び出し元で必ず引数を \[…\]
categories:
    - 18
    - 60
tags:
    - 59
featured_media: 1027
date: 2021-09-30T19:30:00
modified: 2021-10-30T19:01:45
---

こんばんは、しーまんです。

プログラムで関数を作成する時に引数を設定すると思いますが、Golangではデフォルト引数(オプションパラメータ)の設定をすることができません。関数で引数を設定した場合には呼び出し元で必ず引数を渡してあげないとエラーになります。

元々PHPerの私は関数の設計でどうしてもデフォルト引数の関数設計をしてしまいたくなります。

きっと同じように思った方はいるはず。。

そこでGolangでもどうにかデフォルト引数の関数を実装できないかと調べたところ「**Functional Options Pattern**」というデザインパターンで実装することが可能だということが分かりました。

今回は「 **Functional Options Pattern**」を使用してデフォルト引数のある関数の実現方法を紹介したいと思います。

## 実現したいこと

まず今回実現したい内容を確認しておきます。

PHPを使用すると下記のようなコードが実装可能です。

```
function echoJoinStr($str1 = "" , $str2 = "") {
    echo $str1 . $str2;
}

echoJoinStr();                               // ""
echoJoinStr("hello");                   // "hello"
echoJoinStr("hello", "world")     // "helloworld"
```

上記のように引数にデフォルトで代入される値を予め設定しておくことで、呼び出し側は引数を渡しても渡さなくても問題ないように実装することが可能です。

上記のコードでデフォルト引数を設定する意味は殆どありませんが、デフォルト引数のある関数を説明するために簡単な例をあげています。

実際には殆どの場合デフォルト値で問題ないけど変更することも可能なように実装する必要があるものに設定します。

**例. pagerの実装でlimitとoffsetを引数にする場合**

デフォルトlimitで100件と設定しておき、

引数で渡された場合1000件などに変更出来る関数を用意する

Golangでは基本的にはこのデフォルト引数を設定することはできません。

ですので、これを実現するためにはどうすればいいのかをみていきましょう。

## Functional Options Pattern

Golangでデフォルト引数ありの関数を設定したい場合、Dave Cheney氏の記事で取り上げている「 **Functional Options Pattern**」というものを使用すると実現が可能なようです。

 [![](https://s.wordpress.com/mshots/v1/https%3A%2F%2Fdave.cheney.net%2F2014%2F10%2F17%2Ffunctional-options-for-friendly-apis?w=320&h=180)\
\
Functional options for friendly APIs \| Dave Cheney\
\
![](https://www.google.com/s2/favicons?domain=https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)\
\
dave.cheney.net](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis "Functional options for friendly APIs | Dave Cheney")

デザインパターンとかは人によって細かい認識に差異が出て、ああだこうだと議論になる部分なので私個人的には苦手な話題ですが、大切なのは上記であげた問題を解決することです。

ですので、デザインパターンに関して言及することは行いませんが、興味がある方は調べてみてください。

## Golangで実現

では実際にGolangで上記のPHPソースと同じ事ができる、デフォルト引数ありの関数を実現していきたいと思います。コード例は以下になります。

```
type Options struct {
	Str1 string
	Str2 string
}

type Option func(*Options)

func Str1(str1 string) Option {
	return func(opts *Options) {
		opts.Str1 = str1
	}
}

func Str2(str2 string) Option {
	return func(opts *Options) {
		opts.Str2 = str2
	}
}

func echoJoinStr(options ...Option) {
	// デフォルト値設定.
        opts := &Options{
		Str1: "",
		Str2: "",
	}

	for _, option := range options {
		option(opts)
	}

	fmt.Print(opts.Str1 + opts.Str2)
}

func main() {
    echoJoinStr()                                                   // ""
    echoJoinStr(Str1("hello"))                            // "hello"
    echoJoinStr(Str1("hello"), Str2("world"))  // "helloworld"
}
```

デフォルト引数を追加するだけなのにコード量が多いのが気になりますが、個人的にはとてもスッキリした解決ができたかなと思います。

基本的にGolangは余計な機能を削ぎ落として作られているので、デフォルト引数はあってもなくても実際には問題にならないことが多いです。なので機能として存在しないのだと思います。

少なくてもデフォルト引数がないと実装できないということはないはずなので、今回使用した実装もなるべく多用しないようにして、お作法通り実装するのが良さそうです。

## まとめ

今回はGolangでデフォルト引数がある関数の作成方法を紹介しました。

PHPなどと違いデフォルト引数を使用しようとした場合「 **Functional Options Pattern**」を使いひと工夫することで実現ができました。

他の言語からGolangに触れると、他の言語のこんな機能がないのかなと調査することは結構あると思います。今回の「デフォルト引数がある関数」の実装もその一つです。

今回の記事が同じように実装に困っている方の参考に少しでもなれば幸いです。