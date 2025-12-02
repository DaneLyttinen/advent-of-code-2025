package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	ranges, err := readRanges("input.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println("=== EXTREME OPTIMIZATION BENCHMARK ===\n")

	// Baseline: Current Math Method
	startMath := time.Now()
	sumMath := 0
	for _, r := range ranges {
		sumMath += sumValidInRangeMath(r.start, r.end)
	}
	durationMath := time.Since(startMath)

	// Extreme 2: Batch Processing
	startBatch := time.Now()
	sumBatch := batchProcessRanges(ranges)
	durationBatch := time.Since(startBatch)

	fmt.Printf("Math Method:          Sum=%d, Time=%v\n", sumMath, durationMath)
	fmt.Printf("Batch Method:         Sum=%d, Time=%v\n", sumBatch, durationBatch)

	fmt.Println()
	if durationBatch > 0 && durationMath > 0 {
		fmt.Printf("Speedup (Batch vs Math):   %.2fx\n", float64(durationMath)/float64(durationBatch))
	}
}

type rangePair struct {
	start, end int
}

func readRanges(filename string) ([]rangePair, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var result []rangePair
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		for _, r := range parts {
			p := strings.Split(r, "-")
			if len(p) != 2 {
				continue
			}
			start, _ := strconv.Atoi(strings.TrimSpace(p[0]))
			end, _ := strconv.Atoi(strings.TrimSpace(p[1]))
			result = append(result, rangePair{start, end})
		}
	}
	return result, scanner.Err()
}

// Baseline Math Method
func sumValidInRangeMath(start, end int) int {
	minLen := numDigits(start)
	maxLen := numDigits(end)

	sum := 0
	seen := make(map[int]bool, 64)

	pow10 := [11]int{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000}

	for length := minLen; length <= maxLen; length++ {
		for p := 1; p < length; p++ {
			if length%p != 0 {
				continue
			}

			repeats := length / p
			powP := pow10[p]
			K := (intPow(powP, repeats) - 1) / (powP - 1)

			minBase := pow10[p-1]
			maxBase := pow10[p] - 1

			low := (start + K - 1) / K
			high := end / K

			if low < minBase {
				low = minBase
			}
			if high > maxBase {
				high = maxBase
			}

			for base := low; base <= high; base++ {
				val := base * K
				if !seen[val] {
					seen[val] = true
					sum += val
				}
			}
		}
	}

	return sum
}
// Inline bitmap-based deduplication
//go:inline
func addToSum(sum *int, val int, bitmap *[16384]uint64, hasLarge *bool, largeMap map[int]struct{}) {
	if val < 1048576 { // Fits in bitmap
		idx := val >> 6       // val / 64
		bit := uint64(1) << (val & 63) // val % 64
		if bitmap[idx]&bit == 0 {
			bitmap[idx] |= bit
			*sum += val
		}
	} else {
		if !*hasLarge {
			*hasLarge = true
		}
		if _, exists := largeMap[val]; !exists {
			largeMap[val] = struct{}{}
			*sum += val
		}
	}
}

// EXTREME 2: Batch Processing All Ranges Together
func batchProcessRanges(ranges []rangePair) int {
	// Find global min/max
	globalMin := ranges[0].start
	globalMax := ranges[0].end
	for _, r := range ranges[1:] {
		if r.start < globalMin {
			globalMin = r.start
		}
		if r.end > globalMax {
			globalMax = r.end
		}
	}

	minLen := numDigits(globalMin)
	maxLen := numDigits(globalMax)

	var bitmap [16384]uint64
	hasLarge := false
	largeMap := make(map[int]struct{}, 16)
	sum := 0

	// Process all (length, period) combinations once
	for length := minLen; length <= maxLen; length++ {
		for p := 1; p < length; p++ {
			if length%p != 0 {
				continue
			}

			repeats := length / p
			powP := intPowFast(10, p)
			K := (intPowFast(powP, repeats) - 1) / (powP - 1)

			minBase := intPowFast(10, p-1)
			maxBase := intPowFast(10, p) - 1

			// For each range, compute intersection and add
			for _, r := range ranges {
				low := (r.start + K - 1) / K
				high := r.end / K

				if low < minBase {
					low = minBase
				}
				if high > maxBase {
					high = maxBase
				}

				if low > high {
					continue
				}

				for base := low; base <= high; base++ {
					val := base * K
					addToSum(&sum, val, &bitmap, &hasLarge, largeMap)
				}
			}
		}
	}

	return sum
}

// Fast integer power using bit manipulation
func intPowFast(base, exp int) int {
	if exp == 0 {
		return 1
	}
	result := 1
	for exp > 0 {
		if exp&1 == 1 {
			result *= base
		}
		base *= base
		exp >>= 1
	}
	return result
}

func intPow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

func numDigits(n int) int {
	if n < 10 {
		return 1
	}
	if n < 100 {
		return 2
	}
	if n < 1000 {
		return 3
	}
	if n < 10000 {
		return 4
	}
	if n < 100000 {
		return 5
	}
	if n < 1000000 {
		return 6
	}
	if n < 10000000 {
		return 7
	}
	if n < 100000000 {
		return 8
	}
	if n < 1000000000 {
		return 9
	}
	return 10
}

