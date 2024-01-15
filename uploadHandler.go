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

	// convert to pdf
	_, errpdf := processDocumentConvert(fileDest)
	if errpdf != nil {
		http.Error(w, "error processing document"+errpdf.Error(), http.StatusInternalServerError)
		return
	}

	docInfo, errInfo := getPdfData(recentPdf)
	if errInfo != nil {
		http.Error(w, "error getting pdf Data"+errInfo.Error(), http.StatusInternalServerError)
		return
	}

	// return docInfo to the user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docInfo)
	// fmt.Fprintf(w, "File uploaded successfully\n")
	// Schedule file deletion after 30 minutes
	go scheduleFileDeletion(destPath, 5*time.Minute)
}