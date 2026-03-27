package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "github.com/taskfix/taskfix/internal/ai"
)

func main() {
    key := os.Getenv("OPENROUTER_API_KEY")
    if key == "" {
        if len(os.Args) > 1 {
            key = os.Args[1]
        }
    }
    if key == "" {
        fmt.Println("no API key provided")
        return
    }
    c := ai.NewClient("openrouter", key, "openai/gpt-4o-mini")
    models, err := c.FetchModels()
    if err != nil {
        fmt.Println("error:", err)
        return
    }
    fmt.Println("count:", len(models))
    for i, m := range models {
        if i >= 20 {
            break
        }
        fmt.Println(m)
    }

    // Raw HTTP check
    req, _ := http.NewRequest(http.MethodGet, "https://openrouter.ai/api/v1/models", nil)
    req.Header.Set("Authorization", "Bearer "+key)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        fmt.Println("raw http error:", err)
        return
    }
    defer resp.Body.Close()
    fmt.Println("raw status:", resp.Status)
    b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
    s := strings.TrimSpace(string(b))
    if len(s) > 0 {
        fmt.Println("raw body (first 1KB):")
        fmt.Println(s)
    } else {
        fmt.Println("raw body empty")
    }
}
