pipeline {
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
    environment {
        // Docker Hub Repo
        IMAGE_REPO = "vikas1412/antigravity-applicatiions" 
    }
    stages {
        stage('Build Backend') {
            steps {
                container('kaniko') {
                    sh "/kaniko/executor --context dir://${env.WORKSPACE}/apps/backend --dockerfile dir://${env.WORKSPACE}/apps/backend/Dockerfile --destination index.docker.io/${IMAGE_REPO}:backend-${BUILD_NUMBER} "
                }
            }
        }
        stage('Build Frontend') {
            steps {
                container('kaniko') {
                    sh "/kaniko/executor --context dir://${env.WORKSPACE}/apps/frontend --dockerfile dir://${env.WORKSPACE}/apps/frontend/Dockerfile --destination index.docker.io/${IMAGE_REPO}:frontend-${BUILD_NUMBER} "
                }
            }
        }
        stage('Update Manifests') {
            steps {
                script {
                    sh """
                        apt-get update && apt-get install -y git
                        git config user.email "jenkins@example.com"
                        git config user.name "Jenkins CI"
                        sed -i 's|image: .*backend:.*|image: index.docker.io/${IMAGE_REPO}:backend-${BUILD_NUMBER}|' k8s-manifests/backend/deployment.yaml
                        sed -i 's|image: .*frontend:.*|image: index.docker.io/${IMAGE_REPO}:frontend-${BUILD_NUMBER}|' k8s-manifests/frontend/deployment.yaml
                        git add k8s-manifests/
                        git commit -m "Update images to build ${BUILD_NUMBER}"
                        echo "Skipping git push (No credentials configured). Manifests updated locally in running pod."
                    """
                }
            }
        }
    }
}
