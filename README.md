# Infra Poc

This proofs the concept of creating a master control plane that can create multiple control planes.

## Install Crossplane

```bash
helm repo add crossplane-stable https://charts.crossplane.io/stable && helm repo update

helm install crossplane \
--namespace crossplane-system \
--create-namespace crossplane-stable/crossplane

kubectl get deployments -n crossplane-system
```

## GCP

Install the provider

```bash
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-gcp
spec:
  package: xpkg.upbound.io/upbound/provider-gcp:v0.29.0
```

Create the secret

```bash
kubectl create secret \
generic gcp-secret \
-n crossplane-system \
--from-file=creds=./gcp/gcp-credentials.json

kubectl describe secret gcp-secret -n crossplane-system
```

Create a provider config

```bash

apiVersion: gcp.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: gcp-pc
spec:
  projectID: <PROJECT_ID>
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: gcp-secret
      key: creds
```

Create public cluster

```bash
kubectl create -f gcp/gke-cluster.yaml
```

Create private cluster

```bash
kubectl create -f gcp/gke-private-cluster.yaml
```

## AWS

Install the provider

```bash
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-aws
spec:
  package: xpkg.upbound.io/upbound/provider-aws:v0.31.0
```

Create the secret

```bash
kubectl create secret \
generic aws-secret \
-n crossplane-system \
--from-file=creds=./aws/aws-credentials.txt
```

Create a provider config

```bash
apiVersion: aws.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: aws-pc
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: aws-secret
      key: creds
```

Create simple cluster

```bash
kubectl create -f aws/aws-simple-cluster.yaml
```

## AKS

Install the provider

```bash
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-azure
spec:
  package: xpkg.upbound.io/upbound/provider-azure:v0.29.0
```

Create Credentials

```bash
az login


az ad sp create-for-rbac \
--sdk-auth \
--role Owner \
--scopes /subscriptions/<subscription-id></subscription-id>
```

Save output to ./aks/azure-credentials.json

Create the secret

```bash
kubectl create secret \
generic azure-secret \
-n crossplane-system \
--from-file=creds=./aks/azure-credentials.json

kubectl describe secret azure-secret -n crossplane-system
```

Create a provider config

```bash
apiVersion: azure.upbound.io/v1beta1
metadata:
  name: azure-pc
kind: ProviderConfig
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: azure-secret
      key: creds
```

Create simple cluster

```bash
kubectl create -f aks/aks-simple-cluster.yaml
```

## CIVO

Install the provider

```bash
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-civo
spec:
  package: xpkg.upbound.io/upbound/provider-civo:v0.1
```

Create Credentials

Get api key from CIVO account at [Settings > Profile > Security](https://dashboard.civo.com/security).

Create the secret

```bash
kubectl create secret \
generic civo-secret \
-n crossplane-system \
--from-literal=credentials=<api-key>

kubectl describe secret civo-secret -n crossplane-system
```

Create a provider config

```bash
apiVersion: civo.upbound.io/v1beta1
metadata:
  name: civo-pc
kind: ProviderConfig
spec:
  region: lon1
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: civo-secret
      key: credentials
```

Create simple cluster

```bash
kubectl create -f civo/civo-simple-cluster.yaml
```
