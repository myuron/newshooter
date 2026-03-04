package changelog

import "testing"

func TestLatestSection(t *testing.T) {
	input := `# Changelog

## 2.0.0

- Breaking change
- New feature

## 1.0.0

- Initial release
`
	got := LatestSection(input)
	want := "## 2.0.0\n\n- Breaking change\n- New feature"
	if got != want {
		t.Errorf("got:\n%s\n\nwant:\n%s", got, want)
	}
}

func TestLatestSection_SingleVersion(t *testing.T) {
	input := `# Changelog

## 1.0.0

- Only version
`
	got := LatestSection(input)
	want := "## 1.0.0\n\n- Only version"
	if got != want {
		t.Errorf("got:\n%s\n\nwant:\n%s", got, want)
	}
}
