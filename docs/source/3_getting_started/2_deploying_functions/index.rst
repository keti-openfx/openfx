Deploying Functions
===================

본 챕터에서는 OpenFx의 CLI(Command line interface) 명령어를 통해
사용자가 정의한 함수를 빌드, 배포, 호출되기까지의 과정을 다룬다.

OpenFx CLI command
------------------

다음은 OpenFx에서 제공하는 CLI 명령어의 템플릿이다.

-  ``openfx-cli function <COMMAND NAME>``
-  ``openfx-cli fn <COMMAND NAME>``

모든 명령어는 ``--help; -h``\ 를 이용하여 상세 내용을 확인할 수 있다.

-  ``openfx-cli function <COMMAND NAME> --help``

Init Function
~~~~~~~~~~~~~

-  함수 배포를 위한 새 함수 템플릿 생성

``$ openfx-cli function init <FUNCTION NAME> -r <RUNTIME> [-f <YAML FILE NAME>.yaml] [-g <호스트 OS IP>:31113>]   >>    Folder: <FUNCTION NAME> created   Fucntion handler created in folder: <FUNCTION NAME>/src   Rewrite the function handler code in <FUNCTION NAME>/src folder   Config file written: <YAML FILE NAME>.yaml``

    ``<FUNCTION NAME>`` : 생성하고자 하는 함수의 이름; 첫문자는 반드시
    영문 소문자로 작성해야 한다.

    ``--runtime; -r`` : 함수의 런타임 종류; OpenFx에서 지원하는 런타임은
    모두 가능하다.

    -  OpenFx에서 지원하는 런타임의 종류
    -  Golang
    -  Python 2.7 / 3.6
    -  Ruby
    -  CPP
    -  C#
    -  Node.js
    -  Java

    ``--config; -f`` : 함수 설정 파일의 이름; 설정하지 않으면 기본값인
    ``config.yaml``\ 으로 생성된다.

    ``--gateway; -g`` : Gateway의 주소와 포트번호; 함수 호출 시 사용하는
    Gateway의 주소이다. 설정하지 않으면 기본값인
    ``keti.asuscomm.com:31113``\ 으로 설정되며 사용하고자 하는 IP주소의
    도메인 설정이 되어있다면 URL로도 사용이 가능하다.

-  함수 템플릿 디렉토리 생성

새로운 함수를 생성하면 현재 경로 디렉토리 안에 함수 이름과 같은
디렉토리가 생성된다. 디렉토리의 내부 구조는 다음과 같다.

.

├── .yaml
├── Dockerfile
└── src
​ └── handler.

    ``<YAML FILE NAME>.yaml`` : 함수에 대한 정보와 함수의 컨테이너
    설정값

    ``Dockerfile`` : 함수 컨테이너의 베이스 도커 이미지 파일

    ``src/handler.<RUNTIME>`` : 작성해야 하는 함수 파일

-  .yaml

   함수에 대한 정보와 함수의 컨테이너 설정값을 가지고 있다. 이를
   수정하여 함수의 정보와 설정값을 변경할 수 있다.

   ::

       functions:
         <FUNCTION NAME>:
           runtime: <RUNTIME>
           desc: ""
           maintainer: ""
           handler:
             dir: ./src
             file: handler.<RUNTIME>
           docker_registry: <REGISTRY IP>:<PORT>
           image: <REGISTRY IP>:<PORT>/<FUNCTION NAME>
           requests:
             memory: 50Mi
             cpu: 50m
             gpu: ""
         openfx:
         gateway: <호스트 OS IP>:31113

       ``functions`` : 함수 컨테이너 정보 기술

       ``<FUNCTION NAME>`` : 사용자가 설정한 함수 이름

       ``runtime`` : 함수의 런타임 종류; OpenFx에서 지원하는 런타임은
       모두 가능하다.

       ``<RUNTIME>`` : 사용자가 지정한 함수의 런타임

       ``desc`` : 함수에 대한 상세 설명; 함수 작성자가 함수에 대한 설명
       기술

       ``maintainer`` : 함수 작성자 혹은 유지보수 담당자에 대한 정보
       기술; 이메일 혹은 닉네임 등 자유롭게 서술 가능하다.

       ``handler`` : 함수의 엔트리포인트이며, 각각의 항목은 다음과 같다.

       -  ``dir`` : 함수 파일의 디렉토리 위치

       -  ``file`` : 사용자가 작성할 함수 파일

       ``docker_registry`` : 도커 레지스트리의 주소

       ``<REGISTRY IP>:<PORT>`` : 사용자 도커 레지스트리의 IP 주소와
       PORT 번호

       ``image`` : 도커 레지스트리에 전송될 도커 이미지의 이름

       ``<호스트 OS IP>`` : Gateway 설정에서 작성한 호스트 OS IP 주소
       (변경 가능)

       ``requests`` : 사용자가 정의할 함수 컨테이너별 자원 사용량이며,
       각각의 항목은 다음과 같다.

       -  ``memory`` : 함수 컨테이너의 memory 사용량; 최대 200Mi까지
          지정할 수 있으며, 기본값은 50Mi이다.
       -  ``cpu`` : 함수 컨테이너의 cpu 사용량; 최대 80m까지 지정할 수
          있으며, 기본값은 50m이다.
       -  ``gpu`` : 함수 컨테이너의 gpu 사용량; 값이 빈 문자열이면 CPU
          함수로 작동한다.

       ``openfx`` : OpenFx 정보 기술

       ``gateway`` : 함수 호출 시 사용하는 Gateway의 주소

