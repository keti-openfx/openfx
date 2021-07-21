module github.com/keti-openfx/openfx

go 1.16

require (
        github.com/beorn7/perks v1.0.1 // indirect
        github.com/evanphx/json-patch v4.9.0+incompatible // indirect
        github.com/fsnotify/fsnotify v1.4.9 // indirect
        github.com/golang/groupcache v0.0.0-20191227052852-215e87163ea7 // indirect
        github.com/golang/protobuf v1.4.3
        github.com/googleapis/gnostic v0.4.0 // indirect
        github.com/grpc-ecosystem/grpc-gateway v1.16.0
        github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
        github.com/philips/go-bindata-assetfs v0.0.0-20150624150248-3dcc96556217
        github.com/prometheus/client_golang v1.8.0
        github.com/prometheus/client_model v0.0.0-00010101000000-000000000000 // indirect
        github.com/prometheus/common v0.0.0-00010101000000-000000000000 // indirect
        github.com/prometheus/procfs v0.2.0 // indirect
        github.com/soheilhy/cmux v0.1.4
        golang.org/x/net v0.0.0-20201031054903-ff519b6c9102
        google.golang.org/genproto v0.0.0-20201104152603-2e45c02ce95c
        google.golang.org/grpc v1.33.1
        google.golang.org/protobuf v1.25.0
        gopkg.in/yaml.v2 v2.3.0
        k8s.io/api v0.19.0-alpha.1
        k8s.io/apimachinery v0.19.0-alpha.1
        k8s.io/client-go v0.18.6
        k8s.io/gengo v0.0.0-20200413195148-3a45101e95ac // indirect
        k8s.io/klog v1.0.0 // indirect
        k8s.io/klog/v2 v2.2.0 // indirect
        k8s.io/utils v0.0.0-20201104234853-8146046b121e // indirect
        sigs.k8s.io/structured-merge-diff/v4 v4.0.1 // indirect
)

replace (
        github.com/golang/protobuf => github.com/golang/protobuf v1.1.0
        github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.1.0
        github.com/grpc-ecosystem/grpc-gateway => github.com/grpc-ecosystem/grpc-gateway v1.4.1
        github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.1
        github.com/prometheus/client_model => github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910
        github.com/prometheus/common => github.com/prometheus/common v0.0.0-20181020173914-7e9e6cabbd39
        google.golang.org/genproto => google.golang.org/genproto v0.0.0-20180709204101-e92b11657268
        google.golang.org/grpc => google.golang.org/grpc v1.13.0
        k8s.io/api => k8s.io/api v0.18.6
        k8s.io/apimachinery => k8s.io/apimachinery v0.18.6
)
