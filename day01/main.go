package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		f, err = os.Open("day01/input.txt")
		if err != nil {
			fmt.Println("Error opening input.txt:", err)
			return
		}
	}
	defer f.Close()

	curr_val := 50
	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		direction := string(line[0])
		turns, _ := strconv.Atoi(line[1:])

		count += turns / 100

		remainder := turns % 100
		for i := 0; i < remainder; i++ {
			if direction == "L" {
				curr_val--
				if curr_val < 0 {
					curr_val = 99
				}
			} else { // R
				curr_val++
				if curr_val > 99 {
					curr_val = 0
				}
			}

			if curr_val == 0 {
				count++
			}
		}
	}
	fmt.Println(count)
}
