Invoke Internal Service
====================================

OpenFx의 함수는 독립적으로 배포가 가능하며 개별적으로 실행될 수 있는 하나의 서비스로 정의된다. OpenFx에 배포된 서비스는 내부에서 다른 서비스를 호출할 수 있으며, 이를 통해 사용자는 서로 다른 서비스를 엮어 원하는 기능의 애플리케이션 서비스를 완성할 수 있다.

본 챕터에서는 OpenFx에 배포되어 있는 서비스를 호출하는 방법에 대해서 알아본다. 





OpenFx에서 서비스 간 호출을 하기 위해서는 다음과 같은 일련의 과정을 거쳐야 한다.

- 핸들러 파일 안에 사용자가 호출을 원하는 서비스의 이름과 전달할 인자값 설정
- 작성된 함수를 배포하기 위해 도커 이미지로 빌드
- 생성된 도커 이미지를 OpenFx에 배포
- 사용자가 엮은 서비스의 시작점에 있는 서비스를 호출



## Supported Runtimes

서비스 간 호출을 지원하는 런타임 목록은 다음과 같다.

- Golang
- Python2.7 / 3.6
- Ruby
- CPP



서비스 간 호출을 지원하는 런타임들에 대해서 서비스를 호출하기 위한 Handler 코드 작성방법에 대해 설명한다. 새 함수 템플릿 생성 시, Handler 예제 코드의 기본 설정은 사용자가 단일 서비스 호출에 대한 응답을 바로 받을 수 있는 코드로 설정되어 있다. 

다른 서비스 호출에 대한 코드는 주석으로 표시되어 있다. 서비스 호출 기능이 필요할 경우, 전체 주석을 제거하고 single call에 대한 return값만 주석으로 처리해주면 된다.





### Golang

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

  >`<FUNCTION NAME>` : 호출하고자 하는 다른 함수의 이름; OpenFx에 배포되어 있는 함수이여야 한다.
  >
  >`input` : 호출하고자 하는 함수에 넘겨줄 인자값
  >
  >`MeshCall` : 사용자가 설정한 함수를 호출하는 기능 함수; 위에서 설정한 내용이 전달된다. 





### Python2.7 / 3.6

- handler.py

  ```python
  import mesh
  
  def Handler(req):
      # mesh call
      #
      # functionName = "<FUNCTION NAME>"
      # input = req.input
      # result = mesh.mesh_call(functionName, input)
      # return result
  
      # single call
      return req.input
  ```

  >`<FUNCTION NAME>` : 호출하고자 하는 다른 함수의 이름; OpenFx에 배포되어 있는 함수이여야 한다.
  >
  >`input` : 호출하고자 하는 함수에 넘겨줄 인자값
  >
  >`mesh_call` : 사용자가 설정한 함수를 호출하는 기능 함수; 위에서 설정한 내용이 전달된다. 





### Ruby

- handler.rb

  ```ruby
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
  ```

  >`<FUNCTION NAME>` : 호출하고자 하는 다른 함수의 이름; OpenFx에 배포되어 있는 함수이여야 한다.
  >
  >`input` : 호출하고자 하는 함수에 넘겨줄 인자값
  >
  >`mesh_call` : 사용자가 설정한 함수를 호출하는 기능 함수; 위에서 설정한 내용이 전달된다. 





### CPP

- handler.cc

  ```cpp
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
  ```

  >`<FUNCTION NAME>` : 호출하고자 하는 다른 함수의 이름; OpenFx에 배포되어 있는 함수이여야 한다.
  >
  >`input` : 호출하고자 하는 함수에 넘겨줄 인자값
  >
  >`MeshCall` : 사용자가 설정한 함수를 호출하는 기능 함수; 위에서 설정한 내용이 전달된다. 





Handler 코드 작성 후에, 다음 [링크](https://keti-openfx.readthedocs.io/en/latest/3_getting_started/2_deploying_functions/index.html)의 **Build Function**부터 순차적으로 진행하면 다른 서비스를 호출하는 함수가 OpenFx에 배포된다. 호출 시에는 사용자가 엮은 서비스의 시작점에 있는 서비스를 호출해야 정상적으로 호출된다.

