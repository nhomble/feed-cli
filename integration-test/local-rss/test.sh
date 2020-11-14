#!/usr/bin/env sh

docker run -p 8000:8000 --name test -d feed-cli-test:0.1.0

ENDPOINT="http://localhost:8000/rss.xml"
until $(curl --output /dev/null --silent --head --fail $ENDPOINT); do
  curl $ENDPOINT
  docker logs test
  sleep 5
done

echo "$ENDPOINT" | go run feed.go >out.html
tail -n9 out.html >diff1
tail -n9 integration-test/local-rss/output.html >diff2

cat out.html

DIFF=$(diff diff1 diff2)
if [ "$DIFF" != "" ]; then
  echo "THERE WAS A DIFF!"
  echo $DIFF
  docker kill test
  exit 1
fi

docker kill test
