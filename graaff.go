package main

import (
  "os"
  "log"
  "fmt"
  "flag"
  "time"
  "sort"
  "bytes"
  "strings"
  "io/ioutil"
  "text/template"
  "github.com/russross/blackfriday"
)

func changeExtension(filename, newExtension string) string {
  // Takes a filename and changes extension. E.g. post.md => post.html
  parts := strings.SplitAfter(filename, ".")
  fileArray := parts[:len(parts)-1]
  return fmt.Sprintf("%v%v", strings.Join(fileArray, ""), newExtension)
}

func writeFile(outPath, htmlName string, html []byte) error {
  // Writes a file to a specified directory. If the directory does not exist,
  // we create it with the correct permissions.
  _, err := os.Stat(outPath)
  if os.IsNotExist(err) {
    os.Mkdir(outPath, 0766)
  } else if (err != nil) { return err }

  outFile := fmt.Sprintf("%v/%v", outPath, htmlName)
  err = ioutil.WriteFile(outFile, html, 0644)
  if (err != nil) { return err }
  return nil
}

func parseConfig(conf string, params map[string]interface{}, sep string) {
  // Parses a configuration string and adds them to a config-dictionary.
  // Each config-variable is declared using a colon (:)
  vars := strings.Split(conf, sep)
  for _, line := range vars {
    if (line == "") { return }
    l := strings.Split(line, ":")
    params[strings.TrimSpace(l[0])] = strings.TrimSpace(strings.Join(l[1:],":"))
  }
}

func parseTemplate(file string, data interface{}) (out []byte, err error) {
  // Parses a templatefile given a data-interface. Returns the bytes-buffer.
  var buf bytes.Buffer
  t, err := template.ParseFiles(file)
  if err != nil { return nil, err }
  err = t.Execute(&buf, data)
  if err != nil { return nil, err }
  return buf.Bytes(), nil
}

func generateHtml(baseTmpl, tmpl string, params map[string]interface{}) ([]byte, error) {
  // Outputs a byte-array of generated HTML based on a base-template containing
  // a Content block and a child-template.
  base, err := parseTemplate(tmpl, params)
  if (err != nil) { return []byte{}, err }
  params["Content"] = string(base)
  content, err := parseTemplate(baseTmpl, params)
  if (err != nil) { return []byte{}, err }
  return content, nil
}

func truncate(s string, i int) string {
  // Truncates a long string and adds ... if it is too long.
  parts := strings.SplitAfter(s, "</p>")
  if (len(parts) < i) { return s }
  return strings.Join(parts[:i], "")
}

type Post struct {
  // A basic struct for out overview. Wish I didn't need it.
  Title string
  Published string
  Author string
  Abstract string
  Filename string
}

func parseTimeString(ts string) time.Time {
  // Parses various time-strings and returns as a time-object.
  // Surely Golang has a smarter, more dynamic way of doing this?
  format := "2006-01-02"
  long_format := "2006-01-02 15:04"
  t := time.Now()
  if len(ts) == 10 { t, _ = time.Parse(format, ts)
  } else if len(ts) == 16 { t, _ = time.Parse(long_format, ts) }
  return t
}

// We want to be able to sort on publish-date.
type ByPublished []Post
func (b ByPublished) Len() int { return len(b) }
func (b ByPublished) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByPublished) Less(i, j int) bool {
  bi := parseTimeString(b[i].Published)
  bj := parseTimeString(b[j].Published)
  if (bi.Equal(bj)) {
    // If published in same day, sort on title.
    return bytes.Compare([]byte(b[i].Title), []byte(b[j].Title)) < 0
  }
  return bi.After(bj)
}

