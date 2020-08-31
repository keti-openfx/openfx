OpenFx Setup
====================================

다음은 OpenFx 프레임워크를 구성하기 위해 쿠버네티스 클러스터(Kubernetes Cluster) 설치부터 OpenFx 설치까지의 가이드이다. 

 

## Configuration

쿠버네티스 클러스터(Kubernetes Cluster) 위에서 OpenFx 코어를 동작시키기 위해서는 기본적으로 다음과 같은 최소 사양이 만족되어야 한다.

- `프로세서 코어` >= 2
- `메모리` >= 8GB
- `하드디스크` >= 100GB
- `네트워크` 지원 가능



## Kubernetes Cluster 구축

쿠버네티스 클러스터는  다음의 두 가지 방식으로 구축할 수 있다. 

- 하나 혹은 다수의 물리 서버(가상 머신)를 묶어 `ansible playbook`으로 구축
- 미니쿠베(Minikube)를 통해 구축 

다음은 위의 두 가지 방식을 통해 쿠버네티스 클러스터를 구축하는 방법에 대한 가이드이다. 



### Ansible을 통한 쿠버네티스 클러스터 구축

`Ansible`은 파이썬 기반의 IaC(Infrastructure as Code)를 지향하는 자동화 관리 도구이다. yaml 포맷을 기반으로 플레이북을 실행시켜서 여러 머신에 동시에 소프트웨어 패키지를 설치함으로써 인프라 시스템 구축을 자동화 할 수 있다. 그 외에도 에드혹 모드로 모듈을 실행하여 여러 머신의 상태를 조회해 볼 수 있다. 본 가이드에서는 쿠버네티스를 인프라 시스템에 손쉽게 설치할 수 있도록 하는 ansible 기반의 `kubespray`를 활용하여 쿠버네티스 클러스터를 구축하는 방법을 안내한다. 아래는 kubespray를 통해  yaml 파일에 정의된 각각의 롤(role)들을 기반으로 쿠버네티스 클러스터를 구축하는 방법에 대해 설명한다. 



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
  $ cd kubespray
  $ sudo pip install --upgrage setuptools
  $ sudo pip3 install -r requirements.txt
  ```

- Ubuntu

  ```bash
  $ apt-get update
  $ apt-get install python-pip python3-pip
  $ cd kubespray
  $ sudo pip install --upgrage setuptools
  $ sudo pip3 install -r requirements.txt
  ```



#### Configuration

마스터 노드가 될 물리 서버에서 다른 물리 서버로 ssh 접속하기 위해서는 해당 서버의 ip를 알아야 한다. kubespray에서는 이를 `hosts.ini` 파일에 정의하고 있다. 이는 다음과 같이 할 수 있다. 

```bash
$ sudo vi inventory/mycluster/hosts.ini
>>
[all]
node1    ansible_host=<Master Node IP> ip=<Master Node IP>
node2    ansible_host=<Worker Node IPs> ip=<Worker Node IPs>
node3    ansible_host=<Worker Node IPs> ip=<Worker Node IPs>
node4    ansible_host=<Worker Node IPs> ip=<Worker Node IPs>

[kube-master]
node1
node2

[etcd]
node1

[kube-node]
node2
node3

[k8s-cluster:children]
kube-master
kube-node

[calico-rr]
```

> Note
>
> `etcd` 항목에 기재될 노드의 갯수는 반드시 홀수 개가 되어야 하며, `kube-master` 항목에 기재될 마스터 노드의 갯수는 최소 1개 이상이다. 또한 마스터 노드는 동시에 워커 노드가 될 수 있다. 



도커 이미지 개인 저장소(Docker private registry)를 구축한 경우 다음의 경로에 있는 설정 파일을 수정하면 된다. 

```bash
$ sudo vi inventory/mycluster/group_vars/all/docker.yml
>>
# Add here
docker_insecure_registries:
    - <private registry addr:port>
