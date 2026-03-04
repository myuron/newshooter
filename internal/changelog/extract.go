package changelog

import "strings"

// LatestSection extracts the first version section from a CHANGELOG.
// It looks for the first "## " heading and returns everything until the next "## ".
func LatestSection(content string) string {
	lines := strings.Split(content, "\n")
	var start int
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, "## ") {
			if !found {
				start = i
				found = true
			} else {
				return strings.TrimSpace(strings.Join(lines[start:i], "\n"))
			}
		}
	}
	if found {
		return strings.TrimSpace(strings.Join(lines[start:], "\n"))
	}
	return content
}
