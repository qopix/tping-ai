package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("==================================================================")
	fmt.Println("🤖 TPING-AI: ПОЛНОСТЬЮ АВТОНОМНЫЙ ИИ-АГЕНТ С ОБУЧЕНИЕМ НА ХОДУ")
	fmt.Println("Дайте ему ссылку на репозиторий GitHub и поставьте задачу.")
	fmt.Println("Пример: 'Изучи https://github.com и протестируй сайт httpbin.org'")
	fmt.Println("Для завершения введите 'exit'")
	fmt.Println("==================================================================")

	var history []Message
	history = append(history, Message{Role: "system", Content: systemPrompt})

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nUser > ")
		if !scanner.Scan() {
			break
		}
		userInput := strings.TrimSpace(scanner.Text())
		if strings.ToLower(userInput) == "exit" {
			break
		}
		if userInput == "" {
			continue
		}

		history = append(history, Message{Role: "user", Content: userInput})

		// Цикл автономного выполнения задач агентом (Thought -> Action -> Observation)
		for {
			fmt.Println("\n🧠 [ИИ планирует следующее действие]...")
			aiResponse := CallOllama(history)
			history = append(history, Message{Role: "assistant", Content: aiResponse})

			// Проверяем, сгенерировал ли ИИ команду для вызова инструмента
			tool, param := ParseToolCall(aiResponse)
			if tool == "" {
				// Если ИИ не вызывает инструменты, значит он закончил работу и выдал финальный ответ
				break
			}

			fmt.Printf("\n⚙️ [Агент активирует системный метод]: %s %s\n", tool, param)
			var observation string

			// Выполнение действия в зависимости от решения ИИ
			switch tool {
			case "TOOL_CLONE":
				observation = CloneRepo(param)
			case "TOOL_READ_FILES":
				observation = ReadWorkspaceFiles()
			case "TOOL_COMPILE":
				observation = CompileProject(param)
			case "TOOL_RUN":
				observation = RunBinary(param)
			default:
				observation = "Error: Unknown action tool."
			}

			fmt.Println("📥 Результат операции добавлен в контекст ИИ.")
			// Возвращаем реальный результат работы утилиты/компилятора обратно в память нейросети
			history = append(history, Message{
				Role:    "user",
				Content: fmt.Sprintf("[SYSTEM OBSERVATION]:\n%s\nAnalyze this output and decide on your next step.", observation),
			})
		}
	}
}
