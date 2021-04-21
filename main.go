package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fatih/color"
)

type arguments struct {
	Path    string
	Verbose bool
}

var (
	Green   = color.New(color.FgHiGreen).PrintfFunc()
	Cyan    = color.New(color.FgCyan).PrintfFunc()
	Red     = color.New(color.FgRed).PrintfFunc()
	Yellow  = color.New(color.FgYellow).PrintfFunc()
	Magenta = color.New(color.FgHiMagenta).PrintfFunc()
)

// func color(c string) func(...interface{}) string {
// 	return func(args ...interface{}) string {
// 		return fmt.Sprintf(c, fmt.Sprint(args...))
// 	}
// }

type parsedFile struct {
	path         string
	isDir        bool
	timeCreated  string
	timeModified string
	size         int64
	sizeString   string
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

		// set up our parsed file
		file := parsedFile{}

		// determine the absolute path
		path, _ := filepath.Abs(args.Path + sep + f.Name())
		file.path = path

		// find out whether the file is actually a directory
		file.isDir = determineType(path)

		// get the FileInfo for our current path
		stats, _ := os.Stat(file.path)

		// determine the time modified, and format it
		file.timeModified = stats.ModTime().Format(time.RFC822)
		Magenta("%s ", file.timeModified)

		// determine the time created (currently windows only)
		nano := stats.Sys().(*syscall.Win32FileAttributeData).CreationTime.Nanoseconds()
		file.timeCreated = time.Unix(0, nano).Format(time.RFC822)
		// fmt.Printf("%s ", Magenta(file.timeCreated))

		// determine the file size
		file.size = stats.Size()

		// Print
		if file.isDir {
			Green("%s ", "FOLD")
		} else {
			Green("%s ", formatBytes(file.size))
		}

		// Print respective stuff
		if file.isDir {
			Cyan("%s ", f.Name()+"/")
		} else {
			Cyan("%s ", f.Name())
		}

		fmt.Printf("\n")
	}

}

func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	k := 1024
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	i := math.Floor(math.Log(float64(bytes)) / math.Log(float64(k)))
	num := float64(bytes) / math.Pow(float64(k), i)
	if num == math.Round(num) {
		return fmt.Sprintf("%d %s", int(num), sizes[int(i)])
	} else {
		return fmt.Sprintf("%.1f %s", num, sizes[int(i)])
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
