# Kubernetes Deployment Guide

This directory contains Kubernetes manifests for deploying the Passkey Demo application.

## Prerequisites

- Kubernetes cluster (v1.24+)
- kubectl configured
- NGINX Ingress Controller installed
- cert-manager installed (for TLS certificates)
- Metrics Server installed (for HPA)

## Quick Start

### 1. Build and Push Docker Image

```bash
# Build the image
docker build -t your-registry/passkey-demo:latest .

# Push to registry
docker push your-registry/passkey-demo:latest

# Update deployment.yaml with your image
```

### 2. Update Configuration

Edit [`configmap.yaml`](configmap.yaml:9) to set your domain:

```yaml
data:
  RP_ID: "your-domain.com"
  RP_ORIGIN: "https://your-domain.com"
```

Edit [`ingress.yaml`](ingress.yaml:16) to set your domain:

```yaml
spec:
  tls:
  - hosts:
    - your-domain.com
```

### 3. Deploy to Kubernetes

```bash
# Apply all manifests
kubectl apply -f k8s/

# Or apply in order
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
kubectl apply -f k8s/hpa.yaml
```

### 4. Verify Deployment

```bash
# Check namespace
kubectl get namespace passkey-demo

# Check all resources
kubectl get all -n passkey-demo

# Check pods
kubectl get pods -n passkey-demo

# Check service
kubectl get svc -n passkey-demo

# Check ingress
kubectl get ingress -n passkey-demo

# Check HPA
kubectl get hpa -n passkey-demo
```

### 5. View Logs

```bash
# View logs from all pods
kubectl logs -n passkey-demo -l app=passkey-demo --tail=100 -f

# View logs from specific pod
kubectl logs -n passkey-demo <pod-name> -f
```

## Manifest Files

### [`namespace.yaml`](namespace.yaml:1)
Creates the `passkey-demo` namespace for resource isolation.

### [`configmap.yaml`](configmap.yaml:1)
Contains application configuration:
- `RP_ID`: Relying Party ID (your domain)
- `RP_ORIGIN`: Relying Party Origin (your full URL)
- `PORT`: Application port (default: 8080)

### [`deployment.yaml`](deployment.yaml:1)
Defines the application deployment:
- **Replicas**: 3 (for high availability)
- **Resources**: 
  - Requests: 64Mi memory, 100m CPU
  - Limits: 128Mi memory, 200m CPU
- **Probes**: Liveness and readiness checks
- **Security**: Non-root user, dropped capabilities

### [`service.yaml`](service.yaml:1)
Creates a ClusterIP service:
- Exposes port 80 internally
- Routes to container port 8080
- Load balances across pods

### [`ingress.yaml`](ingress.yaml:1)
Configures external access:
- NGINX Ingress Controller
- TLS/HTTPS enabled
- cert-manager integration
- Domain routing

### [`hpa.yaml`](hpa.yaml:1)
Horizontal Pod Autoscaler:
- Min replicas: 3
- Max replicas: 10
- CPU target: 70%
- Memory target: 80%
- Smart scaling policies

## Scaling

### Manual Scaling

```bash
# Scale to 5 replicas
kubectl scale deployment passkey-demo -n passkey-demo --replicas=5

# Check scaling status
kubectl get deployment passkey-demo -n passkey-demo
```

### Automatic Scaling (HPA)

The HPA automatically scales based on:
- CPU utilization (target: 70%)
- Memory utilization (target: 80%)

```bash
# Watch HPA in action
kubectl get hpa -n passkey-demo -w

# Describe HPA for details
kubectl describe hpa passkey-demo-hpa -n passkey-demo
```

## Monitoring

### Check Pod Status

```bash
# Get pod details
kubectl describe pod <pod-name> -n passkey-demo

# Check pod events
kubectl get events -n passkey-demo --sort-by='.lastTimestamp'
```

### Check Resource Usage

```bash
# View resource usage
kubectl top pods -n passkey-demo

# View node resource usage
kubectl top nodes
```

### Access Application

```bash
# Port forward for local testing
kubectl port-forward -n passkey-demo svc/passkey-demo-service 8080:80

# Then access: http://localhost:8080
```

## Troubleshooting

