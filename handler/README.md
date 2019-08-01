[![Build Status](https://travis-ci.org/nad2000/consumer.svg?branch=master)](https://travis-ci.org/nad2000/consumer)
[![Coverage Status](https://coveralls.io/repos/github/nad2000/consumer/badge.svg?branch=master)](https://coveralls.io/github/nad2000/consumer?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nad2000/consumer)](https://goreportcard.com/report/github.com/nad2000/consumer)




# Consumer simple message agnostic event consumer

a sample Lambda...
It can be triggered either by SQS or API Gateway directly.

## Event Message Flow
![ScreenShot](/handler/flow.png?raw=true "Message Flow")


# Building

```
go build -o main . && upx main && zip main.zip main
```

# Testing

```sh
export API_KEY=... CLIENT_ID=... CLIENT_SECRET=...
go test -v .
```

Or create **.env**-file with environment variables.

```ini

# Copy this file to '.env'-file and change values
# ORCID Hub API client credentials:
CLIENT_ID=...
CLIENT_SECRET=...
# UoA API Key:
API_KEY=...

```

