# OpenFx Setup

다음은 OpenFx 프레임워크를 구성하기 위해 쿠버네티스 클러스터(Kubernetes Cluster) 설치부터 OpenFx 설치까지의 가이드이다. 

 

## Configuration

쿠버네티스 클러스터(Kubernetes Cluster) 위에서 OpenFx 코어를 동작시키기 위해서는 기본적으로 다음과 같은 최소 사양이 만족되어야 한다.

- `프로세서 코어` >= 2
- `메모리` >= 8GB
- `하드디스크` >= 100GB
- `네트워크` 지원 가능



## Kubernetes Cluster 구축

쿠버네티스 클러스터는 다음의 두 가지 방식으로 구축할 수 있다. 

- Kubernetes 설치
- Minikube 설치
-> (Kubernetes와 Minikube의 차이 설명)
...
- Kubernetes 설치
    - Ansible을 이용한 Kubernetes 클러스터 구축

다음은 위의 두 가지 방식을 통해 쿠버네티스 클러스터를 구축하는 방법에 대한 가이드이다. 



### Ansible을 통한 쿠버네티스 클러스터 구축

`Ansible`은 파이썬 기반의 IaC(Infrastructure as Code)를 지향하는 자동화 관리 도구이다. `Ansible`은 yaml 파일로 구성된 플레이북을 실행하여 여러 머신에 동시에 소프트웨어 패키지를 설치함으로써 인프라 시스템 구축을 자동화할 수 있다. 그 외에도 `Ansible`은 에드혹 모드로 모듈을 실행하여 여러 머신의 상태를 조회해 볼 수 있다. 본 가이드에서는 쿠버네티스를 인프라 시스템에 손쉽게 설치할 수 있도록 하는 `Ansible`기반의 `kubespray`를 활용하여 OpenFx 프레임워크를 설치하는 방법을 안내한다. 아래는 kubespray를 통해  yaml 파일에 정의된 각각의 롤(role)들을 기반으로 쿠버네티스 클러스터를 구축하는 방법에 대해 설명한다. 



#### Copy ssh key

Ansible을 통해 다수의 물리 서버(가상 머신)를 묶어 쿠버네티스 클러스터를 구축하고자 하는 경우, 하나의 물리 서버에서만 ansible 플레이북을 실행하여 구축할 수 있다. 이는 마스터 노드가 될 물리 서버에서 워커 노드가 될 물리 서버들로의 ssh 접속을 통해 이루어진다. 이를 위해 클러스터를 구성할 물리 서버들에 ssh key를 복사한다.

```bash
$ ssh-keygen -t rsa
$ ssh-copy-id -i ~/.ssh/id_rsa.pub root@<Master Node IP>
$ ssh-copy-id -i ~/.ssh/id_rsa.pub root@<Worker Node IPs>
```



#### Clone kubespray



#### Install required packages

- CentOS

  ```bash
  $ yum --enablerepo=extras install epel-release
  $ yum install python3 python-pip
  $ sudo pip3 install -r requirements.txt
  ```

- Ubuntu

  ```bash
  $ apt-get update
  $ apt-get install python-pip python3-pip
  $ sudo pip3 install -r requirements.txt
  ```




### Minikube

`미니쿠베`는 쿠버네티스처럼 클러스터를 구성하지 않고 단일 컴퓨팅 환경(노트북, 데스크탑 등)에서 쿠버네티스 환경을 만들어준다. 로컬 환경에서 단일 클러스터를 구동시킬 수 있는 도구인 미니쿠베는 단일 노드에 쿠버네티스 클러스터 환경을 구축하기 때문에 접근성이 뛰어나고 클러스터를 관리하기가 수월하며, 이로 인해 더욱 용이해진 디버깅 환경을 사용자에게 제공하여 편의성을 높여준다. 다음은 미니쿠베 설치 방법 및 미니쿠베 환경 위에서 쿠버네티스를 사용할 수 있게 해주는 명령 줄 인터페이스인 `kubectl` 설치 방법이다.

 

#### Install Virtual Machine

