# Slack Bot

シンプルな Slack Bot です。以下の機能があります：

- "hello" というメッセージに "world" と返信
- "image" コマンドで登録された画像を返信
  - `image` - ランダムな画像を返信
  - `image NAME` - 指定した名前の画像を返信
- "updateImage" コマンドで画像を登録
  - 使用方法: `updateImage NAME URL`

## 必要条件

- Go 1.21 以上
- SQLite3
- Slack App の設定
  - Event Subscriptions の有効化
  - Bot Token Scopes: `chat:write`, `channels:history`
  - Subscribe to bot events: `message.channels`

## セットアップ

1. リポジトリをクローン

```bash
git clone https://github.com/n-kurasawa/slack-bot.git
cd slack-bot
```

2. 依存関係のインストール

```bash
go mod tidy
```

3. 環境変数の設定

```bash
cp .env.example .env
```

`.env` ファイルを編集して、以下の値を設定してください：

- `SLACK_BOT_TOKEN` - Slack Bot のトークン
- `SLACK_SIGNING_SECRET` - Slack App の Signing Secret

4. データベースの初期化

```bash
go run cmd/db/main.go
```

## 実行方法

```bash
go run cmd/bot/main.go
```

サーバーは port 3000 で起動します。

## 使い方

Slack でボットとのダイレクトメッセージで以下のコマンドが使えます：

- `hello` - "world" と返信します
- `image` - ランダムな画像を返信します
- `image NAME` - 指定した名前の画像を返信します
- `updateImage NAME URL` - 新しい画像を登録します

## Slack App の設定

1. [Slack API](https://api.slack.com/apps) でアプリを作成
2. "Event Subscriptions" を有効化
3. Request URL に `https://あなたのドメイン/slack/events` を設定
4. Bot Token Scopes に必要な権限を追加
   - `chat:write` - メッセージの送信用
   - `app_mentions:read` - メンションの読み取り用
   - `messages.im:read` - DM の読み取り用
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
│   ├── bot/    - Bot の主要ロジック
│   │   ├── handler.go        - HTTPハンドラ
│   │   ├── messageevent.go   - メッセージイベント処理
│   │   └── slack.go          - Slackクライアントインターフェース
│   └── db/     - データベース操作
│       └── store.go          - SQLiteデータベース操作
├── .env        - 環境変数設定ファイル
├── .gitignore
└── README.md
```

## データベース

画像は SQLite データベースに保存されます。以下のテーブルが作成されます：

```sql
CREATE TABLE images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT NOT NULL,
    name TEXT
);
```

データベースファイル（`images.db`）は自動的に作成され、Git の管理対象外となっています。

## テスト

テストを実行するには以下のコマンドを使用します：

```bash
go test ./... -v
```
