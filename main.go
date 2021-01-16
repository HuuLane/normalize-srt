package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

type Caption struct {
	Timecodes string
	Text      string
}

func (c Caption) String() string {
	return fmt.Sprintf("%s\n%s\n", c.Timecodes, c.Text)
}

func Normalize(filepath string) {
	f, err := os.Open(filepath)
	Must(err)
	defer f.Close()

	captions := make([]Caption, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "-->") {
			// read next line
			scanner.Scan()
			captions = append(captions, Caption{
				Timecodes: line,
				Text:      scanner.Text(),
			})
		}
	}

	// overwrite file
	f, err = os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC, 0666)
	w := bufio.NewWriter(f)
	for i, caption := range captions {
		_, err := w.WriteString(
			fmt.Sprintf("%d\n%s\n", i+1, caption),
		)
		_, err = w.WriteString("\n")
		Must(err)
	}
	err = w.Flush()
	Must(err)
}

func WalkDir(ch chan<- string, done chan<- struct{}, dirpath string) {
	err := filepath.Walk(dirpath,
		func(p string, info os.FileInfo, err error) error {
			Must(err)
			if path.Ext(info.Name()) == ".srt" {
				ch <- p
			}
			return nil
		})
	Must(err)

	done <- struct{}{}
}

func main() {
	dirpath := os.Args[1]

	ch := make(chan string)
	done := make(chan struct{})
	defer close(ch)
	defer close(done)
	go WalkDir(ch, done, dirpath)

	i := 1

loop:
	for {
		select {
		case p := <-ch:
			fmt.Printf("Progression: %d\n", i)
			i++
			go Normalize(p)
		case <-done:
			break loop
		}
	}
}
