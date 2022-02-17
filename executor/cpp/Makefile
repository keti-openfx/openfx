#
# Copyright 2015 gRPC authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#registry=keti.asuscomm.com:5000
###############################################
#### Edit to your private registry address ####
###############################################
registry=10.0.0.255:5000/openfx
version=1.0
docker-build:
	docker build --network=host -t $(registry)/fxwatcher:$(version)-cpp .
docker-push:
	docker push $(registry)/fxwatcher:$(version)-cpp
clean:
	rm -f *.o *.pb.cc *.pb.h fxwatcher
