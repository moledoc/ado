/*
Extract ados recursively from given/current directory.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/moledoc/walks"
)

// search is a variable, that hold regexp of values we are searching from the files.
var search *regexp.Regexp = regexp.MustCompile("TODO:|NOTE:|HACK:|DEBUG:|FIXME:|REVIEW:|BUG:|TEST:|TESTME:|MAYBE:")

// indent is a size variable that dictates how big the spacing is between filepath and ado.
var indent int

// format is a variable, that allows us to control the format of ado output.
var format string

// adosFile is the filename where we write the TODOs, NOTEs etc.
var adosFile string

// adosFileFull is the full path to adosFile.
var adosFileFull string

// dirAction is a dummy function.
// When extracting TODOs, NOTEs etc, we do not want to perform any action on directories.
// So this functions purpose is to be argument for dir.Walk
func dirAction(path string) {}

// Help prints all subcommands and flags for program/tool `ado`.
func Help(checkingFlagCount bool) {

	// HACK: not the most elegant solution, but will do for now.
	order := []string{"file", "indent", "ignore", "search", "add", "depth"}
	orderType := []string{"filename", "int", "filename", "string", "filename", "int"}
	if len(order) != FlagCount || len(orderType) != FlagCount {
		log.Fatal("Incorrect number of flags defined in help funtion")
	}
	if checkingFlagCount {
		return
	}

	fmt.Printf("%-20s\n", "ado [SUBCOMMAND | [FLAG=<value>] [filname(s) | dirname(s)]]\n")
	fmt.Println("SUBCOMMANDS")
	subFormat := "\t%s\n\t\t%s\n"
	fmt.Printf(subFormat, "help", "Print this help.")

	// Modify the printing of flags to my liking.
	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Println("FLAGS")
		for i, name := range order {
			flagLoc := flagSet.Lookup(name)
			fmt.Printf("\t-%s=%s\n", flagLoc.Name, orderType[i])
			fmt.Printf("\t\t")
			fmt.Printf("%s", flagLoc.Usage)
			if flagLoc.DefValue != "" {
				fmt.Printf(" (default %s)\n", flagLoc.DefValue)
			} else {
				fmt.Printf("\n")
			}
		}
	}
	// new line for making space between subcommands and flags.
	fmt.Println()
	flag.Usage()
	fmt.Printf("\nDEFAULT SEARCH KEYWORDS\n\t\t%s\n", search.String())
	fmt.Println("\nEXAMPLES")
	fmt.Printf(subFormat, "ado", "Gets recursively all ados from current working directory.")
	fmt.Printf(subFormat, "ado help", "Prints help for `ado` program.")
	fmt.Printf(subFormat, "ado . ..", "Gets recursively all ados from current and parent directory.")
	fmt.Printf(subFormat, "ado ./ ../", "Gets recursively all ados from current and parent directory.")
	fmt.Printf(subFormat, "ado <file1> ../", "Gets recursively all ados from parent directory and <file1>.")
	fmt.Printf(subFormat, "ado -file=test.txt <file1>", "Gets all ados from <file1> and in addition saves them to test.txt in the current working directory.")
	fmt.Printf(subFormat, "ado -indent=100", "Gets recursively all ados from current directory and prints them so that ados start at column 100.")
	fmt.Printf(subFormat, "ado -ignore=.adoignore", "Gets recursively all ados from current directory, ignoring files and directories mentioned in the given file.")
	fmt.Printf(subFormat, "ado -search=RandomString", "Gets recursively all 'RandomString' mentions from current directory.")
	fmt.Printf(subFormat, "ado -add=.adoadd", "Gets recursively all ados from current directory, including the keywords mentioned in the given file.")
	fmt.Printf(subFormat, "ado -depth=2", "Gets recursively all ados from current directory, until subdirectory depth is 2 (including).")
}

// clean deletes file with path adosFileFull, if it exists.
func clean() {
	// If file deleted, let user know about it
	if _, err := os.Stat(adosFileFull); err == nil {
		errRm := os.RemoveAll(adosFileFull)
		if errRm != nil {
			log.Fatal(errRm)
		}
		// fmt.Printf("Deleted file: %v\n", adosFileFull)
	}
}

// fileAdos is a function that parses TODOs, NOTEs etc (see variable `search`) from given file and writes formatted lines to stdout or adosFile, when it is given.
// Files without file extension are skipped.
func fileAdos(path string) {
	// Will ignore files, that does not have extensions (eg binaries) or have .exe extension.
	if filepath.Ext(path) == "" || filepath.Ext(path) == ".exe" {
		return
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var f *os.File
	if adosFile != "" {
		f, err = os.OpenFile(adosFileFull, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	for lineNr, line := range strings.Split(string(contents), "\n") {
		if search.MatchString(line) {
			// // get the first occurence of `search` and clean the line until that point.
			// adoInd := search.FindStringIndex(line)[0]
			adoInd := 0
			formattedLine := fmt.Sprintf(format, path+":"+fmt.Sprint(lineNr+1)+":", strings.TrimLeft(line[adoInd:], " |\t"))
			// Write formated line to adosFile when its given.
			if adosFile != "" {
				if _, err := f.WriteString(formattedLine); err != nil {
					log.Fatal(err)
				}
			}
			fmt.Print(formattedLine)
		}
	}
}

// addKeywords function adds keywords to internal gloabal variable `search` from given file.
// It's expected, that each line contains regexp compliant string.
func addKeywords(kwFile string) {
	contents, err := os.ReadFile(kwFile)
	if err != nil {
		log.Fatal(err)
	}
	searchLoc := search.String()
	for _, line := range strings.Split(string(contents), "\n") {
		if line == "" {
			break
		}
		searchLoc += "|"
		if line == "." || line == ".." {
			line = "^" + line + "$"
		}
		searchLoc += strings.Replace(line, ".", "\\.", -1)
	}
	search = regexp.MustCompile(searchLoc)
}

// FlagCount is internal global variable to have an assert for ado help function
var FlagCount int = 6

func main() {

	adosFileLoc := flag.String("file", "", "The filename where the TODOs, NOTEs etc are saved in the current working directory.")
	indentLoc := flag.Int("indent", 60, "The size of indentation between filepath and ado.")
	ignoreLoc := flag.String("ignore", ".adoignore", "File, where each line represents one directory or file that is ignored. If when .adoignore exist in the current directory, this flag is not necessary.")
	searchOne := flag.String("search", "", "Search for a specific string (regexp allowed); will overwrite the default search keywords (see keywords).")
	searchFile := flag.String("add", "", "File, where each line represents one additional keyword (regexp allowed).")
	depth := flag.Int("depth", -1, "The depth of directory structure recursion, -1 is exhaustive recursion.")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Check, that flags are presented to our specification.
	for _, arg := range os.Args {
		if strings.Contains(arg, "-") && !strings.Contains(arg, "=") {
			log.Fatal("Invalid flag syntax, expected <flag>=<value>")
		}
	}

	if *searchOne != "" {
		search = regexp.MustCompile(*searchOne)
	}
	if *searchFile != "" {
		addKeywords(*searchFile)
	}
	adosFile = *adosFileLoc
	indent = *indentLoc
	format = "%-" + fmt.Sprint(indent) + "s%s\n"

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// if windows file path structure, then change it to use forward slash (/)
	root = strings.Replace(root, "\\", "/", -1)
	// set walks.Ignore
	walks.SetIgnore(*ignoreLoc)
	// set full path to ados file, when its given
	if adosFile != "" {
		adosFileFull = root + "/" + adosFile
		clean()
	}

	var WaitGroup sync.WaitGroup

	// If no arguments (besides flags or subcommand) is given, then parse current directory.
	if len(os.Args) == 1 || flag.NArg() == 0 {
		WaitGroup.Add(1)
		go func() { defer WaitGroup.Done(); walks.Walk(root, fileAdos, dirAction, *depth) }()
	}

	// make a regexp to handle dot (parent and current) directories
	parentDir := regexp.MustCompile("\\.\\./|\\.\\.$|/\\.\\.$")
	currentDir := regexp.MustCompile("\\./|\\.$|/\\.$")

	for i := flag.NFlag() + 1; i < len(os.Args); i++ {
		rootLoc := root
		arg := os.Args[i]

		if walks.Ignore.MatchString(arg) && walks.Ignore.String() != "" {
			continue
		}
		// For each parent dot directory reference in arg, move local root (rootLoc) one directory higher.
		for i := len(parentDir.FindAllStringIndex(arg, -1)); i > 0; i-- {
			rootLoc = filepath.Dir(rootLoc)
		}
		// Clean arg from dot directories.
		cleanedArg := parentDir.ReplaceAllString(arg, "")
		cleanedArg = currentDir.ReplaceAllString(cleanedArg, "")
		// When arg had a dot directory, then handle the path format correctly.
		var path string
		if cleanedArg == "" {
			path = rootLoc
		} else {
			path = rootLoc + "/" + cleanedArg
		}

		pathType, err := os.Stat(path)
		if err != nil {
			if strings.Contains(path, "help") {
				Help(false)
				os.Exit(0)
			}
			log.Fatal(err)
		}

		switch mode := pathType.Mode(); {
		default:
			log.Fatalf("Unreachable: %s is not a file nor directory\n", path)
		case mode.IsRegular():
			WaitGroup.Add(1)
			go func() { defer WaitGroup.Done(); fileAdos(path) }()
		case mode.IsDir():
			WaitGroup.Add(1)
			go func() { defer WaitGroup.Done(); walks.Walk(path, fileAdos, dirAction, *depth) }()
		}
	}

	// wait until waitgroups are Done
	WaitGroup.Wait()
	walks.WaitGroup.Wait()

	// If adosFile created, let user know about it.
	if _, err := os.Stat(adosFileFull); err == nil {
		fmt.Printf("Ados parsed to: %v\n", adosFileFull)
	}
}
