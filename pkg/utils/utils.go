package utils

import (
	"bytes"
	"net/http"
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

func DoPost(u string, bodyData *bytes.Buffer) (*http.Response, error) {
	if bodyData != nil {
		return http.Post(u, "application/json", bodyData)
	}

	return http.Post(u, "application/json", nil)
}
