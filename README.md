# OpenFx

## Introduction of OpenFx&reg;

> 빠른 응답성, 유연한 확장성, 편의성 등을 통합 지원하는 오픈 소스 서버리스 

![Architecture of the OpenFx](/openfx_architecture.png)



## Intallation

1. [Installing minikube](./documents/1.Installing_Minikube.md)

2. [Building private docker registry](./documents/2.Building_Private_Docker_Registry.md)

3. [Compiling OpenFx](./documents/3.Compile_OpenFx.md)

4. [Deploy OpenFx](./documents/4.Deploy_OpenFx.md)



## Usage

```
git clone https://github.com/keti-openfx/OpenFx.git
cd OpenFx
kubectl apply -f ./namespaces.yml
kubectl apply -f ./yaml
```

## Status

OpenFX는 아직 초기 개발 중으로 향후 오픈소스화 할 예정임.

## Governance

본 프로젝트는 정보통신기술진흥센터(IITP)에서 지원하는 '18년 정보통신방송연구개발사업으로, "API 호출 단위 자원할당 및 사용량 계량이 가능한 서버리스 클라우드 컴퓨팅 기술개발" 임.
