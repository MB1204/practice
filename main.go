package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
)

func serveHTML(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "practice.html") // Ensure practice.html exists in the same directory
}

func submitQuestionHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        log.Println("Error parsing form:", err)
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    question := r.FormValue("question")
    if question == "" {
        http.Error(w, "Question cannot be empty", http.StatusBadRequest)
        return
    }

    fmt.Println("Received question:", question)

    answer, err := getAnswerFromAPI(question)
    if err != nil {
        log.Println("Error retrieving answer from API:", err)
        http.Error(w, "Error retrieving answer", http.StatusInternalServerError)
        return
    }

    fmt.Fprint(w, answer)
}

func main() {
    http.HandleFunc("/", serveHTML)
    http.HandleFunc("/submit-question", submitQuestionHandler)

    port := "8080"
    fmt.Println("Server starting on port", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getAnswerFromAPI(question string) (string, error) {
    apiURL := "https://api.replicate.com/v1/predictions"
    apiKey := os.Getenv("REPLICATE_API_TOKEN") // Load API key from environment variable

    if apiKey == "" {
        return "", fmt.Errorf("API key is not set")
    }

    requestBody, err := json.Marshal(map[string]interface{}{
        "input": map[string]string{
            "question": question,
        },
    })
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "Token "+apiKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("API request failed with status: %s", resp.Status)
    }

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("error decoding API response: %v", err)
    }

    prediction, ok := result["prediction"].(map[string]interface{})
    if !ok {
        return "", fmt.Errorf("unexpected API response structure")
    }

    answer, ok := prediction["answer"].(string)
    if !ok {
        return "", fmt.Errorf("answer not found in API response")
    }

    return answer, nil
}

