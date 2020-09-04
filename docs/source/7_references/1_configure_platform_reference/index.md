Configure Platform Reference
====================================

다음은 ansible로 쿠버네티스 클러스터를 구축하는 경우,  함께 설치되는 도커 및 nvidia 도커 버전 수정 등의 커스터마이징을 위한 가이드이다.  



## Customizing

### Configure nvidia-docker

설치할 nvidia 도커 관련 커스터마이징을 하고자 하는 경우, 다음의 경로에 있는 파일을 수정하여 클러스터를 구축하면 된다. 해당 파일에서는 nvidia 도커를 구동하기 위한 도커 버전 수정, OS별 저장소 수정 등의 작업이 가능하다. 

```bash
$ sudo vim kubespray/roles/nvidia-docker/defaults/main.yml
>>
---
docker_version: 18.09
# CentOS/RedHat nvidia_docker repo                                                                                                     
libnvidia_container_rh_repo_base_url: https://nvidia.github.io/libnvidia-container/centos7/$basearch                                   
libnvidia_container_rh_repo_gpgkey: https://nvidia.github.io/libnvidia-container/gpgkey                                                
nvidia_container_runtime_rh_repo_base_url: https://nvidia.github.io/nvidia-container-runtime/centos7/$basearch                         
nvidia_container_runtime_rh_repo_gpgkey: https://nvidia.github.io/nvidia-container-runtime/gpgkey                                      
nvidia_docker_rh_repo_base_url: https://nvidia.github.io/nvidia-docker/centos7/$basearch                                               
nvidia_docker_rh_repo_gpgkey: https://nvidia.github.io/nvidia-docker/gpgkey

# Ubuntu nvidia_docker repo                                                                                                            
libnvidia_container_ubuntu_repo_base_url: https://nvidia.github.io/libnvidia-container/ubuntu16.04/$(ARCH)                             
libnvidia_container_ubuntu_repo_gpgkey: https://nvidia.github.io/libnvidia-container/gpgkey                                            
nvidia_container_runtime_ubuntu_repo_base_url: https://nvidia.github.io/nvidia-container-runtime/ubuntu16.04/$(ARCH)                   
nvidia_container_runtime_ubuntu_repo_gpgkey: https://nvidia.github.io/nvidia-container-runtime/gpgkey                                  
nvidia_docker_ubuntu_repo_base_url: https://nvidia.github.io/nvidia-docker/ubuntu16.04/$(ARCH)                                         
nvidia_docker_ubuntu_repo_gpgkey: https://nvidia.github.io/nvidia-docker/gpgkey

# Debian nvidia_docker repo                                                                                                            
libnvidia_container_debian_repo_base_url: https://nvidia.github.io/libnvidia-container/debian9/$(ARCH)                                 
libnvidia_container_debian_repo_gpgkey: https://nvidia.github.io/libnvidia-container/gpgkey                                            
nvidia_container_runtime_debian_repo_base_url: https://nvidia.github.io/nvidia-container-runtime/debian9/$(ARCH)                       
nvidia_container_runtime_debian_repo_gpgkey: https://nvidia.github.io/nvidia-container-runtime/gpgkey                                  
nvidia_docker_debian_repo_base_url: https://nvidia.github.io/nvidia-docker/debian9/$(ARCH)                                             
nvidia_docker_debian_repo_gpgkey: https://nvidia.github.io/nvidia-docker/gpgkey
```

뿐만 아니라, 다음의 경로에 있는 파일을 수정하여 OS별 설치할 nvidia 도커 버전 및 버전 패키지 등을 설정할 수 있다.

```bash
$ sudo cd kubespray/roles/nvidia-docker/vars
# choose file that correspond with your OS
>>
---
nvidia_docker_kernel_min_version: '3.10'

nvidia_docker_versioned_pkg:
  'latest': nvidia-docker2
  '1.13': nvidia-docker2=2.0.3-1.docker1.13.1
  '17.03': nvidia-docker2=2.0.3-1.docker17.03.2.ce
  'stable': nvidia-docker2=2.0.3-1.docker17.03.2.ce

nvidia_docker_package_info:
  pkg_mgr: apt
  pkgs:
    - name: "{{ nvidia_docker_versioned_pkg[docker_version | string] }}"
      force: yes

nvidia_docker_repo_key_info:
  pkg_key: apt_key
  url: '{{ nvidia_docker_debian_repo_gpgkey }}'
  repo_keys:
    - "{{ lookup('file', 'nvidia-docker_gpg') }}"

nvidia_docker_repo_info:
  pkg_repo: apt_repository
  repos:
    - deb {{ libnvidia_container_debian_repo_base_url }}
    - deb {{ nvidia_container-runtime_debian_repo_base_url }}
    - deb {{ nvidia_docker_debian_repo_base_url }}
```



### Configure docker

마찬가지로 설치할 도커 관련 커스터마이징을 하고자 하는 경우, 다음의 경로에 있는 파일을 수정하여 클러스터를 구축하면 된다. 해당 파일에서는 도커 버전 수정, OS별 저장소 수정 등의 작업이 가능하다. 

```bash
$ sudo vim kubespray/roles/container-engine/docker/defaults/main.yml
>>
---
docker_version: '18.09'
#docker_selinux_version: '17.03'
...
# Used to override obsoletes=0
yum_conf: /etc/yum.conf
yum_repo_dir: /etc/yum.repos.d
docker_yum_conf: /etc/yum_docker.conf

# CentOS/RedHat docker-ce repo
docker_rh_repo_base_url: 'https://download.docker.com/linux/centos/7/$basearch/stable'
docker_rh_repo_gpgkey: 'https://download.docker.com/linux/centos/gpg'

# Ubuntu docker-ce repo
docker_ubuntu_repo_base_url: "https://download.docker.com/linux/ubuntu"
docker_ubuntu_repo_gpgkey: 'https://download.docker.com/linux/ubuntu/gpg'
...
```

도커 역시 다음의 경로에서 OS에 맞게 도커 버전 및 버전 패키지 등을 수정할 수 있다. 

```bash
$ sudo cd kubespray/roles/container-engine/docker/vars
# choose file that correspond with your OS
>>
docker_kernel_min_version: '3.10'
# https://download.docker.com/linux/debian/
# https://apt.dockerproject.org/repo/dists/debian-wheezy/main/filelist
docker_versioned_pkg:
  'latest': docker-ce
  '1.13': docker-engine=1.13.1-0~debian-{{ ansible_distribution_release|lower }}
  '17.03': docker-ce=17.03.2~ce-0~debian-{{ ansible_distribution_release|lower }}
  '17.06': docker-ce=17.06.2~ce-0~debian
  '17.09': docker-ce=17.09.0~ce-0~debian
  '17.12': docker-ce=17.12.1~ce-0~debian
  '18.03': docker-ce=18.03.1~ce-0~debian
  '18.06': docker-ce=18.06.2~ce~3-0~debian
  '18.09': docker-ce=5:18.09.7~3-0~debian-{{ ansible_distribution_release|lower }}
  'stable': docker-ce=5:18.09.7~3-0~debian-{{ ansible_distribution_release|lower }}
  'edge': docker-ce=5:18.09.7~3-0~debian-{{ ansible_distribution_release|lower }}
...
```



### Create kubernetes cluster

커스터마이징을 완료하였으면 다음의 명령어를 통해 쿠버네티스 클러스터를 구축하면 된다.

```bash
$ ansible-playbook -i inventory/mycluster/hosts.ini --become --become-user=root cluster.yml -e kube_version=v1.15.2
```



