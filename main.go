package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

func serveHTML(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "practice.html")
}

func submitQuestionHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    question := r.FormValue("question")

    answer, err := getAnswerFromAPI(question)
    if err != nil {
        http.Error(w, "Error retrieving answer", http.StatusInternalServerError)
        return
    }

    // Respond with the answer
    fmt.Fprintf(w, answer)
}

func main() {
    http.HandleFunc("/", serveHTML)
    http.HandleFunc("/submit-question", submitQuestionHandler)
    http.ListenAndServe(":8080", nil) // Start the server on port 8080
}

func getAnswerFromAPI(question string) (string, error) {
    apiURL := "https://api.replicate.com/v1/predictions" // Replicate API endpoint
    requestBody, err := json.Marshal(map[string]interface{}{
        "version": "your_model_version", // Replace with the specific model version
        "input": map[string]string{
            "question": question,
        },
    })
    if err != nil {
        return "", err
    }

    resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("API request failed with status: %s", resp.Status)
    }

    var result map[string]string
    body, _ := ioutil.ReadAll(resp.Body)
    json.Unmarshal(body, &result)

    return result["answer"], nil
}

