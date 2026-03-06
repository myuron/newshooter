package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/myuron/news/internal/changelog"
	"github.com/myuron/news/internal/discord"
	"github.com/myuron/news/internal/gemini"
	"github.com/myuron/news/internal/github"
	"github.com/myuron/news/internal/state"
)

func main() {
	geminiKey := requireEnv("GEMINI_API_KEY")
	webhookURL := requireEnv("DISCORD_WEBHOOK_URL")
	ghToken := os.Getenv("GITHUB_TOKEN")

	const stateFile = "state.json"

	type source int
	const (
		sourceChangelog source = iota
		sourceRelease
	)

	type target struct {
		owner  string
		repo   string
		path   string
		title  string
		source source
	}

	targets := []target{
		{"anthropics", "claude-code", "CHANGELOG.md", "Claude Code CHANGELOG Update", sourceChangelog},
		{"openai", "codex", "", "Codex Release Update", sourceRelease},
	}

	st, err := state.Load(stateFile)
	if err != nil {
		log.Fatalf("failed to load state: %v", err)
	}

	changed := false
	ctx := context.Background()

	for _, t := range targets {
		repoKey := t.owner + "/" + t.repo

		var id, content string

		switch t.source {
		case sourceChangelog:
			file, err := github.FetchFile(ghToken, t.owner, t.repo, t.path)
			if err != nil {
				log.Printf("[%s] failed to fetch CHANGELOG: %v", repoKey, err)
				continue
			}
			id = file.SHA
			content = changelog.LatestSection(file.Content)

		case sourceRelease:
			rel, err := github.FetchLatestRelease(ghToken, t.owner, t.repo)
			if err != nil {
				log.Printf("[%s] failed to fetch release: %v", repoKey, err)
				continue
			}
			id = rel.TagName
			content = fmt.Sprintf("# %s\n\n%s", rel.Name, rel.Body)
		}

		if id == st.SHA(repoKey) {
			log.Printf("[%s] No new changes", repoKey)
			continue
		}

		log.Printf("[%s] changed: %s -> %s", repoKey, st.SHA(repoKey), id)

		summary, err := gemini.Summarize(ctx, geminiKey, content)
		if err != nil {
			log.Printf("[%s] failed to summarize: %v", repoKey, err)
			continue
		}

		if err := discord.Send(webhookURL, t.title, summary); err != nil {
			log.Printf("[%s] failed to send to Discord: %v", repoKey, err)
			continue
		}

		log.Printf("[%s] Discord notification sent", repoKey)

		st.SetSHA(repoKey, id)
		changed = true
	}

	if changed {
		if err := state.Save(stateFile, st); err != nil {
			log.Fatalf("failed to save state: %v", err)
		}
	}
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}