func main() {
  // Specify configuration
  var subFolder = flag.String("subfolder", "posts", "The folder containing markdown data")
  var outFolder = flag.String("outfolder", "output", "The folder containing generated html data")
  var layoutFolder = flag.String("layoutfolder", "layouts", "The folder containing different layouts/template files")

  var baseTmpl = flag.String("template", "base.html", "The file defining how your site looks")
  var postTmpl = flag.String("post", "post.html", "The file defining how your post looks")
  var overviewTmpl = flag.String("overview", "overview.html", "The file defining layout of overview/index-page")

  var index = flag.String("index", "index.html", "Index file with overview over posts")
  var overviewTitle = flag.String("overviewtitle", "Index", "Title to use in the overview")
  var truncateLength = flag.Int("truncate", 1, "Number of paragraphs to include in post overview before truncating")
  var seperator = flag.String("seperator", "-----", "The seperator between config and content")
  var config = flag.String("config", "", "Comma seperated list of config variables globally available.")
  var configFile = flag.String("configfile", "variables.conf", "Line seperated file specifying global variables")
  flag.Parse()

  // We assume all folders are based on the current location of the executable.
  baseFolder, err := os.Getwd()
  if (err != nil) { log.Fatal(err) }

  // Parse folders
  outPath := fmt.Sprintf("%v/%v", baseFolder, *outFolder)
  layoutPath := fmt.Sprintf("%v/%v", baseFolder, *layoutFolder)
  postsPath := fmt.Sprintf("%v/%v", baseFolder, *subFolder)

  // Parse layouts
  baseTemplate := fmt.Sprintf("%v/%v", layoutPath, *baseTmpl)
  postTemplate := fmt.Sprintf("%v/%v", layoutPath, *postTmpl)
  overviewTemplate := fmt.Sprintf("%v/%v", layoutPath, *overviewTmpl)

  // Get a list of all posts in the specified input directory.
  files, err := ioutil.ReadDir(postsPath)
  if (err != nil) { log.Fatal(err) }

  // A dictionary holding all variables sent to the templates.
  d := make(map[string]interface{})

  // Parse config from the configuration file
  conff := fmt.Sprintf("%v/%v", baseFolder, *configFile)
  b, _ := ioutil.ReadFile(conff)
  parseConfig(string(b), d, "\n")

  // Parse and load the config from the flag-variable
  // Overwriting eventual variables provided by config-file.
  parseConfig(*config, d, ",")

  // A list of all posts generated
  posts := make([]Post, 0)

  for _, f := range files {
    // Read the specified file.
    b, _ := ioutil.ReadFile(fmt.Sprintf("%v/%v", postsPath, f.Name()))

    // Parse it, seperating the post-config from the content.
    sep := strings.Split(string(b), *seperator)
    if (len(sep) != 2) { log.Fatal("Error in parsing post file. Ensure seperator is present, once") }
    parseConfig(sep[0], d, "\n")

    // Get the compiled markdown for the content.
    d["Post"] = string(blackfriday.MarkdownBasic([]byte(sep[1])))

    // Check if a custom layout is defined
    var content []byte
    if _, ok := d["Layout"]; ok {
      customLayoutFile := fmt.Sprintf("%v/%v", layoutPath, d["Layout"])
      content, _ = generateHtml(baseTemplate, customLayoutFile, d)
      delete(d, "Layout")
    } else {
      content, _ = generateHtml(baseTemplate, postTemplate, d)
    }

    // Write this HTML to a file in the specified output folder. 
    htmlName := changeExtension(f.Name(), "html")
    err = writeFile(outPath, htmlName, content)
    if (err != nil) { log.Fatal(err) }

    // Append to the global list of posts
    posts = append(posts, Post{
      Title: d["Title"].(string),
      Published: d["Published"].(string),
      Author: d["Author"].(string),
      Abstract: truncate(d["Post"].(string), *truncateLength),
      Filename: htmlName,
    })
  }

  // Sort all posts on publish-date and title.
  sort.Sort(ByPublished(posts))
  d["Posts"] = posts
  d["Title"] = *overviewTitle
  // Index-file referencing all posts
  content, _ := generateHtml(baseTemplate, overviewTemplate, d)
  err = writeFile(outPath, *index, content)
}
