package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	dirName   string
	msgSuffix string
	loc       *time.Location
	err       error
	lg        *log.Logger    = log.New(os.Stderr, "", log.Lshortfile)
	prefixRE  *regexp.Regexp = regexp.MustCompile("([[:digit:]]{12})(.*)")
)

const (
	datePrefix = "Date: "
)

func init() {
	var tz string
	flag.StringVar(&dirName, "dir", ".", "A directory name to scan")
	flag.StringVar(&tz, "loc", "America/Chicago", "A location name from the IANA Time Zone database")
	flag.StringVar(&msgSuffix, "ext", ".eml", "A file extension (including dot) to be recognised as a message file")
	flag.Parse()
	loc, err = time.LoadLocation(tz)
	if err != nil {
		lg.Fatal(err)
	}
}

func stampInFile(fn string) (string, error) {

	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := scanner.Text(); strings.HasPrefix(line, datePrefix) {
			t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", strings.TrimPrefix(line, datePrefix))
			if err != nil {
				return "", err
			}
			return t.In(loc).Format("060102150405"), nil
		}
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("Timestamp not found in file: %q", fn)
}

type trans struct {
	tsNew string
	files []string
}

type message_heap map[string]*trans

func scanMessages(dir string) message_heap {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		lg.Fatal(err)
	}

	messages := make(map[string]*trans)

	for _, file := range files {

		if !file.Mode().IsRegular() {
			continue
		}

		fn := file.Name()
		parts := prefixRE.FindStringSubmatch(fn)
		if parts == nil {
			lg.Printf("Skipping %q\n", fn)
			continue
		}

		if msgInfo, ok := messages[parts[1]]; ok {
			msgInfo.files = append(msgInfo.files, parts[2])
		} else {
			msgInfo = &trans{"", []string{parts[2]}}
			messages[parts[1]] = msgInfo
		}

		if strings.HasSuffix(fn, msgSuffix) {
			stamp, err := stampInFile(fn)
			if err != nil {
				lg.Fatal(err)
			}
			messages[parts[1]].tsNew = stamp
		}
	}
	return messages
}

func processMessages(msgs message_heap) {
	for tsOld, msg := range msgs {
		for _, file := range msg.files {
			if tsOld != msg.tsNew {
				oldName := tsOld + file
				newName := msg.tsNew + file
				os.Rename(oldName, newName)
				lg.Printf("%s => %s", oldName, newName)
			}
		}
	}
}

func main() {
	processMessages(scanMessages(dirName))
}
