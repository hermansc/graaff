package main

/*
  Parse Graaff-config
*/

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Globals      map[string]string
	LayoutFolder string
	LayoutFile   string
	OutputFolder string
	OverviewFile string
	Seperator    string
	PostsFolder  string
	BaseTemplate string
	Truncate     int
}

// Global config-variable available to all methods.
var config Config

func parseConfig(filename string) error {
	// Ensure some config-variables are not nil.
	config.Globals = make(map[string]string)

	// Open config file and for each line get configuration variable.
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	// Open file as buffer and read each line one by one.
	fh := bufio.NewScanner(f)
	for fh.Scan() {
		// Get the line.
		line := fh.Text()

		// Strip away whitespaces.
		line = strings.TrimSpace(line)

		// Ensure it is not empty. If it is, we just continue to next line.
		if line == "" {
			continue
		}

		// Ensure it is not a comment.
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Alright, then we parse it.
		err = parseLine(line)
		if err != nil {
			return err
		}
	}

	// Check that the parsed config is non-empty. If some variables are not set by config file,
	// we use the default defined. So after this method is run we are guaranteed a filled configuration.
	validateConfig()

	return nil
}

func parseLine(line string) error {
	/*
	   Parse a config-line, insert to config-struct and fail if something is wrong.
	*/
	// Split the line on spaces.
	params := strings.Split(line, " ")

	// Get the variable-name
	variable := params[0]
	value := params[1]

	// Globals-keyword has a special format.
	if variable == "globals" {
		// On format: Key1='Value1', Key2='Value2'
		// Concat the values.
		val := strings.Join(params[1:], "")

		// Split on commas.
		globals := strings.Split(val, ",")

		for _, glob := range globals {
			// Remove whitespaces.
			glob = strings.TrimSpace(glob)

			// Get Key and Value by splitting on equal sign.
			keyval := strings.Split(glob, "=")

			// Set the config.
			config.Globals[keyval[1]] = strings.Trim(keyval[1], "'")
		}
	}

	// Parse these in a normal way.
	if variable == "layoutfolder" {
		config.LayoutFolder = value
	}
	if variable == "layoutfile" {
		config.LayoutFile = value
	}
	if variable == "outfolder" {
		config.OutputFolder = value
	}
	if variable == "overviewfile" {
		config.OverviewFile = value
	}
	if variable == "postsfolder" {
		config.PostsFolder = value
	}
	if variable == "truncate" {
		// Convert the value to an integer.
		var err error
		config.Truncate, err = strconv.Atoi(value)
		if err != nil {
			return err
		}
	}
	if variable == "seperator" {
		config.Seperator = strings.Trim(variable, "'")
	}
	return nil
}

func validateConfig() {
	// Check that each value in the config is non-null, and set default if it is.
	if len(config.Globals) == 0 {
		config.Globals = map[string]string{"SiteName": "My blog", "Author": "Someone", "Category": "posts"}
	}
	if config.LayoutFolder == "" {
		config.LayoutFolder = "layouts"
	}
	if config.LayoutFile == "" {
		config.LayoutFile = "index.html"
	}
	if config.OutputFolder == "" {
		config.OutputFolder = "generated"
	}
	if config.OverviewFile == "" {
		config.LayoutFile = "overview.html"
	}
	if config.Seperator == "" {
		config.Seperator = "-----"
	}
	if config.PostsFolder == "" {
		config.PostsFolder = "posts"
	}
	if config.BaseTemplate == "" {
		config.PostsFolder = "base.html"
	}
}

func main() {
	err := parseConfig("graaff.config")
	if err != nil {
		log.Fatal(err)
	}
}
