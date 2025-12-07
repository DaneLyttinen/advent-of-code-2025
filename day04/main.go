package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	f, _ := os.Open("input.txt")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	allLines := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		allLines = append(allLines, line)
	}

	direction := [][]int{
		{1, 0},
		{0, 1},
		{1, 1},
		{1, -1},
		{-1, 0},
		{0, -1},
		{-1, -1},
		{-1, 1},
	}

	totalSum := 0
	for{
		removed := false
		for row := 0; row < len(allLines); row++ {
			for col := 0; col < len(allLines[row]); col++ {
				if allLines[row][col] != '@' {
					continue
				}
				
				surroundingSum := 0
				for _, dir := range direction {
					rowOffset := dir[0]
					colOffset := dir[1]
					rowIdx := row + rowOffset
					colIdx := col + colOffset
					if rowIdx < 0 || rowIdx >= len(allLines) || colIdx < 0 || colIdx >= len(allLines[rowIdx]) {
						continue
					}
					if allLines[rowIdx][colIdx] == '@' {
						surroundingSum++

					}
				}
				if surroundingSum < 4 {
					allLines[row] = allLines[row][:col] + "x" + allLines[row][col+1:]
					removed = true
					totalSum++
				}
			}
		}
		if !removed {
			break
		}
	}
	
	fmt.Println(totalSum)
}