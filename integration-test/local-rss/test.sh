#!/usr/bin/env sh

clean() {
  docker kill test
  rm -f diff1
  rm -f diff2
  rm -f out.html
}

docker run -p 8000:8000 --name test -d feed-cli-test:0.1.0

ENDPOINT="http://localhost:8000/rss.xml"
until "$(curl --output /dev/null --silent --head --fail "$ENDPOINT")"; do
  curl "$ENDPOINT"
  docker logs test
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
  echo "$DIFF"
  clean
  exit 1
else
  clean
fi
