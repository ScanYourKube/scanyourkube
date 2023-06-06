# ScanYourKube

This software provides the possibility to automatically update the container version of Deployments and StatefulSets in case that they are affected by a CVE on your Kubernetes cluster. It uses Kubeclarity to scan for vulnerabilities and uses Keel.sh to update the container versions to the newest available one.