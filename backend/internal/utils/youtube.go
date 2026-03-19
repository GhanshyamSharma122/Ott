package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// ExtractYouTubeVideoID extracts the video ID from various YouTube URL formats
func ExtractYouTubeVideoID(url string) (string, error) {
	patterns := []string{
		`(?:youtube\.com\/watch\?v=)([a-zA-Z0-9_-]{11})`,
		`(?:youtube\.com\/embed\/)([a-zA-Z0-9_-]{11})`,
		`(?:youtu\.be\/)([a-zA-Z0-9_-]{11})`,
		`(?:youtube\.com\/v\/)([a-zA-Z0-9_-]{11})`,
		`(?:youtube\.com\/shorts\/)([a-zA-Z0-9_-]{11})`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	// If input is already just a video ID (11 chars)
	url = strings.TrimSpace(url)
	if len(url) == 11 {
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]{11}$`, url)
		if matched {
			return url, nil
		}
	}

	return "", fmt.Errorf("could not extract YouTube video ID from: %s", url)
}

// GetThumbnailURL returns the YouTube thumbnail URL for a video ID
func GetThumbnailURL(videoID string) string {
	return fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)
}
