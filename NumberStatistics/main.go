// 题目A3：众数问题
// 给定含有n个元素的多重集合S，每个元素在S中出现的次数称为该元素的重
// 数，S中重数最大的元素称为众数。例如，S＝{1, 2 ,2 ,2 ,3 ,5}，S的众数是2，该
// 众数的重数为3。要求对于给定的由n个自然数组成的多重集合S，计算S的众数及
// 其重数。
//
// Usage:
// ./NumberStatistics -h //查看帮助
// ./NumberStatistics -n [1,2,4,3,2,4,4,5]
// Result: The mode is: [4], and the number of occurrences is: 3.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
)

var (
	// num 存放多重集合S
	num []uint32
	// mode 存放众数集合（可能存在多个众数）
	mode []uint32
	// frequency key-value形式存放每个数的重数
	frequency = make(map[uint32]uint32)
	// frequencySort 对frequency的value排序
	frequencySort []uint32

	numRead     string
	help        bool
	defaultNums = "[1,2,4,3,2,4,4,5]"
)

func main() {
	flag.StringVar(&numRead, "n", defaultNums, "input `nums`")
	flag.BoolVar(&help, "h", false, "help")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	err := json.Unmarshal([]byte(numRead), &num)
	if err != nil {
		fmt.Println("error format:", err.Error())
		return
	}

	for _, value := range num {
		frequency[value]++
	}

	// 将frequency的value，即每个数出现的次数，使用slice按从小到大排序后判定众数的重数
	for _, value := range frequency {
		frequencySort = append(frequencySort, value)
	}
	bubbleSort(frequencySort)
	modeFrequency := frequencySort[len(frequencySort)-1]

	// 根据重数遍历寻找众数
	for key, value := range frequency {
		if value == modeFrequency {
			mode = append(mode, key)
		}
	}

	fmt.Printf("The mode is: \x1b[31m%d\x1b[0m, and the number of occurrences is: \x1b[31m%d\x1b[0m.\n", mode, modeFrequency)
}

// bubbleSort 冒泡排序
func bubbleSort(data []uint32) {
	for i := 0; i < len(data)-1; i++ {
		for j := 0; j < len(data)-i-1; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}
