package openapi

// SecurityRequirement is the lists of the required security schemes to execute this operation.
// The name used for each property MUST correspond to a security scheme declared in the Security Schemes under the Components Object.
// Security Requirement Objects that contain multiple schemes require that all schemes MUST be satisfied for a request to be authorized.
// This enables support for scenarios where multiple query parameters or HTTP headers are required to convey security information.
// When a list of Security Requirement Objects is defined on the OpenAPI Object or Operation Object,
// only one of the Security Requirement Objects in the list needs to be satisfied to authorize the request.
//
// https://spec.openapis.org/oas/v3.1.1#security-requirement-object
//
// Example:
//
//	api_key: []
type SecurityRequirement map[string][]string

func (o *SecurityRequirement) validateSpec(_ string, validator *Validator) []*validationError { //nolint: unparam // by design
	for k := range *o {
		validator.visited[joinLoc("#", "components", "securitySchemes", k)] = true
	}
	return nil // nothing to validate
}

type SecurityRequirementBuilder struct {
	spec SecurityRequirement
}

func NewSecurityRequirementBuilder() *SecurityRequirementBuilder {
	return &SecurityRequirementBuilder{
		spec: make(SecurityRequirement),
	}
}

func (b *SecurityRequirementBuilder) Build() *SecurityRequirement {
	return &b.spec
}

func (b *SecurityRequirementBuilder) Add(name string, scopes ...string) *SecurityRequirementBuilder {
	b.spec[name] = append(b.spec[name], scopes...)
	return b
}
