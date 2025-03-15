# Slack Bot

シンプルな Slack Bot です。以下の機能があります：

- "hello" というメッセージに "world" と返信
- "image" というメッセージに登録された画像を返信

## 必要条件

- Go 1.21 以上
- SQLite3
- Slack App の設定
  - Event Subscriptions の有効化
  - Bot Token Scopes: `chat:write`, `app_mentions:read`, `messages.im:read`, `files:write`
  - Subscribe to bot events: `message.im`

## セットアップ

1. リポジトリをクローン

```bash
git clone https://github.com/n-kurasawa/slack-bot.git
cd slack-bot
```

2. 依存関係のインストール

```bash
go mod download
```

3. 環境変数の設定

```bash
cp .env.example .env
```

`.env` ファイルを編集して、`SLACK_BOT_TOKEN` に適切な値を設定してください。

4. データベースの初期化

```bash
# データベースの初期化のみ
go run cmd/db/main.go

# 画像の登録
go run cmd/db/main.go -image path/to/your/image.jpg
```

## 実行方法

```bash
go run cmd/bot/main.go
```

サーバーは port 3000 で起動します。

## 使い方

Slack でボットとのダイレクトメッセージで以下のコマンドが使えます：

- `hello` - "world" と返信します
- `image` - 登録された最新の画像を返信します

## Slack App の設定

1. [Slack API](https://api.slack.com/apps) でアプリを作成
2. "Event Subscriptions" を有効化
3. Request URL に `https://あなたのドメイン/slack/events` を設定
4. Bot Token Scopes に必要な権限を追加
   - `chat:write` - メッセージの送信用
   - `app_mentions:read` - メンションの読み取り用
   - `messages.im:read` - DM の読み取り用
   - `files:write` - 画像の送信用
5. アプリをワークスペースにインストール

## 開発環境でのテスト

ローカル環境で開発する場合は、[ngrok](https://ngrok.com/) などのツールを使用してローカルサーバーを公開できます：

```bash
ngrok http 3000
```

生成された URL + "/slack/events" を Slack アプリの Request URL として設定してください。

## プロジェクト構造

```
.
├── cmd
│   ├── bot/    - Slack Bot のメインコード
│   └── db/     - データベース操作用コマンド
├── internal
│   └── db/     - データベース操作用パッケージ
├── .env        - 環境変数設定ファイル
├── .gitignore
└── README.md
```

## データベース

画像は SQLite データベースに保存されます。データベースファイル（`images.db`）は自動的に作成され、Git の管理対象外となっています。
