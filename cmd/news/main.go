package main

import (
	"context"
	"log"
	"os"

	"github.com/myuron/news/internal/discord"
	"github.com/myuron/news/internal/gemini"
	"github.com/myuron/news/internal/github"
	"github.com/myuron/news/internal/state"
)

func main() {
	geminiKey := requireEnv("GEMINI_API_KEY")
	webhookURL := requireEnv("DISCORD_WEBHOOK_URL")
	ghToken := os.Getenv("GITHUB_TOKEN")

	const (
		owner     = "anthropics"
		repo      = "claude-code"
		path      = "CHANGELOG.md"
		stateFile = "state.json"
	)

	st, err := state.Load(stateFile)
	if err != nil {
		log.Fatalf("failed to load state: %v", err)
	}

	file, err := github.FetchFile(ghToken, owner, repo, path)
	if err != nil {
		log.Fatalf("failed to fetch CHANGELOG: %v", err)
	}

	if file.SHA == st.LastSeenSHA {
		log.Println("No new changes in CHANGELOG.md")
		return
	}

	log.Printf("CHANGELOG changed: %s -> %s", st.LastSeenSHA, file.SHA)

	ctx := context.Background()
	summary, err := gemini.Summarize(ctx, geminiKey, file.Content)
	if err != nil {
		log.Fatalf("failed to summarize: %v", err)
	}

	title := "Claude Code CHANGELOG Update"
	if err := discord.Send(webhookURL, title, summary); err != nil {
		log.Fatalf("failed to send to Discord: %v", err)
	}

	log.Println("Discord notification sent successfully")

	st.LastSeenSHA = file.SHA
	if err := state.Save(stateFile, st); err != nil {
		log.Fatalf("failed to save state: %v", err)
	}
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}
