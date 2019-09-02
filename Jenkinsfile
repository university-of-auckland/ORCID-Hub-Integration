pipeline {
    agent {label("uoa-buildtools-small")}
    environment {
        GOPATH = "$WORKSPACE/.go"
	CGO_ENABLED = "0"
	GO111MODULE = "on"
        PATH = "$WORKSPACE/.go/bin:$WORKSPACE/bin:$WORKSPACE/go/bin:$PATH"
    }

    stages {
        stage('Setup') {
            steps {
		sh '.jenkins/install.sh'
		sh 'env'
		sh './go/bin/go env'
            }
        }
        stage('Test') {
            steps {
		sh 'gotest -tags test ./handler/...'
            }
        }
        stage('Build') {
            steps {
		sh 'go vet ./handler'
		sh 'go vet -tags test ./handler'
		sh 'golint ./handler'
		sh 'go get -u golang.org/x/tools/cmd/cover github.com/mattn/goveralls golang.org/x/lint/golint github.com/rakyll/gotest'
		sh 'go build -o main ./handler/ && upx main && zipit'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Stay Tunded....'
            }
        }
    }
}
