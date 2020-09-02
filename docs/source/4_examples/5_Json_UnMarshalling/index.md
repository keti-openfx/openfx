# Json Unmarshalling

본 예제는 OpenFx 함수에서 `언마샬링(Unmarshalling)` 를  통해 입출력 인터페이스를 구성하는 예제이다. 언어는 Golang 이다. 



##### Unmarshalling ?

 언마샬링이란 로우 바이트를 논리적 구조로 변경하는 것을 뜻하며 Decoding 이라 표현한다. 



### Write function

##### Init Handlergo

```
$ openfx-cli fn init unmarshalling --runtime go
>>
Folder: unmarshalling created.
Function handler created in folder: unmarshalling/src
Rewrite the function handler code in unmarshalling/src folder
Config file written: config.yaml

```

##### handler.go

아래와 같이 `handler.go`를 작성한다.

```go
package main                                                                                                                                                                       
import (                                                                                     "encoding/json"                                                                           "fmt"                                                                                     sdk "github.com/keti-openfx/openfx/executor/go/pb"                                   )                                                                                                                                                
type SensorReading struct {                                                                  Name     string `json:"name"`                                                             Capacity int    `json:"capacity"`                                                         Time     string `json:"time"`                                                         }                                                                                                                                     
func Handler(req sdk.Request) string {                                                       var reading SensorReading
    err := json.Unmarshal(req.Input, &reading)                                               if err != nil {                                                                               fmt.Println(err)                                                                     }                                                                                         return fmt.Sprintf("%+v", reading)                                                   }                           
```



### Build function

작성한 함수를 빌드한다

```
$ openfx-cli fn build -f config.yaml 
>>
Building function (unmarshalling) image...
Image: keti.asuscomm.com:5000/unmarshalling built in local environment.
```



### Deploy function

```
$ openfx-cli fn deploy -f config.yaml 
>>
Pushing: unmarshalling, Image: keti.asuscomm.com:5000/unmarshalling in Registry: keti.asuscomm.com:5000 ...
Deploying: unmarshalling ...
Attempting update... but Function Not Found. Deploying Function...
http trigger url: http://keti.asuscomm.com:31113/function/unmarshalling 
```



### Test

```
$ echo '{"name": "battery sensor", "capacity": 40, "time": "2019-01-21T19:07:28Z"}' | openfx-cli fn call unmarshalling
>> 
{Name:battery sensor Capacity:40 Time:2019-01-21T19:07:28Z}
```



### 참고

[ Learn Go: Marshal & Unmarshal JSON in Golang #21 ](https://ednsquare.com/story/learn-go-marshal-unmarshal-json-in-golang------B6LUvY)







