module github.com/hashicorp/terraform/internal/backend/remote-state/cos

go 1.22.0

require (
	github.com/hashicorp/terraform v0.0.0-00010101000000-000000000000
	github.com/hashicorp/terraform/internal/legacy v0.0.0-00010101000000-000000000000
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.588
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts v1.0.588
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag v1.0.233
	github.com/tencentyun/cos-go-sdk-v5 v0.7.42
)

require (
	cloud.google.com/go v0.110.7 // indirect
	cloud.google.com/go/compute v1.23.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.1 // indirect
	cloud.google.com/go/storage v1.30.1 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/apparentlymart/go-versions v1.0.1 // indirect
	github.com/aws/aws-sdk-go v1.44.122 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/clbanning/mxj v1.8.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.11.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.7.3 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-slug v0.14.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl/v2 v2.20.0 // indirect
	github.com/hashicorp/terraform-registry-address v0.2.3 // indirect
	github.com/hashicorp/terraform-svchost v0.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/mozillazg/go-httpheader v0.3.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/zclconf/go-cty v1.14.3 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/oauth2 v0.17.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.126.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

replace github.com/hashicorp/terraform/internal/backend/remote-state/azure => ../azure

replace github.com/hashicorp/terraform/internal/backend/remote-state/consul => ../consul

replace github.com/hashicorp/terraform/internal/backend/remote-state/cos => ../cos

replace github.com/hashicorp/terraform/internal/backend/remote-state/gcs => ../gcs

replace github.com/hashicorp/terraform/internal/backend/remote-state/kubernetes => ../kubernetes

replace github.com/hashicorp/terraform/internal/backend/remote-state/oss => ../oss

replace github.com/hashicorp/terraform/internal/backend/remote-state/pg => ../pg

replace github.com/hashicorp/terraform/internal/backend/remote-state/s3 => ../s3

replace github.com/hashicorp/terraform/internal/legacy => ../../../legacy

replace github.com/hashicorp/terraform => ../../../..
