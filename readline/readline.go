package readline

import (
	"github.com/peterh/liner"
	"os"
	"path/filepath"
)

var (
	historyFile = filepath.Join(os.TempDir(), ".mal_history")
	line        *liner.State
)

func init() {
	line = liner.NewLiner()
	line.SetCtrlCAborts(true)
	// load history from file
	if f, err := os.Open(historyFile); err == nil {
		_, _ = line.ReadHistory(f)
		_ = f.Close()
	}
}

// Close writes this session's history back and closes the liner
// Remember to call this function before exiting!
func Close() {
	// before closing, write history back to file
	if f, err := os.Create(historyFile); err == nil {
		_, _ = line.WriteHistory(f)
	}
	_ = line.Close()
}

// PromptAndRead prints `prompt` in the console and reads one line from user
// Note that '\n' will be removed
func PromptAndRead(prompt string) (string, error) {
	input, err := line.Prompt(prompt)
	if err != nil {
		return "", err
	}
	line.AppendHistory(input)
	return input, nil
}
