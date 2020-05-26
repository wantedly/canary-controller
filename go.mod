module github.com/wantedly/canary-controller

go 1.14

require (
	cloud.google.com/go v0.34.0
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973
	github.com/davecgh/go-spew v1.1.1
	github.com/emicklei/go-restful v2.8.0+incompatible
	github.com/fsnotify/fsnotify v1.4.7
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.0
	github.com/gobuffalo/envy v1.6.11
	github.com/gogo/protobuf v1.2.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20181024230925-c65c006176ff
	github.com/golang/protobuf v1.2.0
	github.com/google/btree v1.0.0
	github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf
	github.com/google/uuid v1.1.0
	github.com/googleapis/gnostic v0.2.0
	github.com/gregjones/httpcache v0.0.0-20181110185634-c63ab54fda8f
	github.com/hashicorp/golang-lru v0.5.0
	github.com/hpcloud/tail v1.0.0
	github.com/imdario/mergo v0.3.6
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.5
	github.com/markbates/inflect v1.0.4
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/petar/GoLLRB v0.0.0-20130427215148-53be0d36a84c
	github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/pkg/errors v0.8.0
	github.com/prometheus/client_golang v0.9.2
	github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910
	github.com/prometheus/common v0.0.0-20181218105931-67670fe90761
	github.com/prometheus/procfs v0.0.0-20181204211112-1dc9a6cbc91a
	github.com/rogpeppe/go-internal v1.0.0
	github.com/spf13/afero v1.2.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	go.uber.org/atomic v1.3.2
	go.uber.org/multierr v1.1.0
	go.uber.org/zap v1.9.1
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9
	golang.org/x/net v0.0.0-20181217023233-e147a9138326
	golang.org/x/oauth2 v0.0.0-20181203162652-d668ce993890
	golang.org/x/sys v0.0.0-20181218192612-074acd46bca6
	golang.org/x/text v0.3.0
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c
	golang.org/x/tools v0.0.0-20181218204010-d4971274fe38
	google.golang.org/appengine v1.3.0
	gopkg.in/inf.v0 v0.9.1
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/api v0.0.0-20181126151915-b503174bad59
	k8s.io/apiextensions-apiserver v0.0.0-20181126155829-0cd23ebeb688
	k8s.io/apimachinery v0.0.0-20181126123746-eddba98df674
	k8s.io/client-go v0.0.0-20181126152608-d082d5923d3c
	k8s.io/code-generator v0.0.0-20181206115026-3a2206dd6a78
	k8s.io/gengo v0.0.0-20181113154421-fd15ee9cc2f7
	k8s.io/klog v0.1.0
	k8s.io/kube-openapi v0.0.0-20181114233023-0317810137be
	sigs.k8s.io/controller-runtime v0.1.8
	sigs.k8s.io/controller-tools v0.1.7
	sigs.k8s.io/testing_frameworks v0.1.0
)

replace gopkg.in/fsnotify.v1 v1.4.7 => github.com/fsnotify/fsnotify v1.4.7
