package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const workspaceDir = "./agent_workspace"

// Инструмент 1: Клонирование репозитория
func CloneRepo(repoURL string) string {
	_ = os.RemoveAll(workspaceDir) // Очищаем старый воркспейс перед клонированием нового

	cmd := exec.Command("git", "clone", repoURL, workspaceDir)
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Failed to clone repository: %v", err)
	}
	return "Repository successfully cloned into workspace. Now you should read files to analyze the source code."
}

// Инструмент 2: Чтение содержимого и исходного кода файлов проекта
func ReadWorkspaceFiles() string {
	var result strings.Builder
	
	err := filepath.Walk(workspaceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Пропускаем папку .git
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}

		// Читаем только файлы кода, скрипты и документацию
		ext := filepath.Ext(path)
		if ext == ".cpp" || ext == ".go" || ext == ".md" || ext == ".h" || ext == ".c" || info.Name() == "Makefile" {
			content, err := os.ReadFile(path)
			if err == nil {
				result.WriteString(fmt.Sprintf("\n--- FILE: %s ---\n", filepath.Base(path)))
				result.WriteString(string(content))
				result.WriteString("\n")
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Sprintf("Error reading files: %v", err)
	}
	if result.Len() == 0 {
		return "No core source files (.cpp, .go, .md, Makefile) found in the workspace."
	}
	return result.String()
}

// Инструмент 3: Компиляция утилиты
func CompileProject(buildCommand string) string {
	args := strings.Fields(buildCommand)
	if len(args) == 0 {
		return "Error: Empty compile command"
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workspaceDir // Выполняем сборку внутри папки проекта

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Compilation failed: %v\nOutput:\n%s", err, string(output))
	}
	return fmt.Sprintf("Compilation successful!\nOutput:\n%s", string(output))
}

// Инструмент 4: Выполнение скомпилированного бинарника
func RunBinary(runCommand string) string {
	args := strings.Fields(runCommand)
	if len(args) == 0 {
		return "Error: Empty run command"
	}

	// Если ИИ пытается запустить локальный бинарник (например, ./tping),
	// мы корректируем путь, так как запуск происходит из рабочей директории
	cmdName := args[0]
	if strings.HasPrefix(cmdName, "./") {
		cmdName = filepath.Join(workspaceDir, cmdName[2:])
	}

	cmd := exec.Command(cmdName, args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Execution failed: %v\nOutput:\n%s", err, string(output))
	}
	return fmt.Sprintf("Execution result:\n%s", string(output))
}
