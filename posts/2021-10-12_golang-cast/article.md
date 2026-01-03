---
id: 1173
title: Go言語の型変換(キャスト)方法についてまとめてみた
slug: golang-cast
status: publish
excerpt: こんばんは、しーまんです。 今回は完全に個人的な備忘録です！ 最近Go言語で開発する場面が増えているのですが、いつも型変換(キャスト)に悩まされます。 なるべくキャストはせずに開発するのがよいということは分かっているので \[…\]
categories:
    - 18
    - 60
tags:
    - 59
    - 76
featured_media: 1175
date: 2021-10-12T19:30:00
modified: 2022-03-19T11:58:35
---

こんばんは、しーまんです。

今回は完全に個人的な備忘録です！

最近Go言語で開発する場面が増えているのですが、いつも型変換(キャスト)に悩まされます。

なるべくキャストはせずに開発するのがよいということは分かっているのですが、使用するモジュールに合わせて型変換をする場面がどうしても出てきてしまします。

全てのモジュールがインターフェースでどの型にも対応なんて都合の良い感じにはできていませんからね。

自身が開発しているシステムとモジュール間で微妙に型が違うなんて場合はないでしょうか。例えば同じ整数でもintとint64とか。

そういった場合の型変換の方法を毎回検索していて、流石に疲れたので一通りまとめようと思いました。

同じように型変換について悩んでいた人にとっては参考になるかもしれません。

## 数値型から数値型

まずは数値型から数値型の型変換についてまとめていきます。

大体の言語では数値型に関しては自動で大きい方というか、32ビットと64ビットだったら64ビットに置き換えてくれたり、intやdoubleなど整数と小数だったら小数に合わせてくれたりを自動で判断して行ってくれる言語が多いですが、Golangではそのあたりの融通は利きません。

しっかりと対応した型を指定しないと **コンパイルエラー** となります。

これは型を意識しないで予期せぬ処理が走ることを防ぐ狙いがあるということは分かりますが、少し窮屈に感じる部分でもあります。

普通の言語では文字列型から数値型とその逆を考えるのがほとんどですが、Golangでは数値型同士の型変換も意識しなくてはいけません。ですのでまずはここからまとめていきましょう。

基本的に数値型から数値型の型変換は「 **型(変数)**」でキャストできます。

ただし桁あふれや、桁数あふれに注意しましょう！

### intからint64

```
var i int = 123
i64 := int64(i)
```

### intからint32

```
var i int = 123
i32 := int32(i)
```

### int64からint

```
var i64 int64 = 123
i := int(i64)
```

※ 桁あふれに注意しましょう。

```
if i64 < math.MinInt32 || i64 > math.MaxInt32 {
        return 0
} else {
        return int(i64)
}
```

などと判定を入れてあげると丁寧

### int32からint

```
var i32 int32 = 123
i := int(i32)
```

### intからuint

```
var i int = 123
ui := uint(i)
```

※ 負の値のキャストに注意しましょう。

### uintからint

```
var ui int = 123
i := int(ui)
```

### intからfloat64

```
var i int = 123
f64 := float64(i)
```

### intからfloat32

```
var i int = 123
f32 := float32(i)
```

### float64からint

```
var f64 float = 123.456
i := int(f64)     // 123
```

※ 小数点以下の切り捨てに注意

### float32からint

```
var f32 float = 123.456
i := int(f32)     // 123
```

※ 小数点以下の切り捨てに注意

## 数値型から文字列型

数値型の型変換ができれば後は他言語と同じです。

数値型と文字列型の型変換をまとめていきましょう。

基本的にはstrconvパッケージを使用し、stringから違う型に変換する場合「 **Parse〇〇**」、逆の変換の場合は「 **Format〇〇**」という関数名になります。

### intからstring

```
import "strconv"

var i int = 123
s := strconv.Itoa(i)
```

### int64からstring

```
import "strconv"

var i64 int64 = 123
s := strconv.FormatInt(i64, 10)
```

※ 64ビットより小さいものをまずは64ビットの数値型に変換する必要があります。

### uint64からstring

```
import "strconv"

var ui64 uint64 = 123
s := strconv.FormatUint(ui64, 10)
```

※ 64ビットより小さいものをまずは64ビットの数値型に変換する必要があります。

### float64からstring

```
import "strconv"

var f64 float64 = 123.456
s := strconv.FormatFloat(f64, 'f', -1, 64)
```

### float32からstring

```
import "strconv"

var f32 float32 = 123.456
s := strconv.FormatFloat(f32, 'f', -1, 32)
```

## 文字列型から数値型

### stringからint

```
import "strconv"

var s string = "123"
i, _ := strconv.Atoi(s)
```

### stringからint64

```
import "strconv"

var s string = "123"
i64, _ := strconv.ParseInt(s, 10, 64)
```

### stringからint32

```
import "strconv"

var s string = "123"
i32, _ := strconv.ParseInt(s, 10, 32)
```

### stringからuint64

```
import "strconv"

var s string = "123"
ui64, _ := strconv.ParseUint(s, 10, 64)
```

### stringからuint32

```
import "strconv"

var s string = "123"
ui32, _ := strconv.ParseUint(s, 10, 32)
```

### stringからfloat64

```
import "strconv"

var s string = "123.456"
f64, _ := strconv.ParseFloat(s, 64)
```

### stringからfloat32

```
import "strconv"

var s string = "123.456"
f32, _ := strconv.ParseFloat(s, 32)
```

## 論理値型から文字列型

最後に論理値型と文字列型の変換についてです。

こちらはそこまでパターンが無いと思いますが、一応まとめておきます。

### boolからstring

```
import "strconv"

var b bool = true
s := strconv.FormatBool(b)
```

## 文字列型から論理値型

### stringからbool

```
import "strconv"

var s string = "true"
b, _ := strconv.ParseBool(s)
```

## まとめ

今回はGo言語における型変換(キャスト)の方法についてまとめてみました。

Go言語のキャスト方法は他言語と違い容易に変更できないようにしているのか、とても覚えきれません。

おそらく私自身も実装の度に今回の記事に戻ってきて、参照するような気がしています。

正確に型を意識して実装することで予期せぬエラーを無くすことは可能だとは思いますが、実装工数との比較的にはどっちがいいか判断するのは難しいと思いますね。

ですので、JSのstrictモードとそうじゃないモードみたいなのがGo言語でもあると個人的にはすごい助かるなと思いました。

Goの融通が利かない型の扱いについてチートシート的な役割になれば幸いです。