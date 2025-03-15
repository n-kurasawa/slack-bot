# Slack Bot

シンプルな Slack Bot です。"hello" というメッセージに "world" と返信します。

## 必要条件

- Go 1.21 以上
- Slack App の設定
  - Event Subscriptions の有効化
  - Bot Token Scopes: `chat:write`, `app_mentions:read`, `messages.im:read`
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

## 実行方法

```bash
go run main.go
```

サーバーは port 3000 で起動します。

## Slack App の設定

1. [Slack API](https://api.slack.com/apps) でアプリを作成
2. "Event Subscriptions" を有効化
3. Request URL に `https://あなたのドメイン/slack/events` を設定
4. Bot Token Scopes に必要な権限を追加
5. アプリをワークスペースにインストール

## 開発環境でのテスト

ローカル環境で開発する場合は、[ngrok](https://ngrok.com/) などのツールを使用してローカルサーバーを公開できます：

```bash
ngrok http 3000
```

生成された URL + "/slack/events" を Slack アプリの Request URL として設定してください。
