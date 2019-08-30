pipeline {
    agent any
    environment {
        PATH = "$WORKSPACE/bin:$WORKSPACE/go/bin:$PATH"
    }

    stages {
        stage('Build') {
            steps {
	    	echo '=== Install/Upgrade'
		sh '.jenkins/install.sh'
                echo '=== Building..'
		sh 'go version'
		sh 'go get -u golang.org/x/tools/cmd/cover github.com/mattn/goveralls golang.org/x/lint/golint github.com/rakyll/gotest'
            }
        }
        stage('Test') {
            steps {
                echo '=== Testing..'
		sh 'gotest -v ./...'
            }
        }
        stage('Deploy') {
            steps {
                echo '=== Deploying....'
                echo 'Stay Tunded....'
            }
        }
    }
}
