# ElasticSearch Slack Indexer

A helper tool for forwarding Slack messages to ElasticSearch backend

## Building

The project uses `gb` tool, install it with: `go get
github.com/constabulary/gb/...`

Once `gb` and `gb-vendor` are available run:

```
gb vendor restore
gb build
```

Binaries will end up in `bin/` directory.

## Configuration

An example config file `eslacki.conf` is provided in the tree. Adjust
it to your liking.

Configuration:

```
[config]
# your slack token
token =
# elastic search URL
url = http://localhost:9200/
# elastic search index to use
index =
```

### Tokens

The app requires an OAuth token to access the RTM API. Generate one
here: https://api.slack.com/docs/oauth-test-tokens

## Running


```
./bin/eslacki -config <path-to-config-file>
./bin/eslacki -config ./eslacki.conf
```

## Other

`cmd/esimport` is an example helper tool for importing Slack log files
collected with other programs (ex. WeeChat).

## License

MIT
