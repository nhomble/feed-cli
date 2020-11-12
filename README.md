feed-cli
===
![CI](https://github.com/nhomble/feed-cli/workflows/CI/badge.svg)

simple utility for formatting (groups of) feeds

# Usage
```bash
$ echo "https://your.feed.com/feed" > go run feed.go
```

## Template Data
You can supply your own golang template with the ```-templateOverride``` parameter. When defining your own template, there
are several template variables at your disposal:

|Field          |Description|
|-----          |-----------|
|.Now           |Get current time (machine dependent)|
|.NowIn (loc)   |Get current time in desired location|
|.Org           |The organization / blog name / title|
|.Feeds         |An array of entries you can perform a range over|

Within a feed object you have

|Field          |Description|
|----           |-----------|
|.Parent        |A reference back to the feed group |
|.Article       |Article name, title of the post |
|.Link          |URL to the source|
|.Published     |last touched time, latest time between published and updated time|

[for example](https://github.com/nhomble/fdmi/blob/master/index.tpl)

## Feed metadata
The goal has been to keep the feed file simple. At a minimum, you can just provide
a list of feeds delimited by new-lines and you are done. To naturally give some more metadata per feed, the
cli recognizes the following format

Then the grammar looks like:

```text
<feedFile>    ::= <line> || <line>\n
<line>        ::= <feed> || <comment>> || ""
<comment>     ::= # <string>
<feed>        ::= <feedLink> <metadata>
<feedLink>    ::= <url>
<metadata>    ::= <key>=<value> <metadata> | ""
<key>         ::= <string>
<value>       ::= <string>
```

#### For example
[My feeds](https://github.com/nhomble/fdmi/blob/master/feeds)
```text
https://nullprogram.com/feed/ daysOld=30
http://xkcd.com/atom.xml
https://blog.codinghorror.com/rss/
```

### Feed metadata fields

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