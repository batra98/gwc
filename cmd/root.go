/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var bytecount bool
var linecount bool
var wordcount bool
var charcount bool
var noflag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gwc",
	Short: "wc written in golang",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		if bytecount == false && linecount == false && wordcount == false && charcount == false {
			noflag = true
		}

		var f *os.File
		var err error

		if len(args) > 0 {
			filepath := args[0]
			f, err = os.Open(filepath)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			defer f.Close()
			defer fmt.Println(filepath)
		} else {
			f = os.Stdin
		}

		buf := make([]byte, 32*1024)
		num_chars := 0
		num_lines := 0
		num_words := 0
		num_bytes := 0

		for {
			c, err := f.Read(buf)

			if linecount || noflag {
				getLines(buf, c, &num_lines)
			}
			if wordcount || noflag {
				getWords(buf, c, &num_words)
			}
			if bytecount || noflag {
				getBytes(buf, c, &num_bytes)
			}

			if charcount {
				getChars(buf, c, &num_chars)
			}

			switch {
			case err == io.EOF:
				printResult(&num_lines, &num_words, &num_bytes, &num_chars)
				return

			case err != nil:
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

	},
}

func printResult(num_lines *int, num_words *int, num_bytes *int, num_chars *int) {
	if linecount || noflag {
		fmt.Printf("%d ", *num_lines)
	}
	if wordcount || noflag {
		fmt.Printf("%d ", *num_words)
	}
	if bytecount || noflag {
		fmt.Printf("%d ", *num_bytes)
	}
	if charcount {
		fmt.Printf("%d ", *num_chars)
	}

}

func getChars(buf []byte, c int, count *int) {
	*count += utf8.RuneCount(buf[:c])
}

func getLines(buf []byte, c int, count *int) {
	lineSep := []byte{'\n'}
	*count += bytes.Count(buf[:c], lineSep)

}
func getBytes(buf []byte, c int, count *int) {
	*count += c
}

func getWords(buf []byte, c int, count *int) {

	text := string(buf[:c])
	words := strings.Fields(text)
	*count += len(words)

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	err = doc.GenMarkdownTree(rootCmd, "docs/")
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gwc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	noflag = false
	rootCmd.Flags().BoolVarP(&bytecount, "bytes", "c", false, "The number of bytes in each input file is written to the standard output")
	rootCmd.Flags().BoolVarP(&linecount, "lines", "l", false, "The number of lines in each input file is written to the standard output.")
	rootCmd.Flags().BoolVarP(&wordcount, "words", "w", false, "The number of words in each input file is written to the standard output.")
	rootCmd.Flags().BoolVarP(&charcount, "chars", "m", false, "The number of chars in each input file is written to the standard output.")
	rootCmd.CompletionOptions.DisableDefaultCmd = false

}
