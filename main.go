package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	OutputDir string `short:"o" long:"output-dir" default:"./" description:"出力するディレクトリ"`
	Quiet     bool   `short:"q" long:"quiet" description:"完了時にFinderを開かない。"`
}

func fatalln(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func countLines(str string) int {
	return strings.Count(str, "\n")
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "stamp-dl"
	parser.Usage = "ここにLINEスタンプページのURLを入れて！"
	args, err := parser.Parse()
	if err != nil {
		if countLines(err.Error()) < 1 {
			parser.WriteHelp(os.Stdout)
		}
		os.Exit(2)
	}

	if len(args) < 1 {
		parser.WriteHelp(os.Stdout)
		os.Exit(2)
	}

	ss, err := FetchStamps(args)
	if err != nil {
		fatalln(err)
	}

	absPath, err := filepath.Abs(opts.OutputDir)
	if err != nil {
		fatalln(err)
	}

	for _, s := range ss {
		err = s.Store(absPath)
		if err != nil {
			fatalln(err)
		}
	}

	if opts.Quiet {
		return
	}

	if err := exec.Command("open", absPath).Run(); err != nil {
		fatalln(err)
	}
}
