module github.com/nodeset-org/hyperdrive

go 1.24

toolchain go1.24.2

require github.com/rocket-pool/node-manager-core v0.5.2-0.20250430074613-76bcf6bb1be0

replace github.com/rocket-pool/node-manager-core => github.com/nodeset-org/node-manager-core v0.6.0

require (
	github.com/alessio/shellescape v1.4.2
	github.com/blang/semver/v4 v4.0.0
	github.com/docker/docker v27.3.1+incompatible
	github.com/dustin/go-humanize v1.0.1
	github.com/ethereum/go-ethereum v1.15.8
	github.com/fatih/color v1.17.0
	github.com/gdamore/tcell/v2 v2.7.4
	github.com/goccy/go-json v0.10.4
	github.com/mholt/archiver/v4 v4.0.0-alpha.8
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nodeset-org/hyperdrive-constellation v1.0.1
	github.com/nodeset-org/hyperdrive-daemon v1.2.2
	github.com/nodeset-org/hyperdrive-stakewise v1.2.2
	github.com/nodeset-org/osha v0.4.0
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58
	github.com/rivo/tview v0.0.0-20230208211350-7dfff1ce7854 // DO NOT UPGRADE
	github.com/rocket-pool/rocketpool-go/v2 v2.0.0-b2.0.20240709170030-c27aeb5fb99b
	github.com/rocket-pool/smartnode/v2 v2.0.0-olddev.0.20240710181452-edcbd6208bdd
	github.com/stretchr/testify v1.10.0
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/urfave/cli/v2 v2.27.5
	golang.org/x/sync v0.11.0
	golang.org/x/term v0.29.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	filippo.io/age v1.2.0 // indirect
	github.com/AdaLogics/go-fuzz-headers v0.0.0-20240716105424-66b64c4bb379 // indirect
	github.com/AdamKorcz/go-118-fuzz-build v0.0.0-20231105174938-2b5cbb29f3e2 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Microsoft/hcsshim v0.12.5 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.17.0 // indirect
	github.com/bodgit/plumbing v1.2.0 // indirect
	github.com/bodgit/sevenzip v1.3.0 // indirect
	github.com/bodgit/windows v1.0.0 // indirect
	github.com/btcsuite/btcd v0.24.2 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.4 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.5 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/compose-spec/compose-go/v2 v2.1.3 // indirect
	github.com/connesc/cipherio v0.2.1 // indirect
	github.com/consensys/bavard v0.1.22 // indirect
	github.com/consensys/gnark-crypto v0.14.0 // indirect
	github.com/containerd/cgroups/v3 v3.0.3 // indirect
	github.com/containerd/containerd v1.7.19 // indirect
	github.com/containerd/containerd/api v1.7.19 // indirect
	github.com/containerd/continuity v0.4.3 // indirect
	github.com/containerd/errdefs v0.1.0 // indirect
	github.com/containerd/fifo v1.1.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/platforms v0.2.1 // indirect
	github.com/containerd/ttrpc v1.2.5 // indirect
	github.com/containerd/typeurl/v2 v2.2.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20240724233137-53bbb0ceb27a // indirect
	github.com/crate-crypto/go-kzg-4844 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/ethereum/c-kzg-4844 v1.0.2 // indirect
	github.com/ethereum/go-verkle v0.2.2 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/ferranbt/fastssz v0.1.3 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/glendc/go-external-ip v0.1.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/herumi/bls-eth-go-binary v1.36.1 // indirect
	github.com/holiman/uint256 v1.3.2 // indirect
	github.com/ipfs/go-cid v0.4.1 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mattn/go-shellwords v1.0.12 // indirect
	github.com/minio/highwayhash v1.0.3 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/patternmatcher v0.6.0 // indirect
	github.com/moby/swarmkit/v2 v2.0.0-20240611172349-ea1a7cec35cb // indirect
	github.com/moby/sys/mount v0.3.3 // indirect
	github.com/moby/sys/mountinfo v0.7.1 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/sys/signal v0.7.0 // indirect
	github.com/moby/sys/symlink v0.2.0 // indirect
	github.com/moby/sys/user v0.1.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multibase v0.2.0 // indirect
	github.com/multiformats/go-multihash v0.2.3 // indirect
	github.com/multiformats/go-varint v0.0.7 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nodeset-org/nodeset-client-go v1.3.1 // indirect
	github.com/nwaples/rardecode/v2 v2.0.0-beta.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/opencontainers/selinux v1.11.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/prysmaticlabs/fastssz v0.0.0-20240620202422-a981b8ef89d3 // indirect
	github.com/prysmaticlabs/go-bitfield v0.0.0-20240328144219-a1caa50c3a1e // indirect
	github.com/prysmaticlabs/gohashtree v0.0.4-beta.0.20240624100937-73632381301b // indirect
	github.com/prysmaticlabs/prysm/v5 v5.1.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rocket-pool/batch-query v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sethvargo/go-password v0.3.1 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/supranational/blst v0.3.14 // indirect
	github.com/therootcompany/xz v1.0.1 // indirect
	github.com/thomaso-mirodin/intmath v0.0.0-20160323211736-5dc6d854e46e // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/vbatts/tar-split v0.11.5 // indirect
	github.com/wealdtech/go-bytesutil v1.2.1 // indirect
	github.com/wealdtech/go-ens/v3 v3.6.0 // indirect
	github.com/wealdtech/go-eth2-types/v2 v2.8.2 // indirect
	github.com/wealdtech/go-eth2-util v1.8.2 // indirect
	github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4 v1.4.1 // indirect
	github.com/wealdtech/go-merkletree v1.0.1-0.20190605192610-2bb163c2ea2a // indirect
	github.com/wealdtech/go-multicodec v1.4.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.etcd.io/etcd/raft/v3 v3.5.14 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.53.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go4.org v0.0.0-20200411211856-f5505b9728dd // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/exp v0.0.0-20240716160929-1d5bc16f04a8 // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20240711142825-46eb208f015d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240711142825-46eb208f015d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240711142825-46eb208f015d // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	lukechampine.com/blake3 v1.2.1 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)

replace github.com/rocket-pool/smartnode/v2 => github.com/nodeset-org/rocketpool-smartnode/v2 v2.0.0-olddev.0.20250130225832-4768ad9c2c82
