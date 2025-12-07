package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	f, _ := os.Open("input.txt")
	defer f.Close()

	scanner := bufio.NewScanner(f)

	ranges := [][]int{}
	totalSum := 0
	isRange := true

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if len(line) == 0 {
			isRange = false
			continue
		}
		
		if isRange {
			rangeParts := strings.Split(line, "-")
			if len(rangeParts) == 2 {
				start, _ := strconv.Atoi(rangeParts[0])
				end, _ := strconv.Atoi(rangeParts[1])
				ranges = append(ranges, []int{start, end})
			}
		} else {
			// number, _ := strconv.Atoi(line)
			// isFresh := false
			// for _, r := range ranges {
			// 	if number >= r[0] && number <= r[1] {
			// 		isFresh = true
			// 		break
			// 	}
			// }
			// if isFresh {
			// 	totalSum++
			// }
		}
	}
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][0] < ranges[j][0]
	})
	
	currentStart := ranges[0][0]
	currentEnd := ranges[0][1]
	
	for i := 1; i < len(ranges); i++ {		
		r := ranges[i]
		if r[0] <= currentEnd+1 {
			currentEnd = max(currentEnd, r[1])
		} else {
			totalSum += currentEnd - currentStart + 1
			currentStart = r[0]
			currentEnd = r[1]
		}
	}
	
	totalSum += currentEnd - currentStart + 1
	
	fmt.Println(totalSum)
}