package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var fileDest string = ""

func setDefaultPrinterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()

	printerName := queryParams.Get("printer-name")

	err := SetDefaultPrinter(printerName)

	if err != nil {
		http.Error(w, "error settings default printer", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "successfully changed the default printer to %s", printerName)
}

func getDefaultPrinterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	printerName, err := GetDefaultPrinter()
	if err != nil {
		http.Error(w, "error getting default printer \n"+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "default printer: %s", printerName)
}

func pdfPreviewHandler(w http.ResponseWriter, r *http.Request) {
	// currentPath, _ := os.Executable()
	fmt.Print(recentPdf)
	http.ServeFile(w, r, recentPdf)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong!")
}

func scheduleFileDeletion(filePath string, delay time.Duration) {
	time.Sleep(delay)
	err := os.Remove(filePath)
	fileDest = ""
	if err != nil {
		fmt.Println("Error deleting file:", err)
	}
	fmt.Printf("File %s deleted after %v\n", filePath, delay)
}
