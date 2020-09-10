# Image Processing 

본 예제는 OpenFx 함수를 통해 이미지 프로세싱을 처리한다.



### Prerequirement

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





### Write function

##### Init Handler.py 

imgprocessing 함수를 생성한다.

```
$ openfx-cli function init imgprocessing --runtime python3
>>
Directory: imgprocessing is created.
Function handler created in directory: imgprocessing/src
Rewrite the function handler code in imgprocessing/src directory
Config file written: config.yaml
```



##### handler.py

아래와 같이 `handler.py`를 작성한다.

```
import numpy as np 
import cv2 

def Handler(req):
    # Bytes -> frame 
    nparr = np.frombuffer(req.input, np.uint8)
    frame = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

    
    """
    frame 데이터 처리 
    """

    # Frame -> Bytes
    res = cv2.imencode('.jpg', frame)[1].tostring()

    return res
```

##### requirements.txt

다음은 데이터 변환에 필요한 패키지 파일을 requirements.txt에 명시한다.

```
opencv-python
opencv-contrib-python
ffmpeg
```



### Build function

작성한 함수를 빌드한다

```
$ cd imgprocessing
$ openfx-cli  function build -f config.yaml -v
>>
Building function (imgprocessing) ...
Sending build context to Docker daemon  8.192kB
Step 1/45 : ARG ADDITIONAL_PACKAGE
Step 2/45 : ARG REGISTRY
Step 3/45 : ARG PYTHON_VERSION
...
```

### Deploy functions

```
$ openfx-cli fn deploy -f config.yaml 
>>
Pushing: crawler, Image: keti.asuscomm.com:5000/imgprocessing in Registry: keti.asuscomm.com:5000 ...
Deploying: imgprocessing ...
Function imgprocessing already exists, attempting rolling-update.
http trigger url: http://keti.asuscomm.com:31113/function/imgprocessing
```



## User Client

#### Init

`User Client`는 Python 언어로 구현하였으며 필요 라이브러리는 다음의 명령어를 통해 설치할 수 있다. 비디오 데이터 변환 및 입력을 위한 라이브러리로 Opencv를 사용하였다.

```
pip install opencv-python
pip install opencv-contrib-python
pip install ffmpeg 

python -m pip install grpcio
python -m pip install grpcio-tool

pip install argparse
```

*"Opencv 외 라이브러리 통해 데이터 인코딩 및 입력이 가능하지만, Handler 함수에서 사용자 라이브러리 설치 및 데이터 디코딩이 필요하다."*



다음은 클라이언트 코드의 작성 예제이다.

```python
import queue
import time
import datetime 
import threading

import argparse 
import numpy as np 
import cv2 

import grpc
import fxgateway_pb2
import fxgateway_pb2_grpc


address = 'keti'
port = 31113

class Client:
    def __init__(self):
        channel = grpc.insecure_channel(address + ':' + str(port))
        self.conn = fxgateway_pb2_grpc.GatewayStub(channel)
        self.dataQueue = queue.Queue()
        self.cap = cv2.VideoCapture(args.video)  

        self.cap.set(3, 960) 
        self.cap.set(4, 640) 

        threading.Thread(target=self.__listen_for_messages).start()
        self.Capture()

    def generator(self):
        while True:
            time.sleep(0.01)
            if self.dataQueue.qsize()>0:
                yield self.dataQueue.get()

    def __listen_for_messages(self):
        time.sleep(5)
        responses = self.conn.Invokes(self.generator())

        try :
            for i in responses:
                nparr = np.frombuffer(i.Output, np.uint8)
                newFrame = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
                cv2.imshow("OpenFx Image processing", newFrame)
                k = cv2.waitKey(1) & 0xff 
                if k == 27: # ESC 키 입력시 종료 
                    break 
                    
            self.cap.release()  
            cv2.destroyAllWindows()     
        except grpc._channel._Rendezvous as err :
            print(err)   
            

    def Capture(self): 
        """
        이 함수는 gRPC 를 위한 정보 입력과 발신 메세지를 처리합니다. 
        """
        time.sleep(1)
        while True:
            ret, frame = self.cap.read() # cap read 
            if cv2.waitKey(1) & 0xFF == ord('q'): 
                break
            res = cv2.imencode('.jpg', frame)[1].tostring()
            msg = gateway_pb2.InvokeServiceRequest(Service= args.Handler, Input=res)
            self.dataQueue.put(msg)

        print("Image Processing END!")

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='This code is written for OpenFx Client about Image Processing')
    parser.add_argument('Handler', type=str,
            metavar='Openfx Function name',
			help='Input to Use OpenFx Function')
    parser.add_argument('--image', type=str, default = int(0),
            metavar='image file Name',
            help='Input to Use image File Name \n')
    args = parser.parse_args()
    c = Client()
```





### Test

Client 를 실행하기 위한 명령어는 다음과 같다.

```
$ python client.py -h
> 

This code is written for OpenFx Client about Image Processing

positional arguments:
  OpenFx Function name  Input to Use OpenFx Function
  Image file Name    Input to Use Image File Name 

optional arguments:
  -h, --help         show this help message and exit
  
$ python3 client.py [$function] --image [$image File]
```

- [$function] : 사용할 OpenFx 함수를 등록한다.
- [$image File] : 사용할 동영상 파일명을 등록한다. 동영상 경로는 현 실행 폴더로 지정해뒀다. 또한 웹 캠으로 동영상 데이터를 입력받을 시 `0`을 입력한다.



```
$ python3 client.py imgprocessing test.jpg
```



