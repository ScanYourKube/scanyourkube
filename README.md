# ScanYourKube
[![SonarCloud](https://sonarcloud.io/images/project_badges/sonarcloud-black.svg)](https://sonarcloud.io/summary/new_code?id=ScanYourKube_scanyourkube)

This software provides the possibility to automatically update the container version of Deployments and StatefulSets in case that they are affected by a CVE on your Kubernetes cluster. It uses Kubeclarity to scan for vulnerabilities and uses Keel.sh to update the container versions to the newest available one.

## Purpose
This project was built during my bachelor thesis. The goal was to implement a system to scan daily the kuberentes cluster for vulnerabilities and automatically update the vulnerable containers to a newer version. As the system under test was a Rancher Kubernetes cluster, a lot of the functionallity is build for it's CRD's. 

## Deployment
### Install using helm
1. Save values.yaml to default file
```
helm show values oci://ghcr.io/scanyourkube/scanyourkube > values.yaml
```

2. Install ScanYourKube on the Kubernetes cluster
```
helm install --values values.yaml --create-namespace scanyourkube oci://ghcr.io/scanyourkube/scanyourkube -n scanyourkube
```


## Roadmap 🚧
- [ ] Refactor code to a more general version, which can be used on more Kubernetes clusters
- [ ] Refactor the updating component and implement own version




