[![Build Status](https://travis-ci.org/university-of-auckland/ORCID-Hub-Integration.svg?branch=master)](https://travis-ci.org/university-of-auckland/ORCID-Hub-Integration)
[![Coverage Status](https://coveralls.io/repos/github/university-of-auckland/ORCID-Hub-Integration/badge.svg?branch=master)](https://coveralls.io/github/university-of-auckland/ORCID-Hub-Integration?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/university-of-auckland/ORCID-Hub-Integration)](https://goreportcard.com/report/github.com/university-of-auckland/ORCID-Hub-Integration)

# NZ ORICD Hub Integration

A flexible and platform agnostic integration solution that can be deployed either as AWS lambda based solution, stand-alone, stand-alone docker based, or hosted virtually with any PAAS provider, e.g., Heroku, Google Cloud App Engine, Cloud Function etc. The solution based on AWS Lambda can be triggered either by SQS or API Gateway directly.

This project can server as a reference for [NZ ORCID Hub](https://github.com/Royal-Society-of-New-Zealand/NZ-ORCID-Hub) integrators.

### Event Message Flow (AWS Lambda based)

![ScreenShot](https://wiki.auckland.ac.nz/rest/gliffy/1.0/embeddedDiagrams/223c4818-415f-4cd2-971d-951f0728ff53.png "Message Flow")

## Building

To deploy on AWS Lambda:

```
go build -o main ./handler/ && upx main && zip main.zip main
```

Docker image:

```sh 

docker build -t handler . 

```

Stand-alone executable:

```sh 

go build -o server -tags 'standalone'  ./handler/

```

## Testing

```sh

export APIKEY=... CLIENT_ID=... CLIENT_SECRET=...
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
APIKEY=...
# PORT on which to server the handler (only for Docker)
PORT=5000

```

## Running Docker

```sh 

# you need to create **.env** file...
docker run -it --env-file .env -p "5050:5050" handler

```

## Running Stand-Alone

```sh 

export PORT=9090
./server

```