### Pods Not Starting

```bash
# Check pod status
kubectl get pods -n passkey-demo

# Describe pod for events
kubectl describe pod <pod-name> -n passkey-demo

# Check logs
kubectl logs <pod-name> -n passkey-demo
```

### Image Pull Errors

```bash
# Check if image exists
docker pull your-registry/passkey-demo:latest

# Create image pull secret if needed
kubectl create secret docker-registry regcred \
  --docker-server=your-registry \
  --docker-username=your-username \
  --docker-password=your-password \
  -n passkey-demo

# Add to deployment.yaml:
# spec:
#   imagePullSecrets:
#   - name: regcred
```

### Ingress Not Working

```bash
# Check ingress status
kubectl describe ingress passkey-demo-ingress -n passkey-demo

# Check ingress controller logs
kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx

# Verify DNS points to ingress IP
kubectl get ingress -n passkey-demo
```

### HPA Not Scaling

```bash
# Check metrics server
kubectl get deployment metrics-server -n kube-system

# Check HPA status
kubectl describe hpa passkey-demo-hpa -n passkey-demo

# Generate load to test
kubectl run -it --rm load-generator --image=busybox -n passkey-demo -- /bin/sh
# Then: while true; do wget -q -O- http://passkey-demo-service; done
```

## Updating the Application

### Rolling Update

```bash
# Update image
kubectl set image deployment/passkey-demo \
  passkey-demo=your-registry/passkey-demo:v2 \
  -n passkey-demo

# Watch rollout
kubectl rollout status deployment/passkey-demo -n passkey-demo

# Check rollout history
kubectl rollout history deployment/passkey-demo -n passkey-demo
```

### Rollback

```bash
# Rollback to previous version
kubectl rollout undo deployment/passkey-demo -n passkey-demo

# Rollback to specific revision
kubectl rollout undo deployment/passkey-demo --to-revision=2 -n passkey-demo
```

## Cleanup

### Delete All Resources

```bash
# Delete all resources in namespace
kubectl delete -f k8s/

# Or delete namespace (removes everything)
kubectl delete namespace passkey-demo
```

### Delete Specific Resources

```bash
# Delete deployment
kubectl delete deployment passkey-demo -n passkey-demo

# Delete service
kubectl delete service passkey-demo-service -n passkey-demo

# Delete ingress
kubectl delete ingress passkey-demo-ingress -n passkey-demo
```

## Production Considerations

### Security

1. **Use Private Registry**: Store images in a private container registry
2. **Network Policies**: Implement network policies to restrict traffic
3. **RBAC**: Configure proper role-based access control
4. **Secrets Management**: Use Kubernetes secrets or external secret managers
5. **Pod Security**: Enable Pod Security Standards

### High Availability

1. **Multiple Replicas**: Run at least 3 replicas
2. **Pod Disruption Budget**: Prevent too many pods from being down
3. **Node Affinity**: Spread pods across different nodes/zones
4. **Resource Limits**: Set appropriate resource requests and limits

### Monitoring & Logging

1. **Prometheus**: Collect metrics
2. **Grafana**: Visualize metrics
3. **ELK/EFK Stack**: Centralized logging
4. **Alerting**: Set up alerts for critical issues

### Backup & Recovery

1. **Database Backups**: If using persistent storage
2. **Configuration Backups**: Version control all manifests
3. **Disaster Recovery**: Document recovery procedures

## Example: Complete Deployment

```bash
#!/bin/bash

# 1. Build and push image
docker build -t myregistry/passkey-demo:v1.0.0 .
docker push myregistry/passkey-demo:v1.0.0

# 2. Update manifests
sed -i 's|passkey-demo:latest|myregistry/passkey-demo:v1.0.0|' k8s/deployment.yaml
sed -i 's|passkey.example.com|myapp.com|g' k8s/configmap.yaml k8s/ingress.yaml

# 3. Deploy
kubectl apply -f k8s/

# 4. Wait for rollout
kubectl rollout status deployment/passkey-demo -n passkey-demo

# 5. Verify
kubectl get all -n passkey-demo

# 6. Test
curl https://myapp.com/health
```

## Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/)
- [cert-manager](https://cert-manager.io/)
- [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)