# golang-songs

## 概要
- Golang+Nuxt.jsのSPAであるオリジナルアプリケーション「Your Songs」のGolangのコード
- 元々PHP/Laravel+Nuxt.jsのSPAで作っていたところを、バックエンドをGolangにリプレイスしました。実際に業務でGolangを使用している方に、コードレビューしてもらいながら実装を進めてきました。
- Spotifyの無料アカウントを作ってお持ちであれば、Spotifyで曲を検索して投稿することができます。

## アプリケーションURL
http://your-songs-laravel.site

## フロント側のコード（Nuxt.js）
- https://github.com/kt-321/nuxt-songs-go


## アプリケーション全体での主な使用技術
- Golang 1.14
- Nuxt.js 2.11
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
- ElastiCache (Redis)
- Terraform
- AWS Secrets Manager

## Golangのコード
- net/httpパッケージでHTTPサーバーの起動
- gorilla/muxを用いてルーティング作成
- ORM用ライブラリGORMを使用
- sql-migrateを用いてマイグレーション
- GolangCI-Lint
- net/http/httptestを用いてテストコード記述
- ルーティングについて必要であるものはJwtMiddlewareでラップ
- Redigoを用いてRedisの使用

APIリクエストがあるとJSON形式でフロントにレスポンスを返しています。

## 実装済みの主な機能
- ユーザー登録・ログイン
- ユーザー情報の取得
- ユーザー情報編集
- SpotifyAPIを用いた曲検索
- 曲の追加
- 曲情報の取得
- 曲の編集
- 曲の削除
- 曲をお気に入りする機能
- Redis(ElastiCache)の利用（曲の取得・追加・更新・削除）
- Clean Architectureを倣ったディレクトリ構成
- テストコード
- Github Actionsを用いた自動テスト
- Github Actionsを用いて、ECR への image の push, ECS(Fargate)への自動デプロイ

## 現在実装中
- 画像をアップロードしてS3に保存する機能
- CloudFrontの導入
- Lambda・API Gatewayの導入

## アプリケーションのTOP画面
<a href="https://gyazo.com/6bc2a6b38c9420c9e41b829ca3a9eba1"><img src="https://i.gyazo.com/6bc2a6b38c9420c9e41b829ca3a9eba1.jpg" alt="Image from Gyazo" width="1511"/></a>

## インフラ構成図
<a href="https://gyazo.com/68eb106befe391785cdfb0e826a61fb0"><img src="https://i.gyazo.com/68eb106befe391785cdfb0e826a61fb0.png" alt="Image from Gyazo" width="700"/></a>
