package main

import (
	"bufio"
	"fmt"
	"os"
	"printer"
	"reader"
	"strings"
	"types"
)

func READ(str string) (types.MalType, error) {
	return reader.ReadStr(str)
}

func EVAL(ast types.MalType, env string) (types.MalType, error) {
	return ast, nil
}

func PRINT(exp types.MalType) (string, error) {
	return printer.PrintStr(exp, true), nil
}

func rep(str string) (string, error) {
	ast, err := READ(str)
	if err != nil {
		return "", err
	}
	exp, err := EVAL(ast, "")
	if err != nil {
		return "", err
	}
	return PRINT(exp)
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("user> ")
		if in.Scan() {
			read := strings.TrimSpace(in.Text())
			result, err := rep(read)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
		} else {
			err := in.Err()
			if err != nil {
				panic(err)
			}
			return
		}
	}
}
