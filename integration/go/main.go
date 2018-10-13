package integration

import (
	"fmt"
	"os"
	"os/exec"
)

// Listen performs a zenaton listen given an app id and api token environment variables
func Listen() error {

	args := []string{"listen", "--app_id=" + os.Getenv("ZENATON_APP_ID"),
		"--api_token=" + os.Getenv("ZENATON_API_TOKEN"),
		"--app_env=go-integration-test",
		"--boot=boot/boot.go"}

	cmd := exec.Command("zenaton", args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("problem listening: ", string(out))
	}

	return err
}