```



#### Setup cluster

다음의 명령어를 실행하여 쿠버네티스 클러스터를 구축한다.

```bash
$ ansible-playbook -i inventory/mycluster/hosts.ini --become --become-user=root cluster.yml -e kube_version=v1.15.2
```



#### Verify 

클러스터 구축이 완료되었으면, 다음의 명령어를 통해 정상적으로 구축이 되었는지를 확인한다. 

```bash
$ kubectl get pods --all-namespaces
kube-system   calico-kube-controllers-b784f96cc-bjzsh   1/1     Running     1          41d
kube-system   calico-node-mhkp8                         1/1     Running     2          41d
kube-system   calico-node-p49md                         1/1     Running     3          41d
kube-system   coredns-74c9d4d795-p74lb                  1/1     Running     1          41d
kube-system   coredns-74c9d4d795-qms9s                  1/1     Running     1          41d
kube-system   dns-autoscaler-576b576b74-95rmk           1/1     Running     1          41d
kube-system   kube-apiserver-node1                      1/1     Running     1          41d
kube-system   kube-controller-manager-node1             1/1     Running     1          41d
kube-system   kube-proxy-hq6zl                          1/1     Running     1          41d
kube-system   kube-proxy-qhkrc                          1/1     Running     1          41d
kube-system   kube-scheduler-node1                      1/1     Running     1          41d
kube-system   kubernetes-dashboard-7c547b4c64-d7ghs     1/1     Running     1          41d
kube-system   local-volume-provisioner-dgftl            1/1     Running     1          41d
kube-system   local-volume-provisioner-dhljn            1/1     Running     1          41d
kube-system   metrics-server-c779857d6-4ss6t            1/1     Running     0          35d
kube-system   nginx-proxy-node2                         1/1     Running     1          41d
```



### Minikube

`미니쿠베`는 쿠버네티스처럼 클러스터를 구성하지 않고 단일 컴퓨팅 환경(노트북, 데스크탑 등)에서 쿠버네티스 환경을 만들어준다. 로컬 환경에서 단일 클러스터를 구동시킬 수 있는 도구인 미니쿠베는 단일 노드에 쿠버네티스 클러스터 환경을 구축하기 때문에 접근성이 뛰어나 클러스터를 관리하기가 수월하다. 이로 인해 더욱 용이해진 디버깅 환경을 사용자에게 제공하여 편의성을 높여준다. 다음은 미니쿠베 설치 방법 및 미니쿠베 환경 위에서 쿠버네티스를 사용할 수 있게 해주는 명령 줄 인터페이스인 `kubectl` 설치 방법이다.

 

#### Install Virtual Machine

미니쿠베를 시작하기 전, 미니쿠베를 통해 쿠버네티스 컴포넌트를 가상 머신(Virtual Machine) 위에서 동작시키기 위해 [버츄얼박스](<https://www.virtualbox.org/>)를 설치한다. 

> Note
>
> 가상 머신이 아닌 호스트 OS 환경(리눅스)이라면, 가상 머신 설치를 생략한다. 



#### Install Minikube

클러스터를 구축할 가상 머신 혹은 호스트 OS 환경이 준비되었다면, 본격적으로 쿠버네티스 클러스터를 로컬 환경에서 구축하기 위한 미니쿠베 설치를 진행하여야 한다. 

- Linux

  ```bash
  $ wget minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube-linux-amd64
  $ sudo mv minikube-linux-amd64 minikube
  $ sudo mv minikube /usr/local/bin
  ```





#### Start Minikube

미니쿠베를 시작하기 전, 미니쿠베는 기본적으로 하이퍼바이저(hypervisor)를 지원하고 있으며 다음의 [링크](<https://kubernetes.io/ko/docs/setup/learning-environment/minikube/#vm-%EB%93%9C%EB%9D%BC%EC%9D%B4%EB%B2%84-%EC%A7%80%EC%A0%95%ED%95%98%EA%B8%B0>)를 통해 지원하고 있는 하이퍼바이저를 확인 후, 이를 설치하여 사용할 수 있다. 하이퍼바이저란 virtualbox, vmware 같이 물리적 호스트에서 다수의 가상머신을 실행할 수 있도록 하여 컴퓨팅 자원을 효과적으로 사용할 수 있게 하는 도구이다. 미니쿠베를 실행하기 위해선 가상 머신이 필요하며, 가상 머신에서 실행하지 않으려면 리눅스 시스템과 도커가 필요하다. 하이퍼바이저 설치까지 완료되었으면, `--driver=<driver_name>` 플래그를 추가해서 미니쿠베를 시작할 수 있다.  뿐만 아니라 쿠버네티스 버전을 명시하여 미니쿠베를 실행할 수 있는데, 현재 OpenFx 코어는 쿠버네티스 버전 `1.15.2`까지 지원하기 때문에 다음과 같이 버전을 지정하여 미니쿠베를 시작해야 한다.  

```bash
$ echo export CHANGE_MINIKUBE_NONE_USER=true >> ~/.bashrc
$ sudo minikube start --driver=<driver_name> --kubernetes-version v1.15.2 --insecure-registry="<IP ADDRESS>:<PORT>"
```



하이퍼바이저 설치를 하지 않았다면 `--driver=none` 플래그를 통해 미니쿠베를 시작할 수도 있다.

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
>>
Starts a local kubernetes cluster

Options:
...
      --apiserver-port=8443: The apiserver listening port
      --cpus=2: Number of CPUs allocated to the minikube VM
      --docker-opt=[]: Specify arbitrary flags to pass to the Docker daemon. (format: key=value)
      --insecure-registry=[]: Insecure Docker registries to pass to the Docker daemon.  The default service CIDR range
will automatically be added.
      --memory='2000mb': Amount of RAM allocated to the minikube VM (format: <number>[<unit>], where unit = b, k, m or
g)
      --vm-driver='virtualbox': VM driver is one of: [virtualbox parallels vmwarefusion kvm2 vmware none]
      --wait=true: Wait until Kubernetes core services are healthy before exiting

Usage:
  minikube start [flags] [options]

Use "minikube start options" for a list of global command-line options (applies to all commands).
```



