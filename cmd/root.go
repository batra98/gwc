/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var bytecount bool
var linecount bool
var wordcount bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gwc",
	Short: "wc written in golang",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		filepath := args[0]

		f, err := os.Open(filepath)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer f.Close()
		if linecount == true {
			num_lines, err := getLines(f)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Printf("%d ", num_lines)
		}
		if wordcount {
			num_words, err := getWords(f)

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Printf("%d ", num_words)

		}

		if bytecount == true {
			num_bytes, err := getBytes(f)

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Printf("%d ", num_bytes)
		}

		fmt.Println(filepath)
	},
}

func getLines(f *os.File) (int, error) {
	f.Seek(0, 0)
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := f.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}

}
func getBytes(f *os.File) (int64, error) {

	fstat, err := f.Stat()

	if err != nil {
		return 0, err
	}

	return fstat.Size(), nil

}

func getWords(f *os.File) (int, error) {
	f.Seek(0, 0)

	data, err := io.ReadAll(f)
	if err != nil {
		return 0, err
	}

	text := string(data)
	words := strings.Fields(text)
	wordCount := len(words)

	return wordCount, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gwc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&bytecount, "bytes", "c", false, "The number of bytes in each input file is written to the standard output")
	rootCmd.Flags().BoolVarP(&linecount, "lines", "l", false, "The number of lines in each input file is written to the standard output.")
	rootCmd.Flags().BoolVarP(&wordcount, "words", "w", false, "The number of words in each input file is written to the standard output.")
}
