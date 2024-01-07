package main

import (
	"math"
	"math/rand"
	"sort"
	"strings"
	"unicode/utf8"
)

func hide(srcSlice []string, percent int) []string {

	var cnt int
	var numSlice []int
	var flag int
	var resSlice []string
	var nCounter int
	var tempSlice []string

	for _, v := range srcSlice {
		cnt += len(strings.Fields(v))
	}

	nCntPercent := math.Floor(float64(cnt * percent / 100))

	for {
		if len(numSlice) >= int(nCntPercent) {
			break
		}
		flag = 1

		n := rand.Intn(cnt) + 1
		for _, v := range numSlice {
			if v == n {
				flag = 0
				break
			}
		}

		if flag == 0 || n == 0 {
			continue
		} else {
			numSlice = append(numSlice, n)
		}
	}

	sort.Ints(numSlice)

	for _, v := range srcSlice {
		for _, vv := range strings.Fields(v) {
			nCounter += 1
			flag = 0

			for _, vCnt := range numSlice {
				if nCounter == vCnt {
					flag = 1
					break
				}

				if vCnt > nCounter {
					flag = 0
					break
				}
			}

			if flag == 0 {
				tempSlice = append(tempSlice, vv)
			} else {
				tempSlice = append(tempSlice, strings.Repeat("*", utf8.RuneCountInString(vv)))
			}
		}

		resSlice = append(resSlice, strings.Join(tempSlice, " "))

		tempSlice = nil
	}

	return resSlice
}
