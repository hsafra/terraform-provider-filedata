/*
This file has api functions for working with the data file structure we define.
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
lines numbering starts from 1.
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
lines numbering starts from 1.
*/
func WriteLine(filePath string, n int, text string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
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
		for i := len(lines); i < n; i++ {
			lines = append(lines, "")
		}
	}

	lines[n-1] = text

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("error while writing to the file: %v", err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error while writing to the file: %v", err)
	}
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error while writing to the file: %v", err)
		}
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error while writing to the file: %v", err)
	}

	return nil
}

/*
LineCount returns the number of lines in the file.
*/
func LineCount(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("unable to open the file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}

	if scanner.Err() != nil {
		return 0, fmt.Errorf("error while reading the file: %v", scanner.Err())
	}

	return lines, nil
}

/*
RemoveFile removes the file at the given path.
*/
func RemoveFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("unable to remove the file: %v", err)
	}
	return nil
}

/*
TrimFile removes all lines after the n-th line from the file.
*/
func TrimFile(filePath string, n int) error {
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
		return fmt.Errorf("file doesn't contain %d lines", n)
	}

	lines = lines[:n]

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("error while writing to the file: %v", err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error while writing to the file: %v", err)
	}
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error while writing to the file: %v", err)
		}
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error while writing to the file: %v", err)
	}
	return nil
}
