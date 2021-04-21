package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type arguments struct {
	Path    string
	Verbose bool
}

var (
	Green   = color("\033[38;5;10m%s\033[0m")
	Teal    = color("\033[38;5;14m%s\033[0m")
	Orange  = color("\033[38;5;215m%s\033[0m")
	Yellow  = color("\033[38;5;221m%s\033[0m")
	Magenta = color("\033[38;5;5m%s\033[0m")
)

func color(c string) func(...interface{}) string {
	return func(args ...interface{}) string {
		return fmt.Sprintf(c, fmt.Sprint(args...))
	}
}

type parsedFile struct {
	Path         string
	isDir        bool
	timeCreated  string
	timeModified string
	size         string
}

func main() {
	args := arguments{}

	for _, arg := range os.Args[1:] {
		if arg[0] != '-' {
			args.Path = arg
		} else {
			switch arg {
			case "-v", "--verbose":
				args.Verbose = true
			default:
				fmt.Printf("unknown flag %s\n", arg)
				os.Exit(1)
			}
		}
	}

	if args.Path == "" {
		args.Path = "."
	}

	dir, err := os.Open(args.Path)
	if err != nil {
		log.Fatalf("Failed to open directory: %s", err)
	}
	defer dir.Close()

	sep := string(os.PathSeparator)

	list, _ := dir.ReadDir(0)
	for _, f := range list {

		// Set up our parsed file
		file := parsedFile{}

		// Determine the absolute path
		path, _ := filepath.Abs(args.Path + sep + f.Name())
		file.Path = path

		// Find out whether the file is actually a directory
		file.isDir = determineType(path)

		stats, _ := os.Stat(file.Path)

		file.timeModified = stats.ModTime().Format(time.RFC822)
		fmt.Printf("%s ", Magenta(file.timeModified))

		// file.timeCreated = stats.Sys().(*syscall.Stat_t)

		// Print respective stuff
		if file.isDir {
			fmt.Printf("%s ", Teal(f.Name()+"/"))
		} else {
			fmt.Printf("%s ", Green(f.Name()))
		}

		fmt.Printf("\n")
	}

}

func determineType(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Failed to read file or directory: %s", err)
	}

	switch mode := file.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	default:
		return false
	}
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}
