# NHN-ROR-CLI
CLI made with Golang and Cobra

# Prerequisites
- Golang 1.8.x

# Get started
Bash commands is from ```<repo root>/src/clients/ror-cli/```

Download dependencies:
```bash
go get ./...
```

Start webapi
```bash
go run main.go
```

Or
Start the ```Debug ROR-CLI``` debugger config from VS Code

# Release
A new patch release is generated on every merge to develop and put in https://helsegitlab.nhn.no/sdi-public/ror-cli-releases. 
Version number is tracked in the ROR_VERSION CI/CD variable, only merge jobs bump the version