#### Further progress

미니쿠베를 시작한 후, 다음의 명령어를 통해 `~/.kube`, `~/.minikube` 디렉토리의 권한을 `$USER`로 변경해야 한다. 이는 현재 사용자가 쿠버네티스 및 미니쿠베 관련 설정 파일들을 수정할 수 있게 하기 위함이다.

```bash
$ sudo chown -R $USER ~/.kube ~/.minikube
```



쿠버네티스에서 자동 스케일링을 하기 위해서는 각 노드의 자원 사용량 정보를 알아야 한다. 이를 위해 쿠버네티스에서는 각 노드 별 메트릭 데이터를 수집하는 `heapster`와  `metrics-server`를 제공하고 있다. 쿠버네티스에서 기본적으로 설정된 도구는 `heapster`인데 `metrics-server`가 조금 더 나은 성능을 보이고 있다. 이 두 개의 도구를 동시에 실행하면 충돌이 일어나기 때문에 미니쿠베에서는 수동으로 `heapster`를 종료하고 `merics-server`를 실행해주어야 한다. 

```bash
$ sudo minikube addons disable heapster
$ sudo minikube addons enable metrics-server
```



#### Install kubectl

`kubectl `은 쿠버네티스를 제어하기 위한 명령 줄 인터페이스이다. 미니쿠베를 통해 구축된 로컬 환경에서의 쿠버네티스 클러스터를 사용하기 위해선 설치해야만 하는 필수 요소이다. 이는 아래와 같은 명령어로 설치를 진행할 수 있다. 

