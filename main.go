package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

// Handler for serving static files from "dist" and "assets" directories
func staticFileHandler() http.Handler {
	fileServer := http.FileServer(http.Dir("."))
	return http.StripPrefix("/static/", fileServer)
}

// Handler for serving the HTML template
func serveTemplate(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil) // Data can be passed to the template if needed
}

// Handler for serving the PDF file
func servePDF(w http.ResponseWriter, r *http.Request) {
	// Define the file path for the PDF
	pdfPath := filepath.Join("assets", "CV Alfia Meilani Putri_NEW.pdf")

	// Open the PDF file
	file, err := os.Open(pdfPath)
	if err != nil {
		http.Error(w, "Error opening PDF file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get file information
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Error getting PDF file info", http.StatusInternalServerError)
		return
	}

	// Set headers for the PDF file
	w.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name())
	w.Header().Set("Content-Type", "application/pdf")

	// Serve the file using ServeContent to handle the headers properly
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func main() {
	// Serve static files from the "dist/" and "assets/" directories
	http.Handle("/static/", staticFileHandler())

	// Serve the template HTML
	http.HandleFunc("/", serveTemplate)

	// // Serve the PDF file for download
	// http.HandleFunc("/download-cv", servePDF)

	// Start the server on port 8080 with a log notification
	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)

	// Run the server
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
