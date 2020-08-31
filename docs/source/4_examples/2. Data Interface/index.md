## 2. Client Data Interface 

본 예제는 사용자 클라이언트에서 OpenFx의 함수를 호출하기 위한 Data Interface를 정리하였다.



#### Prerequirement

OpenFx는 gRPC 프로토콜로 설계된 서버리스 프레임워크이다. OpenFx는 gRPC 프로토콜 사용을 장려한다. 이는 기능적으로 HTTP 프로토콜도 지원하나 gRPC Gateway를 통해 변환이 필요하여 속도 지연이 생길 우려가 있기 때문이다.  gRPC은 통신 구조를 정의하기 위한 Protobuf이 필요하며, 정의한 데이터로만 통신이 가능하다. 현재 OpenFx의 정의된 Streaming Protobuf 의 통신 구조는 다음과 같다.

```protobuf
rpc Invoke(InvokeServiceRequest) returns(Message) {} 
message InvokeServiceRequest {                                                             string Service = 1;                                                                       bytes Input = 2;                                                                       }     

message Message {                                                                           string Msg = 1;                                                                         }     

```

입력은 Bytearray를 입력받고 출력은 String 타입으로 데이터로 반환된다. python과 같은 동적인 경우 자동으로 타입 변환이 되지만 Go, C, C++, Java의 같은 정적 언어인 경우 타입 변환에 신경을 써야한다.



또한, OpenFx 통신을 위한 gRPC Protobuf 정의가 필요하다. 다음의 명령을 통해 `Pb` 폴더의 `fxgateway.proto` 을 컴파일한다. 컴파일 언어는 `python` 이다.

```
python -m grpc_tools.protoc -I${GOPATH}/src/github.com/digitalcompanion-keti/pb \ 
            --python_out=. \
             --grpc_python_out=. \
            ${GOPATH}/src/github.com/digitalcompanion-keti/pb/gateway.proto
```

컴파일 후 실행 폴더 내 `fxgateway_pb2.py` 와 `fxgateway_pb2_gprc.py` 이 생성된다.



`Golang` 같은 경우 기본적으로 컴파일 파일이 제공된다.  뿐만 아니라 필요에 따라서는 `pb` 폴더의 Makefile 을 통해 컴파일가능하다.

```
$ make fxgateway
```

컴파일 후 실행 폴더 내 `fxgateway.pb.gw.go` , `fxgateway.swagger.json`,  `fxgateway.swagger.json` 이 생성된다. 





* Image

* ###### Json(Marshalling & UnMarshalling)



* ### Image 

  본 예제에서는 Python 사용자 클라이언트를 생성하여 동영상 데이터를 이미지로 짤라 Openfx에 전송하는 

* ### Json (Marshalling & Unmarshalling)

  개발 언어는 `Golang`으로 진행하였다.  

  
  
  
  
  ```
  $ openfx-cli function init json-echo --runtime python3
  >>>
  Folder: json-echo created.
  Function handler created in folder: json-echo/src
  Rewrite the function handler code in json-echo/src folder
  Config file written: config.yaml
  ```

  
  
  #### Write function
  
  ##### hander.py
  
  아래의 코드는 입력받은 json 데이터중 `key1`  의 키를 가진 값을 반환하는 코드이다. 
  
  
  
  ```python
  import json                                                                                                                                      
  def Handler(req):                                                                    
      info = json.loads(req.input)                                                       
      return info["key1"] 
  ```
  
  
  
  #### Build function
  
  ```
  $ openfx-cli fn build -f config.yaml -v
  >>
  Building function (json-echo) image...
  Sending build context to Docker daemon  5.632kB
  Successfully tagged keti.asuscomm.com:5000/json-echo:latest
  Image: keti.asuscomm.com:5000/json-echo built in local environment.
  ```
  
  
  
  #### Test Function
  
  ```
  $ vi in.json
  >> 
  {
  	"key1": "value1",
  	"key2": "value2"
  }
  
  $ root@ubuntu:~/go/src/github.com/keti-openfx# cat in.json | openfx-cli fn call echo
  >> 
  value1 
  ```
  
  

