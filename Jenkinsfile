pipeline {
    // install golang 1.14 on Jenkins node
    agent any
    tools {
        go 'golang'
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
    }
	 
    stages {
        stage("build") {
            steps {
                echo 'BUILD EXECUTION STARTED'
                sh 'go version'
                sh 'go get ./...'
                sh 'go build -o myapp'
            }
        }
		stage('Copy build directory') {
            steps {
			    //sh './myapp'
			    sh 'cp -r build /'

            }
        }
    }
}

