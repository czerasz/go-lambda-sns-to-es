# AWS Lambda: SNS to ElasticSearch

## Requirements

- Golang >= 1.11

## Development

Install dependencies:

```bash
go get
```

## Tests

Run tests with:

```bash
go test -cover ./...
```

## Build

Generate artifact `./dist/deployment.zip`:

```bash
./scripts/build
```

## Environment Variables

| Name | Default Value | Description |
| ---- | ------------- | ----------- |
| `ES_URL` | `http://127.0.0.1:9200` | The URL to ElasticSearch |
| `ES_INDEX_TEMPLATE` | `{{ .Prefix }}-{{ .Date.Year }}.{{ .Date.Month }}.{{ .Date.Day }}` | ElasticSearch index template.<br/><br/>Available objects:<br/>- `Env`<br/>- `Date`<br/>- `Prefix` |
| `ES_DOC_TYPE_NAME` | `notification` | Document type used to index the event message |
| `DEBUG` | `false` | Enable debugging mode |

## Testing/Debugging Localy

- start HTTP debugging server:

  ```bash
  docker run --rm \
            --name=http-debugger \
            --publish 8080:80 \
            mendhak/http-https-echo
  ```

- execute lambda function with sample event:

  ```bash
  cat ./test-data/event.json | LOCAL_TEST=true DEBUG=true ES_URL='http://localhost:8080' go run ./src/main.go
  ```