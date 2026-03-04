# 実装計画: Changelog Discord Notifier

## Context
GitHubリポジトリのCHANGELOG.mdを取得し、Google Geminiで日本語要約し、Discordに通知するGoアプリケーションを新規作成する。GitHub Actionsで定期実行する。

## プロジェクト構成

```
├── .github/workflows/notify.yml   # GitHub Actions ワークフロー
├── cmd/news/main.go               # エントリポイント
├── internal/
│   ├── github/client.go           # GitHub API でCHANGELOG取得
│   ├── gemini/client.go           # Gemini APIで日本語要約
│   ├── discord/webhook.go         # Discord Webhook送信
│   └── state/state.go             # 状態管理 (重複防止)
├── state.json                     # 最後に処理したSHA記録
├── go.mod
└── docs/architecture.md           # 既存
```

## 実装ステップ

### Step 1: `go.mod` 作成
- `module github.com/myuron/news`, Go 1.23
- 依存: `google.golang.org/genai` のみ (GitHub API・Discord Webhookは標準ライブラリで実装)

### Step 2: `internal/state/state.go`
- `state.json` の読み書き。`last_seen_sha` でCHANGELOGの変更を検知
- ファイル未存在時は空のStateを返す

### Step 3: `internal/github/client.go`
- `GET /repos/{owner}/{repo}/contents/{path}` で CHANGELOG.md とblob SHAを取得
- `net/http` + `encoding/json` + `encoding/base64` で実装
- `GITHUB_TOKEN` で認証

### Step 4: `internal/gemini/client.go`
- `google.golang.org/genai` SDK使用、モデル: `gemini-2.0-flash`
- プロンプト: 最新バージョンの変更を日本語で箇条書き要約

### Step 5: `internal/discord/webhook.go`
- Discord Webhook URLにEmbed形式でPOST
- 4096文字超過時は末尾を切り詰め

### Step 6: `cmd/news/main.go`
- 環境変数から設定読み込み (`GEMINI_API_KEY`, `DISCORD_WEBHOOK_URL`, `GITHUB_TOKEN`)
- フロー: state読込 → GitHub取得 → SHA比較 → 変更あれば要約 → Discord送信 → state保存
- エラーは `log.Fatalf` でfail fast

### Step 7: `.github/workflows/notify.yml`
- `schedule: cron '0 9 * * *'` (UTC 9:00 = JST 18:00) + `workflow_dispatch`
- Go build → 実行 → `state.json` に変更あればcommit & push
- `permissions: contents: write`

### Step 8: `state.json` 初期ファイル
- `{"last_seen_sha": ""}` で初回実行時にフル要約を送信

## 環境変数 (GitHub Secrets)
| 変数名 | 用途 |
|--------|------|
| `GEMINI_API_KEY` | Gemini API キー |
| `DISCORD_WEBHOOK_URL` | Discord Webhook URL |
| `GITHUB_TOKEN` | GitHub Actions 自動提供 |

## 検証方法
1. `go build ./cmd/news` でビルド確認
2. 環境変数をセットしてローカル実行: `GEMINI_API_KEY=... DISCORD_WEBHOOK_URL=... GITHUB_TOKEN=... go run ./cmd/news`
3. Discord チャンネルに要約が投稿されることを確認
4. 再実行して「No new changes」で終了することを確認
5. GitHub Actions で `workflow_dispatch` から手動実行して動作確認