-  Dockerfile

함수 컨테이너의 베이스 도커 이미지 파일이다. 해당 파일을 기반으로 함수
컨테이너를 빌드한다.

-  src/handler.

   런타임에 따라 실제 함수를 작성할 수 있는 기본 템플릿이 제공된다.

Write Function
~~~~~~~~~~~~~~

-  Handler 코드 작성

다음은 런타임별 함수 내에서 실행되는 예제 코드이다. 사용자는 다음의 예제
코드를 응용하여 함수의 실행 코드를 작성할 수 있다.

#### Golang

-  handler.go

   .. code:: go

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

#### Python 2.7 / 3.6

-  handler.py

   .. code:: python

       import mesh

       def Handler(req):
           # mesh call
           #
           # functionName = "<FUNCTIONNAME>"
           # input = req.input
           # result = mesh.mesh_call(functionName, input)
           # return result

           # single call
           return req.input

#### Ruby

-  handler.rb

   .. code:: ruby

       #!/usr/bin/env ruby

       module FxWatcher
         def FxWatcher.Handler(argStr)
             # mesh call
             # functionName = "<FUNCTIONNAME>"
             # input = argStr
             # result = FxWatcher.mesh_call(functionName, input)
             # return result 
             #
             # single call
             return argStr
         end
       end

#### CPP

-  handler.cc

   .. code:: cpp

       #include <iostream>

       using namespace std;
       //extern string MeshCall(string functionName, string input);

       string Handler(const string req) {

         // mesh call
         // string functionName = "<FUNCTIONNAME>";
         // string input = req;
         // string result = MeshCall(functionName, input);
         // return result 
         //
         // single call
         return req;
       }

#### C#

-  handler.cs

   .. code:: csharp

       namespace Fx
       {
           class Function
           {
               public byte[] Handler(byte[] Input)
               {
                   return Input; 
               }
           }
       }

#### Node.js

-  handler.js

   .. code:: js

       // handler.js

       function Handler(argStr) {
           return argStr;
       }

       module.exports = Handler;

#### Java

-  Handler.java

   .. code:: java

       package io.grpc.fxwatcher;

       import com.google.protobuf.ByteString;

       public class Handler {

         public static String reply(ByteString input) {
           return input.toStringUtf8() + "test";
         }

       }

Build Function
~~~~~~~~~~~~~~

-  작성한 함수를 OpenFx에 배포하기 위한 도커 이미지 생성

yaml 파일에 작성된 함수 컨테이너의 설정값을 기준으로 함수 이미지를
빌드한다.

``$ openfx-cli function build -v [-f <YAML FILE NAME>.yaml] [--nocache] [-g <호스트 OS IP>:31113>]   >>    Building function (<FUNCTION NAME>) image ...   Image: <REGISTRY IP>:<PORT>/<FUNCTION NAME> built in local environment.``

    ``--config; -f`` : 함수 설정 파일의 이름; 함수 생성 시 설정했다면
    해당 파일로 옵션을 추가한다.

    ``--buildverbose; -v`` : 이미지 빌드 과정을 로그로 출력

    ``--nocache`` : 이미지 빌드에 캐시 미사용

    ``--gateway; -g`` : Gateway의 주소; 함수 생성 시 설정한 Gateway의
    주소를 입력한다.

Test Function
~~~~~~~~~~~~~

-  작성한 함수를 배포 전, 로컬 환경에서 테스트 진행

함수 배포에 앞서 작성한 함수가 정상적으로 동작하는지 확인하기 위해 로컬
환경에서 함수 이미지를 실행한다.

