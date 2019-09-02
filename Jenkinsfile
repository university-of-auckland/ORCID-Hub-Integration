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
      }
    }
    stage('Test') {
      steps {
        sh 'gotest -tags test ./handler/...'
          sh 'gotestsum --junitfile tests.xml -- -v -tags test ./handler/...'
      }
    }
    stage('Build') {
      steps {
        sh 'go vet ./handler'
          sh 'go vet -tags test ./handler'
          sh 'golint ./handler'
          sh 'go build -o main ./handler/ && upx main && zipit'
      }
    }
    stage('Deploy') {
      steps {
        echo 'Stay Tunded....'
      }
    }
  }

  post {
    always {
    archiveArtifacts artifacts: 'mail.zip', fingerprint: true
      junit 'tests.xml'
    }
  }
}
