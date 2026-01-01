package main

import (
	"fmt"
	"os"
	"sugu/repl"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		// 引数なしの場合はREPLを起動
		repl.Start(os.Stdin, os.Stdout)
	} else if len(args) == 1 {
		// 引数がある場合はファイルを実行
		filename := args[0]
		if err := repl.RunFile(filename, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Usage: sugu [filename]")
		os.Exit(1)
	}
}
