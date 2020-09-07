Deploying Functions
====================================

본 챕터에서는 OpenFx의 CLI(Command Line Interface) 명령어를 통해 사용자가 정의한 함수를 빌드, 배포, 호출되기까지의 과정을 다룬다.



## Get Started

`openfx-cli`를 설치했다면 `openfx-cli` 에서 제공하는 명령어를 사용하여 함수 생성 및 배포를 시작할 수 있다. go언어를 사용하여  `echo`함수를 만드는 예제를 통해 OpenFx에 함수를 만들어 보자.



### Init Function

- 함수 배포를 위한 새 함수 템플릿 생성

  ```
  $ openfx-cli function init echo -r go
  >> 
  Folder: echo created
  Fucntion handler created in folder: echo/src
  Rewrite the function handler code in echo/src folder
  Config file written: config.yaml
  ```




- 함수 템플릿 디렉토리 생성

  새로운 함수를 생성하면 현재 경로 디렉토리 안에 함수 이름과 같은 디렉토리가 생성된다. 디렉토리의 내부 구조는 다음과 같다.

  ```
echo
  └── config.yaml
  └── Dockerfile
  └── src
       └──handler.go
  ```




### Write Function

- Handler 코드 작성 

  다음은 기본적으로 생성되는 go언어의 함수 템플릿이다. 템플릿 코드는 입력된 데이터를 그대로 반환해주는 `echo`함수이다. 사용자는 다음의 예제코드를 응용하여 함수의 실행 코드를 작성할 수 있다.

  - handler.go

    ```go
    package main
    
    import sdk "github.com/keti-openfx/openfx/executor/go/pb"
    
    //import mesh "github.com/keti-openfx/openfx/executor/go/mesh"
    
    func Handler(req sdk.Request) string {
    	// mesh call
    	//
    	// functionName := "<FUNCTION NAME>"
    	// input := string(req.Input)
    	// result := mesh.MeshCall(functionName, []byte(input))
    	// return result
    
    	// single call
    	return string(req.Input)
    }
    ```





### Build Function

- 작성한 함수를 OpenFx에 배포하기 위한 도커 이미지 생성

  yaml 파일에 작성된 함수 컨테이너의 설정값을 기준으로 함수 이미지를 빌드한다.
  
  ```
  $ openfx-cli function build -v 
  >> 
  Building function (echo) image ...
  Image: ketiasuscomm.com:5000/echo built in local environment.
  ```
  



### Test Function

- 작성한 함수 배포 전, 로컬 환경에서 테스트 진행

  함수 배포에 앞서 작성한 함수가 정상적으로 동작하는지 확인하기 위해 로컬 환경에서 함수 이미지를 실행한다.

  ```
  $ echo "Hello" | openfx-cli function run echo
  >>
  Running image (ketiasuscomm.com:5000/echo) in local
  Starting FxWatcher Server ...
  Call echo in user's local
  Handler request: Hello
  
  Handler reply: Hello
  [1]+  Stopped                 echo "Hello" | openfx-cli function run echo
  ```

  > *Ctrl + Z를 통해 함수 실행을 중지할 수 있다.*



### Deploy Function

- 생성된 도커 이미지를 통해 OpenFx에 함수 배포

  yaml 파일을 통해 생성된 도커 이미지를 도커 레지스트리에 푸시하고 함수 컨테이너를 OpenFx에 배포한다.

  ```
    $ openfx-cli function deploy -f config.yaml -v 
    >> 
    Pushing: echo, Image: ketiasuscomm.com:5000/echo in Registry: ketiasuscomm.com:5000...
    ...
    Deploying: echo ...
    Attempting update... but Function Not Found. Deploying Function...
    http trigger url: http://ketiasuscomm.com:31113/function/echo
  ```




### Confirm Function

- OpenFx에 배포가 완료된 함수의 목록 확인

  ```
   $ openfx-cli function list
   >> 
   Function            Image                         Maintainer    Invocations    Replicas    Status    Description
   <FUNCTION NAME>     $(repo)/echo                                0              1           Ready
  ```




### Call Function

- OpenFx에 배포된 함수 호출

  ```
 $ echo "Hello" | openfx-cli function call echo
   >>
   Hello
  ```
  




### Function Info

- OpenFx에 배포된 특정 함수의 정보 확인

  ```
 $ openfx-cli function info echo
   >>
   name: echo
   image: ketiasuscomm.com:5000/echo
   invocationcount: 4
   replicas: 1
   annotations: {}
   availablereplicas: 1
   labels:
     openfx_fn: echo
  ```
  



### Function Log

- OpenFx에 배포된 특정 함수의 로그 확인

  ```
 $ openfx-cli function log echo
   >>
   ---
   Name: echo
   Log: 
   ...
  ```
  
