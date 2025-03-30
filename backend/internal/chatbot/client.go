package chatbot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

// rule of thumb for gemini is 4 characters per token
// https://ai.google.dev/gemini-api/docs/tokens?lang=go
// decided to made gemini-1.5-flash our default since it has more capabilities
// https://ai.google.dev/gemini-api/docs/models#gemini-1.5-flash
type GeminiClient struct {
	Client *genai.Client
	config *genai.GenerateContentConfig
}

const model = "gemini-1.5-flash"

const basePrompt = `
    You are an assistant for a study group app used by undergraduate students at Trinity College Dublin.
    Your role is to help them study by giving clear, accurate answers and explanations.
    Always reply in plain text only. No markdown, no bullet points, no formatting. Keep it simple.
    You may receive study files as context before the user's question — use the data from these files
    to enhance your answer when applicable, but if the question is outside their scope,
    answer it using your general knowledge.
    Be friendly and encouraging, but stay focused on helping the user understand the topic.
`

func NewClient() (GeminiClient, error) {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")

	var config *genai.GenerateContentConfig = &genai.GenerateContentConfig{
		Temperature: genai.Ptr[float32](0.5)}

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
	}, nil
}

func (gc *GeminiClient) SetTemperature(temp float32) {
	gc.config.Temperature = genai.Ptr[float32](temp)

}

func (gc *GeminiClient) Prompt(ctx context.Context, input string) (string, error) {
	result, err := gc.Client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(fmt.Sprintf("%s\nUser question: %s", basePrompt, input)),
		gc.config)

	if err != nil {
		return "", err
	}

	return result.Text(), err
}

func (gc *GeminiClient) PromptWithContext(ctx context.Context, input string, memory map[string]string) (string, error) {

	sb := &strings.Builder{}

	for filename, data := range memory {
		sb.WriteString(fmt.Sprintf("[Start of %s]:\n%s[End of %s]\n", filename, data, filename))
	}

	result, err := gc.Client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(
			fmt.Sprintf("%s\nUse the following files as context:\n\n%s\nUser question: %s",
				basePrompt,
				sb.String(),
				input)),
		gc.config)

	if err != nil {
		return "", err
	}

	return result.Text(), err
}
