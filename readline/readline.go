package readline

import (
	"github.com/peterh/liner"
	"os"
	"path/filepath"
)

var (
	history_file = filepath.Join(os.TempDir(), ".mal_history")
	line         *liner.State
)

func init() {
	line = liner.NewLiner()
	line.SetCtrlCAborts(true)
	// load history from file
	if f, err := os.Open(history_file); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}

func Close() {
	// before closing, write history back to file
	if f, err := os.Create(history_file); err == nil {
		line.WriteHistory(f)
	}
	line.Close()
}

func PromptAndRead(prompt string) (string, error) {
	input, err := line.Prompt(prompt)
	if err != nil {
		return "", err
	}
	line.AppendHistory(input)
	return input, nil
}
