/*
This file has api functions for working with the data file structure we define.
The file has two primitives: ReadLine and WriteLine
*/

package api

import (
	"bufio"
	"fmt"
	"os"
)

/*
ReadLine reads the n-th line from the file and returns it as a string.
If the files doesn't contain n lines, it returns an empty string.
If the files doesn't exist it returns an error.
*/
func ReadLine(filePath string, n int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to open the file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line := 1
	for scanner.Scan() {
		if line == n {
			return scanner.Text(), nil
		}
		line++
	}

	if scanner.Err() != nil {
		return "", fmt.Errorf("error while reading the file: %v", scanner.Err())
	}

	return "", nil
}

/*
WriteLine writes the given string to the n-th line of the file.
If the files doesn't contain n lines it appends empty lines till n-1 and then the string at line n
Other lines are not modified.
*/
func WriteLine(filePath string, n int, text string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("unable to open the file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if scanner.Err() != nil {
		return fmt.Errorf("error while reading the file: %v", scanner.Err())
	}

	if n > len(lines) {
		for i := len(lines); i < n-1; i++ {
			lines = append(lines, "")
		}
	}

	lines = append(lines[:n-1], text)
	lines = append(lines, lines[n-1:]...)

	file.Truncate(0)
	file.Seek(0, 0)
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error while writing to the file: %v", err)
		}
	}
	writer.Flush()

	return nil
}
