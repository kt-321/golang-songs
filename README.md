# golang-songs
golang project

## オリジナルで作成したアプリケーション「Your Songs」のバックエンド

プロジェクト3に記載のオリジナルアプリケーション（PHP/Laravel+Nuxt.jsのSPA）のバックエンドをGolangにリプレイスしました。Spotifyの無料アカウントを作ってお持ちであれば、Spotifyで曲を検索して投稿することができます。

【アプリケーションURL】
http://your-songs-laravel.site

【Github】
- https://github.com/kt-321/golang-songs　→　Golang
- https://github.com/kt-321/nuxt-songs-go　→　Nuxt.js

　実際に業務でGolangを使用している方に、コードレビューしてもらいながら実装を進めています。

【主な使用技術】
- Golang
- Nuxt.js
- TypeScript
- AWS
- VPC
- EC2
- Route53
- RDS for MySQL
- S3
- ALB
- ECS
- ECR
- Terraform
- AWS Secrets Manager

【Golangのコード】
- net/httpパッケージでHTTPサーバーの起動
- gorilla/muxを用いてルーティング作成
- ORM用ライブラリGORMを使用
- sql-migrateを用いてマイグレーション
- GolangCI-Lint
- net/http/httptestを用いてテストコード記述
- ルーティングについて必要であるものはJwtMiddlewareでラップ

APIリクエストがあるとJSON形式でフロントにレスポンスを返しています。

【実装済みの主な機能】
- ユーザー登録
- ログイン
- ユーザー情報編集
- 指定したユーザーの情報参照
- SpotifyAPIを用いた曲検索
- 曲の追加
- 曲情報の参照
- 曲の編集
- 曲の削除
- ユーザーフォロー機能
- 曲のお気に入り機能
- Clean Architectureを倣ったディレクトリ構成
- テストコード
- Github Actionsを用いた自動テスト
- Github Actionsを用いて、ECR への image の push, ECS(Fargate)への自動デプロイ

【現在実装中】
- 画像をアップロードしてS3に保存する機能


<a href="https://gyazo.com/ce7bf25c667275e805debd2a415009f6"><img src="https://i.gyazo.com/ce7bf25c667275e805debd2a415009f6.png" alt="Image from Gyazo" width="700"/></a>
