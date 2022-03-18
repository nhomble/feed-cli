FROM python:3.10.3

WORKDIR /integration-test
COPY rss.xml .

EXPOSE 8000:8000
CMD ["python", "-m", "http.server", "8000" ]