// +build tools


package ent

import (
	"sifu-box/ent/provider"
	"sifu-box/ent/ruleset"
	"sifu-box/ent/schema"
	"sifu-box/ent/template"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	providerFields := schema.Provider{}.Fields()
	_ = providerFields
	// providerDescName is the schema descriptor for name field.
	providerDescName := providerFields[0].Descriptor()
	// provider.NameValidator is a validator for the "name" field. It is called by the builders before save.
	provider.NameValidator = func() func(string) error {
		validators := providerDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// providerDescPath is the schema descriptor for path field.
	providerDescPath := providerFields[1].Descriptor()
	// provider.PathValidator is a validator for the "path" field. It is called by the builders before save.
	provider.PathValidator = func() func(string) error {
		validators := providerDescPath.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(_path string) error {
			for _, fn := range fns {
				if err := fn(_path); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// providerDescDetour is the schema descriptor for detour field.
	providerDescDetour := providerFields[2].Descriptor()
	// provider.DetourValidator is a validator for the "detour" field. It is called by the builders before save.
	provider.DetourValidator = providerDescDetour.Validators[0].(func(string) error)
	rulesetFields := schema.RuleSet{}.Fields()
	_ = rulesetFields
	// rulesetDescTag is the schema descriptor for tag field.
	rulesetDescTag := rulesetFields[0].Descriptor()
	// ruleset.TagValidator is a validator for the "tag" field. It is called by the builders before save.
	ruleset.TagValidator = func() func(string) error {
		validators := rulesetDescTag.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(tag string) error {
			for _, fn := range fns {
				if err := fn(tag); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// rulesetDescType is the schema descriptor for type field.
	rulesetDescType := rulesetFields[1].Descriptor()
	// ruleset.TypeValidator is a validator for the "type" field. It is called by the builders before save.
	ruleset.TypeValidator = func() func(string) error {
		validators := rulesetDescType.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(_type string) error {
			for _, fn := range fns {
				if err := fn(_type); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// rulesetDescPath is the schema descriptor for path field.
	rulesetDescPath := rulesetFields[2].Descriptor()
	// ruleset.PathValidator is a validator for the "path" field. It is called by the builders before save.
	ruleset.PathValidator = func() func(string) error {
		validators := rulesetDescPath.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(_path string) error {
			for _, fn := range fns {
				if err := fn(_path); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// rulesetDescFormat is the schema descriptor for format field.
	rulesetDescFormat := rulesetFields[3].Descriptor()
	// ruleset.FormatValidator is a validator for the "format" field. It is called by the builders before save.
	ruleset.FormatValidator = func() func(string) error {
		validators := rulesetDescFormat.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(format string) error {
			for _, fn := range fns {
				if err := fn(format); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// rulesetDescLabel is the schema descriptor for label field.
	rulesetDescLabel := rulesetFields[4].Descriptor()
	// ruleset.LabelValidator is a validator for the "label" field. It is called by the builders before save.
	ruleset.LabelValidator = func() func(string) error {
		validators := rulesetDescLabel.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(label string) error {
			for _, fn := range fns {
				if err := fn(label); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// rulesetDescDownloadDetour is the schema descriptor for download_detour field.
	rulesetDescDownloadDetour := rulesetFields[5].Descriptor()
	// ruleset.DownloadDetourValidator is a validator for the "download_detour" field. It is called by the builders before save.
	ruleset.DownloadDetourValidator = rulesetDescDownloadDetour.Validators[0].(func(string) error)
	// rulesetDescUpdateInterval is the schema descriptor for update_interval field.
	rulesetDescUpdateInterval := rulesetFields[6].Descriptor()
	// ruleset.UpdateIntervalValidator is a validator for the "update_interval" field. It is called by the builders before save.
	ruleset.UpdateIntervalValidator = rulesetDescUpdateInterval.Validators[0].(func(string) error)
	// rulesetDescNameServer is the schema descriptor for name_server field.
	rulesetDescNameServer := rulesetFields[7].Descriptor()
	// ruleset.NameServerValidator is a validator for the "name_server" field. It is called by the builders before save.
	ruleset.NameServerValidator = rulesetDescNameServer.Validators[0].(func(string) error)
	templateFields := schema.Template{}.Fields()
	_ = templateFields
	// templateDescName is the schema descriptor for name field.
	templateDescName := templateFields[0].Descriptor()
	// template.NameValidator is a validator for the "name" field. It is called by the builders before save.
	template.NameValidator = func() func(string) error {
		validators := templateDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
}
