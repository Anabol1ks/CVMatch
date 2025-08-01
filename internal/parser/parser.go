package parser

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ledongthuc/pdf"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
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
	  "start_date": "2016-09-01" // дата начала обучения,
	  "end_date": "2020-06-30 // дата окончания обучения (обычно указана только дата окончания)"
	}
  ]
}

Текст резюме:
%s
`, text)
}

// Парсинг резюме через LLM (Ollama)
func ParseResumeWithLLM(pdfPath string, modelName string) (string, error) {
	start := time.Now()
	fmt.Println("[LLM] Начинаем парсинг резюме через LLM...")
	resumeText, err := ExtractTextFromPDF(pdfPath)
	if err != nil {
		fmt.Println("[LLM] Ошибка извлечения текста из PDF:", err)
		return "", err
	}
	fmt.Println("[LLM] Текст резюме успешно извлечён, длина:", len(resumeText))
	llm, err := ollama.New(ollama.WithModel(modelName))
	if err != nil {
		fmt.Println("[LLM] Ошибка инициализации Ollama:", err)
		return "", err
	}
	ctx := context.Background()
	prompt := BuildPrompt(resumeText)
	fmt.Println("[LLM] Prompt сформирован, длина:", len(prompt))
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Ты — парсер резюме. Возвращай только JSON в указанной структуре."),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}
	resp, err := llm.GenerateContent(ctx, content)
	if err != nil {
		fmt.Println("[LLM] Ошибка генерации ответа LLM:", err)
		return "", err
	}
	fmt.Println("[LLM] Ответ LLM получен, длина:", len(resp.Choices[0].Content))
	fmt.Println(resp.Choices[0].Content)
	elapsed := time.Since(start)
	fmt.Printf("[LLM] Время парсинга резюме: %s\n", elapsed)
	return resp.Choices[0].Content, nil
}
