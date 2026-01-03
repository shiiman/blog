---
id: 1287
title: 【Golang】OAuth2.0を実装する方法/Google認証
slug: golang-oauth2
status: publish
excerpt: こんにちは、しーまんです。 今回はいつもと違い、サムザップ Advent Calendar 2021の12/3の記事です。毎年恒例のアドベントカレンダー用の記事になります。 とはいえ、特に変わったところはなく普通に技術記 \[…\]
categories:
    - 18
    - 60
tags:
    - 85
    - 86
featured_media: 1290
date: 2021-12-03T12:00:00
modified: 2022-11-04T17:41:25
---

こんにちは、しーまんです。

今回はいつもと違い、 [サムザップ Advent Calendar 2021](https://qiita.com/advent-calendar/2021/sumzap) の12/3の記事です。

毎年恒例のアドベントカレンダー用の記事になります。

とはいえ、特に変わったところはなく普通に技術記事を上げていきたいと思います。

今回はGolangでOAuth2.0を実装する方法をGoogle認証を例にして紹介していきたいと思います。

それでは早速やっていきましょう。

## OAuth2.0とは

まずはOAuth2.0をみていきましょう。

### OAuth2.0の概要

通常認証処理というと **ID** と **Password** のようなクレデンシャル情報をやりとりして認証するログイン認証を思い浮かぶと思います。

ここで問題になるのはクレデンシャル情報を渡してしまうと、全ての情報にアクセスできてしまうので情報漏洩が問題になります。

例えば皆さんGoogleのサービスを使いますよね。Googleドキュメントやスプレッドシートなどです。

スプレッドシートの情報を何かしらのWebサービスに渡して何かしらの処理を行う場合、ご自身の **ID** と **Password** を渡してしまったらどうでしょうか？

するとスプレッドシートの情報だけ渡したかったのに、そのWebサービスではご自身のGmailだったり、Google Photosだったり、その他機密情報を持ったドキュメントだったりアクセスし放題になってしまいます。

そこで限られた情報だけアクセスを許可する必要があります。

上記の例では特定にスプレッドシートの情報だけ渡せればいいですよね。

そんな時に使用するのがOAuth認証です。

こちらを使用すれば、ユーザが許可したものだけにアクセスを制限させることが可能になります。

### OAuth2.0のシーケンス

OAuth2.0には [RFC 6749](https://tools.ietf.org/html/rfc6749) (The OAuth 2.0 Authorization Framework) で定義されている4つの認可フローがあります。

ただこちらを全て説明し出すと途方もない労力になってしまうので、詳しい説明は今回省かせていただきます。 [@TakahikoKawasaki](https://qiita.com/TakahikoKawasaki) がそのあたりをとても分かりやすくまとめた記事を上げてくださっているので、詳細を知りたい方はぜひご覧になってみてください。

[一番分かりやすい OAuth の説明](https://qiita.com/TakahikoKawasaki/items/e37caf50776e00e733be)

[OAuth 2.0 全フローの図解と動画](https://qiita.com/TakahikoKawasaki/items/200951e5b5929f840a1f)

今回は実際に認証アプリを作成する際によく使用される、RFC 6749, [4.1. Authorization Code Grant](https://tools.ietf.org/html/rfc6749#section-4.1) で定義されている認可フローついてを紹介します。

[![OAuth2.0 - シーケンス図](https://shiimanblog.com/wp-content/uploads/2021/11/1bd64ca7670b2417de6b176be337397a-1-800x762.png)](https://shiimanblog.com/wp-content/uploads/2021/11/1bd64ca7670b2417de6b176be337397a-1.png)

図にあるように登場人物は「 **ブラウザ**」「 **クライアントアプリ**」「 **サーバアプリ**」「 **認可サーバ**」の4人です。そして今回Golangの実装を紹介するのは「サーバアプリ」の部分になります。アプリを実際に作成する際は「クライアントアプリ」の実装も必要です。

「認可サーバ」に関しては各サービス側で用意してくださっています。例えばGoogleの認可サーバだったり、Githubの認可サーバがこれにあたります。

## golang/oauth2

上記で紹介したシーケンスをもとに実装をしていきましょう！

Golangでは [golang/oauth2](https://github.com/golang/oauth2) というパッケージがありますので、こちらを使用すると実装が簡単になります。

とはいえ認証の実装にはかなり骨が折れる作業になるので細かくステップ毎に実装内容をみていきましょう。

ステップ1\. OAuth2.0用のconfを作成

ステップ2\. RedirectURLの作成

ステップ3\. ユーザの承認&リダイレクト

ステップ4\. Auth CodeからToken取得

ステップ5\. Tokenの使用

ステップ6\. Tokenリフレッシュ

### ステップ1\. OAuth2.0用のconfを作成

まずは認可サーバの情報などを設定する必要があります。

こちらは各サービス毎に取得方法が異なりますのでそれぞれのサービスのドキュメントを参照ください。

Googleの場合は [こちら](https://developers.google.com/identity/protocols/oauth2) になります。

手順どおり進めていくとクレデンシャルJSONが発行されます。

こちらを使用してconfファイルを作成していきます。

JSONの中身は「 **client\_id**」「 **client\_secret**」「 **auth\_uri**」「 **token\_uri**」「 **redirect\_uri**」などが定義されています。OAuth2.0のconfファイルはこれらのパラメータがセットされた構造体ですので、JSONを使用しなくても直接セットすることで動作可能です。

Google認証以外でconfファイルを作成する場合、上記項目を直接指定することが多いです。

```
import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/sheets/v4"
)

func NewGoogleAuthConf() *oauth2.Config {
        // 実際にはSecretManagerなどに保存して、そこから取得する.
	credentialsJSON := `jsonの中身`

        // 第2引数に認証を求めるスコープを設定します.
        // 今回はスプレッドシートのリード権限スコープを指定.
	config, err := google.ConfigFromJSON(credentialsJSON, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config
}
```

### ステップ2\. RedirectURLの作成

次はリダイレクトURLの作成です。

`AuthCodeURL` メソッドを使用することでリダイレクトURLの作成が可能です。こちらをクライアントに渡し、リダイレクトすることでユーザに認証画面を表示させます。

CSRF対策のため検証用のstate文字列を付与しましょう

```
func Auth() error {
	conf := NewGoogleAuthConf()
	state := `CSRF攻撃を防ぐためにstateパラメータをつける.コールバック後のトークン取得時に検証する.`

	// stateをsessionなどに保存.

	// リダイレクトURL作成.
	redirectURL := conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce),

	// redirectURLをクライアントに返す.
}
```

googleの認可サーバにリダイレクトするURLは実際には以下のような感じで作成されます。

```
https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id={client_id}prompt=consent&redirect_uri={redirect_uri}&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fspreadsheets.readonly&state={state}
```

### ステップ3\. ユーザの承認&リダイレクト

クライアントから認証画面にリダイレクトさせるとユーザにスプレッドシートへのアクセスを許可するか確認する画面が開きます。許可がされたら、アプリ側にリダイレクトされるのでクライアントアプリ側でそのアクセスポイントを作成しておきましょう。

[![Google認可](https://shiimanblog.com/wp-content/uploads/2021/11/370389726c4cf5ee3904949f4f72edc7.png)](https://shiimanblog.com/wp-content/uploads/2021/11/370389726c4cf5ee3904949f4f72edc7.png)

認可サーバから返却されるcallbackは下記のようなURLになります。

```
https://{call_back_uri}?state={state}&code={code}&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fspreadsheets.readonly
```

アプリ側のリダイレクトURLにcodeとstateがクエリパラメータとして付与されていますので、こちらをサーバアプリに送ります。

### ステップ4\. Auth CodeからToken取得

クライアントから送られてきた **code** を使用してサーバで **Token** を取得します。

`Exchange` メソッドを使用することでToken情報が取得可能です。

取得したトークン情報は「 **アクセストークン**」「 **リフレッシュトークン**」「 **有効期限**」の3要素あります。こちらをDBなどに保存しておきましょう。

```
func Link() error {

	// クライアントからcode, stateを取得.
	code := `code`
	state := `state`

	// stateが正しいか検証.

	conf := helper.NewGoogleAuthConf()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return err
	}

	// token.AccessToken, token.RefreshToken, token.ExpiryをDBに保存.
}
```

データベースへの保存は下記のようなテーブルにそれぞれ保存しておけばよいでしょう。

**google\_authテーブル**

[![認可情報保存テーブル](https://shiimanblog.com/wp-content/uploads/2021/11/table.png)](https://shiimanblog.com/wp-content/uploads/2021/11/table.png)

ここまでのステップ4までで認証処理は終わりです。

保存したTokenを使用することでスコープに則った情報にアクセスすることが可能です。

### ステップ5\. Tokenの使用

では実際にトークンを使用してスプレッドシートにアクセスしてみましょう。

DBに保存したアクセストークンを取得し、スプレッドシートにアクセスします。

```
func NewClient(ctx context.Context) *http.Client {
	conf := NewGoogleAuthConf()

	// DBに保存したトークン情報取得.
	accessToken := `アクセストークン`
	refreshToken := `リフレッシュトークン`
	expiry := `有効期限`

	token := &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    "bearer",
		RefreshToken: refreshToken,
		Expiry:              expiry,
	}

	return conf.Client(ctx, token)
}

func SpreadsheetSheetGet() error {

	// クライアントから取得したいスプレッドシートIDを受け取る.
	spreadsheetID := `スプレッドシートID`
	readRange : = `読み込み範囲`

	// クライアント取得.
	ctx := context.Background()
	client := NewClient(ctx)

	// シート情報取得.
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	// https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets.values/get .
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return err
	}

	// クライアントにスプレッドシート情報を渡す.
}
```

これでトークンを使用してスプレッドシートにアクセスできました。

しかし、このままの実装ではバグがあります。

それは折角DBにトークンを保存しているのに、有効期限が切れた場合スプレッドシートにアクセスできなくなってしまうことです。アクセスのたびに都度認証を入れる(トークンを保存しない)ならよいですが、それでは使い勝手が悪いですよね。

ということで最後はトークンの有効期限が切れた場合にリフレッシュトークンを使用し、トークン情報を更新する方法を紹介します。

### ステップ6\. Tokenリフレッシュ

リフレッシュトークンを使用するには `TokenURL` に `grant_type=refresh_token` オプションをつけてリクエストすればokです。これは `oauth2.NewClient` を使用すれば内部で自動的に有効期限を判定して、トークンを更新してくれます。 その際にDBに保存されているトークン情報も更新させます。

そこで先程定義した `NewClient` を修正していきます。

```
func NewClient(ctx context.Context) *http.Client {
	conf := NewGoogleAuthConf()

	// DBに保存したトークン情報取得.
	accessToken := `アクセストークン`
	refreshToken := `リフレッシュトークン`
	expiry := `有効期限`

	token := &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    "bearer",
		RefreshToken: refreshToken,
		Expiry:              expiry,
	}

	// token取得.
	tokenSource := conf.TokenSource(ctx, token)

	// token更新.
	mySrc := &MyTokenSource{
		src:               tokenSource,
		f:                   TokenRefresh,
		dbID:            `更新するDBのレコードID`,
	}
	reuseSrc := oauth2.ReuseTokenSource(token, mySrc)
	client := oauth2.NewClient(ctx, reuseSrc)
	return client
}

type MyFunc func(*oauth2.Token, uint) error

func TokenRefresh(t *oauth2.Token, dbID uint) error {
	// 更新されたtoken情報と対象のDBレコードIDをもとにDBのToken情報を更新.

	return nil
}

type MyTokenSource struct {
	src               oauth2.TokenSource
	f                 MyFunc
	dbID              uint
}

func (s *MyTokenSource) Token() (*oauth2.Token, error) {
	t, err := s.src.Token()
	if err != nil {
		return nil, err
	}
	if err = s.f(t, s.dbID); err != nil {
		return t, err
	}
	return t, nil
}
```

上記の様に修正することでトークンの有効期限のチェックをパッケージにまかせつつ、トークンが更新された際にDBに保存されたトークン情報も更新してくれるようになりました。これで1度取得したトークンを使い回すことができるようになり利便性が向上しますね。

## まとめ

今回はGolangでOAuth2.0を実装する方法を紹介しました。

今回はGoogle認証を作りましたが、最近のwebサービスは結構OAuth2.0に対応していますので、いろいろなサービスに流用することが可能です。

認証周りの実装に困っている方はぜひ今回の記事を参考に実装してみてください。

今回はアドベントカレンダーの1記事ということですので、明日は [@RyotoKitajima](https://qiita.com/RyotoKitajima) さんの記事になります。年末恒例のイベントですので、他の方の記事も読んでこの機会にいっぱいインプットしましょう！