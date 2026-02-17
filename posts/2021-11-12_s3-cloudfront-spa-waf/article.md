---
id: 1256
title: 【AWS】S3+CloudFrontでSPA対応する方法とWAFのIP制限干渉を解決する方法！
slug: s3-cloudfront-spa-waf
status: publish
date: 2021-11-12T19:30:00
modified: 2021-11-11T17:24:34
excerpt: AWSでS3+CloudFrontを使ったSPA構築時のルーティング問題と、WAFのIP制限との干渉を解決する方法を解説します。
categories:
    - 19
    - 18
tags:
    - 79
    - 81
    - 82
    - 83
    - 84
featured_media: 1262
---

こんばんは、しーまんです。

今回はAWSでSPAのWebサイトを構築する際にハマった点を解説していきたいと思います。

サーバを立てずにサッと簡単なWebサイトを構築する際はよく用いる方法ですので、参考にしていただければと思います。

## AWSで簡単なSPAサイトの構築

まずSPAはご存知でしょうか。

SPAとは Single Page Application の略で文字通り単一のページで構成されるWebアプリケーションです。

通常だとページの遷移のたびに新しいページを構成する要素をサーバからダウンロードして表示しますが、SPAの場合1つのページをあたかも遷移したページのように見せかけページを構成する要素だけを入れ替えて表示します。

こうすることによってページの切り替えが早く、ユーザビリティが向上します。

* * *

こちらのSPAをAWSで構築する場合、よくS3+CloudFontの構成で実現したりします。

