package main

import (
	"fmt"
	"os/exec"
	"strings"
)

var defaultPrinter = ""

// GetAllPrintersUsingLp returns a list of all available printers from CUPS using the "lpstat -p" command.
func GetAllPrinters() ([]string, error) {
	cmd := exec.Command("lpstat", "-p")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get printers: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var printerNames []string
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 4 && fields[0] == "printer" { // Check for "printer" keyword
			printerNames = append(printerNames, fields[1]) // Extract the second field
		}
	}
	return printerNames, nil
}

// PrintDocumentUsingLp prints a document to the specified printer with options using the "lp" command.
//
//	options := map[string]string{
//		"copies": "2",
//		"page-ranges": "1-3,5",
//		"orientation-requested": "4", // 4 for landscape
//	}
func PrintDocument(filename string, options map[string]string) error {
	var args []string
	if len(defaultPrinter) == 0 {
		_, err := GetDefaultPrinter()
		if err != nil {
			return fmt.Errorf("failed getting default printer: %s", err.Error())
		}
	}

	args = append(args, "-d", defaultPrinter)

	for key, value := range options {
		args = append(args, fmt.Sprintf("-o %s=%s", key, value))
	}

	args = append(args, filename)

	cmd := exec.Command("lp", args...)
	fmt.Println(args)
	//TODO: err is the output so maybe test this later
	// troubleshoot note: maybe the arguments are the one making the error
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("failed to print file: %s", err.Error())
	}

	fmt.Printf("File %s sent to printer %s with options\n", filename, defaultPrinter)
	return nil
}

func SetDefaultPrinter(printerName string) error {

	cmd := exec.Command("lpoptions", "-d", printerName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to set default printer: %w", err)

	}
	defaultPrinter = printerName
	return nil
}

func GetDefaultPrinter() (string, error) {
	cmd := exec.Command("lpstat", "-d")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get default printer: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "system default destination:") {
			printerName := strings.TrimSpace(strings.Fields(line)[3])
			defaultPrinter = printerName
			return printerName, nil // Extract printer name

		}
	}

	return "", fmt.Errorf("default printer not found")
}
