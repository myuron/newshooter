package state

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

type RepoState struct {
	LastSeenSHA string `json:"last_seen_sha"`
}

type State struct {
	Repos map[string]RepoState `json:"repos"`
}

func Load(path string) (State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return State{Repos: make(map[string]RepoState)}, nil
		}
		return State{}, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, err
	}
	if s.Repos == nil {
		s.Repos = make(map[string]RepoState)
	}
	// Migrate from legacy single-SHA format
	var legacy struct {
		LastSeenSHA string `json:"last_seen_sha"`
	}
	if json.Unmarshal(data, &legacy) == nil && legacy.LastSeenSHA != "" {
		if _, ok := s.Repos["anthropics/claude-code"]; !ok {
			s.Repos["anthropics/claude-code"] = RepoState{LastSeenSHA: legacy.LastSeenSHA}
		}
	}
	return s, nil
}

func Save(path string, s State) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
}

func (s State) SHA(repoKey string) string {
	return s.Repos[repoKey].LastSeenSHA
}

func (s *State) SetSHA(repoKey, sha string) {
	if s.Repos == nil {
		s.Repos = make(map[string]RepoState)
	}
	s.Repos[repoKey] = RepoState{LastSeenSHA: sha}
}
