package openapi

import "github.com/santhosh-tekuri/jsonschema/v6"

type validationOptions struct {
	allowExtensionNameWithoutPrefix bool
	allowRequestBodyForGet          bool
	allowRequestBodyForHead         bool
	allowRequestBodyForDelete       bool
	allowUndefinedTagsInOperation   bool
	allowUnusedComponents           bool
	doNotValidateExamples           bool
	doNotValidateDefaultValues      bool
	validateDataAsJSON              bool
	updateCompiler                  []func(*jsonschema.Compiler)
}

// ValidationOption is a type for validation options.
type ValidationOption func(*validationOptions)

// AllowExtensionNameWithoutPrefix is a validation option to allow extension name without `x-` prefix.
func AllowExtensionNameWithoutPrefix() ValidationOption {
	return func(v *validationOptions) {
		v.allowExtensionNameWithoutPrefix = true
	}
}

// AllowRequestBodyForGet is a validation option to allow request body for GET operation.
func AllowRequestBodyForGet() ValidationOption {
	return func(v *validationOptions) {
		v.allowRequestBodyForGet = true
	}
}

// AllowRequestBodyForHead is a validation option to allow request body for HEAD operation.
func AllowRequestBodyForHead() ValidationOption {
	return func(v *validationOptions) {
		v.allowRequestBodyForHead = true
	}
}

// AllowRequestBodyForDelete is a validation option to allow request body for DELETE operation.
func AllowRequestBodyForDelete() ValidationOption {
	return func(v *validationOptions) {
		v.allowRequestBodyForDelete = true
	}
}

// AllowUndefinedTagsInOperation is a validation option to allow undefined tags in operation.
func AllowUndefinedTagsInOperation() ValidationOption {
	return func(v *validationOptions) {
		v.allowUndefinedTagsInOperation = true
	}
}

// AllowUnusedComponents is a validation option to allow unused components.
func AllowUnusedComponents() ValidationOption {
	return func(v *validationOptions) {
		v.allowUnusedComponents = true
	}
}

// DoNotValidateExamples is a validation option to skip examples validation.
func DoNotValidateExamples() ValidationOption {
	return func(v *validationOptions) {
		v.doNotValidateExamples = true
	}
}

// DoNotValidateDefaultValues is a validation option to skip default values validation.
func DoNotValidateDefaultValues() ValidationOption {
	return func(v *validationOptions) {
		v.doNotValidateDefaultValues = true
	}
}

func ValidateStringDataAsJSON() ValidationOption {
	return func(v *validationOptions) {
		v.validateDataAsJSON = true
	}
}

// UpdateCompiler is a type to modify the jsonschema.Compiler.
func UpdateCompiler(f func(*jsonschema.Compiler)) ValidationOption {
	return func(v *validationOptions) {
		v.updateCompiler = append(v.updateCompiler, f)
	}
}
