// main.go
package main

import (
    "flag"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"

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

    color.Cyan("Starting server on port %s", port)
    http.HandleFunc("/", handler)
    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        color.Red("Failed to start server: %v", err)
        os.Exit(1)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        color.Yellow("Missing Authorization header from %s", r.RemoteAddr)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Expecting header in the format: "Bearer <token>"
    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) != 2 || parts[0] != "Bearer" {
        color.Yellow("Invalid Authorization header format from %s", r.RemoteAddr)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    token := parts[1]
    if token != authKey {
        color.Yellow("Unauthorized access attempt from %s", r.RemoteAddr)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    defer r.Body.Close()
    body, err := io.ReadAll(r.Body)
    if err != nil {
        color.Red("Error reading request body: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = os.WriteFile(filePath, body, 0644)
    if err != nil {
        color.Red("Error writing to file: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    color.Green("Request from %s processed successfully", r.RemoteAddr)
    fmt.Fprintln(w, "OK")
}
