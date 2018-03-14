package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"log"
)

func READ(str string) string {
	return str
}

func EVAL(str string, env string) string {
	return str
}

func PRINT(str string) string {
	return str
}

func rep(str string) string {
	return PRINT(EVAL(READ(str), ""))
}

func main() {
	in := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("user> ")
		if in.Scan() {
			read := strings.TrimSpace(in.Text())
			fmt.Println(rep(read))
		} else {
			err := in.Err()
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
	}
}
