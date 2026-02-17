---
id: 1034
title: 【Golang】HTTPリクエストのレスポンスを構造体定義なしでjsonに吐き出す方法
slug: output-json
status: publish
date: 2021-10-01T19:30:00
modified: 2021-10-01T13:29:14
excerpt: GolangでHTTPレスポンスのJSONを構造体定義なしでファイルに出力する方法を紹介。JSON-to-Goの活用法も解説します。
categories:
    - 18
    - 60
tags:
    - 59
    - 61
    - 62
featured_media: 1035
---

こんばんは、しーまんです。

Golangでデータを扱う際は構造体を定義して扱うのが一般的です。HTTPリクエストのレスポンスデータに対しても構造体を定義し受け取ります。

その際の構造体を作成するwebサービスとして [JSON-to-Go](https://mholt.github.io/json-to-go/) はよくお世話になっております。

 [![](https://s.wordpress.com/mshots/v1/https%3A%2F%2Fmholt.github.io%2Fjson-to-go%2F?w=320&h=180)\
\
JSON-to-Go: Convert JSON to Go instantly\
\
![](https://www.google.com/s2/favicons?domain=https://mholt.github.io/json-to-go/)\
\
mholt.github.io](https://mholt.github.io/json-to-go/ "JSON-to-Go: Convert JSON to Go instantly")

しかし、外部サービス連携などでapiをリクエストした際に構造体を定義せずにjson出力したい場合などがあります。そのような場合どのように実装すればよいのか検証したので、その方法を解説致します。

## 構造体を定義しない理由

構造体を定義しない状況っていうのはかなり特殊な状態だと思いますので、私の事例を元に解説致します。

最近クラウドサービスの普及により、REST api で情報を取得できるSaaSが増えてきています。

そんな中でapiで取得したデータをjsonに吐き出し、そのデータをBigQueryへ投入したい要件がありました。

イメージとしては下図のような感じですね。

[![saas api](https://shiimanblog.com/wp-content/uploads/2021/10/saas_api-800x294.png)](https://shiimanblog.com/wp-content/uploads/2021/10/saas_api.png)

### 2つの構造データ作成を省略

上記の実装の際、BigQueryにデータを投入する時にAutodetectという設定があり、これはテーブルのカラム定義を自動で推論してくれる機能です。Autodetect機能を使用しない場合ではBigQuery側のテーブルの構造データも作成する必要があります。

つまり、Golang側では

APIの取得に構造体を作成し、BigQueryのデータ投入にテーブルの構造データ作成

という2つほど、データの構造を定義しないといけません。

それはとても面倒で工数もかかりますし、APIの数も大量にありますので、構造の定義を省きたいというわけです。

BigQueryに関しては上記で取り上げたAutodetect機能を使用すれば解決しますが、一応制限事項もありますので、使用する際はマニュアルを参照ください。

[GCP Bigquery JSONデータ読み込みマニュアル](https://cloud.google.com/bigquery/docs/loading-data-cloud-storage-json?hl=ja#limitations)

### SaaSのAPI取得で構造体定義したくない理由

またSaaSのAPIに関して構造体定義したくない理由が3点ほどあります。

1. SaaS側で提供されているマニュアル、テストレスポンスに誤りがある場合が存在する

実際に本番リクエストを送って確認するまで安心できません。

2. レスポンスが頻繁に更新される

レスポンスの更新が行われるともちろん構造体の更新も必要になります。

3. 必須値以外が漏れる可能性がある

返って来ない場合があるレスポンスの値は構造体定義時に漏れる可能性があります。

上記のような理由からSaaSのAPI取得は可能であれば構造体を定義しないで、動的に扱いたいと思います。

## HTTPレスポンスを動的に扱う方法

では実際にHTTPレスポンスを動的に扱う方法をみていきましょう。

まずは構造体を定義する通常パターン

エラーを全て無視して簡易的に書くとこのようなコードになります。

```
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func main() {
	url := "https://xxxx"

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	client := new(http.Client)
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// JSONを構造体にエンコード
	var response Response
	json.Unmarshal(body, &response)

	file, _ := os.Create("response.json")
	defer file.Close()

	json.NewEncoder(file).Encode(response)
}

type Response struct {
	Column []struct {
		ID   int    `json:"ID"`
		Name string `json:"NAME"`
	} `json:"Column"`
}

```

上記コードで json.Unmarshal のところで、構造体を渡しています。

しかし、json.Unmarshalの第2引数には \*interface{} を渡すことが可能です。また json.EncoderのEncodeも interface{} を渡すことが可能です。

こちらを利用して動的にレスポンスを扱うことが可能になります。

実際に動的にレスポンスを取得してJSONを出力するコードをみていきましょう。

```
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func main() {
	url := "https://xxxx"

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	client := new(http.Client)
	resp, _ := client.Do(req)

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// JSONを構造体にエンコード
	var response interface{}
	json.Unmarshal(body, &response)

	file, _ := os.Create("response.json")
	defer file.Close()

	json.NewEncoder(file).Encode(response)
}
```

変更点は下記ですね。

var response Response が var response interface{} に変わっています。

そうです、interfaceをうまく使用してあげると構造体を定義しなくてもレスポンスをjsonに吐き出すことが可能になります。

ただし、今回はjsonを吐き出すことが目的でしたので、こちらの方法を取りましたが、実際に値を使うことを考えるととても大変です。

下記のようにキャストの嵐で値をとりだすことになりコード量が半端じゃなく増えます。

```
response.(map[string]interface{})["Column"].([]interface{})[0].(map[string]interface{})["NAME"].(string)
```

また、そもそも使用する値が決まっているのだったら構造体も定義できますよね。

ですので、今回の要件のように特殊なケースを除いては、ちゃんと構造体を定義した方が処理を楽に書くことができます。

## まとめ

今回はHTTPリクエストのレスポンスを動的にJSON化する方法を解説しました。

なかなかニッチなケースではありますが、元々PHPerである私としては構造体定義なしでもっと気軽にレスポンスを扱いたいって思いもあります。まぁ余計なreflectionなどの処理がなくなるので処理自体は高速なんでしょうね。

ということで、このようなパターンで困っている方がどれだけいるかは分かりませんが、同じように困っている方の参考になれば幸いです。