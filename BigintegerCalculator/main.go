// 大整数算术运算器
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"regexp"
	"strings"
)

var (
	expr        string
	help        bool
	exprSlice   []string
	resultSlice []string
	outputSlice []string
	inputFile   string
	outputFile  string
	sqrtNum     string
)

func main() {
	flag.StringVar(&expr, "e", "", "input arithmetic `expression`")
	flag.StringVar(&inputFile, "f", "", "`file` to read by line")
	flag.StringVar(&outputFile, "o", "", "`file` to output")
	flag.StringVar(&sqrtNum, "sqrt", "", "`num` to sqrt")
	flag.BoolVar(&help, "h", false, "help")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if expr == "" && inputFile == "" && sqrtNum == "" {
		fmt.Println("\x1b[31m[!] Please use -h to view help, like `./BigintegerCalculator -h`\x1b[0m")
		return
	}

	if expr != "" && inputFile != "" {
		fmt.Println("\x1b[31m[!] too many parameters\x1b[0m")
		return
	}

	if expr != "" {
		// 限制输入格式
		ok := checkExpr(expr)
		if !ok {
			fmt.Println("\x1b[31m[!] invalid format\x1b[0m")
			return
		}
		result := calc(expr)
		fmt.Printf("\x1b[32m[*]\x1b[0m %s = %s\n", expr, result)
		return
	}

	if inputFile != "" {
		exprSlice = readLine(inputFile)
		for _, v := range exprSlice {
			resultSlice = append(resultSlice, calc(v))
		}

		for k, v := range exprSlice {
			for i, j := range resultSlice {
				if k == i {
					outputSlice = append(outputSlice, fmt.Sprintf("%s = %s", v, j))
				}
			}
		}

		fmt.Println("\x1b[32m[*] The calculation results of batch processing are as follows:\x1b[0m")
		for _, v := range outputSlice {
			fmt.Printf("\x1b[32m[*]\x1b[0m %s\n", v)
		}

		if outputFile != "" {
			writeLine(outputFile, outputSlice)
			fmt.Printf("\x1b[32m[*] The calculation result is saved to %s successfully!\x1b[0m\n", outputFile)
			return
		}
	} else if outputFile != "" {
		fmt.Println("\x1b[31m[!] require -f\x1b[0m")
		return
	}

	if sqrtNum != "" {
		pattern := `^[\d+]*$`
		reg := regexp.MustCompile(pattern)
		if reg == nil {
			fmt.Println("\x1b[31m[!] regexp err\x1b[0m")
			return
		}
		isNum := reg.MatchString(sqrtNum)
		if !isNum {
			fmt.Println("\x1b[31m[!] invalid format\x1b[0m")
			return
		}
		big, ok := new(big.Int).SetString(sqrtNum, 10)
		if !ok {
			fmt.Println("\x1b[31m[!] invalid format\x1b[0m")
			return
		}
		big.Sqrt(big)
		fmt.Printf("\x1b[32m[*]\x1b[0m sqrt(%s) = %s\n", sqrtNum, big)
	}

}

// calc 计算一个表达式，传入表达式，返回计算结果
func calc(expr string) string {
	var result big.Int
	// 获取两个操作数
	bigNum := regexpBigNum(expr)

	big1, ok1 := new(big.Int).SetString(bigNum[0][0], 10)
	big2, ok2 := new(big.Int).SetString(bigNum[1][0], 10)
	if !ok1 || !ok2 {
		fmt.Println("\x1b[31m[!] invalid format\x1b[0m")
		os.Exit(1)
	}

	operator := regexpOperator(expr)
	switch {
	case operator == "+":
		result.Add(big1, big2)
	// 因为正则匹配到的数带有负号，所以操作符为负号的也是相加
	case operator == "-":
		result.Add(big1, big2)
	case operator == "*":
		result.Mul(big1, big2)
	case operator == "/":
		result.Div(big1, big2)
	}
	return result.String()
}

// regexpBigNum 利用正则匹配处理运算表达式，获取两个大整数，返回匹配到操作数的二维数组
func regexpBigNum(expr string) [][]string {
	pattern := `[\-]?([0-9]{1,}\.?[0-9]*)`
	reg := regexp.MustCompile(pattern)
	if reg == nil {
		fmt.Println("\x1b[31m[!] regexp err\x1b[0m")
		os.Exit(1)
	}
	bigNumArr := reg.FindAllStringSubmatch(expr, 2)
	return bigNumArr
}

// regexpOperator 利用正则匹配获取两个数之间的运算符
func regexpOperator(expr string) string {
	pattern := `([+\-*/])`
	reg := regexp.MustCompile(pattern)
	if reg == nil {
		fmt.Println("\x1b[31m[!] regexp err\x1b[0m")
		os.Exit(1)
	}
	operator := reg.FindAllStringSubmatch(expr, -1)
	if strings.HasPrefix(expr, "-") {
		return operator[1][0]
	} else {
		return operator[0][0]
	}

}

// checkExpr 对输入的算术表达式限制格式，正确返回true，错误返回false
func checkExpr(expr string) bool {
	// 判断传入表达式是否为空
	if expr == "" {
		return false
	}
	// 判断操作数是否纯数字
	exprSub := strings.Replace(expr, "+", "", -1)
	exprSub = strings.Replace(exprSub, "-", "", -1)
	exprSub = strings.Replace(exprSub, "*", "", -1)
	exprSub = strings.Replace(exprSub, "/", "", -1)
	exprSub = strings.Replace(exprSub, "(", "", -1)
	exprSub = strings.Replace(exprSub, ")", "", -1)
	pattern := `^[\d+]*$`
	reg := regexp.MustCompile(pattern)
	if reg == nil {
		fmt.Println("\x1b[31m[!] regexp err\x1b[0m")
		os.Exit(1)
	}
	isNum := reg.MatchString(exprSub)
	if !isNum {
		return false
	}
	// 判断操作符数量是否正确
	if (strings.Count(expr, "-") > 2) || (strings.Count(expr, "+")+strings.Count(expr, "*")+strings.Count(expr, "/") > 1) {
		return false
	}
	// 判断操作符是否格式正确
	operator := regexpOperator(expr)
	if len(operator) > 1 {
		return false
	}
	// 判断括号是否匹配
	if strings.Contains(expr, "(") && !strings.Contains(expr, ")") {
		return false
	}
	// 判断括号数量是否正确
	if strings.Contains(expr, "(") && strings.Contains(expr, ")") {
		if (strings.Count(expr, "(")+strings.Count(expr, ")"))%2 != 0 {
			return false
		}
		if !strings.Contains(expr, "(-") {
			return false
		}
	}
	return true
}

// readLine 按行读取给定文件
func readLine(filename string) []string {
	var slice []string
	fi, err := os.Open(filename)
	if err != nil {
		fmt.Printf("\x1b[31m[!] %s\x1b[0m", err.Error())
		os.Exit(1)
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		ok := checkExpr(string(a))
		if !ok {
			fmt.Printf("\x1b[31m[!] invalid format\n[*] Incorrect expression:\n[*]\x1b[0m %s", string(a))
			os.Exit(1)
		}
		slice = append(slice, string(a))
	}
	return slice
}

// writeLine 按行写入给定文件
func writeLine(filename string, slice []string) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("\x1b[31m[!] %s\x1b[0m", err.Error())
		os.Exit(1)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, v := range slice {
		fmt.Fprintln(w, v)
	}
	w.Flush()
}
