# Ninetails

Ninetails is tail-like monitoring application written in Go. It's purpose is to allow for easy identifcation of important fields, data or lines using Regex and a yaml configuration file.

# Installation
The easiest approach is to use something like `wget`: 

```
curl -o ninetails "https://github.com/mikemackintosh/ninetails/releases/download/v1.0.2/ninetails-$(uname -o |  awk '{print tolower($0)}')-$(uname -m)"
chmod +x ninetails
```

# Configuration
Each project I am working on has a `.ninetails.yml` configuration file. An example can be found below:

Tails define each of your search strings and the mappings of what you want to occur. You can do full line highlighting by passing `color`. Alternatively, if you pass `format`, include match groups in your `search` string and you'll be able to reference their positional matches within your `format` string.

Colors can be customized, along with other escapes. These strings are auto-escaped when being output, `\033[` is the implied prefix.
```
---
tails:

  - search: '"level":"(.*?)","msg":"(.*?)"'
    format: '"level":"\YELLOW$1\CLEAR","msg":"\BABYBLUE$2"'

  - search: '"payload":(\[.*?\]),'
    format: '"payload":\BABYBLUE$1\CLEAR,'

  - search: '"error":"(.*?)"'
    format: '"error":"\RED$1\CLEAR'
    exit_on_match: true

  - search: "INFO"
    color: PURPLE

  - search: "DEBUG"
    color: ORANGE

colors:
  PURPLE: "38;5;129m"
  PINK: "38;5;162m"
  RED: "38;5;196m"
  ORANGE: "38;5;208m"
  YELLOW: "38;5;184m"
  GREEN: "38;5;154m"
  BLUE: "38;5;32m"
  GREY: "38;5;242m"
  DARKGREY: "38;5;239m"
  LIGHTGREY: "38;5;249m"
  BABYBLUE: "38;5;123m"
  LIGHTPINK: "38;5;212m"
  WHITE: "38;5;7m"
  CLEAR: "0m"
  ```
