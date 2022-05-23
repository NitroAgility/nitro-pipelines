pipeline {
    agent any
    options { disableConcurrentBuilds() }
    stages {
        stage('Build') {
            steps {
                sh 'docker build -t nitroagility/nitro-bitbucket-pipelines:2.$BUILD_NUMBER -f ./bitbucket/Dockerfile .'
            }
        }
        stage('Deploy') {
            environment {
                NITRO_PIPELINES_DEV_GITHUB_USERNAME = credentials('jenkins/nitro/dev/docker-hub/username')
                NITRO_PIPELINES_DEV_GITHUB_PASSWORD = credentials('jenkins/nitro/dev/docker-hub/password')
            }
            steps {
                sh 'echo "$NITRO_PIPELINES_DEV_GITHUB_PASSWORD" | docker login --username "$NITRO_PIPELINES_DEV_GITHUB_USERNAME" --password-stdin'
                sh 'docker push nitroagility/nitro-bitbucket-pipelines:2.$BUILD_NUMBER'
            }
        }
    }
    post {
        success {
            slackSend (color: '#00FF00', message: "SUCCESSFUL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL}) triggered by ${currentBuild.getBuildCauses('hudson.model.Cause$UserIdCause')}")
        }
        failure {
            slackSend (color: '#FF0000', message: "FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL}) triggered by ${currentBuild.getBuildCauses('hudson.model.Cause$UserIdCause')}")
        }
    }
}