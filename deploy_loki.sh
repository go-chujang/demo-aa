#!/bin/bash
set -e

echo "applying persistentVolumeClaims..."
kubectl apply -f ./.microk8s/loki/pvcs.yaml

echo "applying configMaps..."
kubectl apply -f ./.microk8s/loki/configmap.yaml

echo "applying rbac configurations..."
kubectl apply -f ./.microk8s/loki/rbac.yaml

echo "deploying loki..."
kubectl apply -f ./.microk8s/loki/loki.yaml

echo "deploying promtail..."
kubectl apply -f ./.microk8s/loki/promtail.yaml

echo "deploying grafana..."
kubectl apply -f ./.microk8s/loki/grafana.yaml

echo "deployment loki-promtail-grafana completed."