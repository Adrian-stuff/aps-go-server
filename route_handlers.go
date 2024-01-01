package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var fileDest string = ""

func getAllPrintersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	printers, err := GetAllPrinters()
	if err != nil {
		http.Error(w, "Error getting printers", http.StatusInternalServerError)
		return
	}

	data := map[string][]string{"printers": printers}

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonData)

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer file.Close()
	if strings.Split(handler.Filename, ".")[1] != "docx" && strings.Split(handler.Filename, ".")[1] != "doc" && strings.Split(handler.Filename, ".")[1] != "pdf" {
		http.Error(w, "Wrong file extension "+strings.Split(handler.Filename, ".")[1], http.StatusBadRequest)
		return
	}
	// Generate a unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)

	// Create the file
	destPath := filepath.Join(uploadDir, filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// set global file destination
	fileDest = destPath

	// Copy the file content to the destination file
	_, err = io.Copy(destFile, file)
	if err != nil {
		http.Error(w, "Error copying file content", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully\n")

	// Schedule file deletion after 30 minutes
	go scheduleFileDeletion(destPath, 5*time.Minute)
}

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

func printDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if len(fileDest) == 0 {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	// TODO: implement converting the document first
	pdfPath, errpdf := processDocumentConvert(fileDest)
	if errpdf != nil {
		http.Error(w, "error processing document"+errpdf.Error(), http.StatusInternalServerError)
		return

	}

	err := PrintDocument(pdfPath, map[string]string{"copies": "1"})

	if err != nil {
		http.Error(w, "error printing"+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "successfully printed")

	// delete file
	// go scheduleFileDeletion(fileDest, 1)
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
