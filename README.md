feed-cli
===
simple utility for formatting (groups of) feeds

# Usage
```bash
$ echo "https://your.feed.com/feed" > go run feed.go
```

## Options
```text
-templateOverride [relative path to go template]
    point to your own go templates for formatting
    $ go run feed.go -templateOverride ./index.tpl 
```