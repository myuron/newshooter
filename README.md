# newshooter

開発ツールの更新情報を追跡し、Gemini 2.5 Flash で日本語に要約して Discord に通知する Go アプリケーションです。

GitHub Actions で定期実行されます。

## 追跡対象

| プロジェクト | ソース種別 | 取得方法 |
|---|---|---|
| [anthropics/claude-code](https://github.com/anthropics/claude-code) | CHANGELOG.md | GitHub Contents API |
| [openai/codex](https://github.com/openai/codex) | GitHub Release | GitHub Releases API |
| [Rork](https://rorkapp.notion.site/Changelog-for-Docs-and-Discord-2c76979e738b806abbb8dd3238507bff) | Notion ページ | Jina Reader API |

## 仕組み

1. 各ソースから最新の変更内容を取得
2. `state.json` に保存された前回の状態と比較し、差分があるか判定
3. 更新がある場合、Gemini 2.5 Flash で日本語に要約
4. Discord Webhook で通知（Embed 形式、4096 文字制限）
5. `state.json` を更新し、GitHub Actions が自動コミット