- Linux

  ```bash
  $ curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
  $ curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.14.0/bin/linux/amd64/kubectl
  $ chmod +x ./kubectl
  $ sudo mv ./kubectl /usr/local/bin/kubectl
  $ kubectl version
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

CoreDns에는 Loop 플러그인이라는 서브 모듈이 존재한다. Loop 플러그인이란 임의의 probe query를 자신에게 보내고 이를 몇번 반환받게 되는지를 추적한다. 위와 같은 에러는 loop 감지 플러그인이 업스트림 dns 서버 중 하나에서 무한 전달 루프를 감지해서 발생하는 에러이다. 이는 다음과 같이 해결할 수 있다. 

- Solution #1

  **CoreDns configmap**을 수정한다. 이는 아래와 같은 명령어를 실행 후, `loop`이라는 단어를 삭제한다.

  ```bash
  $ kubectl -n kube-system edit configmap coredns
  ```

  `loop`이라는 단어를 삭제한 후,  새로운 설정이 적용된 Pod를 생성하기 위해 기존의 Pod를 삭제한다.

  ```bash
  $ kubectl -n kube-system delete pod -l k8s-app=kube-dns
  ```

- Solution #2 

  **Solution #1**의 방법으로 에러가 해결이 안되면 이는 방화벽 규칙의 문제일 수 있다. 쿠버네티스 클러스터 구동 시, 기본적으로 추가되는 방화벽 규칙들이 있다. 하지만 쿠버네티스 클러스터 구동 중, 방화벽 규칙이 제대로 추가되지 않거나 기존의 규칙들과 충돌이 일어날 수 있다. 이와 같은 경우, 기존의 규칙들을 모두 제거하고 쿠버네티스 및 도커 관련 방화벽 규칙들을 재정의 해주어야 하며, 이는 아래와 같은 명령어로 실행할 수 있다. 

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
  > **Solution #1** 의 방법으로 에러 해결 시, **Solution #2** 는 진행하지 않아도 된다. 



## OpenFx 설치

쿠버네티스 클러스터가 구동 중이고, 도커 이미지를 담을 개인 저장소를 구축하여 로그인까지 완료하였으면 본격적으로 OpenFx를 배포하여야 한다. 이를 위해 먼저 소스들을 컴파일하여 도커 이미지로 빌드하여 저장소에 저장하여야 한다. 



### Requirements

- [go version >= 1.12](<https://golang.org/dl/>)
- [docker version >= 18.06](<https://docs.docker.com/get-docker/>)



### Package dependencies

OpenFx API 게이트웨이를 사용하기 위한 의존 패키지들은 다음과 같다.



- **google.golang.org/grpc**
- **github.com/protocolbuffers/protobuf**
- **github.com/golang/protobuf/protoc-gen-go**
- **github.com/grpc-ecosystem/grpc-gateway**
- **grpcio-tools**

OpenFx API 게이트웨이는 본 가이드에 명시된 버전의 의존 패키지들로 최적화 되어있다.  아래는 각각의 의존 패키지들에 대한 설치 방법이다.

#### google.golang.org/grpc

```bash
$ go get -u google.golang.org/grpc
$ go get -u golang.org/x/sys/unix
```

#### github.com/protocolbuffers/protobuf

```bash
$ curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
$ unzip protoc-3.7.1-linux-x86_64.zip -d protoc3

