package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var recentPdf = ""

func processDocumentConvert(filePath string) (string, error) {
	// improve this soon TODO:
	fmt.Println(filePath)
	if strings.Split(filePath, ".")[1] != "pdf" {
		// convert to pdf
		pdfPath, err := ConvertToPdf(filePath)
		if err != nil {
			return "", err
		}

		return pdfPath, nil
	}
	return filePath, nil
}

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
		recentPdf = outputPath
		return outputPath, nil
	case err := <-errorChan:
		return "", err
	}
}

type DocumentInfo struct {
	Paper       string
	Orientation int
	Pages       int
}

func parseDocumentInfo(stdout string) (*DocumentInfo, error) {
	// Compile regular expressions for efficient reuse
	pageSizeRegex := regexp.MustCompile(`Page size:\s*(\d+)\s*x\s*(\d+)\s+pts\s*(?:\(([^)]*)\))?`)
	pageCountRegex := regexp.MustCompile(`Pages:\s+(\d+)`)

	matches := pageSizeRegex.FindStringSubmatch(stdout)
	if matches == nil {
		return nil, fmt.Errorf("unable to extract page size information")
	}

	widthStr, heightStr := matches[1], matches[2]
	width, err := strconv.Atoi(widthStr)
	if err != nil {
		return nil, fmt.Errorf("invalid width value: %w", err)
	}
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		return nil, fmt.Errorf("invalid height value: %w", err)
	}

	matches = pageCountRegex.FindStringSubmatch(stdout)
	if matches == nil {
		return nil, fmt.Errorf("unable to extract page count information")
	}

	pageCountStr := matches[1]
	pages, err := strconv.Atoi(pageCountStr)
	if err != nil {
		return nil, fmt.Errorf("invalid page count value: %w", err)
	}

	// orientation := "portrait"
	orientation := 3
	//  3 for portrait as it is the convention in cups
	if width > height {
		orientation = 4
	}

	identifiedPaper := "Unknown"
	if width == 612 && height == 792 || width == 792 && height == 612 {
		identifiedPaper = "Letter"
	} else if width == 612 && height == 1008 || width == 1008 && height == 612 {
		identifiedPaper = "Legal"
	} else if width == 595 && height == 842 || width == 842 && height == 595 {
		identifiedPaper = "A4"
	}

	return &DocumentInfo{
		Paper:       identifiedPaper,
		Orientation: orientation,
		Pages:       pages,
	}, nil
}

func getPdfData(filePath string) (*DocumentInfo, error) {
	cmd := exec.Command("pdfinfo", filePath)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing pdfinfo: %s", err.Error())
	}

	result, err := parseDocumentInfo(string(stdout))
	if err != nil {
		return nil, fmt.Errorf("error parsing document information: %s", err.Error())
	}

	return result, nil
}
