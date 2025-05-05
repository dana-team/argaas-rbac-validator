package config_fetcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	clusterTokensBasePath = "/mnt/cluster-tokens"
)

// FetchClusterToken reads the token file for the destination cluster.
func FetchClusterToken(server string) (string, error) {
	safeName := sanitizeClusterName(server)
	path := filepath.Join(clusterTokensBasePath, safeName+"-token")
	token, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read cluster token failed: %w", err)
	}
	return strings.TrimSpace(string(token)), nil
}

// sanitizeClusterName safely converts server URL to file-safe name.
func sanitizeClusterName(server string) string {
	safe := strings.ReplaceAll(server, "https://", "")
	safe = strings.ReplaceAll(safe, "http://", "")
	safe = strings.ReplaceAll(safe, ".", "-")
	safe = strings.ReplaceAll(safe, ":", "-")
	safe = strings.ReplaceAll(safe, "/", "-")
	return safe
}
