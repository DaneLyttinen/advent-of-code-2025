package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func findMaxKDigits(line string, k int) int {
	n := len(line)

	result := 0
	currentIdx := 0

	for i := 0; i < k; i++ {
		remaining := k - i - 1
		searchEnd := n - remaining

		maxDigit := 0
		maxPos := currentIdx
		
		for j := currentIdx; j < searchEnd; j++ {
			digit := int(line[j] - '0')
			if digit > maxDigit {
				maxDigit = digit
				maxPos = j
			}
		}
		
		result = result*10 + maxDigit
		currentIdx = maxPos + 1
	}
	
	return result
}

func main() {
	f, _ := os.Open("input.txt")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	totalSum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		maxJoltage := findMaxKDigits(line, 12)
		totalSum += maxJoltage
	}
	
	fmt.Println(totalSum)
}