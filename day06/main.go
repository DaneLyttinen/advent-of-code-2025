package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	f, _ := os.Open("input.txt")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	rows := []string{}
	
	for scanner.Scan() {
		rows = append(rows, scanner.Text())
	}
	
	numRows := rows[:4]
	opRow := rows[4]
	maxWidth := len(opRow)
	
	grandTotal := 0
	problemNums := []int{}
	var problemOp rune
	
	for col := maxWidth - 1; col >= 0; col-- {
		isAllSpaces := true
		for row := 0; row < 5; row++ {
			if col < len(rows[row]) && rows[row][col] != ' ' {
				isAllSpaces = false
				break
			}
		}
		
		if isAllSpaces && len(problemNums) > 0 {
			result := problemNums[0]
			for i := 1; i < len(problemNums); i++ {
				if problemOp == '+' {
					result += problemNums[i]
				} else {
					result *= problemNums[i]
				}
			}
			grandTotal += result
			problemNums = []int{}
			continue
		}
		
		opChar := rune(opRow[col])
		
		numStr := ""
		for row := 0; row < 4; row++ {
			if col < len(numRows[row]) && numRows[row][col] != ' ' {
				numStr += string(numRows[row][col])
			}
		}
		
		if numStr != "" {
			num, _ := strconv.Atoi(numStr)
			problemNums = append(problemNums, num)
		}
		
		if opChar == '+' || opChar == '*' {
			problemOp = opChar
		}
	}
	

	if len(problemNums) > 0 {
		result := problemNums[0]
		for i := 1; i < len(problemNums); i++ {
			if problemOp == '+' {
				result += problemNums[i]
			} else {
				result *= problemNums[i]
			}
		}
		grandTotal += result
	}
	
	fmt.Println(grandTotal)
}
