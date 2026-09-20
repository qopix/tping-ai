package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("==================================================================")
	fmt.Println("🤖 TPING-AI: АВТОНОМНЫЙ ИИ-АГЕНТ С ПОЛНЫМ ДОСТУПОМ К ФАЙЛАМ")
	fmt.Println("ИИ может читать, создавать, переписывать, удалять файлы и собирать софт.")
	fmt.Println("==================================================================")

	// Создаем историю диалога и закладываем туда системный промпт из config.go
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
			fmt.Println("Выход из программы...")
			break 
		}
		if userInput == "" { 
			continue 
		}

		// Записываем команду пользователя в память ИИ
		history = append(history, Message{Role: "user", Content: userInput})

		// Бесконечный цикл автономных размышлений (Thought -> Action -> Observation)
		for {
			fmt.Println("\n🧠 [ИИ планирует следующее действие]...")
			
			// Вызываем Ollama (функция лежит в agent.go)
			aiResponse := CallOllama(history)
			
			// Запоминаем мысли ИИ
			history = append(history, Message{Role: "assistant", Content: aiResponse})

			// Парсим ответ на наличие вызова системного инструмента
			tool, param := ParseToolCall(aiResponse)
			if tool == "" { 
				// Если ИИ не вызвал инструмент, значит он закончил задачу или выдал финальный ответ
				break 
			}

			fmt.Printf("\n⚙️ [Агент активирует метод]: %s со значением: %s\n", tool, param)
			var observation string

			// Выполнение физического действия на твоем ПК через функции из tools.go
			switch tool {
			case "TOOL_CLONE": 
				observation = CloneRepo(param)
			case "TOOL_LIST_DIR": 
				observation = ListDirectory(param)
			case "TOOL_READ_FILE": 
				observation = ReadSpecificFile(param)
			case "TOOL_WRITE_FILE": 
				observation = WriteSpecificFile(param)
			case "TOOL_DELETE_FILE": 
				observation = DeleteSpecificFile(param)
			case "TOOL_COMPILE": 
				observation = CompileProject(param)
			case "TOOL_RUN": 
				observation = RunBinary(param)
			default: 
				observation = "Error: Unknown tool name used by AI."
			}

			fmt.Println("📥 Результат операции добавлен в память ИИ.")
			
			// Возвращаем реальный ответ операционной системы обратно в контекст нейросети
			history = append(history, Message{
				Role:    "user",
				Content: fmt.Sprintf("[SYSTEM OBSERVATION]:\n%s\nAnalyze this output and decide on your next step.", observation),
			})
		}
	}
}
