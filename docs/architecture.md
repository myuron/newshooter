# architecture.md

アプリケーションのアーキテクチャについて記載する。

## アプリケーション概要

情報をDiscordに通知するスキーム。

## 実現方式

例)Claude Codeの場合
1. GitHubのanthropic/claude-codeのCHANGELOG.mdを取得する。
2. Google Geminiを使用して日本語で要約する。
3. 要約結果をWebhookを用いてDiscordに送信する。
4. 1-3までは全てGitHub Actionsで定期実行を行う。

## プログラミング言語

Go
