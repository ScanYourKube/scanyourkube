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
For the purpose of understanding how the components of the software are communicating with each other, the following C4-container diagram should be used.
![C4-container diagram](https://kroki.io/plantuml/svg/eNp1VMtu2zAQvOsrtj65QJte8gFxjQB9IK1huwjai0BTa4kxRQrkKo7_vkuRsqTYvplLz-zOzFIPnoSjttbZB2WkbguE5X2-tIbwje6aaxdCGXTxKiNFmi_6KhRKlE7UsLcOfrY7XGrhFJ1AmAIOiPrOV1m2QuetmXtiUkOfYLaJv2b8cwGpDIIITaFMCaun-9nHLDt3mXspzMm27sAd8he7CxRc-sul0BR-2F3kks4avodjpWQF5FRZcm-gCqHTDXYfD4zuZqyFESXGf7RNIQhD583JE9b54xvNQ0sZRXGHscQCG21PddKxZYKhEvqMkMw5opS9rtxhqTy5wDxYuk7FTpEBjgWdERrOMOhhUFlPwbAw_XDda4y-evB2T0fhMAj-wjmxE6YQrgAr1QDz0yFroXTuuzPP8fg5nCGee7nKpMl8TQ1s0L3G5v--LZ6nZMEKZ5DQ56JRycdYgMXqe084rYLQ2h5jNNKhIGVjZkJz33hMNieQ3b2gpHdCwhaGjnEZ-1ZpN1lDcEOm8HsvG1sElmyNetjaCxHLMBS3FeMZHHreShnWaIK-Gvqq9RVjhuhUzdsIZCM6_yrkYaCYRrLhx-IBu2D8APmFqqx29tqbGeOn8ODK5mm76oe-gF6dfo38wPCV02mdC0s_0mH4g1B3Gd3kTMFs4yNNry8kGtznGJxtywqecVdZe7jNMnmfPVn0sXvl4ctkRI2-ERJ94pmgbkSj9UUwe2fr9z5HGRfL8SfKGa307_N2ZtkDu8_f0_9Cvvsu)

Further more information about the different components of the job are described in the following C4-Component diagramm.
![C4-component diagramm](https://kroki.io/plantuml/svg/eNq1VdtuEzEQfd-vGPLUSqW88AENoRJQCVBShAChleN1EhNfVr4kjRD_znjX3ntCX-hLNTM-Z65nc2cdMc5Lkb3gigpfMFi8zhdaOfbkbsupAOGKmamQLLViKqIyx53AQPLCW062hkgJG21gRYn6pr158GsGH_Q6y7KGOn-jvSqIOV11X-GjmxHs-ncGbYor4p2WxHGaW3x4A7N5clTI9uUMYwutTcEVccyC2zFowBDAs-sec3DlNJmI_jdfgCiutqA3lb3Hmo1iIczVxhDrjKfOGzZI1b7rJXxo4RfTIkZ6xSl2ohUcudtBCx2kUtrxTXzaS_axE7icrkuROvWWmWEqXxaI6SX5Urku09ewRFzqAmn_ZKuTdUzm90_1tKgghrtTHNKitqBgpdAnGYkfEd56AmEHiaQdStocomFbjms6VYVFJyyjM7DOFaBMcLhEQAODBIOdti5cQL2XFI7NWIfCUc6C1Rt3JAYvUBWvUBuoyHD-BWjKW5jtFykJF7mtbKzj_mWwobZTu1zFyqx0JayYOdTJv7-bf-2TdU6OlLx_bPPP7xNh3wtECH2MR2dYfQFYOfoxb3MQnbPX61-MukEje8ZEyIj_bu0updrXJvYQpkHjIaRZ4h0ElmzJxEjzI6U-Gr7dYuehNMOc4exQixMt60VVzxTRxMUmqniV4UN28AKXQ9ais-FIOKzkjLCXqSb8Kmgjm8F1pv2pGVzgHRf2_5jPfiFWDMfZjeIh6Eb559pvlfrcrYxLmpRn0yr1xgSBt3LrNH9-gPURDvaLoxru9Nw8-nKshsMqSf4I4vsZ0dN7GqnvGXtrpYR_d5gOf3L_AoUf0zc=)

## Roadmap 🚧
- [ ] Refactor code to a more general version, which can be used on more Kubernetes clusters
- [ ] Refactor the updating component and implement own version




