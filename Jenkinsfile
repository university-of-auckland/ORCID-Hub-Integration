pipeline {
    agent any
    environment {
        GOPATH = "$WORKSPACE/.go"
	CGO_ENABLED = "0"
        PATH = "$WORKSPACE/.go/bin:$WORKSPACE/bin:$WORKSPACE/go/bin:$PATH"
    }

    stages {
        stage('Build') {
            steps {
	    	echo '=== Install/Upgrade'
		sh '.jenkins/install.sh'
		sh 'go version'
		sh 'go get -u golang.org/x/tools/cmd/cover github.com/mattn/goveralls golang.org/x/lint/golint github.com/rakyll/gotest'
		sh 'go build -o main ./handler/ && upx main && zipit'
            }
        }
        stage('Test') {
            steps {
		sh 'gotest -tags test -v ./handler/...'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Stay Tunded....'
            }
        }
    }
}
