#!/usr/bin/env sh

docker run feed-cli-test:0.1.0 &
sleep 5

echo "http://localhost:8000/rss.xml" | go run feed.go >out.html
tail -n9 out.html >diff1
tail -n9 integration-test/local-rss/output.html >diff2

cat out.html

DIFF=$(diff diff1 diff2)
if [ "$DIFF" != "" ]; then
  echo "THERE WAS A DIFF!"
  echo $DIFF
  exit 1
fi
