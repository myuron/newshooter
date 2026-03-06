package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const maxDescriptionLen = 2000

type embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color,omitempty"`
}

type webhookPayload struct {
	Embeds []embed `json:"embeds"`
}

func Send(webhookURL, title, description string) error {
	if len(description) > maxDescriptionLen {
		description = description[:maxDescriptionLen-3] + "..."
	}

	payload := webhookPayload{
		Embeds: []embed{
			{
				Title:       title,
				Description: description,
				Color:       0x5865F2, // Discord Blurple
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord webhook returned %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
