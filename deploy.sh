#!/bin/bash
set -e

echo "applying namespace..."
kubectl apply -f ./.microk8s/namespace.yaml

echo "applying configmap..."
kubectl apply -f ./.microk8s/configmap.yaml

echo "applying persistentVolumeClaims..."
kubectl apply -f ./.microk8s/pvcs.yaml

echo "applying stack deployments..."
kubectl apply -f ./.microk8s/dp_stack.yaml

echo "waiting for stack deployments to be ready..."
kubectl rollout status deployment/dp-stack-mongodb -n demo
kubectl rollout status deployment/dp-stack-ganache -n demo
kubectl rollout status deployment/dp-stack-zookeeper -n demo
kubectl rollout status deployment/dp-stack-kafka -n demo

echo "applying stack services..."
kubectl apply -f ./.microk8s/svc_stack.yaml

echo "applying ingress..."
kubectl apply -f ./.microk8s/ingress.yaml

echo "wait for seconds"
sleep 30

echo "applying setup job..."
kubectl apply -f ./.microk8s/setup.yaml

echo "waiting for setup job to complete..."
if kubectl wait --for=condition=complete --namespace=demo job/setup-once --timeout=30s; then
  echo "setup completed."
else
  echo "setup failed."
  exit 1
fi


echo "applying application deployments..."
kubectl apply -f ./.microk8s/dp_app.yaml
kubectl apply -f ./.microk8s/svc_app.yaml

echo "deployment completed."
