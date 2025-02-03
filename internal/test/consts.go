package internal_test

const (
	AdapterTag                  string = "nodeset/hyperdrive-example-adapter:v0.2.0"
	ServiceTag                  string = "nodeset/hyperdrive-example-service:v0.2.0"
	ProjectName                 string = "hde-test"
	GlobalAdapterContainerName  string = "hd_em_adapter"
	ProjectAdapterContainerName string = "hd_" + ProjectName + "_em_adapter"
	ServiceContainerName        string = ProjectName + "_example"
	LogDir                      string = "/tmp/hde-adapter-test/log"
	SystemDir                   string = "/tmp/hde-adapter-test/sys"
	CfgDir                      string = "/tmp/hde-adapter-test/cfg"
	UserDir                     string = "/tmp/hde-adapter-test/user"
	KeyPath                     string = UserDir + "/secrets/adapter.key"
	UserDataPath                string = "/tmp/hde-adapter-test/data"
	TestKey                     string = "test-key"
)
