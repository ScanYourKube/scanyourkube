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
![C4 container diagram](https://kroki.io/plantuml/svg/eNp1VNuOmzAQfecrpvuUSu32pR-wabpSL2obJalW3RdkzAS8GA-yzWbz9x1jIJDLUzKYc5k5Yx6cF9a3tU7eKSN1myOsPqcrMh7f_H1z7UAogzYeJV55zQfDU8iVKKyoYU8WtlKYf9Tan22GSbJG68gscnxFTQ3aD3D3dfh_x8USxiMoyXllClDegaO9PwiLQAZ8iRDYrEGPDtiV84x-nySjg4Vj1SOrVvxe-kIZU0-NwA_Kopy0ZPgcDqWSJXirioItdhrdTID2sWA0CJNDLYwoML7RNrnwGJS3R_ZQp49vfhEkpRZW-SMrBLVVrLi1RtOxRuOD9o4JTk-CzgTJnBNKOfSVWiyU8zYwn8a96R92HRngyHg2QsMIgwE2zjS4Px0PPXrO10ynzQ1_4gx5EiYXNgeS6gRzc5O1UDp1Xc0-Hj-GGmI9tKtM78zVvoEt2tco_vxt-TQnq8Z8U9Gofo594Mv194Fw_hSE1nSI0UiLwiuKmQnNurHsx9yDKHtB6c8aqRB1UOSfe1cOUlUsuYcwDdmHP8yyoTywJBvU0-W-aGMVbLGwmLqw6HgvZVikM_zV4NetK8Pej_GpmjcSPEV8-kXIakoyD2aLJneAXTzuBPqNqigzunZzpvg5PMxm-2u37o1fQq_63yBfM7YHsrU2rP6kE8OfjLpL6iZnH88uXtX-DoZcQwYchqW2KOEJs5Kous0yu6UDWZxkd9fDt8uIGl0jJLqeZ4a6EY7WF9HsLdXnc45tXCzI39jOZLH_jDuaJA88ff7i_gcKowj2)

Further more information about the different components of the job are described in the following C4 component diagram.
![C4 component diagram](https://kroki.io/plantuml/svg/eNq1VdtuEzEQfd-vGPLUSqW88AENoRJQCVBShAChleN1EhNfVr4kjRD_znjX3ntCX-hLNTM-Z65nc2cdMc5Lkb3gigpfMFi8zhdaOfbkbsupAOGKmamQLLViKqIyx53AQPLCW062hkgJG21gRYn6pr158GsGH_Q6y7KGOn-jvSqIOV11X-GjmxHs-ncGbYor4p2WxHGaW3x4A7N5clTI9uUMYwutTcEVccyC2zFowBDAs-sec3DlNJmI_jdfgCiutqA3lb3Hmo1iIczVxhDrjKfOGzZI1b7rJXxo4RfTIkZ6xSl2ohUcudtBCx2kUtrxTXzaS_axE7icrkuROvWWmWEqXxaI6SX5Urku09ewRFzqAmn_ZKuTdUzm90_1tKgghrtTHNKitqBgpdAnGYkfEd56AmEHiaQdStocomFbjms6VYVFJyyjM7DOFaBMcLhEQAODBIOdti5cQL2XFI7NWIfCUc6C1Rt3JAYvUBWvUBuoyHD-BWjKW5jtFykJF7mtbKzj_mWwobZTu1zFyqx0JayYOdTJv7-bf-2TdU6OlLx_bPPP7xNh3wtECH2MR2dYfQFYOfoxb3MQnbPX61-MukEje8ZEyIj_bu0updrXJvYQpkHjIaRZ4h0ElmzJxEjzI6U-Gr7dYuehNMOc4exQixMt60VVzxTRxMUmqniV4UN28AKXQ9ais-FIOKzkjLCXqSb8Kmgjm8F1pv2pGVzgHRf2_5jPfiFWDMfZjeIh6Eb559pvlfrcrYxLmpRn0yr1xgSBt3LrNH9-gPURDvaLoxru9Nw8-nKshsMqSf4I4vsZ0dN7GqnvGXtrpYR_d5gOf3L_AoUf0zc=)

## Roadmap 🚧
- [ ] Refactor code to a more general version, which can be used on more Kubernetes clusters
- [ ] Refactor the updating component and implement own version




