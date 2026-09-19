package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const workspaceDir = "./agent_workspace"

// 1. Клонирование
func CloneRepo(repoURL string) string {
	_ = os.RemoveAll(workspaceDir)
	cmd := exec.Command("git", "clone", repoURL, workspaceDir)
	if err := cmd.Run(); err != nil { return fmt.Sprintf("Clone failed: %v", err) }
	return "Repository cloned successfully."
}

// 2. Просмотр папки
func ListDirectory(path string) string {
	targetPath := filepath.Join(workspaceDir, path)
	files, err := os.ReadDir(targetPath)
	if err != nil { return fmt.Sprintf("Failed to list dir: %v", err) }
	
	var result strings.Builder
	for _, file := range files {
		typeStr := "File"
		if file.IsDir() { typeStr = "Dir" }
		result.WriteString(fmt.Sprintf("[%s] %s\n", typeStr, file.Name()))
	}
	return result.String()
}

// 3. Чтение конкретного файла
func ReadSpecificFile(path string) string {
	targetPath := filepath.Join(workspaceDir, path)
	content, err := os.ReadFile(targetPath)
	if err != nil { return fmt.Sprintf("Failed to read file: %v", err) }
	return string(content)
}

// 4. Запись/Создание/Перезапись файла
func WriteSpecificFile(params string) string {
	parts := strings.SplitN(params, "|", 2)
	if len(parts) < 2 { return "Error: Invalid write parameters. Use path|content" }
	
	targetPath := filepath.Join(workspaceDir, parts[0])
	// Создаем подпапки, если их нет
	_ = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
	
	err := os.WriteFile(targetPath, []byte(parts[1]), 0644)
	if err != nil { return fmt.Sprintf("Failed to write file: %v", err) }
	return fmt.Sprintf("File %s successfully written/updated.", parts[0])
}

// 5. Удаление файла или папки
func DeleteSpecificFile(path string) string {
	targetPath := filepath.Join(workspaceDir, path)
	err := os.RemoveAll(targetPath)
	if err != nil { return fmt.Sprintf("Failed to delete: %v", err) }
	return fmt.Sprintf("Successfully deleted %s", path)
}

// 6. Компиляция
func CompileProject(buildCommand string) string {
	args := strings.Fields(buildCommand)
	if len(args) == 0 { return "Error: Empty compile command" }
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workspaceDir
	out, err := cmd.CombinedOutput()
	if err != nil { return fmt.Sprintf("Compile failed: %v\n%s", err, string(out)) }
	return fmt.Sprintf("Compilation successful!\n%s", string(out))
}

// 7. Запуск бинарников
func RunBinary(runCommand string) string {
	args := strings.Fields(runCommand)
	if len(args) == 0 { return "Error: Empty run command" }
	cmdName := args[0]
	if strings.HasPrefix(cmdName, "./") {
		cmdName = filepath.Join(workspaceDir, cmdName[2:])
	}
	cmd := exec.Command(cmdName, args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil { return fmt.Sprintf("Execution failed: %v\n%s", err, string(out)) }
	return string(out)
}
