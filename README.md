# golang-songs

## 概要
- Golang+Nuxt.jsのSPAであるオリジナルアプリケーション「Your Songs」のGolang(バックエンド)のコード
- 元はPHP/Laravel+Nuxt.jsのSPAで作っていたところを、バックエンドをGolangにリプレイスしました。実際に業務でGolangを使用している方に、コードレビューしてもらいながら実装を進めてきました。
- Spotifyの無料アカウントを作ってお持ちであれば、Spotifyで曲を検索して投稿することができます。

## アプリケーションURL
http://your-songs-laravel.site

## フロント側のコード（Nuxt.js）
https://github.com/kt-321/nuxt-songs-go


## アプリケーション全体での主な使用技術
- Golang 1.14
- Nuxt.js 2.11
- TypeScript 3.9
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
- パッケージ管理にGOMODULE使用
- testing パッケージを用いてテストコード記述
- go-jwt-middlewareパッケージを用いてJWT認証の実装
- Redigoを用いてRedisの使用
- HTTPサーバーのgraceful shutdown
- encoding/jsonを用いてjsonのエンコード/デコード
- bcryptを用いてパスワードをハッシュ化
- GolangCI-Lintの使用
- デバッグにdelveを使用

構造体にjsonタグを付与することで、APIリクエストがあるとJSON形式でフロントにレスポンスを返しています。

## 実装済みの主な機能
- ユーザー登録・ログイン
- ユーザー情報の取得
- ユーザー情報編集
- SpotifyAPIを用いた曲検索
- 曲の追加
- 曲情報の取得
- 曲の編集
- 曲の削除
- 曲をお気に入り登録する機能
- Redis(ElastiCache)の利用（曲の取得・追加・更新・削除）
- Clean Architectureを倣ったディレクトリ構成
- テストコード
- Github Actionsを用いた自動テスト
- Github Actionsを用いて、ECR へのimageの自動push, ECS(Fargate)でのコンテナ作成

## 現在実装中
- 多層キャッシュ構造
- 画像をアップロードしてS3に保存する機能
- CloudFrontの導入
- Lambda・API Gatewayの導入

## アプリケーションのTOP画面
<a href="https://gyazo.com/6bc2a6b38c9420c9e41b829ca3a9eba1"><img src="https://i.gyazo.com/6bc2a6b38c9420c9e41b829ca3a9eba1.jpg" alt="Image from Gyazo" width="1511"/></a>

## インフラ構成図
<a href="https://gyazo.com/68eb106befe391785cdfb0e826a61fb0"><img src="https://i.gyazo.com/68eb106befe391785cdfb0e826a61fb0.png" alt="Image from Gyazo" width="700"/></a>
