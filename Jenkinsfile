pipeline {
  agent {label("uoa-buildtools-small")}
  environment {
    GOPATH = "$WORKSPACE/.go"
      CGO_ENABLED = "0"
      GO111MODULE = "on"
      PATH = "$WORKSPACE/.go/bin:$WORKSPACE/bin:$WORKSPACE/go/bin:$PATH"
  }

  stages {
    /*stage('SETUP') {
      steps {
        sh '.jenkins/install.sh'
	sh 'go version'
	sh 'go env'
	sh 'env'
      }
    }
    stage('TEST') {
      steps {
        // sh 'gotest -tags test ./handler/...'
        sh 'gotestsum --junitfile tests.xml -- -v -tags test ./handler/...'
        junit 'tests.xml'
      }
    }
    stage('BUILD') {
      steps {
        sh 'go vet ./handler'
        sh 'go vet -tags test ./handler'
        sh 'golint ./handler'
        sh 'go build -o main ./handler/ && upx main && zipit'
        archiveArtifacts artifacts: 'main.zip', fingerprint: true
      }
    }
    */
    stage('AWS Credential Grab') {
      steps{
        print "☯ Authenticating with AWS"
        withCredentials([usernamePassword(credentialsId:"aws-user-sandbox", passwordVariable: 'password', usernameVariable: 'username'), string(credentialsId: "aws-token-sandbox", variable: 'token')]) {
          sh "python3 /home/jenkins/aws_saml_login.py --idp iam.auckland.ac.nz --user $USERNAME --password $PASSWORD --token $TOKEN --profile 'default'"
        }
      }
    }
    stage('DEPLOY') {
      steps {
      	script {
	  // "destroy" provisioned environment 
	  // if (env.PROVISION == 'true') {
             sh 'terraform version'
             sh '.jenkins/terraform.sh'
	     dir("deployment") {
               sh 'terraform version'
               sh "terraform init"
               // sh "terraform plan -no-color"
               sh "terraform workspace new ${ENV} || terraform select ${ENV}"
               sh "terraform refresh -no-color"
               sh "terraform plan -no-color -out ${ENV}.plan"
	      // if (env.RECREATE == 'true') {
	          sh "terraform destroy -no-color -force -refresh=true"
	      // }
	      // Provision and deploy the handler
	      // sh "terraform apply ${ENV}.plan -no-color -auto-approve -refresh=true"
	      sh "terraform apply ${ENV}.plan -no-color"
	     }
	  // } else {
	    // // Deploy the handler to already provisioned environment
	    // sh "aws lambda update-function-code --function-name ORCIDHUB_INTEGRATION --publish --zip-file 'fileb://$WORKSPACE/main.zip' --region=ap-southeast-2"
	  // }
	}
      }
    }
  }
}
