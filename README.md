# 🌱 GoGardener [WIP]

A blazing fast cli for generating digital garden from markdown files. Supports markdown files exported from [Roam Research](https://roamresearch.com/), [Obsidian](https://obsidian.md/) etc.

Demo at https://vishnuks.com/garden/

## 🎉 Features

* Generates HTML pages from the markdown files.
* Supports backlinks.
* Supports tags (coming soon).


## 🛠️ Installation

```bash
git clone https://github.com/hackerkid/gogardener
cd gogardener
go install .
```

## ✍️ Usage

```
~/go/bin/gogardener --input markdown-files-dir --output output-dir --base-template base.html

```

* `input` - The directory containing your markdown files.
* `output` - The directory where the generated HTML files should be stored (⚠️ All the files in this directory would be removed.)
* `base-template` - The file that should be used as the base template for the HTML files. You can use the one included in the repo for starting up.

## 🚧 Todo

* Make tags work
* Make images work.
