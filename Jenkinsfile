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
      - name: registry-config
        mountPath: /kaniko/.docker/
  volumes:
    - name: registry-config
      configMap:
        name: local-registry-hosting
"""
        }
    }
    environment {
        // Point to the service name of the registry in k8s if possible, or use the node port/localhost mapping
        // Inside the cluster, we can reach the registry at "registry:5000" if we created a service, 
        // OR we use the localhost mapping if Kaniko supports it (often tricky).
        // Best approach for Kind + Local Registry + Kaniko: 
        // Use the host IP or the special registry service name if available.
        // For this demo, let's assume we push to "localhost:5001" which is mapped on nodes, 
        // BUT Kaniko inside pod needs to reach it. 
        // Often "registry.kube-public.svc.cluster.local:5000" is best if we expose it.
        // Let's create a Service for the registry first or use the container IP.
        
        // Simpler: Use the Kind node's IP for the registry (host.docker.internal equivalent?)
        // Let's rely on the registry-configmap we created? 
        // Actually, for simplicity in this specific "Kind + Local Registry" setup:
        // We push to `registry.default.svc.cluster.local:5000` (we need to create a service for the registry container).
        
        IMAGE_REPO = "localhost:5001" 
    }
    stages {
        stage('Build Backend') {
            steps {
                container('kaniko') {
                    sh '/kaniko/executor --context dir://${env.WORKSPACE}/apps/backend --dockerfile dir://${env.WORKSPACE}/apps/backend/Dockerfile --destination localhost:5001/backend:${BUILD_NUMBER} --insecure --skip-tls-verify'
                }
            }
        }
        stage('Build Frontend') {
            steps {
                container('kaniko') {
                    sh '/kaniko/executor --context dir://${env.WORKSPACE}/apps/frontend --dockerfile dir://${env.WORKSPACE}/apps/frontend/Dockerfile --destination localhost:5001/frontend:${BUILD_NUMBER} --insecure --skip-tls-verify'
                }
            }
        }
        stage('Update Manifests') {
            steps {
                script {
                    // Quick & Dirty: simple sed replacement and commit back
                    // In real world: Use a separate git-ops repo or clean PR flow
                    sh """
                        apt-get update && apt-get install -y git
                        git config user.email "jenkins@example.com"
                        git config user.name "Jenkins CI"
                        sed -i 's|image: localhost:5001/backend:.*|image: localhost:5001/backend:${BUILD_NUMBER}|' k8s-manifests/backend/deployment.yaml
                        sed -i 's|image: localhost:5001/frontend:.*|image: localhost:5001/frontend:${BUILD_NUMBER}|' k8s-manifests/frontend/deployment.yaml
                        git add k8s-manifests/
                        git commit -m "Update images to build ${BUILD_NUMBER}"
                        // We need credentials to push. For this local demo, we might skip pushing 
                        // if we don"t have credentials mounted. 
                        // OPTION: Just verify builds for now since we don't have user's git password/token.
                        echo "Skipping git push (No credentials configured). Manifests updated locally in running pod."
                    """
                }
            }
        }
    }
}
