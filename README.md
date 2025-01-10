# skectl

`skectl` is a command-line tool similar to OpenShift CLI (oc), providing cluster login and basic Kubernetes resource management functionality.

## Features

- Cluster login (similar to `oc login`)
- Support username/password login
- Support token-based login
- Support interactive input
- Get resource information (similar to `kubectl get`)
- Support skipping TLS verification
- Support multi-cluster configuration management and context switching

## Installation

```bash
go install github.com/withlin/oc-demo@latest
```

## Usage

### Login to cluster

```bash
# Login with username and password (interactive input)
skectl login https://api.cluster.example.com:6443
```

## Development

### Requirements

- Go 1.20 or higher

### Build

```bash
go build -o skectl
```

### Run Tests

```bash
go test ./...
```

## License

MIT License
