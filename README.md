# Mate

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->  

- [Mate](#mate)
  - [Introduction](#introduction)
    - [What is a static site generator?](#what-is-a-static-site-generator)
    - [Pros and cons of static site generators](#pros-and-cons-of-static-site-generators)
    - [Performance](#performance)
  - [How to use](#how-to-use)
      - [Project structure](#project-structure)
    - [Layouts](#layouts)
      - [Example of a layout](#example-of-a-layout)
      - [Post root object properties](#post-root-object-properties)
      - [Layout functions](#layout-functions)
      - [Special layouts](#special-layouts)
      - [Creating a homepage](#creating-a-homepage)
      - [Assets](#assets)
  - [Posts](#posts)
      - [What is Markdown](#what-is-markdown)
      - [Front matter](#front-matter)
        - [Front matter parameters](#front-matter-parameters)
    - [Public folder](#public-folder)
  - [To-do](#to-do)

<!-- /code_chunk_output -->

## Introduction

**Mate** is a **static site generator** built in **Go** as a personal learning project of mine.

It has very basic features but it is blazing fast thanks to the implementation in Go, and can be further optimized and extended the more I get confidence with the language, its concurrency patterns (which I am not used to coming from PHP), and so on.

### What is a static site generator?

To provide easily content for a website generally we rely on **frameworks** or even more often **CMS applications**; while extremely flexible and powerful, these applications are **dynamic** and therefore might depend a lot on their infrastructural dependencies and issues, other than being inherently more complex to manage and extend, involving many interacting parts.

This kind of complexity, though, is not needed for simpler websites, where we could benefit from a different approach; instead of relying on **databases**, **services** and **backend-generated content**, we **generate pure HTML files** using text files as source (written as **posts** in **Markdown** format), static assets and some HTML layouts.

### Pros and cons of static site generators

**Pros:**

* **no backend** nor backend skills needed, the resulting site is pure HTML. This means also:

   * **no overhead at all** and zero computational capabilities needed to serve the pages, you are literally using your server speed at 100% with the fastest possible pages

   * the **resources** neeeded on the server will also be **minimal** 

   * the pages can be easily cached for even more speed

   * the site's `public` folder will work **regardless of any backend installed or not**, so it will be de facto **100% portable and functional anywhere**, on any server capable of serving HTML and static assets

* **super-fast**, almost real-time generation (see [Performance](#performance) section)

* easy to understand **pure text format**

  * being text-based, the content can easily be:

    * **versioned under source control** (f.ex. with `git`)

    * **universally viewed and edited** with any text editor, even including simple notepad apps. 
    
    * opening and editing a text file is also usually blazing fast on any OS 

**Cons:**

* **not dynamic** at all, it's pure HTML + assets (images and JavaScript code)
* therefore it cannot perform by itself any backend stuff such as authentication, using user imput, using querystrings, ... 
* problematic anywhere complexity and/or high dynamicity is needed (f. ex. use of querystrings like as pagination), especially when the volume of the content grows
* no preview nor any WYSIWYG possibility

### Performance

This is a **typical benchmark** for this version on my current development machine (Intel® Core™ i7-1065G7 CPU @ 1.30GHz × 8). It processes 3000 randomly generated text files in around 0.5 seconds.

```
No. of posts:       3000 posts
Aggregated .md:            3ms            3ms
Parsed to HTML:          264ms          261ms
Written HTML:            474ms          210ms
Assets copied:           474ms             0s
```

## How to use

If you are using the sources, build with:

```bash
$ go build main.go
```

If the `mate` binary is, or is moved into, an **executable path** (see `$PATH` variable in your system), it can be launched directly.

#### Project structure

Prepare your project using this structure:

```
myProject/
    |--- layouts/
    |       |--- index.tpl.html
    |       |--- post.tpl.html
    |       |--- tag.tpl.html
    |       |--- ...
    |       └--- assets/
    |               |--- my.css
    |               |--- my.js
    |               └--- ...
    |--- posts/
    |      |--- index.md
    |      |--- post.md
    |      └--- ...
    └--- public/
```

Then `cd` into your project:

```bash
$ cd myProject
$ mate
```

The last command might need a path (f.ex. `./mate` or `./somedir/mate`) if the `mate` binary is in not a folder included into `$PATH`.

### Layouts

A **layout** is the "blueprint" that is used to generate HTML files based on a post's content and metadata.

You will put your **layouts** in the `/layouts` folder, in files with the layout name and a `.tpl.html` extension (f. ex. `myLayout.tpl.htm` for the `myLayout` layout).

Technically, layout is a **Go HTML template**: i.e. essentially a HTML file with some dynamic content and flow control possibilities (enclosed by special delimiters `{{` and `}}` and following Go's template syntax).
This template will be paired each time with a post that will determine the values of its dynamic parts, and for each post a static HTML file will be generated, preserving the folder structure of the original posts.
In other words, you use layouts to "skin" and "style" your textual content and transform it in a usable HTML page.

You can have as many layouts as you want: one default layout posts, one for tags, one for the homepage, special layouts for certain selected content ...

#### Example of a layout

**layouts/post.tpl.html**

```html
<!DOCTYPE html>
<html>
    <head>
        <title>{{ .Title }}</title>
    </head>
    <body>
        <h1>{{ .Title }}</h1>
        {{ .Body }}

        {{ range tags -}}
            <a href="{{ linktag . }}">{{ . }}</a>
        {{ end }}
    
        {{ range .Tags }}
            <a href="{{ linktag . }}">{{ . }}</a>
        {{ end }}
    </body>
</html>
```

Notable things to know:

* All dynamic parts are enclosed by `{{` and `}}`

* In case of postlayouts the root object (`.`) is a `Post` structure whose properties (like `.Body` or `.Title`) are used to generate content.

* You can spot some **built-in functions** of the Go templates (like `range`) ...

* ... and **custom template functions** are defined to help generating certain parts of the HTML code; they follow the syntax rules of Go template functions.
  
  This example shows a first function that does not require parameters (`tags`) and a second that does (`linktag`, where the parameter is `.`, which in turn will hold in this particular case the value of the current tag in the `range` loop):

  ```html
  <ul>
    {{ range tags }}
      <li><a href="{{ linktag . }}">{{ . }}</a></li>
    {{ end }}
  </ul>
  ```
* while `.Tags` is a **property of Post** containing the **tags for the current post**, `tags .` is a **function** giving back **all the tags used in the project**.

#### Post root object properties 

Property     | Content
-------------|------------------------------------
`.Title`     | The post's title; it's taken from the file name but can be overridden in the front matter
`.Body`      | The post, parsed as HTML

#### Layout functions

Function  | Parameter(s) | Description
----------|--------------|---------------
`tags`    | -            | Returns a `[]string` slice with all the tags collected from all the posts
`linktag` | `tag`        | Helps generating a URL for the specified `tag`

#### Special layouts

Some layouts are treated specially:

* `post.tpl.html` is the **default layout for posts**; it can be overridden per single post via the post's front matter's `layout` property

* `tag.tpl.html` is the **default layout for a specific tag page**

#### Creating a homepage

Although not mandatory, usually it's also advisable to have **a homepage for your website**; anytime we visit a website's root URL, most servers will search for a `index.html` file to serve.

You can easily accomplish this:

* first create an `index.md` Markdown file at the root of `/posts`; this source file will make generate a `index.html` for the **website's homepage** in the `/public` root

* you might also most likely want a **dedicated layout for your homepage**.
  This can be done via a normal front matter override, it's perfectly legit to create a special `home.tpl.html` layout (or `index.tpl.html` or whichever name you prefer) that will work just for the homepage:

  ```yaml
  ---
  layout: home
  ---
  # This is the homepage!
  ```

#### Assets

Your site needs to include **publicly accessible static front-end assets**, most notably CSS stylesheet and JavaScript files.

To this end, the folder `/layouts/assets` will be copied as-is and it's the ideal place to put your assets into. 

The `/public/assets` folder will mirror its contents, placing `/assets` folder at the website's root.

## Posts

**Posts** are the heart of the content of your site; they are **plain text files**   written in **Markdown format** and saved with a `.md` extension.

#### What is Markdown

**Markdown** is a **rich text content format**; more specificially, it's a **purely textual format** that aims at two goals: 

1. provide a super-easy way to specify end edit formatting in a text file (usually converted into HTML)

2. choose a format that would make sense and could convey the formatting also when viewed in pure text format, with simple applications like notepads

As an example, this Markdown snipped uses the syntax to specify boldface (`**...**`) and level 1 title (`#`), and is suitably converted in html:

```
# Hello world!
Time to learn **Markdown**!
```

and will be rendered in a way similar to:

```
<h1>Hello world!</h1>
<p>Time to learn <strong>Markdown</strong></p>
```

#### Front matter

A post can optionally begin with a **front matter**, which is a short preamble written in **YAML**. 

Front matter it is used to provide **additional metadata** for the post, in the form of **key-value pairs**. It's always at the beginning of the file and enclosed between lines with the `---` delimiter.

An example of a post with front matter would be:

```yaml
---
layout: otherLayout
title: "Custom title"
tags: foo, bar
---
# Actual content

Some **real** content.
```

This front matter sets the values for `layout`, `title` and `tags` - the first two will also override the defaults. 

##### Front matter parameters

The accepted parameters are:

parameter  | value 
-----------|---------------------------
`layout`   | Specify an alternative layout. It will need a `<name>.tpl.html` file (*)
`title`    | Custom title
`tags`     | A comma-separated list of tags; whitespace around tags will be trimmed
`author`   | Custom author

(*) **hierarchical tags** are not implemented, but folders can used to an extent to emulate this (f.ex. `mytag` and `mytag/subtag` -> `mytag.tpl.html` and `mytag/subtag.tpl.html`

### Public folder

In `/public` you will find all the generated HTML files together with the assets (under `/public/assets`).

### Other notes/caveats

* the `/public` folder should **not be edited**, as it will also be **wiped** on every generation

## To-do

- [ ] support `author`
- [ ] support concurrency in parsing
- [ ] support dates
- [ ] write proper tests
- [ ] support command line parameters
- [ ] support watcher with incremental build
- [ ] support pagination
- [ ] support more functions like next, previous post ...
- [ ] add static and dynamic server
