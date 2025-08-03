package parser

import (
	"CVMatch/internal/config"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ledongthuc/pdf"
	"github.com/sheeiavellie/go-yandexgpt"
)

func ExtractTextFromPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var textBuilder bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	_, err = textBuilder.ReadFrom(b)
	if err != nil {
		return "", err
	}
	return textBuilder.String(), nil
}

func BuildPrompt(text string) string {
	return fmt.Sprintf(`
Ты — помощник по анализу резюме. Проанализируй текст и верни результат в формате JSON со следующей структурой:

{
  "full_name": "ФИО",
  "email": "email@example.com",
  "phone": "+7 000 000-00-00",
  "location": "Город",
  "skills": ["Go", "PostgreSQL"],
  "experience": [
	{
	  "company": "Компания",
	  "position": "Должность",
	  "start_date": "2023",
	  "end_date": "2027",
	  "description": "Описание работы"
	}
  ],
  "education": [
	{
	  "institution": "Университет",
	  "degree": "Степень",
	  "field": "Специальность",
	  "start_date": "2016-09-01" // дата начала обучения (если не указано, остваить пустым),
	  "end_date": "2020-06-30 // дата окончания обучения (обычно указана только дата окончания)"
	}
  ]
}

Текст резюме:
%s
`, text)
}

// Парсинг резюме через LLM (Ollama)
func ParseResumeWithLLM(pdfPath string, cfg *config.Config) (string, error) {
	start := time.Now()
	fmt.Println("[LLM] Начинаем парсинг резюме через LLM...")
	resumeText, err := ExtractTextFromPDF(pdfPath)
	if err != nil {
		fmt.Println("[LLM] Ошибка извлечения текста из PDF:", err)
		return "", err
	}
	fmt.Println("[LLM] Текст резюме успешно извлечён, длина:", len(resumeText))
	prompt := BuildPrompt(resumeText)
	fmt.Println("[LLM] Prompt сформирован, длина:", len(prompt))

	client := yandexgpt.NewYandexGPTClientWithAPIKey(cfg.YandexGPTIAM)
	request := yandexgpt.YandexGPTRequest{
		ModelURI: yandexgpt.MakeModelURI(cfg.YandexGPTCatalog, yandexgpt.YandexGPTModelLite),
		CompletionOptions: yandexgpt.YandexGPTCompletionOptions{
			Stream:      false,
			Temperature: 0.7,
			MaxTokens:   2000,
		},
		Messages: []yandexgpt.YandexGPTMessage{
			{
				Role: yandexgpt.YandexGPTMessageRoleSystem,
				Text: "Ты — парсер резюме. Возвращай только JSON в указанной структуре.",
			},
			{
				Role: yandexgpt.YandexGPTMessageRoleUser,
				Text: prompt,
			},
		},
	}
	response, err := client.GetCompletion(context.Background(), request)
	if err != nil {
		fmt.Println("Request error")
		return "", err
	}
	fmt.Println("[LLM] Ответ LLM получен, длина:", len(response.Result.Alternatives[0].Message.Text))
	// fmt.Println("[LLM] Ответ LLM:", response.Result.Alternatives[0].Message.Text)
	elapsed := time.Since(start)
	fmt.Printf("[LLM] Время парсинга резюме: %s\n", elapsed)
	return response.Result.Alternatives[0].Message.Text, nil
}