[![s3+cloudfront構成図](https://shiimanblog.com/wp-content/uploads/2021/11/1266bf9daf9ea60649edb5ff772bbeec-800x290.png)](https://shiimanblog.com/wp-content/uploads/2021/11/1266bf9daf9ea60649edb5ff772bbeec.png)

この構成にする理由は簡単で、 **S3** を使用することで自身で専用のサーバを用意する必要がなくなります。また **CloudFront** を使用することでグローバルにCacheを持つことができ、軽快にアクセスを捌くことが可能だからです。

間に **WAF** がありますが、これはCloudFrontへのアクセスに対して何かしらの制限をかけたい場合に使用します。よくあるのが開発中の環境に対して、特定のIPアドレスからのみにアクセスを制限するなどの要件です。

## CloudFrontでErrorPages設定

### SPAリロードのエラー

通常のWebサーバだと上記のような設計で問題ないのですが、SPAの場合は先程も少し触れましたが、ページ遷移する際は仮想的なページ遷移であって実際にそのURLのhtmlが存在するわけではありません。

つまりWebブラウザでリロードを走らせた場合、デフォルトのルーティング設定だと403エラーが返ってきてしまいます。

[![403エラー](https://shiimanblog.com/wp-content/uploads/2021/11/403.png)](https://shiimanblog.com/wp-content/uploads/2021/11/403.png)

こんな感じのエラーをみたことありませんか。よくあるエラーですよね。

これはSPAサイトでは例えば `example.com/login` というページに遷移した際に `/login` という部分は実際には存在していないページで、仮想的にuriを作成しているために起こる現象です。

### SPAリロードエラー対応

これを解決する方法は **CloudFront** でErrorPagesを設定する方法がよく取り上げられています。

まずCloudFrontのページから対象のディストリビューションを選択しエラーページタブを選択します。

そのページから「カスタムエラーレスポンスを作成」を選択します。

[![CloudFront - エラーページ1](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting1-800x196.png)](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting1.png)

そしてレスポンスコード403, 404に対して200に変更しつつ「/index.html」を返すように設定します。

[![CloudFront - エラーページ2](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting2-800x687.png)](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting2.png)

[![CloudFront - エラーページ3](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting3-1-800x697.png)](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting3-1.png)

最終的に下記の画像のような設定ができていればOKです。

これでリロードした際はシングルページのindex.htmlにアクセスすることができ、エラーが表示されてしまう問題は解決できます。

[![CloudFront - エラーページ4](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting4-800x131.png)](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting4.png)

## WAFのIP制限時の403をどう扱うか

### 403でハマったポイント

多くのサイトだと上記の設定で説明が終わっているものがほとんどです。

しかし今回は上記環境でWAFによるIP制限をかけた際に問題が発生しました。

上記設定は403, 404のレスポンスを200にしてレスポンスページを変更する対応でした。

しかしこれによってWAFによるIP制限のデフォルトレスポンスである403も200に変更され、/index.htmlにアクセスできてしまう事象が起きました。

つまりIP制限したはずなのでindex.htmlだけは表示されてしまうという問題が発生したのです。

これではIP制限の意味が無いのでどうにかできないか思考錯誤しました。

最初は「 **Lamda@Edge**」や「 **CloudFront Functions**」を使用してリクエストuriを書き換えるなどを試みてみたのですが、結局うまくいかず。

半日ほどいろいろ書き換えを考えていたのですが、そもそも「IP制限の403」と「ファイルが無かった際に返される403」の見分けがつかないとこが原因だよなと思い、じゃあWAFによるIP制限のレスポンスを変えちゃえばいいじゃんという結論にたどり着きました。

### WAFによるIP制限のレスポンスをカスタマイズ

ということでWAFによるIP制限時のレスポンスを403ではなく401に変更します。

そうすれば200が返ることなくちゃんとIP制限をエラーとして返すことが可能になります。

まずはエラー時にクライアントに返却するカスタムレスポンスボディを作成します。

WAFのページから対象のWebACLを選択し、「 **Custom respons bodies**」タブを開き「 **Create custom response body**」を選択します。

カスタムレスポンスボディに関しては作成しなくても問題ありません。また使用する場合、ご自身でお好きなようにカスタマイズしてご利用ください。

[![WAF- 設定1](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting5-800x220.png)](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting5.png)

レスポンスボディを設定し、名前をつけて保存します。

[![WAF- 設定2](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting6-800x682.png)](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting6.png)

次に「 **Rules**」タブを開き「 **Default web ACL action for requests that don’t match any rules**」を設定していきます。これは指定したIPをAllowにしてそれ以外のアクセスをBlockする時の設定です。いわゆるホワイトリスト方式の場合ですね。

ブラックリスト方式の場合は制限するIPルール側にカスタムレスポンスを設定します。

ページ右下にある「 **Edit**」を選択しましょう。

[![WAF- 設定3](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting7-800x345.png)](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting7.png)

そしてResponse codeを401に設定し、先程上記で設定したレスポンスボディを設定します。

ここでResponse codeを401に設定することが今回の一番の目的です。この設定をすることで、CloudFront側のエラーページ対応からIP制限の処理を除くことができます。

[![WAF- 設定4](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting8-800x706.png)](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting8.png)

保存が正常に完了すると下記のように表示されます。

[![WAF- 設定5](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting9-800x120.png)](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting9.png)

実際に許可していないIPでブラウザからアクセスした場合は下記のような表示がされることが確認できます。

[![WAF- 設定6](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting10-1.png)](https://shiimanblog.com/wp-content/uploads/2021/11/error_page_setting10-1.png)

## まとめ

今回はS3+CloudFrontでSPAのWebサイトを構築する方法と、リロードエラーの対応を解説しました。そのあとWAFによるIP制限の問題についても対応方法を解説しています。

最初はなかなか良い方法が思いつかなく、見当違いな方法を検証していましたが、WAFのカスタムレスポンスを設定するというシンプルな方法で問題解決することができ、スッキリしました。

AWSの記事は気づくと初めてだったですね。普段から触っている分当たり前の領域になってきてしまっていてあまり記事出せていなかったですが、今後はCloudに関しての記事もどんどん出していこうと思います。

以上、どなたかの参考になれば幸いです。