$ sudo mv protoc3/bin/* /usr/local/bin/
$ sudo mv protoc3/include/* /usr/local/include/

$ sudo chown $USER /usr/local/bin/protoc
$ sudo chown -R $USER /usr/local/include/google

$ export PATH=$PATH:/usr/local/bin
```

#### github.com/golang/protobuf/protoc-gen-go

```bash
$ go get -u github.com/golang/protobuf/protoc-gen-go
$ cd $GOPATH/src/github.com/golang/protobuf/protoc-gen-go
$ git checkout tags/v1.2.0 -b v1.2.0
$ go install
```

#### github.com/grpc-ecosystem/grpc-gateway

```bash
$ cd $GOPATH/src/github.com
$ git clone https://github.com/grpc-ecosystem/grpc-gateway.git
$ cd grpc-ecosystem/grpc-gateway
$ git checkout v1.4.1

// Install protoc-gen-grpc-gateway
$ cd protoc-gen-grpc-gateway
$ go get -u github.com/golang/glog
$ go install

// Install protoc-gen-swagger
$ cd ../protoc-gen-swagger
$ go install
```

#### grpcio-tools

```bash
$ pip install grpcio-tools
$ pip3 install grpcio-tools
```



#### Troubleshooting

Mac 사용자들은 다음과 같은 오류가 발생할 수 있다.

```bash
Traceback (most recent call last):
  File "/usr/bin/pip", line 11, in <module>
    sys.exit(main())
  File "/usr/lib/python2.7/dist-packages/pip/__init__.py", line 215, in main
    locale.setlocale(locale.LC_ALL, '')
  File "/usr/lib/python2.7/locale.py", line 581, in setlocale
    return _setlocale(category, locale)
locale.Error: unsupported locale setting
```

이는 지역 및 언어 선택이 정상적을 설정이 되어 있지 않아서 발생하는 오류이다.  이는 다음과 같이 해결할 수 있다.

```bash
$ export LC_ALL="en_US.UTF-8"
$ export LC_CTYPE="en_US.UTF-8"
$ sudo dpkg-reconfigure locales
```

locales를 재설정하게 되면 `Configuring locales`라는 화면이 표시되고, `en_US.UTF-8 UTF-8`이 체크되어 있는지 확인한 후 OK를 눌러 설정을 마치면 된다.



### Compile OpenFx

먼저 아래와 같이 `keti-openfx`라는 폴더를 생성하여 OpenFx 소스코드를 복제할 위치를 지정한다.

```bash
$ mkdir $GOPATH/src/github.com/keti-openfx
$ cd $GOPATH/src/github.com/keti-openfx
```

다음은 OpenFx 프레임워크 위 서비스들을 관리하는 gateway의 이미지를 생성하는 방법에 대한 가이드이다. 먼저, `openfx` 저장소를 클론하여 openfx 디렉토리로 이동한다.

```bash
$ git clone https://github.com/keti-openfx/openfx.git
$ cd openfx
```

openfx 디렉토리 내의 `Makefile`에 `REGISTRY`란을 도커 레지스트리 서버에 맞춰 변경한다.

```bash
$ sudo vim Makefile
REGISTRY=<REGISTRY IP ADDRESS> : <PORT>
...
```

`make` 명령어를 이용해서 `openfx-gateway`를 컴파일하고, 이미지를 생성한 뒤, 개인 도커 레지스트리에 저장한다.

```bash
$ make build
$ make push
```

다음은 OpenFx 프레임워크 위 서비스들을 실행하기 위한 gRPC 서버인 executor의 이미지를 생성하는 방법에 대한 가이드이다. 앞서 클론한 openfx 디렉토리 내의 executor 디렉토리로 이동한다.

```bash
$ cd executor
```

OpenFx executor는 다음과 같이 총 7개의 runtime 버전이 존재한다.

- go
- python
- nodejs
- ruby
- java
- cpp
- csharp

각각의 runtime 폴더에 있는 `Makefile` 의 `registry` 를 도커 레지스트리 서버에 맞춰 변경한다.

```bash
$ cd go
$ sudo vim Makefile
registry=<REGISTRY IP ADDRESS>:<PORT>
...

$ cd ../python
$ sudo vim Makefile
registry=<REGISTRY IP ADDRESS>:<PORT>
...

$ cd ../nodejs
$ sudo vim Makefile
registry=<REGISTRY IP ADDRESS>:<PORT>
...

$ cd ../ruby
$ sudo vim Makefile
registry=<REGISTRY IP ADDRESS>:<PORT>
...

$ cd ../java
$ sudo vim Makefile
registry=<REGISTRY IP ADDRESS>:<PORT>
...

$ cd ../cpp
$ sudo vim Makefile
registry=<REGISTRY IP ADDRESS>:<PORT>
...

$ cd ../csharp
$ sudo vim Makefile
registry=<REGISTRY IP ADDRESS>:<PORT>
...
```

`executor` 폴더로 돌아와서 `make` 명령을 실행하여 runtime별 executor를 컴파일 한 후, 각각의 이미지를 생성하여 개인 도커 레지스트리에 저장한다.

```bash
$ cd ..
$ make
```

컴파일 완료 후, `docker images`와 레지스트리에 있는 이미지를 확인했을 때, 아래와 같이 결과나 나오면 성공적으로 컴파일이 완료된 것이다.

```bash
$ docker images
>>
REPOSITORY                       TAG                 IMAGE ID            CREATED       SIZE
<REGISTRY IP>:<PORT>/fxwatcher   0.1.0-csharp        5ab2175321bd        37 minutes ago
690MB
<REGISTRY IP>:<PORT>/fxwatcher   0.1.0-cpp           50bc61dc1545        35 minutes ago
6.66GB
<REGISTRY IP>:<PORT>/fxwatcher   0.1.0-java          0bec44a16eec        30 minutes ago
548MB
<REGISTRY IP>:<PORT>/fxwatcher   0.1.0-ruby          56fa4e607ac6        26 minutes ago 
490MB
<REGISTRY IP>:<PORT>/fxwatcher   0.1.0-nodejs        b26348908044        28 minutes ago
335MB
<REGISTRY IP>:<PORT>/fxwatcher   0.1.0-python3       5779598d8ad0        25 minutes ago  413MB
<REGISTRY IP>:<PORT>/fxwatcher   0.1.0-python2       b91ef13cede0        32 minutes ago  401MB
<REGISTRY IP>:<PORT>/fxwatcher   0.1.0-go            3cd97230054d        39 minutes ago  793MB
<REGISTRY IP>:<PORT>/fxgateway   0.1.0               89bc10ce43ec        3 hours ago    255MB
<none>                           <none>              3d1f57588f3f        3 hours ago    986MB
python                           2.7-alpine          ee70cb11da0d        13 days ago    61.3MB
python                           3.4-alpine          c06adcf62f6e        2 months ago    72.9MB
registry                         2                   f32a97de94e1        2 months ago    25.8MB
alpine                           3.7                 6d1ef012b567        2 months ago    4.21MB
golang                           1.9.7               ef89ef5c42a9        10 months ago  750MB
golang                           1.10.1              1af690c44028        12 months ago  780MB

$ curl -k -X GET https://<ID>:<PASSWD>@<REGISTRY IP ADDRESS>:<PORT>/v2/_catalog
>>
{"repositories":["fxgateway","fxwatcher"]}
```



### Deploy OpenFx

OpenFx 컴파일을 완료하였다면, 이제 OpenFx를 배포해보자.

먼저, 쿠버네티스에서 개인 도커 레지스트리로부터 도커 이미지를 다운받으려면 **도커 인증(Docker credential)**이 필요하다. 이를 위해서 도커 레지스트리 타입의 **Secret**을 사용하여 레지스트리에 인증을 받는다. 그리고 도커 인증을 생성하고 배포하기 위해 **yaml** 파일에 이를 설정한다. 이는 다음의 절차를 통해 진행된다.

- Secret 생성
- 도커 인증파일을 base64로 변환한 **.dockerconfigjson** 내용 확인
- **keti-openfx/openfx/deploy/yaml**폴더의 **docker-registry-secret.yaml**에 **.dockerconfigjson**내용을 이전에 확인한 **.dockerconfigjson**으로 변경

#### Create Secret

다음의 명령어를 통해 secret을 생성

```bash
$ kubectl create secret docker-registry regcred --docker-server=<REGSTRY IP>:<PORT> --docker-username=<your-name> --docker-password=<your-password>
```

- `<REGSTRY IP>:<PORT>` : 개인 레지스트리 주소와 포트
- `<your-name>` : 도커 로그인을 위한 아이디
- `<your-password>` : 도커 로그인을 위한 비밀번호

#### Inspecting the Secret `regcred`

다음의 명령어를 통해 도커 인증파일을 bash64로 변환한 **.dockerconfigjson** 내용을 확인

```bash
$ kubectl get secret regcred --output=yaml
>>
kind: Secret
metadata:
  ...
  name: regcred
  ...
data:
  .dockerconfigjson: eyJodHRwczovL2luZGV4L ... J0QUl6RTIifX0=
type: kubernetes.io/dockerconfigjson
```

- `.dockerconfigjson`에 나오는 정보는 위와 상이할 수 있다.

#### Configure docker-registry-secret.yml

**keti-openfx/openfx/deploy/yaml** 폴더의 **docker-registry-secret.yaml** 파일의 **.dockerconfigjson**의 내용을 위에서 확인한 **Secret** `regcred`의 **.dockerconfigjson** 내용으로 변경

```bash
apiVersion: v1
kind: Secret
metadata:
  name: regcred
  namespace: openfx
data:
  .dockerconfigjson: eyJodHRwczovL2luZGV4L ... J0QUl6RTIifX0=
type: kubernetes.io/dockerconfigjson
---
apiVersion: v1
kind: Secret
metadata:
  name: regcred
  namespace: openfx-fn
data:
  .dockerconfigjson: eyJodHRwczovL2luZGV4L ... J0QUl6RTIifX0=
type: kubernetes.io/dockerconfigjson
```

- `.dockerconfigjson`에 나오는 정보는 위와 상이할 수 있다.

위와 같이 도커 인증까지 완료되었으면 이제 본격적으로 쿠버네티스 클러스터에 OpenFx 게이트웨이를 배포하여야 한다. 이는 아래와 같이 진행할 수 있다. 

#### Configure gateway-dep.yml

**openfx/deploy/yaml** 폴더의 **gateway-dep.yml**파일의 **image**란의 레지스트리 IP와 Port를 변경한다.

```bash
$ sudo vim $GOPATH/scr/github.com/keti-openfx/openfx/deploy/yaml/gateway-dep.yml
>>
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: fxgateway
  namespace: openfx
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: fxgateway
    spec:
      serviceAccountName: fxgateway
      imagePullSecrets:
        - name: regcred
      containers:
      - name: fxgateway
        image: <REGISTRY IP>:<REGISTRY PORT>/fxgateway:0.1.0
        imagePullPolicy: Always

        env:
        - name: FUNCTION_NAMESPACE
          value: openfx-fn
        - name: IMAGE_PULL_POLICY
          value: "Always"

        ports:
        - containerPort: 10000
          protocol: TCP

        resources:
          requests:
            memory: 250Mi
          limits:
            memory: 250Mi
```

게이트웨이 이미지를 `pull`할 레지스트리 주소 변경까지 완료하였으면, 다음의 명령어를 통해 OpenFx 컴포넌트들을 배포한다.

```bash
$ cd $GOPATH/scr/github.com/keti-openfx/openfx/deploy
$ kubectl apply -f ./namespaces.yml
$ kubectl apply -f ./yaml
$ kubectl get pods --all-namespaces
>>
NAMESPACE     NAME                               READY   STATUS             RESTARTS   AGE
kube-system   coredns-fb8b8dccf-4bq7x            1/1     Running            0          113s
kube-system   coredns-fb8b8dccf-jw6j2            1/1     Running            0          113s
kube-system   etcd-minikube                      1/1     Running            0          4m19s
kube-system   kube-addon-manager-minikube        1/1     Running            0          4m22s
kube-system   kube-apiserver-minikube            1/1     Running            0          4m17s
kube-system   kube-controller-manager-minikube   1/1     Running            0          4m6s
kube-system   kube-proxy-h8q7p                   1/1     Running            0          5m11s
kube-system   kube-scheduler-minikube            1/1     Running            0          4m16s
kube-system   kubernetes-dashboard-7b8ddcb5d6..  1/1     Running            0
4m18s
kube-system   metrics-server-89cd44dc7-d8jvj     1/1     Running            0
4m18s
kube-system   storage-provisioner                1/1     Running            0          5m7s
openfx        fxgateway-755df6464f-6zrqw         1/1     Running            0          6m28s
openfx        prometheus-5c8f7f7c7d-zhpbb        1/1     Running            0          6m30s
openfx        grafana-core-7d6b476bb9-dj9bl      1/1     Running            0          6m30s
openfx        grafana-import-dashboards-tnskf    1/1     Running            0          6m30s
...
```

- `STATUS`가 **Running**이 아닌 경우에는 [링크](https://kubernetes.io/ko/docs/reference/kubectl/cheatsheet/)를 참조하여 포드의 로그를 확인한다.
- `kubectl apply`를 통해 배포되는 pod 중 prometheus는 컴퓨팅 자원들의 metric 정보를 수집하는 모니터링 툴이며, grafana는 prometheus를 통해 수집된 metric 정보들을 시각화해주는 툴이다.

