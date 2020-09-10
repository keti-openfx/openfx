# CLI Setup

OpenFx-cli는 OpenFx를 사용하기 위한 Command Line Interface 도구이다. OpenFx 서비스들은 OpenFx-cli를 활용하여 생성, 빌드, 배포, 테스트가 가능하다. Openfx-cli 설치 방법은 다음의 과정을 통해 진행한다. 




# Requirements

- [go version >= 1.12](<https://golang.org/dl/>)

- [docker version >= 18.06](<https://docs.docker.com/get-docker/>) 

  > Note
  > OS에 맞게 설치 요망

- [OpenFx-Gateway & Executor 구동]()



# Get Started

### Setting insecure registries

도커 레지스트리는 SSL 인증서 없이 이용할 수 없다. SSL 인증서 없이 도커 레지스트리를 사용하기 위해서는 `insecure-registries`에 대한 설정이 필요하다. `insecure-registries` 에 대한 설정은 아래와 같이 진행한다. 

```bash
$ sudo vim /etc/docker/daemon.json
>>
{"insecure-registries": ["YOUR PRIVATE REGISTRY SERVER IP:PORT"]}

$ systemctl daemon-reload
$ service docker restart
```



### Compile OpenFx-CLI

- `openfx-cli` 저장소를 __keti-openfx__ 디렉토리 밑에 클론하여 컴파일을 진행한다. 


  ```bash
  $ cd $GOPATH/src/github.com/keti-openfx
  $ git clone https://github.com/keti-openfx/openfx-cli.git
  $ cd openfx-cli
  ```

- `make` 명령을 실행하여 컴파일을 진행한다.   

  ```bash
  $ make build
  >>
  go build
  go install
  ```

- `$GOPATH/bin`을 확인해보면 `openfx-cli`가 컴파일 되어 있는 것을 확인할 수 있다. 

  ```bash
  $ cd $GOPATH/bin
  $ ls
  openfx-cli
  ```
