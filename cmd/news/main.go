package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"os"

	"github.com/myuron/news/internal/changelog"
	"github.com/myuron/news/internal/discord"
	"github.com/myuron/news/internal/gemini"
	"github.com/myuron/news/internal/jina"
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
		sourceURL
	)

	type target struct {
		owner  string
		repo   string
		path   string
		url    string
		title  string
		source source
	}

	targets := []target{
		{"anthropics", "claude-code", "CHANGELOG.md", "", "Claude Code CHANGELOG Update", sourceChangelog},
		{"openai", "codex", "", "", "Codex Release Update", sourceRelease},
		{"rork", "changelog", "", "https://rorkapp.notion.site/Changelog-for-Docs-and-Discord-2c76979e738b806abbb8dd3238507bff", "Rork Changelog Update", sourceURL},
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

		case sourceURL:
			markdown, err := jina.FetchMarkdown(t.url)
			if err != nil {
				log.Printf("[%s] failed to fetch URL via Jina: %v", repoKey, err)
				continue
			}
			id = fmt.Sprintf("%x", sha256.Sum256([]byte(markdown)))
			if id == st.SHA(repoKey) {
				log.Printf("[%s] No new changes", repoKey)
				continue
			}
			summary, err := gemini.Summarize(ctx, geminiKey, markdown)
			if err != nil {
				log.Printf("[%s] failed to summarize: %v", repoKey, err)
				continue
			}
			log.Printf("[%s] changed: %s -> %s", repoKey, st.SHA(repoKey), id)
			if err := discord.Send(webhookURL, t.title, summary); err != nil {
				log.Printf("[%s] failed to send to Discord: %v", repoKey, err)
				continue
			}
			log.Printf("[%s] Discord notification sent", repoKey)
			st.SetSHA(repoKey, id)
			changed = true
			continue
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
