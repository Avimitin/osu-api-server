pipeline {
	agent {
		docker {
			image 'golang:1.15-buster'
		}
	}
	stages{
		stage('Setup') {
			steps {
				sh 'go get ./...'
			}
		}

		stage('Test') {
			steps {
				sh 'go test -v ./...'
			}
		}
	}
}
