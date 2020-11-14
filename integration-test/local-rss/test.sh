#!/usr/bin/env sh

echo "http://localhost:8000/rss.xml" | go run feed.go >out.html
tail -n10 out.html >diff1
tail -n10 integration-test/local-rss/output.html >diff2

DIFF=$(diff diff1 diff2)
if [ "$DIFF" != "" ]; then
  echo "THERE WAS A DIFF!"
  echo $DIFF
  exit 1
fi
