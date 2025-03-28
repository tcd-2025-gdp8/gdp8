package chatbot

import (
	"context"
	"errors"
	"os"

	"google.golang.org/genai"
)

type GeminiClient struct {
    model string
    Client *genai.Client
}

func NewClient(model string) (GeminiClient, error) {
    ctx := context.Background()
    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        return GeminiClient{}, errors.New("API_KEY not set")
    }

    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey:  apiKey,
        Backend: genai.BackendGeminiAPI,
    })

    if err != nil {
        return GeminiClient{}, err
    }

    return GeminiClient{
        model: model,
        Client: client,
    }, nil
}
