package main

import (
	"fmt"
	"os/exec"
	"regexp"
)

func ConvertToPdf(filePath string) (string, error) {
	outputChan := make(chan string)
	errorChan := make(chan error)

	go func() {
		cmd := exec.Command("lowriter", "--headless", "--print-to-file", "--outdir", uploadDir, filePath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			errorChan <- fmt.Errorf("failed to execute command: %w\n%s", err, string(output))
			return
		}

		regex := regexp.MustCompile(`->\s([^\s]+\.pdf)\susing`)
		match := regex.FindStringSubmatch(string(output))
		if len(match) > 1 {
			outputPath := match[1]
			outputChan <- outputPath
		} else {
			errorChan <- fmt.Errorf("no output path found in command output: %s", string(output))
		}
	}()

	select {
	case outputPath := <-outputChan:
		return outputPath, nil
	case err := <-errorChan:
		return "", err
	}
}
