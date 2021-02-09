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

		stage('Build') {
			steps {
				sh 'go build -ldflags="-s -w" -o bin/osuapi-linux-amd64 cmd/cmd.go'
				sh 'tar -zcvf bin/osuapi-linux-amd64.tar.gz bin/osuapi-linux-amd64'
			}
		}

		stage('Archive') {
			steps {
				archiveArtifacts 'bin/*.tar.gz'
			}
		}
	}
}
