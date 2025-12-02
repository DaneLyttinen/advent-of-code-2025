package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isValid(num int) bool {
	s := strconv.Itoa(num)
	n := len(s)
	return strings.Contains((s + s)[1 : 2*n-1], s)
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		f, err = os.Open("day02/input.txt")
		if err != nil {
			fmt.Println("Error opening input.txt:", err)
			return
		}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	
	totalInvalidSum := 0

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		ranges := strings.Split(line, ",")
		for _, r := range ranges {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				continue
			}
			
			start, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
			end, _ := strconv.Atoi(strings.TrimSpace(parts[1]))

			for i := start; i <= end; i++ {
				if isValid(i) {
					totalInvalidSum += i
				}
			}
		}
	}

	fmt.Println(totalInvalidSum)
}
