feed-cli
===
![CI](https://github.com/nhomble/feed-cli/workflows/CI/badge.svg)

simple utility for formatting (groups of) feeds

# Usage
```bash
$ echo "https://your.feed.com/feed" > go run feed.go
```

## Feed metadata
The goal has been to keep the feed file simple. At a minimum, you can just provide
a list of feeds delimited by new-lines and you are done. To naturally give some more metadata per feed, the
cli recognizes the following format

Given a feed file formatted like:
```text
[feed]
[feed]
```

Then the grammar looks like:

```text
feed        := [feedLink] [metadata]
feedLink    := string
metadata    := [key]=[value] [metadata]
key         := string
value       := string
```

#### For example
[My feeds](https://github.com/nhomble/fdmi/blob/master/feeds)
```text
https://nullprogram.com/feed/ daysOld=30
http://xkcd.com/atom.xml
https://blog.codinghorror.com/rss/
```

### Feed metadata

|Key|Description|
|---|-----------|
|daysOld|positive number to indicate number of days old to respect this feed|
|limit|positive number to indicate a limit to the number of entries to pull for a feed|

## Options
```text
-templateOverride [relative path to go template]
    // point to your own go templates for formatting
    $ go run feed.go -templateOverride ./index.tpl 

-numWorkers [natural number]
    // control the amount of parallelism in fetching feeds
    $ go run feed.go -numWorkers 5
```