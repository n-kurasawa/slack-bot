package bot_test

import (
	"os"
	"testing"

	"github.com/n-kurasawa/slack-bot/internal/bot"
	"github.com/n-kurasawa/slack-bot/internal/db"
	"github.com/slack-go/slack/slackevents"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *db.Store {
	t.Helper()

	// テスト用のデータベースファイルを作成
	dbPath := "test.db"
	t.Cleanup(func() {
		os.Remove(dbPath)
	})

	store, err := db.NewStore(dbPath)
	require.NoError(t, err)

	err = store.CreateTable()
	require.NoError(t, err)

	return store
}

func TestMessageEventHandler_HandleMessage(t *testing.T) {
	tests := []struct {
		name    string
		event   *slackevents.MessageEvent
		setup   func(*db.Store)
		want    string
		wantErr bool
	}{
		{
			name: "hello コマンドが成功する",
			event: &slackevents.MessageEvent{
				Text: "hello",
			},
			setup: func(store *db.Store) {},
			want:  "world",
		},
		{
			name: "image コマンドが成功する（名前指定なし）",
			event: &slackevents.MessageEvent{
				Text: "image",
			},
			setup: func(store *db.Store) {
				err := store.SaveImage("test", "http://example.com/test.jpg")
				require.NoError(t, err)
			},
			want: "test\nhttp://example.com/test.jpg",
		},
		{
			name: "image コマンドが成功する（名前指定あり）",
			event: &slackevents.MessageEvent{
				Text: "image test",
			},
			setup: func(store *db.Store) {
				err := store.SaveImage("test", "http://example.com/test.jpg")
				require.NoError(t, err)
			},
			want: "test\nhttp://example.com/test.jpg",
		},
		{
			name: "image コマンドが失敗する（画像が存在しない）",
			event: &slackevents.MessageEvent{
				Text: "image test",
			},
			setup:   func(store *db.Store) {},
			wantErr: true,
		},
		{
			name: "updateImage コマンドが成功する",
			event: &slackevents.MessageEvent{
				Text: "updateImage test http://example.com/test.jpg",
			},
			setup: func(store *db.Store) {},
			want:  "画像を保存しました :white_check_mark:",
		},
		{
			name: "updateImage コマンドが失敗する（引数が不足）",
			event: &slackevents.MessageEvent{
				Text: "updateImage test",
			},
			setup: func(store *db.Store) {},
			want:  "不正なコマンド形式です。使用方法: updateImage NAME URL",
		},
		{
			name: "未対応のコマンドの場合はエラーを返す",
			event: &slackevents.MessageEvent{
				Text: "unknown",
			},
			setup:   func(store *db.Store) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := setupTestDB(t)
			tt.setup(store)

			handler := bot.NewMessageEventHandler(store)
			got, err := handler.HandleMessage(tt.event)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
