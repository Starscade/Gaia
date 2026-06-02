package utils

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Starscade/Gaia/internal/text"
)

func Confirm(message string) bool {
	fmt.Print(message + text.PrintYesOrNo)
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "y", "yes":
			return true
		default:
			return false
		}
	}

}

func GetEnv(env_var_name string, default_value string) string {
	env_var := os.Getenv(env_var_name)

	if env_var == "" {
		env_var = default_value
	}

	os.Setenv(env_var_name, env_var)

	return env_var
}

func GetFile(glob_path string) (string, error) {
	files, err := filepath.Glob(glob_path)
	if err != nil {
		return "", err
	}
	var bodies strings.Builder
	for _, file_path := range files {
		file, err := os.ReadFile(file_path)
		if err != nil {
			return "", err
		}
		eof_start := "EOF " + file_path + "\n\n"
		eof_end := "\n\nEOF " + file_path
		bodies.WriteString(eof_start + string(file) + eof_end)
	}
	return bodies.String(), nil
}

func GetStdin() string {
	stat, _ := os.Stdin.Stat()
	var input []byte
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		input, _ = io.ReadAll(os.Stdin)
	}
	return string(input)
}

func PrintErr(err string) {
	fmt.Println(text.PrintErr, err)
	os.Exit(1)
}

func PrintHelp() {
	fmt.Print(text.PrintGaia)
	flag.PrintDefaults()
	fmt.Print(text.PrintHelpUsage)
}
