# OpenFx-WorkShop

이 문서는 OpenFx로 서버리스 함수를 작성하는 방법을 학습하기 위한 워크숍이다.



## Requirements

해당 워크숍에서는 쿠버네티스 클러스터에 OpenFx를 배포하고 OpenFx가 제공하는 CLI를 설치하는 것으로 시작한다. Openfx 배포 및 CLI 설치를 위한 문서는 [다음의 링크]()에서 제공한다.



## Workshop composition

| Name               | Details                                                      |
| ------------------ | ------------------------------------------------------------ |
| 감정 분석          | python, `textblob` 를 활용한 텍스트 기반의 감정 분석 함수    |
| 이미지 프로세싱    | python, `opencv`, `ffmpeg` 를 활용한 이미지 프로세싱 함수와 <br />이를 위한 사용자 클라이언트 구성 |
| MQTT Broker 구성   | OpenFx 함수와 IOT 기기 통신을 위한 MQTT Broker 구성          |
| 크롤링 봇          | python 을 활용한 크롤링 함수                                 |
| Json Unmarshalling | Golang  Json 형식의 데이터 처리 함수                         |