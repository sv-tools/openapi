package openapi

import (
	"fmt"
	"strings"
)

// OpenAPI is the root object of the OpenAPI document.
//
// https://spec.openapis.org/oas/v3.1.1#openapi-object
//
// Example:
//
//	openapi: 3.1.1
//	info:
//	  title: Minimal OpenAPI example
//	  version: 1.0.0
//	paths: { }
type OpenAPI struct {
	// An element to hold various schemas for the document.
	Components *Extendable[Components] `json:"components,omitempty" yaml:"components,omitempty"`
	// REQUIRED
	// Provides metadata about the API. The metadata MAY be used by tooling as required.
	Info *Extendable[Info] `json:"info" yaml:"info"`
	// Additional external documentation.
	ExternalDocs *Extendable[ExternalDocs] `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// Holds the relative paths to the individual endpoints and their operations.
	// The path is appended to the URL from the Server Object in order to construct the full URL.
	// The Paths MAY be empty, due to Access Control List (ACL) constraints.
	Paths *Extendable[Paths] `json:"paths,omitempty" yaml:"paths,omitempty"`
	// The incoming webhooks that MAY be received as part of this API and that the API consumer MAY choose to implement.
	// Closely related to the callbacks feature, this section describes requests initiated other than by an API call,
	// for example by an out of band registration.
	// The key name is a unique string to refer to each webhook, while the (optionally referenced) PathItem Object describes
	// a request that may be initiated by the API provider and the expected responses.
	WebHooks Webhooks `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
	// The default value for the $schema keyword within Schema Objects contained within this OAS document.
	// This MUST be in the form of a URI.
	JsonSchemaDialect string `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`
	// REQUIRED
	// This string MUST be the version number of the OpenAPI Specification that the OpenAPI document uses.
	// The openapi field SHOULD be used by tooling to interpret the OpenAPI document.
	// This is not related to the API info.version string.
	OpenAPI string `json:"openapi" yaml:"openapi"`
	// A declaration of which security mechanisms can be used across the API.
	// The list of values includes alternative security requirement objects that can be used.
	// Only one of the security requirement objects need to be satisfied to authorize a request.
	// Individual operations can override this definition.
	// To make security optional, an empty security requirement ({}) can be included in the array.
	Security []SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	// A list of tags used by the document with additional metadata.
	// The order of the tags can be used to reflect on their order by the parsing tools.
	// Not all tags that are used by the Operation Object must be declared.
	// The tags that are not declared MAY be organized randomly or based on the tools’ logic.
	// Each tag name in the list MUST be unique.
	Tags []*Extendable[Tag] `json:"tags,omitempty" yaml:"tags,omitempty"`
	// An array of Server Objects, which provide connectivity information to a target server.
	// If the servers property is not provided, or is an empty array, the default value would be a Server Object with a url value of /.
	Servers []*Extendable[Server] `json:"servers,omitempty" yaml:"servers,omitempty"`
}

func checkUnusedComponent[T any](name string, m map[string]T, validator *Validator) []*validationError {
	var errs []*validationError
	for k := range m {
		id := joinLoc("#", "components", name, k)
		if !validator.visited[id] {
			errs = append(errs, newValidationError(id, ErrUnused))
		}
	}
	return errs
}

func (o *OpenAPI) validateSpec(location string, validator *Validator) []*validationError {
	var errs []*validationError
	if o.OpenAPI == "" {
		errs = append(errs, newValidationError(joinLoc(location, "openapi"), ErrRequired))
	} else {
		if !strings.HasPrefix(o.OpenAPI, "3.1.") {
			errs = append(errs, newValidationError(joinLoc(location, "openapi"), fmt.Errorf("unsupported version: %s", o.OpenAPI)))
		}
	}
	if o.Info == nil {
		errs = append(errs, newValidationError(joinLoc(location, "info"), ErrRequired))
	} else {
		errs = append(errs, o.Info.validateSpec(joinLoc(location, "info"), validator)...)
	}

	// validate tags first to memorize them for later checking
	if o.Tags != nil {
		for i, tag := range o.Tags {
			errs = append(errs, tag.validateSpec(joinLoc(location, "tag", i), validator)...)
		}
	}

	if err := checkURL(o.JsonSchemaDialect); err != nil {
		errs = append(errs, newValidationError(joinLoc(location, "jsonSchemaDialect"), err))
	}
	if o.Servers != nil {
		for i, server := range o.Servers {
			errs = append(errs, server.validateSpec(joinLoc(location, "servers", i), validator)...)
		}
	}
	if o.Paths != nil {
		errs = append(errs, o.Paths.validateSpec(joinLoc(location, "paths"), validator)...)
	}
	if o.WebHooks != nil {
		for name, webhook := range o.WebHooks {
			errs = append(errs, webhook.validateSpec(joinLoc(location, "webhooks", name), validator)...)
		}
	}
	if o.Components != nil {
		errs = append(errs, o.Components.validateSpec(joinLoc(location, "components"), validator)...)
	}
	if o.Security != nil {
		for i, security := range o.Security {
			errs = append(errs, security.validateSpec(joinLoc(location, "security", i), validator)...)
		}
	}
	if o.ExternalDocs != nil {
		errs = append(errs, o.ExternalDocs.validateSpec(joinLoc(location, "externalDocs"), validator)...)
	}
	if o.Paths == nil && o.WebHooks == nil && o.Components == nil {
		errs = append(errs, newValidationError(joinLoc(location, "paths||webhooks||components"), ErrRequired))
	}

	// check for unused
	for i, t := range o.Tags {
		if !validator.visited[joinLoc("tags", t.Spec.Name, "used")] {
			errs = append(errs, newValidationError(joinLoc(location, "tags", i), fmt.Errorf("'%s': %w", t.Spec.Name, ErrUnused)))
		}
	}
	if o.Components != nil && !validator.opts.allowUnusedComponents {
		errs = append(errs, checkUnusedComponent("schemas", o.Components.Spec.Schemas, validator)...)
		errs = append(errs, checkUnusedComponent("responses", o.Components.Spec.Responses, validator)...)
		errs = append(errs, checkUnusedComponent("parameters", o.Components.Spec.Parameters, validator)...)
		errs = append(errs, checkUnusedComponent("examples", o.Components.Spec.Examples, validator)...)
		errs = append(errs, checkUnusedComponent("requestBodies", o.Components.Spec.RequestBodies, validator)...)
		errs = append(errs, checkUnusedComponent("headers", o.Components.Spec.Headers, validator)...)
		errs = append(errs, checkUnusedComponent("securitySchemes", o.Components.Spec.SecuritySchemes, validator)...)
		errs = append(errs, checkUnusedComponent("links", o.Components.Spec.Links, validator)...)
		errs = append(errs, checkUnusedComponent("callbacks", o.Components.Spec.Callbacks, validator)...)
		errs = append(errs, checkUnusedComponent("paths", o.Components.Spec.Paths, validator)...)
	}

	for k, v := range validator.linkToOperationID {
		if !validator.visited[joinLoc("operations", v)] {
			errs = append(errs, newValidationError(k, "'%s' not found", v))
		}
	}
	return errs
}

