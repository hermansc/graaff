# Graaff - A simple static generator in Go.

**Graaff**, named after the [Van de Graaff electrostatic
generator](http://en.wikipedia.org/wiki/Van_de_Graaff_generator) is a super
simple static generator for HTML - not energy. It takes a folder of Markdown
documents, converts and outputs the result in a blog-aware highly customizable
layout.

The goals for this project were:
* To be extemely simple. I want Apache/Nginx to serve my files, not
Go/Python/Foobar.
* I don't need automatic reloading, site-generation, building nor any other
"magic".
* I use markdown, not Textile or *X*.  * *Fast*, this is much thanks to a great
[markdown library](https://github.com/russross/blackfriday) but also due to
Graaff's simplicity.

## Getting started with Graaff

First install Golang. Instructions for this procedure is given
[here](http://golang.org/doc/install).

In order to compile all markdown `.md` files in the `posts`-folder provided in
this repository, you only need to run the binary as so:

    $ go run graaff.go

Now your HTML-files should be generated and written to the `output`-folder.
These can then be opened with your web-browser or simply served by your
favorite static file server, such as Apache or Nginx.

## Writing your first blog-post

A simple blog-post looks like this:

    Title: Nemo enim ipsam
    Author: Mark Antony
    Published: 2014-03-17
    -----
    Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium
    doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore
    veritatis et quasi architecto beatae vitae dicta sunt explicabo.
    
Then by editing the layout we can access the blog-variables:

    <h1>{{.Title}}</h1>
    <p>
        <a href="index.html">← Overview</a> - 
        Published: <em>{{.Published}}</em> - 
        Author: <em>{{.Author}}</em>
    </p>
    <p>{{.Post}}</p>

Let's say you want to add a custom variable for your own site, this is as
simple as adding it before the seperator `-----`:

    (...)
    Published: 2014-03-17
    Tags: new-blog,sunny
    -----
    (...)

And it can be accessed in the template file with `{{.Tags}}`.

## Customization

Ignoring all the customizations available in the template files with regards to
how the site looks, there are equally many customizations available concerning
how *Graaff* works. By issuing `--help` we see the options available and their
default values:

    $ ./graaff --help
    Usage of ./graaff:
      -config="": Comma seperated list of config variables globally available.
      -configfile="variables.conf": Line seperated file specifying global variables
      -index="index.html": Index file with overview over posts
      -layoutfolder="layouts": The folder containing different layouts/template files
      -outfolder="output": The folder containing generated html data
      -overview="overview.html": The file defining layout of overview/index-page
      -overviewtitle="Index": Title to use in the overview
      -post="post.html": The file defining how your post looks
      -seperator="-----": The seperator between config and content
      -subfolder="posts": The folder containing markdown data
      -template="base.html": The file defining how your site looks
      -truncate=1: Number of paragraphs to include in post overview before truncating

Special mention will here be given to the `config` and `configfile` variables.
A configuration-file or string can contain globally available variables for
your site. The two flags have the same purpose and effect, but the variables
passed in `config` have higher presedence. A configfile `variables.conf` can
look like:

    SiteName: My awesome generated site
    Author: My name

And now `{{.SiteName}}` and `{{.Author}}` are available throughout all your
templates. This is equal to saying:

    ./graaff -config="SiteName: My awesome generated site, Author: My Name"

You can even specify your own **custom layout** for one blog-post:

    Layout: mycustomlayout.html

And then by creating this file in your `layouts`-folder the blog-post gets a
unique design.

## Caveats

As with everything simple, there are some things that fits for me, but perhaps
not for you. One is the generation of the overview-page (`index.html`) and the
other handling of dates.

The first, the overview page is "special" in the way that only the global
variables and some variables specified in the code, are available to the
template. These are:
`{{.Author}}`,`{{.Published}}`,`{{.Title}}`,`{{.Abstract}}` and
`{{.Filename}}`. Thus if you want custom variables in this overview you need to
fork the repository and customize the code for yourself. It should be as simple
as changing the `Post`-struct and appending the right values in the
`main()`-function.

The second is handling of dates. As I wanted to sort the posts based on
publish-date, the code accepts two formats:
[ISO-8601](http://en.wikipedia.org/wiki/ISO_8601) on the form `2014-03-17` and
the second being `2014-03-17 18:00`. There is perhaps a smarter way to do this,
but I didn't really care.

## Licence

[GNU General Public License v2
(GPL-2)](https://tldrlegal.com/license/gnu-general-public-license-v2#summary)

**You can**:
 * Commercial Use 
 * Modify 
 * Distribute 
 * Place Warranty

**You cannot**:
 * Sublicense 
 * Hold Liable

**You must**:
 * Include Original 
 * Disclose Source

Full licence available [here](http://www.gnu.org/licenses/gpl-2.0.html)
