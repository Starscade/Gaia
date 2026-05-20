package tools

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Starscade/Gaia/internal/text"
)

func ExitOnErr(err error) {
	if err != nil {
		PrintErr(err.Error())
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

func GetFile(glob_path string) string {
	files, err := filepath.Glob(glob_path)
	ExitOnErr(err)
	var bodies strings.Builder
	for _, file_path := range files {
		file, err := os.ReadFile(file_path)
		ExitOnErr(err)
		eof_start := "EOF " + file_path + "\n\n"
		eof_end := "\n\nEOF " + file_path
		bodies.WriteString(eof_start + string(file) + eof_end)
	}
	return bodies.String()
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
	log.Fatal(text.PRINT_ERR, err)
}

func PrintHelp() {
	fmt.Print(text.PRINT_GAIA)
	flag.PrintDefaults()
	fmt.Print(text.PRINT_HELP_USAGE)
}
