package config

import (
	"fmt"
	"strconv"

	"github.com/nodeset-org/hyperdrive/modules/config"
)

type enableParamInstance struct {
	info    *config.ModuleInfo
	intance *config.ModuleInstance
}

func NewEnableParamInstance(info *config.ModuleInfo, instance *config.ModuleInstance) *enableParamInstance {
	return &enableParamInstance{
		info:    info,
		intance: instance,
	}
}

func (e *enableParamInstance) GetID() config.Identifier {
	return ""
}

func (e *enableParamInstance) GetName() string {
	return "Enable " + string(e.info.Descriptor.Name)
}

func (e *enableParamInstance) GetDescription() config.DynamicProperty[string] {
	return config.DynamicProperty[string]{
		Default:  "Enable the " + string(e.info.Descriptor.Name) + " module. Check this to configure the module, and start it when you start the Hyperdrive service.",
		Template: "",
	}
}

func (e *enableParamInstance) GetType() config.ParameterType {
	return config.ParameterType_Bool
}

func (e *enableParamInstance) GetDefault() any {
	return false
}

func (e *enableParamInstance) GetAdvanced() bool {
	return false
}

func (e *enableParamInstance) GetDisabled() config.DynamicProperty[bool] {
	return config.DynamicProperty[bool]{
		Default:  false,
		Template: "",
	}
}

func (e *enableParamInstance) GetHidden() config.DynamicProperty[bool] {
	return config.DynamicProperty[bool]{
		Default:  false,
		Template: "",
	}
}

func (e *enableParamInstance) GetOverwriteOnUpgrade() bool {
	return false
}

func (e *enableParamInstance) GetAffectedContainers() []string {
	return []string{}
}

func (e *enableParamInstance) Serialize() map[string]any {
	return nil
}

func (e *enableParamInstance) Deserialize(data map[string]any) error {
	return nil
}

func (e *enableParamInstance) CreateSetting() config.IParameterSetting {
	return nil
}

func (e *enableParamInstance) GetMetadata() config.IParameter {
	return e
}

func (e *enableParamInstance) GetValue() any {
	return e.intance.Enabled
}

func (e *enableParamInstance) SetValue(value any) error {
	val, ok := value.(bool)
	if !ok {
		return fmt.Errorf("invalid value type for module [%s] enable flag: %T", e.info.Descriptor.GetFullyQualifiedModuleName(), value)
	}
	e.intance.Enabled = val
	return nil
}

func (e *enableParamInstance) String() string {
	return strconv.FormatBool(e.intance.Enabled)
}

func (e *enableParamInstance) Validate() []error {
	return nil
}
