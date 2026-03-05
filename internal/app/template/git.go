package template

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type GitClient struct {
	cacheDir string
}

func NewGitClient(cacheDir string) *GitClient {
	if cacheDir == "" {
		homeDir, _ := os.UserHomeDir()
		cacheDir = filepath.Join(homeDir, ".scaffold", "cache")
	}
	return &GitClient{cacheDir: cacheDir}
}

func (g *GitClient) Clone(repoURL string) (string, error) {
	return g.CloneWithRef(repoURL, "")
}

func (g *GitClient) CloneWithRef(repoURL, ref string) (string, error) {
	repoName := extractRepoName(repoURL)
	suffix := ""
	if ref != "" && ref != "latest" {
		suffix = "-" + ref
	}
	targetPath := filepath.Join(g.cacheDir, repoName+suffix)

	if _, err := os.Stat(targetPath); err == nil {
		if err := os.RemoveAll(targetPath); err != nil {
			return "", err
		}
	}

	if err := os.MkdirAll(g.cacheDir, 0755); err != nil {
		return "", err
	}

	var cmd *exec.Cmd
	if ref != "" && ref != "latest" {
		cmd = exec.Command("git", "clone", "--depth", "1", "--branch", ref, repoURL, targetPath)
	} else {
		cmd = exec.Command("git", "clone", "--depth", "1", repoURL, targetPath)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return targetPath, nil
}

func (g *GitClient) ListRemoteTags(repoURL string) ([]string, error) {
	cmd := exec.Command("git", "ls-remote", "--tags", "--refs", repoURL)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tags []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) == 2 {
			ref := parts[1]
			tag := strings.TrimPrefix(ref, "refs/tags/")
			tags = append(tags, tag)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(tags)))

	return tags, nil
}

// TagInfo contains tag name and annotation information
type TagInfo struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (g *GitClient) ListTagsWithAnnotations(repoURL string) ([]TagInfo, error) {
	repoName := extractRepoName(repoURL)
	targetPath := filepath.Join(g.cacheDir, repoName+"-tags-meta")

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		if err := os.MkdirAll(g.cacheDir, 0755); err != nil {
			return nil, err
		}
		cmd := exec.Command("git", "clone", "--bare", "--filter=blob:none", repoURL, targetPath)
		if err := cmd.Run(); err != nil {
			return nil, err
		}
	} else {
		cmd := exec.Command("git", "fetch", "--tags", "--force")
		cmd.Dir = targetPath
		cmd.Run()
	}

	cmd := exec.Command("git", "tag", "-l", "--format=%(refname:short)|%(contents:subject)")
	cmd.Dir = targetPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var tags []TagInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		tag := TagInfo{Name: parts[0]}
		if len(parts) == 2 {
			tag.Message = strings.TrimSpace(parts[1])
		}
		tags = append(tags, tag)
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name > tags[j].Name
	})

	return tags, nil
}

func (g *GitClient) Pull(repoName string) error {
	targetPath := filepath.Join(g.cacheDir, repoName)
	cmd := exec.Command("git", "pull")
	cmd.Dir = targetPath
	return cmd.Run()
}

func (g *GitClient) IsCached(repoName string) bool {
	targetPath := filepath.Join(g.cacheDir, repoName)
	_, err := os.Stat(targetPath)
	return err == nil
}

func (g *GitClient) GetCachePath(repoName string) string {
	return filepath.Join(g.cacheDir, repoName)
}

func (g *GitClient) ClearCache() error {
	return os.RemoveAll(g.cacheDir)
}

func extractRepoName(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	name := parts[len(parts)-1]
	return strings.TrimSuffix(name, ".git")
}
