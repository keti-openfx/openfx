## 1. Function Test By Runtime Languages

다음은 CLI를 통해 함수 초기화시 제공되는 예제로 입출력 기능을 가진다.  런타임 언어별 초기화를 통해 제공되는 언어가 다르며 코드는 아래의 코드와 같다.

ex) 입출력 기능

```
$ echo Hello DCF | openfx-cli function call echo 
>> Hello DCF
```



- Golang example

  handler.go

  ```
  package main
  
  import sdk "github.com/keti-openfx/openfx/executor/go/pb"
  
  func Handler(req sdk.Request) string {
      return string(req.Input)
  }
  ```

- Python 2.7 / 3.4 example

  handler.py

  ```
  def Handler(req):
      return req.input
  ```

- Node Js example

  handler.js

  ```
  function Handler(argStr) {
      return argStr;
  }
  
  module.exports = Handler;
  ```

- Ruby example

  handler.rb

  ```
  #!/usr/bin/env ruby
  
  module FxWatcher
    def FxWatcher.Handler(argStr)
      return argStr
    end
  end
  ```

- C++ example

  handler.cc

  ```
  #include <iostream>
  
  using namespace std;
  
  string Handler(const string req) {
    return req;
  }
  ```

- Java example

  Handler.java

  ```
  package io.grpc.fxwatcher;
  
  import com.google.protobuf.ByteString;
  
  public class Handler {
  
    public static String reply(ByteString input) {
      return input.toStringUtf8() + "test";
    }
  
  }
  ```

- C# example

  handler.cs

  ```
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
  ```









