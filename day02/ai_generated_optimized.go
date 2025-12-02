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

	// Extreme 1: Zero-Allocation Direct Arithmetic
	startExtreme := time.Now()
	sumExtreme := 0
	for _, r := range ranges {
		sumExtreme += sumValidExtreme(r.start, r.end)
	}
	durationExtreme := time.Since(startExtreme)

	// Extreme 2: Batch Processing
	startBatch := time.Now()
	sumBatch := batchProcessRanges(ranges)
	durationBatch := time.Since(startBatch)

	fmt.Printf("Math Method:          Sum=%d, Time=%v\n", sumMath, durationMath)
	fmt.Printf("Extreme Method:       Sum=%d, Time=%v\n", sumExtreme, durationExtreme)
	fmt.Printf("Batch Method:         Sum=%d, Time=%v\n", sumBatch, durationBatch)

	fmt.Println()
	if durationExtreme > 0 && durationMath > 0 {
		fmt.Printf("Speedup (Extreme vs Math): %.2fx\n", float64(durationMath)/float64(durationExtreme))
	}
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

// EXTREME 1: Zero-Allocation Direct Arithmetic with Inline Dedup
func sumValidExtreme(start, end int) int {
	minLen := numDigits(start)
	maxLen := numDigits(end)

	// Pre-allocated bitmap for deduplication (covers most cases)
	// Using fixed-size array instead of map for cache locality
	var bitmap [16384]uint64 // Covers values up to 16384*64 = 1,048,576
	hasLarge := false
	largeMap := make(map[int]struct{}, 8)

	sum := 0

	// Pre-compute all powers of 10 up to 10^10
	pow10_0 := 1
	pow10_1 := 10
	pow10_2 := 100
	pow10_3 := 1000
	pow10_4 := 10000
	pow10_5 := 100000
	pow10_6 := 1000000
	pow10_7 := 10000000
	pow10_8 := 100000000
	pow10_9 := 1000000000
	pow10_10 := 10000000000

	for length := minLen; length <= maxLen; length++ {
		// Unroll common cases for p (most patterns are p=1,2,3,5)
		
		// p=1 (repeating single digit: 11, 111, 1111, etc)
		if length >= 2 {
			K := (intPowFast(10, length) - 1) / 9 // 10^n - 1 / 9 = 111...1
			minBase := 1
			maxBase := 9

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
				addToSum(&sum, val, &bitmap, &hasLarge, largeMap)
			}
		}

		// p=2
		if length >= 4 && length%2 == 0 {
			repeats := length / 2
			K := (intPowFast(100, repeats) - 1) / 99
			minBase := 10
			maxBase := 99

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
				addToSum(&sum, val, &bitmap, &hasLarge, largeMap)
			}
		}

		// p=3
		if length >= 6 && length%3 == 0 {
			repeats := length / 3
			K := (intPowFast(1000, repeats) - 1) / 999
			minBase := 100
			maxBase := 999

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
				addToSum(&sum, val, &bitmap, &hasLarge, largeMap)
			}
		}

		// p=4
		if length >= 8 && length%4 == 0 {
			repeats := length / 4
			K := (intPowFast(10000, repeats) - 1) / 9999
			minBase := 1000
			maxBase := 9999

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
				addToSum(&sum, val, &bitmap, &hasLarge, largeMap)
			}
		}

		// p=5
		if length == 10 {
			K := (intPowFast(100000, 2) - 1) / 99999
			minBase := 10000
			maxBase := 99999

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
				addToSum(&sum, val, &bitmap, &hasLarge, largeMap)
			}
		}

		// Handle remaining periods (6,7,8,9) - rare cases
		for p := 6; p < length; p++ {
			if length%p != 0 {
				continue
			}

			repeats := length / p
			var powP int
			switch p {
			case 1:
				powP = pow10_1
			case 2:
				powP = pow10_2
			case 3:
				powP = pow10_3
			case 4:
				powP = pow10_4
			case 5:
				powP = pow10_5
			case 6:
				powP = pow10_6
			case 7:
				powP = pow10_7
			case 8:
				powP = pow10_8
			case 9:
				powP = pow10_9
			case 10:
				powP = pow10_10
			default:
				powP = intPowFast(10, p)
			}

			K := (intPowFast(powP, repeats) - 1) / (powP - 1)

			var minBase, maxBase int
			switch p {
			case 1:
				minBase = pow10_0
			case 2:
				minBase = pow10_1
			case 3:
				minBase = pow10_2
			case 4:
				minBase = pow10_3
			case 5:
				minBase = pow10_4
			case 6:
				minBase = pow10_5
			case 7:
				minBase = pow10_6
			case 8:
				minBase = pow10_7
			case 9:
				minBase = pow10_8
			case 10:
				minBase = pow10_9
			default:
				minBase = intPowFast(10, p-1)
			}

			switch p {
			case 1:
				maxBase = pow10_1 - 1
			case 2:
				maxBase = pow10_2 - 1
			case 3:
				maxBase = pow10_3 - 1
			case 4:
				maxBase = pow10_4 - 1
			case 5:
				maxBase = pow10_5 - 1
			case 6:
				maxBase = pow10_6 - 1
			case 7:
				maxBase = pow10_7 - 1
			case 8:
				maxBase = pow10_8 - 1
			case 9:
				maxBase = pow10_9 - 1
			case 10:
				maxBase = pow10_10 - 1
			default:
				maxBase = intPowFast(10, p) - 1
			}

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
				addToSum(&sum, val, &bitmap, &hasLarge, largeMap)
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

