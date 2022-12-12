#!/usr/bin/env sh

ENDPOINT="http://localhost:8000/rss.xml"
until $(curl --output /dev/null --silent --head --fail $ENDPOINT); do
  curl $ENDPOINT
  sleep 5
done
echo "Done waiting for mock rss feed"

echo "$ENDPOINT" | go run feed.go >out.html

echo "Data retrieved"
cat out.html

tail -n9 out.html >diff1
tail -n9 integration-test/local-rss/output.html >diff2

cat out.html

DIFF=$(diff diff1 diff2)
if [ "$DIFF" != "" ]; then
  echo "THERE WAS A DIFF!"
  echo $DIFF
  exit 1
fi

exit 0
