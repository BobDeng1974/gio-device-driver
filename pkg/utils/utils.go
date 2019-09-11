package utils

import (
	"os/exec"
	"strings"
)

func GetHostIP() (string, error) {
	// route | awk '/default/ { print $2 }'

	cmd := exec.Command("ash", "-c", "route | awk '/default/ { print $2 }'")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	ip := strings.TrimSuffix(strings.TrimSpace(string(out)), "\n")

	return ip, nil
}
