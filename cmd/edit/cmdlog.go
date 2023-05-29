package edit

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const logFilePath = "cmd.log"

func ApplyLog(key, valueType string) {
	line := fmt.Sprintf("%s %s", key, valueType)
	appendToFile(line)
}

func DelLog(key, valueType string) {
	line := fmt.Sprintf("%s %s", key, valueType)
	removeLineFromFile(line)
}

func ResetLog() []string {
	file, err := os.OpenFile(logFilePath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err := os.Truncate(logFilePath, 0); err != nil {
		log.Fatal(err)
	}

	return lines
}

func appendToFile(line string) {
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := file.WriteString(line + "\n"); err != nil {
		log.Fatal(err)
	}
}

func removeLineFromFile(line string) {
	file, err := os.OpenFile(logFilePath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currLine := scanner.Text()
		if currLine != line {
			lines = append(lines, currLine)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err := file.Truncate(0); err != nil {
		log.Fatal(err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			log.Fatal(err)
		}
	}
	if err := writer.Flush(); err != nil {
		log.Fatal(err)
	}
}
