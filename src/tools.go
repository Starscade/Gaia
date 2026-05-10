package main

import (
	"flag"
	"fmt"
	"log"
)

func printErr(err string) {
	log.Fatal(" \033[1;31mERR:\033[0m ", err)
}

func printHelp() {
	fmt.Println(PRINT_GAIA)
	flag.PrintDefaults()
	fmt.Println(PRINT_HELP_USAGE)
}
