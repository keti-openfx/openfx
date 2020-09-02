# Crawling Bot

본 예제는 OpenFx 함수를 통해 구현한 크롤링 봇 예제이다.  본 예제에서의 함수는 네이버 홈페이지의 뉴스 헤드 이슈를 크롤링한다. 언어는 파이썬이며  감정 분석을 위해 [BeautifulSoup](https://www.crummy.com/software/BeautifulSoup/bs4/doc/) 사용하였다.



### Write function

##### Init Handler.py 

```
$ openfx-cli fn init crawler --runtime python3
>>
Folder: crawler created.
Function handler created in folder: crawler/src
Rewrite the function handler code in crawler/src folder
Config file written: config.yaml
```

##### handler.py

아래와 같이 `handler.py`를 작성한다.

```python
import requests                                                                           from bs4 import BeautifulSoup                                                                                                              
def Handler(req):                                                                             source = requests.get("http://www.naver.com").text                                       soup = BeautifulSoup(source, "html.parser")                                               hotkeys = soup.select("a.issue")                                                                                                             
    hot = []                                                                                 
    index = 0                                                                                 for key in hotkeys:                                                                           index += 1                                                                               hot.append(str(index) + "," + key.text)                                                   if index >= 20:                                                                               break                                                                         
    return '\n'.join(hot)                         
```

##### requirements.txt

함수 구성에 필요한 라이브러리를 requirements.txt에 명시한다.

```
bs4                               
requests 
```



### Build function

작성한 함수를 빌드한다

```
$ openfx-cli fn build -f config.yaml 
>>
Building function (crawler) image...
Image: keti.asuscomm.com:5000/crawler built in local environment.
```



### Deploy function

```
$ openfx-cli fn deploy -f config.yaml 
>>
Pushing: crawler, Image: keti.asuscomm.com:5000/crawler in Registry: keti.asuscomm.com:5000 ...
Deploying: crawler ...
Function crawler already exists, attempting rolling-update.
http trigger url: http://keti.asuscomm.com:31113/function/crawler 
```



### Test

```
$ echo "" | openfx-cli fn call crawler
>>
1,태풍 '마이삭' 시속 23㎞로 한반도 접근 중…자정께 부산 근접
2,정은경 "코로나 폭발적 급증은 억제…이번주가 안정·확산 기로"
3,2주간 코로나19 사망자 20명, 모두 60대 이상…'사망후 확진'도
4,서울 실내운동시설 3곳서 잇단 집단감염…사랑제일교회 1천117명
5,카카오게임즈 1억원 넣어도 수익은 19만원…경쟁률 1500대1 기준
6,노영민 "문대통령 사저부지에 건물 들어서면 기존 집 처분"
7,[1보] 미래통합당, '국민의힘'으로 당명 교체 확정
8,정부 "국회-의료계 합의 결과 존중"…의정갈등 풀리나
9,16일만에 퇴원한 전광훈 '사기극' 운운하며 문대통령 비난
10,野 "보좌관 전화 왔었다" 녹취공개…추미애·보좌관 고발
```



##### 참고

[ [python] 파이썬 크롤링(네이버 실시간 검색어) ](https://blockdmask.tistory.com/385)







