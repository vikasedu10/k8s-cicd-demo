pipeline {
    agent none
    environment {
        IMAGE_REPO = "vikas1412/antigravity-applicatiions"
    }
    stages {
        stage('Build Backend') {
            agent {
                kubernetes {
                    yaml """
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: kaniko
    image: gcr.io/kaniko-project/executor:debug
    command:
    - cat
    tty: true
    volumeMounts:
      - name: docker-config
        mountPath: /kaniko/.docker/
  volumes:
    - name: docker-config
      secret:
        secretName: docker-hub-config
"""
                }
            }
            steps {
                container('kaniko') {
                    sh "/kaniko/executor --context ${env.WORKSPACE}/apps/backend --dockerfile ${env.WORKSPACE}/apps/backend/Dockerfile --destination index.docker.io/${IMAGE_REPO}:backend-${BUILD_NUMBER} --destination index.docker.io/${IMAGE_REPO}:backend-latest"
                }
            }
        }
        stage('Build Frontend') {
            agent {
                kubernetes {
                    yaml """
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: kaniko
    image: gcr.io/kaniko-project/executor:debug
    command:
    - cat
    tty: true
    volumeMounts:
      - name: docker-config
        mountPath: /kaniko/.docker/
  volumes:
    - name: docker-config
      secret:
        secretName: docker-hub-config
"""
                }
            }
            steps {
                container('kaniko') {
                    sh "/kaniko/executor --context ${env.WORKSPACE}/apps/frontend --dockerfile ${env.WORKSPACE}/apps/frontend/Dockerfile --destination index.docker.io/${IMAGE_REPO}:frontend-${BUILD_NUMBER} --destination index.docker.io/${IMAGE_REPO}:frontend-latest"
                }
            }
        }
    }
}
