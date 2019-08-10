[![Build Status](https://travis-ci.org/university-of-auckland/ORCID-Hub-Integration.svg?branch=master)](https://travis-ci.org/university-of-auckland/ORCID-Hub-Integration)
[![Coverage Status](https://coveralls.io/repos/github/university-of-auckland/ORCID-Hub-Integration/badge.svg?branch=master)](https://coveralls.io/github/university-of-auckland/ORCID-Hub-Integration?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/university-of-auckland/ORCID-Hub-Integration)](https://goreportcard.com/report/github.com/university-of-auckland/ORCID-Hub-Integration)

# AWS Lambda Based Event Handlder

a sample Lambda...
It can be triggered either by SQS or API Gateway directly.

### Event Message Flow
![ScreenShot](/handler/flow.png?raw=true "Message Flow")


## Building

To deploy on AWS Lambda:

```
go build -o main ./handler/ && upx main && zip main.zip main
```

Docker image:

```sh 

docker build -t handler . 

```


## Testing

```sh

export API_KEY=... CLIENT_ID=... CLIENT_SECRET=...
gotest -v .

# more verbose:
gotest -v . -args -verbose

# with 'live server' instead of using the mock:
gotest -v . -args -live

# to get coverage report (user "-tags test" to exclude Lambda specific bits from the coverage):
gotest ./... -tags test -cover -coverprofile coverage.out  ; go tool cover -html=coverage.out -o coverage.html

```

Or create **.env**-file with environment variables.

```ini

# Copy this file to '.env'-file and change values
# ORCID Hub API client credentials:
CLIENT_ID=...
CLIENT_SECRET=...
# UoA API Key:
API_KEY=...
# PORT on which to server the handler (only for Docker)
PORT=5000

```

## Runing Docker

```sh 

# you need to create **.env** file...
docker run -it --env-file .env -p "5050:5050" handler

```
