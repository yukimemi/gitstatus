package main

import ( // {{{
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
) // }}}

// for debug {{{
type debugT bool

var debug = debugT(false)

func (d debugT) Println(args ...interface{}) {
	if d {
		log.Println(args...)
	}
} // }}}

func trace(s string) string { // {{{
	debug.Println("entering:", s)
	return s
} // }}}

func un(s string) { // {{{
	debug.Println("leaving:", s)
} // }}}

func failOnError(e error) { // {{{
	if e != nil {
		log.Fatal("Error:", e)
	}
} // }}}

func main() { // {{{

	root, e := os.Getwd()
	failOnError(e)

	var searchName string
	if len(os.Args) > 1 {
		searchName = os.Args[1]
	} else {
		searchName = "."
	}

	fmt.Println("root", root)

	e = printGitStatus(root, searchName)
	failOnError(e)

} // }}}

func printGitStatus(root string, searchName string) error { // {{{

	walkFn := func(path string, info os.FileInfo, err error) error {

		m, e := regexp.MatchString(`\.git$`, info.Name())
		if m && info.IsDir() {
			os.Chdir(path)
			os.Chdir("..")
			cmd := exec.Command("git", "remote", "-v")
			out, e := cmd.Output()
			failOnError(e)

			if m, _ := regexp.MatchString(searchName, string(out)); m {
				cmd := exec.Command("git", "status")
				out, e := cmd.Output()
				failOnError(e)
				if n, _ := regexp.MatchString("nothing to commit, working directory clean", string(out)); !n {
					fmt.Println("--------------------------------------------------------------------------------")
					path, _ := os.Getwd()
					fmt.Println("★ ", path)
					fmt.Println(string(out))
					fmt.Println("--------------------------------------------------------------------------------\n")
				}
			}
		}
		return e
	}
	e := filepath.Walk(root, walkFn)

	return e

} // }}}
