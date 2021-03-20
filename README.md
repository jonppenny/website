# My personal website
A personal blog website written in Go. Mainly used for learning Go.

## Setup
1. Create database, DB user and grant privileges. Create the DB schema based on the readme in the `/docs` directory.
2. Rename `config.yaml.example` to `config.yaml` or use the `-config` flag when running the app to specify a location/filename.
3. Build the binary by running `go build ./cmd/webapp`.
4. Run the binary with `./webapp`.
5. In your browser, go to `http://localhost:9990` unless you used the `-addr` flag to specify a port number.

## Flags
| Flag | Description
|---|---|
| -addr | specify a port number |
| -config | directory/filename.yaml containing the config options |