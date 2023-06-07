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


## Architecture
For the purpose of understanding how the components of the software are communicating with each other, the following C4 container diagram should be used.
![C4 container diagram](https://kroki.io/plantuml/svg/eNp1VMtu2zAQvOsrtjm5QJNe-gFJnRzaoq1huyh6EihqLTGiuAJJxfHfd2lSsuTHyV5RszM7s9Sj88L6vtXZB2Wk7kuE5Zd8Scbju3_orh0IZdDGo8wrr_lgeAqlEpUVLezIwkYK8496-6MvMMtWaB2ZRYlvqKlD-wnunof_d1w8wXgENTmvTAXKO3C083thEciArxFCN2vQowNW5TyjP2bZqGDhmPXArA2_l79Swa2nQuA7FZFOWjJ8DvtayRq8VVXFEo8cR0-AdrFgNAhTQiuMqDC-0Xel8BiYNwfW0OYv734RKKUWVvkDMwS2Zax4tE7ToUXjA_eWG5yeBJ4JkntOWsphrtxipZy3ofPJ7nV6eJzIAEfG3ggNIwwG2OhpUH86TjOO1k_85pE_c4rshSmFLYGkOgHdXGYrlM7dsWYlL_ehhlgnaa71HTi0b0zqif-xoRhec7BXvp63a8aMc9Gp5GUK_Wn1bfBw_hSE1rSP8UiLwiuKuQnNrsQyWZ1AVLyi9GejNIg6MPLPg6sHqiaWoEzwQ6YFGPzsqAxdsjXq6YJfjLEMsphYTFVYdLybMizTGf5q-Kve1WH3xwhVy1vJnkZ8_lXIZtpkHs2GfXeA99H5EfQLVVUXdO32TPFzePBm83O7SsIvoVf1r5GvGssD2Vsb1n8yieHPRntM6mbPFM82Xtd0D0OuIQMOw1Jf1fAXi5qoud1ldlOHZtHJ430P3y8jWnSdkOhSnxnqRjhaX0Szs9Se-xzHuFiQP3GcyWL_Hnc0yx7Zff7q_gevfArb)

Further more information about the different components of the job are described in the following C4 component diagram.
![C4 component diagram](https://kroki.io/plantuml/svg/eNq1Vclu2zAQvesrpj4lQJpe-gFxnRzSAG1hp4eiKASaom3WFClwsWMU_fcOJVKiFru5JJdgZvjebG_kO2OJtq4U2TsuqXAFg8XHfKGkZS_2tpoKEC6ZngqVlZJMBlRmuRUYiF6452SrSVnCRmlYUSJ_KKef3JrBZ7XOsqylzj8pJwuiT1fpK3x0M4Jd_8mgS3FFnFUlsZzmBh_ewGweHTWyeznD2EIpXXBJLDNgdwxaMHjw7LrH7F05jSai_8_nIZLLLahNbe-xZi2ZD3O50cRY7ah1mg1Sde96CZ86-MW0iCmd5BQ7URKO3O6ggw5SSWX5JjztJfuSBC6nSylip84wPUzlqgIxvSTfa9dl-gYWiStVIO3fbHUylpX5w0szLSqI5vYUhrRoLChYJdSpDMTPCO88njBBImlCSVsharbluKZTXVhwwjI4PetcAp4JDpcIaGEQYbBTxnoFNHuJ4dBMwQ5MqApHBUZt7JFo1KAsPuB14E36AyhAUd4BTb_MknCRm9rGSh7eexsaOzZsSlsBLuOAWXEpBViFpiyA1Y-9OPqcifZIxfuqm397jLx9LxAh1DGoT7NGCtgA-nE0rTIS_av1b0btoJ89Y8JnxH-3ZhdT7RsTT8YPhQZFxKGiIDxLtmRidPyjk33WfLvFSfjSNLOa4wLqV2gZJ-p6pogmpBupgjz9F-3gBO6IrEWy6kA4rOTMhS9jTfh5ULpsB5dM-2s7OM87LuztmM9-KlZeT2nUeJXFT8C59ruTfe1WxiVN3mnbKnVa-0vv7i5p_vwAGxEO9oujGu703Dz6V7nqju2nP8ZfAT29p9H1vWJv3Snh3x2m87-9_wCG49b2)

## Roadmap 🚧
- [ ] Refactor code to a more general version, which can be used on more Kubernetes clusters
- [ ] Refactor the updating component and implement own version




