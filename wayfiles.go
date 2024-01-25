package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/exp/slices"
)

// global variables must be declared without the colon :
var extensions = []string{
	"pem",
	"env",
	"sql",
	"cfg",
	"config",
	"apk",
	"json",
	"yml",
	"yaml",
	"xml",
	"log",
	"git",
	"enc",
	"key",
	"ini",
	"ps1",
	"sh",
	"bat",
	"exe",
	"cgi",
	"msi",
	"jar",
	"py",
	"db",
	"mdb",
	"bak",
	"bkp",
	"bkf",
	"inc",
	"asa",
	"old",
	"iso",
	"bin",
	"swf",
	"pl",
	"htm",
	"txt",
	"doc",
	"docx",
	"xls",
	"xlsx",
	"ppt",
	"pptx",
	"pdf",
	"eml",
	"email",
	"msg",
	"gadget",
	"tmp",
	"temp",
	"xz",
	"dll",
	"bz2",
	"do",
	"zst",
	"bz",
	"gz",
	"ovpn",
	"vpn",
}

func banner() {
	banner := ` 

	╦ ╦┌─┐┬ ┬┌─┐┬ ┬  ┌─┐┌─┐
	║║║├─┤└┬┘├┤ │ │  ├┤ └─┐
	╚╩╝┴ ┴ ┴ └  ┴ ┴─┘└─┘└─┘											  
	
				with <3 by Lihaft
	`
	fmt.Println(banner)
}

func verifyOS() {
	if runtime.GOOS == "windows" {
		fmt.Println("Wayfiles only runs in Linux systems. Exiting...")
		os.Exit(1)
	}
}

// Open the file list and recreate the extensions slice
func useCustomList(flList string) {
	file, err := os.Open(flList)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	sc := bufio.NewScanner(file)
	extensions = make([]string, 0) // empty the extensions slice using make

	for sc.Scan() {
		extensions = append(extensions, sc.Text())
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
}

// Append extensions
func includeExtensions(flInc string) {
	includes := strings.Split(flInc, ",")
	extensions = append(extensions, includes...)
}

// TODO: Refactor this in the future with the optimal functions
func excludeExtensions(flExc string) {
	newExtensions := make([]string, 0)
	excludes := strings.Split(flExc, ",")

	for _, ext := range extensions {
		if !slices.Contains(excludes, ext) {
			newExtensions = append(newExtensions, ext)
		}
	}
	extensions = newExtensions
}

// Include or exclude certain extensions based on flags.
func updateExtensions(flInc string, flExc string, flList string) {
	if flList != "" {
		useCustomList(flList) // use custom list of extensions if there is one specified
	}
	if flExc != "" {
		excludeExtensions(flExc)
	}
	if flInc != "" {
		includeExtensions(flInc)
	}

}

// Core function to search for juicy extensions
func searchFiles(flPath string, flSilent bool, flInc string, flExc string, flList string) {

	updateExtensions(flInc, flExc, flList)

	for idx, _ := range extensions {
		ext := extensions[idx]

		// If silent flag is not specified, be verbose
		if !flSilent {
			yellow := color.New(color.FgYellow).SprintFunc()
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("\n%s Searching for files that have the extension: .%s \n", yellow("[*]"), red(ext))
		}

		out, err := exec.Command("grep", "-ihR", ".*\\."+ext+"$", flPath).Output()

		if err != nil && !flSilent {
			fmt.Println("No results found.")
			fmt.Println("")
		}

		output := string(out[:])
		fmt.Print(output)
	}

	// DEBUGGING PURPOSES
	//fmt.Println(extensions)
	//fmt.Println(len(extensions))
}

func main() {
	version := "1.0.2"
	verifyOS()

	var (
		flPath    = flag.String("f", "", "[REQUIRED] File/path for either a file with URLs or a directory with Wayback Machine results. (E.g.: wayfiles -f urls.txt | wayfiles -f ~/waymore/results | wayfiles -f .)")
		flSilent  = flag.Bool("s", false, "Silent/Pipable mode. Not verbose mode. Just print the URLs to stdout")
		flInc     = flag.String("i", "", "Include extra extensions to search for in the format of: ext1,ext2,ext3 E.g.: -i php,js,aspx")
		flExc     = flag.String("e", "", "Exclude certain extensions to search for in the format of: ext1,ext2,ext3 E.g.: -e db,pdf,doc")
		flList    = flag.String("l", "", "Use a custom list with extensions instead of using the default ones. The list must include one extension per line. E.g: php (newline) js (newline) aspx (newline) etc")
		flExt     = flag.Bool("ext", false, "Print the list of extensions by default")
		flVersion = flag.Bool("version", false, "Print version number")
		flHelp    = flag.Bool("h", false, "Prints help menu") // flag package automatically adds the --help flag too
	)

	flag.Parse()

	if *flVersion {
		fmt.Println("VERSION: ", version)
		os.Exit(0)
	}

	if *flHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *flExt {
		fmt.Println(extensions)
		os.Exit(0)
	}

	if *flPath == "" {
		fmt.Println("\n ERROR: required argument -f must be specified with a file/path of a file or folder.")
		fmt.Println("")
		flag.Usage()
		os.Exit(1)
	}

	if !*flSilent {
		banner()
	}

	searchFiles(*flPath, *flSilent, *flInc, *flExc, *flList)

}