미니쿠베를 시작하기 전, 미니쿠베를 통해 쿠버네티스 컴포넌트를 가상 머신(Virtual Machine) 위에서 동작시키기 위해 [버츄얼박스](<https://www.virtualbox.org/>)를 설치한다. 

> Note
>
> 가상 머신이 아닌 호스트 OS 환경(Mac, 리눅스 등)이라면, 가상 머신 설치를 생략한다. 



#### Install Minikube

클러스터를 구축할 가상 머신 혹은 호스트 OS 환경이 준비되었다면, 본격적으로 쿠버네티스 클러스터를 로컬 환경에서 구축하기 위한 미니쿠베 설치를 진행하여야 한다. 

- Mac OS

  ```bash
  $ brew cask install minikube
  ```

- Linux

  ```bash
  $ wget minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube-linux-amd64
  $ sudo mv minikube-linux-amd64 minikube
  $ sudo mv minikube /usr/local/bin
  ```

- Window

  1. Windows PowerShell을 관리자 모드로 실행

  2. 다음의 명령어를 통해 Chocolatey 설치 전 환경 설정 확인

     ```bash
     PS C:\WINDOWS\system32> Get-ExecutionPolicy
     Unrestricted
     ```

     > Note
     >
     > 위 명령어의 실행 결과 값이 `Restricted` 인 경우, `Get-ExecutionPolicy`  대신에 `Set-ExecutionPolicy AllSigned ` 혹은 `Set-ExecutionPolicy Bypass -Scope Process` 명령어를 입력한다. 

  3. 다음의 명령어를 통해 Chocolatey 설치

     ```bash
     Set-ExecutionPolicy Bypass -Scope Process -Force;[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
     ```

  4. Chocolatey가 제대로 설치되었는지 확인

     ```bash
     PS C:\WINDOWS\system32> choco -v
     0.10.15
     ```

  5. 미니쿠베 설치

     ```bash
     PS C:\WINDOWS\system32> choco install minikube
     ```

#### Start Minikube

미니쿠베를 시작하기 전, 미니쿠베는 기본적으로 VM 드라이버를 지원하고 있으며 다음의 [링크](<https://kubernetes.io/ko/docs/setup/learning-environment/minikube/#vm-%EB%93%9C%EB%9D%BC%EC%9D%B4%EB%B2%84-%EC%A7%80%EC%A0%95%ED%95%98%EA%B8%B0>)를 통해 지원하고 있는 VM 드라이버를 확인 후, 이를 설치하여 사용할 수 있다.  VM 드라이버 설치까지 완료되었으면, `--driver=<driver_name>` 플래그를 추가해서 미니쿠베를 시작할 수 있다.  뿐만 아니라 쿠버네티스 버전을 명시하여 미니쿠베를 실행할 수 있는데, 현재 OpenFx 코어는 쿠버네티스 버전 `1.15.2`까지 지원하기 때문에 다음과 같이 버전을 지정하여 미니쿠베를 시작해야 한다.  

```bash
$ echo export CHANGE_MINIKUBE_NONE_USER=true >> ~/.bashrc
$ sudo minikube start --driver=<driver_name> --kubernetes-version v1.15.2 --insecure-registry="<IP ADDRESS>:<PORT>"
```



VM 드라이버 설치를 하지 않았다면 `--driver=none` 플래그를 통해 미니쿠베를 시작할 수도 있다.

```bash
$ echo export CHANGE_MINIKUBE_NONE_USER=true >> ~/.bashrc
$ sudo minikube start --driver=none --kubernetes-version v1.15.2 --insecure-registry="<IP ADDRESS>:<PORT>"
```

> Note
>
> `<IP ADDRESS>:<PORT>` 는 도커 레지스트리 서버의 주소와 포트번호를 적어주어야 한다. 도커 이미지 개인 저장소(Docker private registry)를 구축하는 방법은 다음의 [링크]()를 통해 진행하면 된다. 



추가로 `minikube start`의 다양한 플래그가 궁금하다면 다음의 명령어를 입력하여 원하는 정보를 얻을 수 있다.

```bash
$ sudo minikube start --help
```



#### Further progress

미니쿠베를 시작한 후, 다음의 명령어를 통해 `~/.kube`, `~/.minikube` 디렉토리의 권한을 `$USER`로 변경해야 한다. 

```bash
$ sudo chown -R $USER ~/.kube ~/.minikube
```



그 후, 미니쿠베에서 Horizontal Pod Autoscaling이 가능하게 하기 위해 다음과 같이 설정을 변경해주어야 한다.

```bash
$ sudo minikube addons disable heapster
$ sudo minikube addons enable metrics-server
```



#### Install kubectl

`kubectl `은 쿠버네티스를 제어하기 위한 명령 줄 인터페이스이다. 미니쿠베를 통해 구축된 로컬 환경에서의 쿠버네티스 클러스터를 사용하기 위해선 설치해야만 하는 필수 요소이다. 이는 아래와 같은 명령어로 설치를 진행할 수 있다. 

- MacOS

  ```bash
  $ brew install kubernetes-cli
  $ kubectl version
  ```

- Linux

  ```bash
  $ curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
  $ curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.14.0/bin/linux/amd64/kubectl
  $ chmod +x ./kubectl
  $ sudo mv ./kubectl /usr/local/bin/kubectl
  $ kubectl version
  ```

- Window

  1. Windows PowerShell을 관리자 모드로 실행

  2. 다음의 명령어를 통해 `kubectl` 설치

     ```bash
     PS C:\WINDOWS\system32> choco install kubernetes-cli
     PS C:\WINDOWS\system32> kubectl version
     ```



#### Verify installed minikube

가상 머신, 미니쿠베, 그리고 kubectl 설치까지 모두 완료하였으면, 쿠버네티스 클러스터가 정상적으로 동작하는지를 확인하여야 한다. 이는 아래와 같은 명령어를 통해 확인할 수 있다. 

```bash
$ kubectl get pods --all-namespaces
```



#### Troble Shooting (CrashLoopBackOff Error)

호스트 OS 혹은 가상 머신에서 미니쿠베를 실행할 경우, 아래와 같은 에러가 발생할 수 있다. 

```bash
$ kubectl get pods --all-namespaces
>>
NAMESPACE     NAME                               READY   STATUS             RESTARTS   AGE
kube-system   coredns-fb8b8dccf-mtn7d            0/1     CrashLoopBackOff   5          3m54s
kube-system   coredns-fb8b8dccf-t584j            0/1     CrashLoopBackOff   5          3m54s
kube-system   etcd-minikube                      1/1     Running            0          2m46s
kube-system   kube-addon-manager-minikube        1/1     Running            0          4m1s
kube-system   kube-apiserver-minikube            1/1     Running            0          2m51s
kube-system   kube-controller-manager-minikube   1/1     Running            0          2m52s
kube-system   kube-proxy-rtswf                   1/1     Running            0          3m54s
kube-system   kube-scheduler-minikube            1/1     Running            0          2m51s
kube-system   storage-provisioner                1/1     Running            0          3m53s
```

이는 다음과 같이 해결할 수 있다. 

- Solution #1

  __CoreDns configmap__을 수정한다. 이는 아래와 같은 명령어를 실행 후, `loop`이라는 단어를 삭제한다.

  ```bash
  $ kubectl -n kube-system edit configmap coredns
  ```

  `loop`이라는 단어를 삭제한 후,  새로운 설정이 적용된 Pod를 생성하기 위해 기존의 Pod를 삭제한다.

  ```bash
  $ kubectl -n kube-system delete pod -l k8s-app=kube-dns
  ```

- Solution #2 

  __Solution #1__의 방법으로 에러가 해결이 안되면 이는 방화벽 규칙의 문제일 수 있다. 쿠버네티스 클러스터 구동 시, 기본적으로 추가되는 방화벽 규칙들이 있다. 하지만 쿠버네티스 클러스터 구동 중, 방화벽 규칙이 제대로 추가되지 않거나 기존의 규칙들과 충돌이 일어날 수 있다. 이와 같은 경우, 기존의 규칙들을 모두 제거하고 쿠버네티스 및 도커 관련 방화벽 규칙들을 재정의 해주어야 하며, 이는 아래와 같은 명령어로 실행할 수 있다. 

  ```bash
  $ iptables -t nat -F
  $ iptables -t mangle -F
  $ iptables -F
  $ iptables -X
  $ iptables -P INPUT ACCEPT
  $ iptables -P FORWARD ACCEPT
  $ iptables -P OUTPUT ACCEPT
  $ iptables -N DOCKER
  $ iptables -N DOCKER-ISOLATION
  ```

  > Note 
  >
  > __Solution #1__ 의 방법으로 에러 해결 시, __Solution #2__ 는 진행하지 않아도 된다. 
