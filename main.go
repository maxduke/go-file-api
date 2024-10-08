// main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	authKey  string
	port     string
	filePath string
)

func main() {
	// Parse command line arguments
	flag.StringVar(&authKey, "auth_key", "", "Authentication key for API calls")
	flag.StringVar(&port, "port", "8080", "Port number to listen on")
	flag.StringVar(&filePath, "file", "", "File path to write request data to")
	flag.Parse()

	if authKey == "" || filePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	logInfo("Starting server on port %s", port)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logError("Failed to start server: %v", err)
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		logWarning("Missing Authorization header from %s", r.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Expecting header in the format: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		logWarning("Invalid Authorization header format from %s", r.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := parts[1]
	if token != authKey {
		logWarning("Unauthorized access attempt from %s", r.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logError("Error reading request body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(filePath, body, 0644)
	if err != nil {
		logError("Error writing to file: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logSuccess("Request from %s processed successfully", r.RemoteAddr)
	fmt.Fprintln(w, "OK")
}

func logInfo(format string, a ...interface{}) {
	color.Cyan("%s [INFO] "+format, append([]interface{}{timestamp()}, a...)...)
}

func logWarning(format string, a ...interface{}) {
	color.Yellow("%s [WARN] "+format, append([]interface{}{timestamp()}, a...)...)
}

func logError(format string, a ...interface{}) {
	color.Red("%s [ERROR] "+format, append([]interface{}{timestamp()}, a...)...)
}

func logSuccess(format string, a ...interface{}) {
	color.Green("%s [SUCCESS] "+format, append([]interface{}{timestamp()}, a...)...)
}

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}