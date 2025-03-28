package chatbot

import (
	"context"
	"errors"
	"os"

	"google.golang.org/genai"
)

// rule of thumb for gemini is 4 characters per token
// https://ai.google.dev/gemini-api/docs/tokens?lang=go
// decided to made gemini-1.5-flash our default since it has more capabilities
// https://ai.google.dev/gemini-api/docs/models#gemini-1.5-flash
type GeminiClient struct {
    Client *genai.Client
    config *genai.GenerateContentConfig
    context context.Context
}

func NewClient() (GeminiClient, error) {
    ctx := context.Background()
    apiKey := os.Getenv("GEMINI_API_KEY")

    var config *genai.GenerateContentConfig = &genai.GenerateContentConfig{Temperature: genai.Ptr[float32](0.5)}

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
        Client: client,
        config: config,
        context: ctx,
    }, nil
}

func (gc *GeminiClient) SetTemperature(temp float32) {
    gc.config.Temperature = genai.Ptr[float32](temp)

}

func (gc *GeminiClient) Prompt(input string) (string, error) {
	result, err := gc.Client.Models.GenerateContent(
        gc.context,
        "gemini-1.5-flash",
        genai.Text(input),
        gc.config)

        if err != nil {
            return "", err
        }

        return result.Text(), err;
}
