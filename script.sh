# Hello Dapr from Dagger
dagger -m github.com/shykes/daggerverse/hello@v0.1.2 call hello --giant --greeting="Let's Start" --name="Dapr Workflow Demo"

# Run Dapr Locally
dapr run -f dapr.yaml

# Launch Dapr Debugging Processes
# (vscode)

http://localhost:3001/

# Build and Push Container Images
dagger call build-push --dapr-dir="." --service-dir="./services/" --path=dapr-k8s.yaml

# Deploy to Kubernetes
dapr run -f dapr-k8s.yaml -k