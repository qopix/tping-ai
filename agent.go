package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Отправка истории диалога в Ollama с выводом текста в реальном времени
func CallOllama(messages []Message) string {
	reqBody := OllamaRequest{
		Model:    modelName,
		Stream:   true,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Sprintf("Error encoding JSON: %v", err)
	}

	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "\n[Ошибка]: Не удалось подключиться к Ollama. Убедитесь, что сервер запущен командой 'ollama serve'\n"
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var fullResponse strings.Builder

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF || err != nil {
			break
		}

		var chunk OllamaResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			continue
		}

		fmt.Print(chunk.Message.Content)
		fullResponse.WriteString(chunk.Message.Content)

		if chunk.Done {
			break
		}
	}
	fmt.Println()
	return fullResponse.String()
}

// Парсер текста для извлечения триггеров инструментов, сгенерированных ИИ
func ParseToolCall(content string) (string, string) {
	lines := strings.Split(content, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "TOOL_CLONE:") {
			return "TOOL_CLONE", strings.TrimPrefix(line, "TOOL_CLONE:")
		}
		if strings.HasPrefix(line, "TOOL_READ_FILES") {
			return "TOOL_READ_FILES", ""
		}
		if strings.HasPrefix(line, "TOOL_COMPILE:") {
			return "TOOL_COMPILE", strings.TrimPrefix(line, "TOOL_COMPILE:")
		}
		if strings.HasPrefix(line, "TOOL_RUN:") {
			return "TOOL_RUN", strings.TrimPrefix(line, "TOOL_RUN:")
		}
	}
	return "", ""
}
