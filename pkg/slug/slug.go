package slug

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

func Generate(name string) string {
	slug := strings.ToLower(name)
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	
	if slug == "" {
		slug = "project"
	}
	
	return slug
}

func Validate(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}
	
	if len(slug) < 2 {
		return fmt.Errorf("slug must be at least 2 characters long")
	}
	
	if len(slug) > 50 {
		return fmt.Errorf("slug must be no more than 50 characters long")
	}
	
	if !slugRegex.MatchString(slug) {
		return fmt.Errorf("slug must contain only lowercase letters, numbers, and hyphens")
	}
	
	if strings.Contains(slug, "--") {
		return fmt.Errorf("slug cannot contain consecutive hyphens")
	}
	
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return fmt.Errorf("slug cannot start or end with hyphens")
	}
	
	return nil
}
