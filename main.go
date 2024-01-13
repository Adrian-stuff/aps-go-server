package main

import (
	"fmt"
	"net/http"
	"os"
)

const uploadDir = "./uploads"

func main() {
	// Create the upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Println("Error creating upload directory:", err)
		return
	}
	// static html
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/", http.StripPrefix("/", fs))

	// printer handlers
	http.HandleFunc("/printers", getAllPrintersHandler)
	http.HandleFunc("/setDefaultPrinter", setDefaultPrinterHandler)
	http.HandleFunc("/getDefaultPrinter", getDefaultPrinterHandler)

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/print", printDocumentHandler)

	// pdf preview
	http.HandleFunc("/pdfPreview", pdfPreviewHandler)

	http.HandleFunc("/ping", pingHandler)

	// websocket
	http.HandleFunc("/ws", websocketHandler)
	fmt.Println("listening at port 8080")
	http.ListenAndServe(":8080", nil)
}
