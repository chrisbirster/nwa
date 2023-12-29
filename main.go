package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const dailyNoteFilename = `\d{4}-\d{2}-\d{2}.md$`
const weeklyNoteFilename = `\d{4}-W\d{2}.md$`

type Note string

// / Notes With Attitude
type NWA struct {
	noteDir string
	notes   []Note
	dailies []Note
	weeklies []Note
}

func newNWAWithDir(noteDir string) *NWA {
	_, err := os.Stat(noteDir)
	if err != nil {
		panic(errors.New("Could not get file info for note directory"))
	}
	return &NWA{
		noteDir: noteDir,
	}
}

func (n NWA) compileDailyRegex() *regexp.Regexp {
	dailyNoteFullpath := filepath.Join(n.noteDir, "dailies", dailyNoteFilename)
	dailyRegex, err := regexp.Compile(dailyNoteFullpath)
	if err != nil {
		panic(errors.New("Could not compile regex for daily note: " + dailyNoteFullpath))
	}
	return dailyRegex
}

func (n NWA) compileWeeklyRegex() *regexp.Regexp {
	weeklyNoteFullpath := filepath.Join(n.noteDir, "weeklies", weeklyNoteFilename)
	weeklyRegex, err := regexp.Compile(weeklyNoteFullpath)
	if err != nil {
		panic(errors.New("Could not compile regex for weekly note: " + weeklyNoteFullpath))
	}
	return weeklyRegex
}

func (n *NWA) initNotes() error {
	dailyRegex := n.compileDailyRegex()
	weeklyRegex := n.compileWeeklyRegex()
	return filepath.Walk(n.noteDir, func(path string, info os.FileInfo, err error) error {
		// add to daily cache
		if dailyRegex.MatchString(path) {
			n.dailies = append(n.dailies, Note(path))
			return nil
		}

		// add to weekly cache
		if weeklyRegex.MatchString(path) {
			n.weeklies = append(n.weeklies, Note(path))
			return nil
		}

		// add to note cache
		if strings.Contains(path, ".md") {
			n.notes = append(n.notes, Note(path))
			return nil
		}
		return nil
	})
}

func (n *NWA) init() {
	err := n.initNotes()
	if err != nil {
		panic(err)
	}
}

func (n NWA) readNote(noteName string) string {
	return "called readNote"
}

func newNWA() *NWA {
	defaultNoteDir, err := os.UserHomeDir()
	if err != nil {
		panic(errors.New("Could not get user home directory"))
	}

	return &NWA{
		noteDir: path.Join(defaultNoteDir, "notes"),
	}
}

func main() {
	nwa := newNWAWithDir("/Users/gm/note-vault")
	nwa.init()
	for _, daily := range nwa.dailies {
		fmt.Println(daily)
	}
}
