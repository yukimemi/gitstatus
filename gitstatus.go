package main

import ( // {{{
	"bytes"
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

	fmt.Println("root", root)

	e = filepath.Walk(root, walkFn)
	failOnError(e)

} // }}}

func walkFn(path string, info os.FileInfo, err error) error {

	m, e := regexp.MatchString(`\.git$`, info.Name())
	if m && info.IsDir() {
		os.Chdir(path)
		os.Chdir("..")
		cmd := exec.Command("git", "remote", "-v")
		out, e := cmd.Output()
		failOnError(e)
		buf := bytes.NewBuffer(out)
		stdout := buf.String()
		if m, _ := regexp.MatchString("yukimemi", stdout); m {
			cmd := exec.Command("git", "status")
			out, e := cmd.Output()
			failOnError(e)
			fmt.Println(path)
			fmt.Printf("%s\n", out)
		}
	}
	return e
}
