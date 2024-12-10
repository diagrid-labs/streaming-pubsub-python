#!/bin/bash

# Default context
DEFAULT_CONTEXT="kind-kind"

# Set context from argument or use default
CONTEXT=${1:-$DEFAULT_CONTEXT}

# Set namespace to the specified context
echo "Setting kubectl context to $CONTEXT..."
kubectl config use-context $CONTEXT

echo "Checking Dapr status in the cluster..."
DAPR_STATUS=$(dapr status -k 2>&1)

set -e

if echo "$DAPR_STATUS" | grep -q "No status returned. Is Dapr initialized in your cluster?"; then
  echo "Dapr is not installed. Installing Dapr..."
  dapr init -k --enable-ha --runtime-version 1.14.4
   # Wait for all pods in the dapr-system namespace to be ready
  echo "Waiting for Dapr pods to be ready..."
  while [[ $(kubectl get pods -n dapr-system -o 'jsonpath={..status.conditions[?(@.type=="Ready")].status}' | grep -o "True" | wc -l) -ne $(kubectl get pods -n dapr-system --no-headers | wc -l) ]]; do
    echo "Waiting for Dapr pods to be ready..."
    sleep 2
  done
else
  echo "Dapr is already installed"
fi

# Build app images and load them in the kind cluster
echo "Building and loading app images..."

docker build -t publisher:latest -f publisher/Dockerfile publisher
kind load docker-image publisher:latest

docker build -t regular-subscriber:latest -f regular-pubsub/Dockerfile regular-pubsub
kind load docker-image regular-subscriber:latest

docker build -t streaming-subscriber:latest -f streaming-pubsub/Dockerfile streaming-pubsub
kind load docker-image streaming-subscriber:latest

# Install Redis
echo "Installing Redis..."
if helm ls --all --short | grep -q '^redis$'; then
  echo "Redis is already installed"
else
  helm install redis bitnami/redis
fi


# Create secret with the password if it doesn't exist
echo "Creating Redis secrets..."
REDIS_PASSWORD=$(kubectl get secret --namespace default redis -o jsonpath="{.data.redis-password}" | base64 -d)
if kubectl get secret redis-secret -n $ns > /dev/null 2>&1; then
  echo "Secret redis-secret already exists in namespace $ns"
else
  kubectl create secret generic redis-secret --from-literal=redis-password=$REDIS_PASSWORD
fi

# Deploy the apps
echo "Deploying the apps..."
kubectl apply -f components/pubsub-component.yaml
kubectl apply -f deploy/publisher.yaml
kubectl apply -f deploy/regular-subscriber.yaml
kubectl apply -f deploy/streaming-subscriber.yaml

kubectl rollout restart deployment