\`\`\` $ echo "Hello" \| openfx-cli function run [-f .yaml] >> Running
image (:/) in local Starting FxWatcher Server ... Call in user's local
Handler request: Hello

Handler reply: Hello [1]+ Stopped echo "Hello" \| openfx-cli function
run \`\`\`

    ``--config; -f`` : 함수 설정 파일의 이름; 함수 생성 시 설정했다면
    해당 파일을 작성한다.

    *Ctrl + Z를 통해 함수 실행을 중지할 수 있다.*

Deploy Function
~~~~~~~~~~~~~~~

-  생성된 도커 이미지를 통해 OpenFx에 함수 배포

yaml 파일을 통해 생성된 도커 이미지를 도커 레지스트리에 푸시하고 함수
컨테이너를 OpenFx에 배포한다.

``$ openfx-cli function deploy -f <YAML FILE NAME>.yaml -v -g <호스트 OS IP>:31113> [--min <NUMBER>] [--max <NUMBER>] [--registry <REGISTRY IP>:<PORT>] [--replace=<TRUE OR FALSE] [--update=<TRUE OR FALSE>]     >>      Pushing: <FUNCTION NAME>, Image: <REGISTRY IP>:<PORT>/<FUNCTION NAME> in Registry: <REGISTRY IP>:<PORT>...     ...     Deploying: <FUNCTION NAME> ...     Attempting update... but Function Not Found. Deploying Function...     http trigger url: http://<호스트 OS IP>:31113/function/<FUNCTION NAME>``

    ``--config; -f`` : 함수 설정 파일의 이름; 함수 생성 시 설정했다면
    해당 파일을 작성하고 아닐 시에는 기본값인 ``config.yaml``\ 을
    입력한다.

    ``--buildverbose; -v`` : 이미지 빌드 과정을 로그로 출력

    ``--gateway; -g`` : Gateway의 주소; 함수 생성 시 설정한 Gateway의
    주소를 입력한다.

    ``--min`` : 함수 레플리카의 최솟값 (default = 1)

    ``--max`` : 함수 레플리카의 최댓값 (default = 1)

    ``--registry`` : 함수를 배포하고자 하는 도커 레지스트리 주소

    ``--replace`` : 존재하는 같은 이름의 함수를 제거하고 재생성

    ``--update`` : 존재하는 같은 이름의 함수에 롤링 업데이트를 수행
    (기본 설정=true)

Confirm Function
~~~~~~~~~~~~~~~~

-  OpenFx에 배포가 완료된 함수의 목록 확인

``$ openfx-cli function list [-g <호스트 OS IP>:31113>]    >>     Function            Image                         Maintainer    Invocations    Replicas    Status    Description    <FUNCTION NAME>     $(repo)/<FUNCTION NAME>                     0              1           Ready``

    ``--gateway; -g`` : Gateway의 주소; 함수 생성 시 설정한 Gateway의
    주소를 입력한다. 입력하지 않을 시에는 Gateway의 기본값인
    ``keti.asuscomm:31113``\ 에 배포된 함수의 목록이 나타난다.

Call Function
~~~~~~~~~~~~~

-  OpenFx에 배포된 함수 호출

``$ echo "Hello" | openfx-cli function call <FUNCTION NAME> [-g <호스트 OS IP>:31113>]    >>    Hello``

    ``--gateway; -g`` : Gateway의 주소; 함수 생성 시 설정한 Gateway의
    주소를 입력한다. 입력하지 않을 시에는 Gateway의 기본값인
    ``keti.asuscomm:31113``\ 에 배포된 함수가 호출된다.

Function Info
~~~~~~~~~~~~~

-  OpenFx에 배포된 특정 함수의 정보 확인

``$ openfx-cli function info <FUNCTION NAME> [-g <호스트 OS IP>:31113>]    >>    name: <FUNCTION NAME>    image: <REGISTRY IP>:<PORT>/<FUNCTION NAME>    invocationcount: 4    replicas: 1    annotations: {}    availablereplicas: 1    labels:      openfx_fn: <FUNCTION NAME>``

    ``--gateway; -g`` : Gateway의 주소; 함수 생성 시 설정한 Gateway의
    주소를 입력한다. 입력하지 않을 시에는 Gateway의 기본값인
    ``keti.asuscomm:31113``\ 에 배포된 함수의 정보가 나타난다.

Function Log
~~~~~~~~~~~~

-  OpenFx에 배포된 특정 함수의 로그 확인

``$ openfx-cli function log <FUNCTION NAME> [-g <호스트 OS IP>:31113>]    >>    ---    Name: <FUNCTION NAME>    Log:     ...``

    ``--gateway; -g`` : Gateway의 주소; 함수 생성 시 설정한 Gateway의
    주소를 입력한다. 입력하지 않을 시에는 Gateway의 기본값인
    ``keti.asuscomm:31113``\ 에 배포된 특정 함수의 로그가 나타난다.
