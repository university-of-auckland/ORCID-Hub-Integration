pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
	    	echo '=== Install/Upgrade'
		sh '.jenkins/install.sh'
                echo '=== Building..'
		sh 'go version'
            }
        }
        stage('Test') {
            steps {
                echo '=== Testing..'
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
