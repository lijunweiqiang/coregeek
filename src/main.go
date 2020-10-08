package main

import "os"

const (
	COMMAND_GAME = iota
	COMMAND_COMPRESS
	COMMAND_DECOMPRESS
	INVALID
)

func main() {
	args := os.Args
	command, inputFile, outPutfile := parseArgs(args)

}

func parseArgs(args []string) (int, string, string) {
	if len(args) != 3 {
		return INVALID, "", ""
	}
	switch args[0] {
	case "-game":
		return COMMAND_GAME, args[1], args[2]
	case "-compress":
		return COMMAND_COMPRESS, args[1], args[2]
	case "-decompress":
		return COMMAND_DECOMPRESS, args[1], args[2]
	default:
		return INVALID, "", ""
	}
}
