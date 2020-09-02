# 감정 분석

본 예제는 OpenFx 함수를 통해 구현한 감정 분석 예제이다.  언어는 파이썬이며  감정 분석을 위해 [TextBlob project](http://textblob.readthedocs.io/en/dev/) 라이브러리를 사용하였다.  



### Write function

##### Init Handler.py 

```
$ openfx-cli function init sentiment-analysis --runtime python3
>>
Directory: sentiment-analysis is created.
Function handler created in directory: sentiment-analysis/src
Rewrite the function handler code in sentiment-analysis/src directory
Config file written: config.yaml
```



##### handler.py

아래와 같이 `handler.py`를 작성한다.

```python
import json                                    
from textblob import TextBlob                                                                                                      
def Handler(req):   
	input_decode = req.input.decode('utf-8')                                                 blob = TextBlob(input_decode)
	
    output = "Sentiment(polarity={}, {}) \n".format(blob.polarity, blob.subjectivity)     
    return output                          
```



##### requirements.txt

함수 구성에 필요한 라이브러리를 requirements.txt에 명시한다.

```
textblob
```



### Build function

작성한 함수를 빌드한다

```
$ cd sentiment-analysis
$ openfx-cli  function build -f config.yaml -v
>>
Building function (sentiment-analysis) ...
Sending build context to Docker daemon  8.192kB
Step 1/45 : ARG ADDITIONAL_PACKAGE
Step 2/45 : ARG REGISTRY
Step 3/45 : ARG PYTHON_VERSION
...
```

### Deploy function

```
$ openfx-cli fn deploy -f config.yaml 
>>
Pushing: sentiment-analysis, Image: keti.asuscomm.com:5000/sentiment-analysis in Registry: keti.asuscomm.com:5000 ...
Deploying: sentiment-analysis ...
Attempting update... but Function Not Found. Deploying Function...
http trigger url: http://keti.asuscomm.com:31113/function/sentiment-analysis 
```



### Test

```
$ echo "Have a nice day" | openfx-cli function call sentiment-analysis
>>
Sentiment(polarity=0.6, 1.0)
```



### 참고

[textblob 공식 사이트](https://textblob.readthedocs.io/en/dev/quickstart.html)







