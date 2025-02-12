package internal_test

const (
	AdapterTag                  string = "nodeset/hyperdrive-example-adapter:v0.2.0"
	ServiceTag                  string = "nodeset/hyperdrive-example-service:v0.2.0"
	ProjectName                 string = "hde-test"
	GlobalAdapterContainerName  string = "hd-em_adapter"
	ProjectAdapterContainerName string = "hd_" + ProjectName + "_em_adapter"
	ServiceContainerName        string = ProjectName + "-em_example"
	LogDir                      string = "/tmp/hde-adapter-test/log"
	SystemDir                   string = "/tmp/hde-adapter-test/sys"
	CfgDir                      string = "/tmp/hde-adapter-test/cfg"
	UserDir                     string = "/tmp/hde-adapter-test/user"
	RuntimeDir                  string = "/tmp/hde-adapter-test/runtime"
	DataDir                     string = "/tmp/hde-adapter-test/data"
	AdapterKeyPath              string = UserDir + "/secrets/adapter.key"
	UserDataPath                string = "/tmp/hde-adapter-test/data"
	TestKey                     string = "test-key"
)

type InternalTestTemplateSource struct {
}

func (t *InternalTestTemplateSource) AdapterTag() string {
	return AdapterTag
}

func (t *InternalTestTemplateSource) ServiceTag() string {
	return ServiceTag
}

func (t *InternalTestTemplateSource) ProjectName() string {
	return ProjectName
}

func (t *InternalTestTemplateSource) GlobalAdapterContainerName() string {
	return GlobalAdapterContainerName
}

func (t *InternalTestTemplateSource) ProjectAdapterContainerName() string {
	return ProjectAdapterContainerName
}

func (t *InternalTestTemplateSource) ServiceContainerName() string {
	return ServiceContainerName
}

func (t *InternalTestTemplateSource) LogDir() string {
	return LogDir
}

func (t *InternalTestTemplateSource) SystemDir() string {
	return SystemDir
}

func (t *InternalTestTemplateSource) CfgDir() string {
	return CfgDir
}

func (t *InternalTestTemplateSource) UserDir() string {
	return UserDir
}

func (t *InternalTestTemplateSource) RuntimeDir() string {
	return RuntimeDir
}

func (t *InternalTestTemplateSource) DataDir() string {
	return DataDir
}

func (t *InternalTestTemplateSource) KeyPath() string {
	return AdapterKeyPath
}

func (t *InternalTestTemplateSource) UserDataPath() string {
	return UserDataPath
}

func (t *InternalTestTemplateSource) TestKey() string {
	return TestKey
}
