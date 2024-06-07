package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "os"
)

type StatusResponse struct {
    StatusCode int    `json:"status_code"`
    Status     string `json:"status"`
    Message    string `json:"message"`
}


func checkWebsiteStatus(url string) (int, string, error) {
    client := http.Client{
        Timeout: 10 * time.Second,
    }

    resp, err := client.Get(url)
    if err != nil {
        return 0, "DOWN", err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        return resp.StatusCode, "UP", nil
    }
    return resp.StatusCode, "DOWN", fmt.Errorf("received status code: %d", resp.StatusCode)
}

// GET localhost:3000/status
func getControlPanelState(w http.ResponseWriter, r *http.Request) {
    websiteURL := os.Getenv("CONTROLPANEL_URL")

    fmt.Println(fmt.Sprintf("URL: %s status get", websiteURL))
    statusCode, status, err := checkWebsiteStatus(websiteURL)
    message := "Website is " + status

    if err != nil {
        message = err.Error()
    }

    response := StatusResponse{
        StatusCode: statusCode,
        Status:     status,
        Message:    message,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}