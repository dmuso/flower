package user

import (
	"fmt"
	"strings"

	"flower/api/internal/platform/auth"
)

func inferUsername(email string) string {
	local := email
	if at := strings.Index(email, "@"); at >= 0 {
		local = email[:at]
	}
	local = strings.ToLower(local)
	var b strings.Builder
	for _, r := range local {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		}
	}
	name := b.String()
	if name == "" {
		name = "user"
	}
	if len(name) > 100 {
		name = name[:100]
	}
	return name
}

func uniquifyUsername(base string, taken func(string) (bool, error)) (string, error) {
	if base == "" {
		return "", fmt.Errorf("username: empty base")
	}
	ok, err := taken(base)
	if err != nil {
		return "", err
	}
	if !ok {
		return base, nil
	}
	root := base
	if len(root) > 95 {
		root = root[:95]
	}
	for i := 0; i < 32; i++ {
		suffix, err := auth.RandomUnambiguous(4)
		if err != nil {
			return "", err
		}
		candidate := root + "-" + suffix
		inUse, err := taken(candidate)
		if err != nil {
			return "", err
		}
		if !inUse {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("username: could not allocate a unique username")
}
