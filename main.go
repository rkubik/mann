package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
	"time"
)

var (
	defaultConfigDir       = ".mann"
	defaultPagesDir        = "doc"
	defaultLastCommandFile = "last_command"
	defaultComment         = "Untitled"
)

var (
	listFlag = flag.Bool("l", false, "List pages")
	addFlag  = flag.Bool("a", false, "Add to page")
)

type baseConfig struct {
	Dir             string
	PagesDir        string
	LastCommandFile string
}

func getBaseConfig() (c *baseConfig, err error) {
	usr, err := user.Current()
	if err != nil {
		return c, err
	}
	c = &baseConfig{"", "", ""}
	c.Dir = path.Join(usr.HomeDir, defaultConfigDir)
	c.PagesDir = path.Join(c.Dir, defaultPagesDir)
	c.LastCommandFile = path.Join(c.Dir, defaultLastCommandFile)
	return c, err
}

func doSaveLastCommand(c *baseConfig, comment string) error {
	b, err := ioutil.ReadFile(c.LastCommandFile)
	if err != nil {
		return fmt.Errorf("Failed to find last command file (%s)", c.LastCommandFile)
	} else if len(b) == 0 {
		return fmt.Errorf("Last command was empty")
	}
	cmd := string(b)
	args := strings.Split(cmd, " ")
	return doSaveCommand(c, args[0], cmd, comment)
}

func doSaveCommand(c *baseConfig, prog string, cmd string, comment string) error {
	page := path.Base(prog)
	pageFile := path.Join(c.PagesDir, page)

	f, err := os.OpenFile(pageFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("# %s\n%s\n", comment, cmd))
	if err != nil {
		return err
	}

	return nil
}

func doPrintPages(c *baseConfig) error {
	pages, err := getPages(c.PagesDir)
	if err != nil {
		return err
	}

	for _, page := range pages {
		fmt.Println(page)
	}

	return nil
}

func doPrintPage(c *baseConfig, page string) (err error) {
	output, err := getPage(path.Join(c.PagesDir, page))
	if err != nil {
		return err
	}

	fmt.Println(output)

	return err
}

func getPage(f string) (string, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getPages(d string) (pages []string, err error) {
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return pages, err
	}

	for _, f := range files {
		pages = append(pages, path.Base(f.Name()))
	}

	return pages, nil
}

func main() {
	var err error
	flag.Parse()

	c, err := getBaseConfig()

	if *listFlag {
		if flag.NArg() > 0 {
			err = doPrintPage(c, flag.Arg(0))
		} else {
			err = doPrintPages(c)
		}
	} else if *addFlag {
		if flag.NArg() > 0 {
			err = doSaveLastCommand(c, strings.Join(flag.Args(), " "))
		} else {
			err = doSaveLastCommand(c, fmt.Sprintf("%s (%s)", defaultComment, time.Now().Format("2006-01-02 15:04:05")))
		}
	} else {
		fmt.Printf("Usage: %s [-l] [-a]\n", path.Base(os.Args[0]))
		return
	}

	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}
}