type OpenAPIBuilder struct {
	spec *Extendable[OpenAPI]
}

func NewOpenAPIBuilder() *OpenAPIBuilder {
	return &OpenAPIBuilder{spec: NewExtendable(&OpenAPI{
		OpenAPI:           "3.1.1",
		JsonSchemaDialect: "https://spec.openapis.org/oas/3.1/dialect/base",
	})}
}

func (b *OpenAPIBuilder) Build() *Extendable[OpenAPI] {
	return b.spec
}

func (b *OpenAPIBuilder) Extensions(v map[string]any) *OpenAPIBuilder {
	b.spec.Extensions = v
	return b
}

func (b *OpenAPIBuilder) AddExt(name string, value any) *OpenAPIBuilder {
	b.spec.AddExt(name, value)
	return b
}

func (b *OpenAPIBuilder) OpenAPI(openAPI string) *OpenAPIBuilder {
	b.spec.Spec.OpenAPI = openAPI
	return b
}

func (b *OpenAPIBuilder) Info(info *Extendable[Info]) *OpenAPIBuilder {
	b.spec.Spec.Info = info
	return b
}

func (b *OpenAPIBuilder) Components(components *Extendable[Components]) *OpenAPIBuilder {
	b.spec.Spec.Components = components
	return b
}

func (b *OpenAPIBuilder) AddComponent(name string, component any) *OpenAPIBuilder {
	if b.spec.Spec.Components == nil {
		b.spec.Spec.Components = NewComponents()
	}
	b.spec.Spec.Components.Spec.Add(name, component)
	return b
}

func (b *OpenAPIBuilder) ExternalDocs(externalDocs *Extendable[ExternalDocs]) *OpenAPIBuilder {
	b.spec.Spec.ExternalDocs = externalDocs
	return b
}

func (b *OpenAPIBuilder) Paths(paths *Extendable[Paths]) *OpenAPIBuilder {
	b.spec.Spec.Paths = paths
	return b
}

func (b *OpenAPIBuilder) AddPath(path string, item *RefOrSpec[Extendable[PathItem]]) *OpenAPIBuilder {
	if b.spec.Spec.Paths == nil {
		b.spec.Spec.Paths = NewPaths()
	}
	b.spec.Spec.Paths.Spec.Add(path, item)
	return b
}

func (b *OpenAPIBuilder) WebHooks(webHooks Webhooks) *OpenAPIBuilder {
	b.spec.Spec.WebHooks = webHooks
	return b
}

func (b *OpenAPIBuilder) AddWebHook(name string, path *RefOrSpec[Extendable[PathItem]]) *OpenAPIBuilder {
	if b.spec.Spec.WebHooks == nil {
		b.spec.Spec.WebHooks = NewWebhooks()
	}
	b.spec.Spec.WebHooks[name] = path
	return b
}

func (b *OpenAPIBuilder) JsonSchemaDialect(jsonSchemaDialect string) *OpenAPIBuilder {
	b.spec.Spec.JsonSchemaDialect = jsonSchemaDialect
	return b
}

func (b *OpenAPIBuilder) Security(security ...SecurityRequirement) *OpenAPIBuilder {
	b.spec.Spec.Security = security
	return b
}

func (b *OpenAPIBuilder) AddSecurity(v ...SecurityRequirement) *OpenAPIBuilder {
	b.spec.Spec.Security = append(b.spec.Spec.Security, v...)
	return b
}

func (b *OpenAPIBuilder) Tags(tags ...*Extendable[Tag]) *OpenAPIBuilder {
	b.spec.Spec.Tags = tags
	return b
}

func (b *OpenAPIBuilder) AddTags(tags ...*Extendable[Tag]) *OpenAPIBuilder {
	b.spec.Spec.Tags = append(b.spec.Spec.Tags, tags...)
	return b
}

func (b *OpenAPIBuilder) Servers(servers ...*Extendable[Server]) *OpenAPIBuilder {
	b.spec.Spec.Servers = servers
	return b
}

func (b *OpenAPIBuilder) AddServers(servers ...*Extendable[Server]) *OpenAPIBuilder {
	b.spec.Spec.Servers = append(b.spec.Spec.Servers, servers...)
	return b
}
