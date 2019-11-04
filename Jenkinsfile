pipeline {
  agent {label("uoa-buildtools-small")}
  environment {
    GOPATH = "$WORKSPACE/.go"
    CGO_ENABLED = "0"
    GO111MODULE = "on"
    PATH = "$WORKSPACE/.go/bin:$WORKSPACE/bin:$WORKSPACE/go/bin:$PATH"
    AWS_DEFAULT_REGION = "ap-southeast-2"
    TF_INPUT = "0"
    TF_CLI_ARGS = "-no-color"
    TF_CLI_ARGS_input = "false"
    TF_CLI_ARGS_refresh = "true"
    JAVA_OPTS = "-Dsun.jnu.encoding=UTF-8 -Dfile.encoding=UTF-8"
    LANG= "en_US.UTF-8"
    LANGUAGE = "en_US"
    LC_ALL = "en_US.UTF-8"
  }

  stages {
    // Imports artifacts if build was previously successful
    stage('Import Artifacts') {
      steps {
        copyArtifacts filter: '*.tar.gz', fingerprintArtifacts: true, optional: true, projectName: 'integration-orcidhub-build-deploy', selector: lastSuccessful()
	sh 'tar xf ./binaries.tar.gz || true'
	sh 'tar xf ./terraform.tar.gz || true'
      }
    }
    stage('SETUP') {
      steps {
        sh '.jenkins/install.sh'
	sh 'go version; go env; env'
	sh 'locale'
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
             // sh 'terraform version'
             sh '.jenkins/terraform.sh'
	     dir("deployment") {
	       // workaround to remove a role if it exists:
               sh "terraform init || true"
               // sh "terraform plan"
               sh "terraform workspace new ${ENV} || terraform workspace select ${ENV}"
               // sh "terraform refresh"
               // sh "terraform plan -out ${ENV}.plan"
	
	      // if (env.RECREATE == 'true') {
	      sh './purge.sh'
	        sh 'terraform destroy -auto-approve'
	        sh 'if [ "${RECREATE}" == "true" ] || (git log -1 --pretty=%B | grep -iq \'\\[RECREATE\\]\') ; the ./purge.sh ; terraform destroy -auto-approve; fi'
	      // }
	      // Provision and deploy the handler
	      sh "terraform apply -auto-approve"
	     }
	  // } else {
	    // // Deploy the handler to already provisioned environment
	    // sh "aws lambda update-function-code --function-name ORCIDHUB_INTEGRATION --publish --zip-file 'fileb://$WORKSPACE/main.zip' --region=ap-southeast-2"
	  // }
	}
      }
    }
    // Archive what was achieved, even if unsuccessful so the next run understands even partial components
    stage('Archive Artifacts') {
      steps {
        sh 'tar czf binaries.tar.gz ./.go ./go ./bin'
        sh 'tar czf terraform.tar.gz ./deployment/terraform.tfstate* ./deployment/.terraform'
        archiveArtifacts artifacts: 'archive.tar.gz', onlyIfSuccessful: false
        // archiveArtifacts artifacts:  '.go/**,go/**,bin/**,**/terraform.tfstate.d/**,**/terraform.tfstate,deployment/.terraform/**', onlyIfSuccessful: false
      }
    }
  }
}
