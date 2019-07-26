[![Build Status](https://travis-ci.com/nad2000/consumer.svg?branch=master)](https://travis-ci.com/nad2000/consumer)
[![Coverage Status](https://coveralls.io/repos/github/nad2000/consumer/badge.svg?branch=master)](https://coveralls.io/github/nad2000/consumer?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nad2000/consumer)](https://goreportcard.com/report/github.com/nad2000/consumer)




# Consumer simple message agnostic event consumer

a sample Lambda...
It can be triggered either by SQS or API Gateway directly.

## Event Message Flow
![ScreenShot](/flow.png?raw=true "Message Flow")


# Building

```
go build . && upx consumer && zip consumer.zip consumer
```

# Testing

```
export API_KEY=... CLIENT_ID=... CLIENT_SECRET=...
go test -v .
```
