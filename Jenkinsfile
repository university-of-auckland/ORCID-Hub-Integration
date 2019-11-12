pipeline {
  agent {label("uoa-buildtools-small")}
  environment {
    LANG= "en_US.UTF-8"
    LANGUAGE = "en_US"
    // LC_ALL = "en_US.UTF-8"
    AWS_DEFAULT_REGION = (EVN == "test" ? "us-east-1" : "ap-southeast-2")
    CGO_ENABLED = "0"
    COMMIT_MESSAGE = sh([ script: 'git log -1 --pretty=%B', returnStdout: true ]).trim()
    GO111MODULE = "on"
    GOPATH = "$WORKSPACE/.go"
    JAVA_OPTS = "-Dsun.jnu.encoding=UTF-8 -Dfile.encoding=UTF-8"
    PATH = "$WORKSPACE/.go/bin:$WORKSPACE/bin:$WORKSPACE/go/bin:$PATH"
    TF_CLI_ARGS = "-no-color"
    TF_CLI_ARGS_input = "false"
    TF_CLI_ARGS_refresh = "true"
    TF_INPUT = "0"
  }

  stages {
    // Imports artifacts if build was previously successful
    stage('Import Artifacts') {
      steps {
        copyArtifacts filter: '*.tar.gz,main.zip', fingerprintArtifacts: true, optional: true, projectName: 'integration-orcidhub-build-deploy' // , selector: lastSuccessful()
        // copyArtifacts filter: 'terraform.tar.gz,binary.tar.gz', fingerprintArtifacts: true, optional: false, projectName: 'integration-orcidhub-build-deploy' // , selector: lastSuccessful()
      }
    }
    /* stage('SETUP') {
      steps {
	// sh 'tar xf ./binaries.tar.gz || true'
        sh '.jenkins/install.sh'
	// sh 'go version; go env; env; locale'
        // sh 'tar czf binaries.tar.gz ./.go ./go ./bin'
        // archiveArtifacts artifacts: 'binaries.tar.gz', onlyIfSuccessful: false
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
        sh 'staticcheck ./handler'
        sh 'go build -o main ./handler/ && upx main && zip -0 main.zip main'
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
	  if (env.PROVISION == 'true' || COMMIT_MESSAGE.toUpperCase().contains("[PROVISION]")) {
	     sh 'tar xf ./terraform.tar.gz || true'
             // sh 'terraform version'
             sh '.jenkins/terraform.sh'
	     dir("deployment") {
               sh "terraform init || true"
               sh "terraform workspace new ${ENV} || terraform workspace select ${ENV}"
	       sh "terraform plan"
	       // Destruction if checked RECREATE or the commit message contains '[RECREATE]'
	       if (env.RECREATE == 'true' || COMMIT_MESSAGE.toUpperCase().contains("[RECREATE]")) {
	         sh '"$WORKSPACE/deployment/purge.sh"'
                 sh 'aws sts get-caller-identity --query Account --output text'
		 sh 'terraform destroy -auto-approve'
	         sh '"$WORKSPACE/deployment/destroy.sh"'
	       }
	       // Provision and deploy the handler
	       sh "terraform apply -auto-approve"
	       sh "terraform output"
	       sh '"$WORKSPACE/deployment/create.sh"'
	     }
             sh 'tar czf terraform.tar.gz ./deployment/terraform.tfstate* ./deployment/.terraform'
	  } else {
	    // Deploy the handler to already provisioned environment
	    sh "aws lambda update-function-code --function-name ORCIDHUB_INTEGRATION_${ENV} --publish --zip-file 'fileb://$WORKSPACE/main.zip'"
	  }
          archiveArtifacts artifacts: 'terraform.tar.gz', onlyIfSuccessful: false
          archiveArtifacts artifacts: 'main.zip', fingerprint: true
	}
      }
    }
  }
}
