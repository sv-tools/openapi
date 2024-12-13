package openapi

import (
	"fmt"
	"strings"
)

// OpenAPI is the root object of the OpenAPI document.
//
// https://spec.openapis.org/oas/v3.1.0#openapi-object
//
// Example:
//
//	openapi: 3.1.0
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
	WebHooks map[string]*RefOrSpec[Extendable[PathItem]] `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
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

func checkUnusedComponent[T any](name string, m map[string]T, opts *specValidationOptions) []*validationError {
	var errs []*validationError
	for k := range m {
		id := joinLoc("#", "components", name, k)
		if !opts.visited[id] {
			errs = append(errs, newValidationError(id, ErrUnused))
		}
	}
	return errs
}

func (o *OpenAPI) validateSpec(loc string, opts *specValidationOptions) []*validationError {
	var errs []*validationError
	if o.OpenAPI == "" {
		errs = append(errs, newValidationError(joinLoc(loc, "openapi"), ErrRequired))
	} else {
		if !strings.HasPrefix(o.OpenAPI, "3.1.") {
			errs = append(errs, newValidationError(joinLoc(loc, "openapi"), fmt.Errorf("unsupported version: %s", o.OpenAPI)))
		}
	}
	if o.Info == nil {
		errs = append(errs, newValidationError(joinLoc(loc, "info"), ErrRequired))
	} else {
		errs = append(errs, o.Info.validateSpec(joinLoc(loc, "info"), opts)...)
	}

	// validate tags first to memorize them for later checking
	if o.Tags != nil {
		for i, tag := range o.Tags {
			errs = append(errs, tag.validateSpec(joinLoc(loc, "tag", i), opts)...)
		}
	}

	if err := checkURL(o.JsonSchemaDialect); err != nil {
		errs = append(errs, newValidationError(joinLoc(loc, "jsonSchemaDialect"), err))
	}
	if o.Servers != nil {
		for i, server := range o.Servers {
			errs = append(errs, server.validateSpec(joinLoc(loc, "servers", i), opts)...)
		}
	}
	if o.Paths != nil {
		errs = append(errs, o.Paths.validateSpec(joinLoc(loc, "paths"), opts)...)
	}
	if o.WebHooks != nil {
		for name, webhook := range o.WebHooks {
			errs = append(errs, webhook.validateSpec(joinLoc(loc, "webhooks", name), opts)...)
		}
	}
	if o.Components != nil {
		errs = append(errs, o.Components.validateSpec(joinLoc(loc, "components"), opts)...)
	}
	if o.Security != nil {
		for i, security := range o.Security {
			errs = append(errs, security.validateSpec(joinLoc(loc, "security", i), opts)...)
		}
	}
	if o.ExternalDocs != nil {
		errs = append(errs, o.Components.validateSpec(joinLoc(loc, "externalDocs"), opts)...)
	}
	if o.Paths == nil && o.WebHooks == nil && o.Components == nil {
		errs = append(errs, newValidationError(joinLoc(loc, "paths||webhooks||components"), ErrRequired))
	}

	// check for unused
	for i, t := range o.Tags {
		if !opts.visited[joinLoc("tags", t.Spec.Name, "used")] {
			errs = append(errs, newValidationError(joinLoc(loc, "tags", i), fmt.Errorf("'%s': %w", t.Spec.Name, ErrUnused)))
		}
	}
	if o.Components != nil && !opts.allowUnusedComponents {
		errs = append(errs, checkUnusedComponent("schemas", o.Components.Spec.Schemas, opts)...)
		errs = append(errs, checkUnusedComponent("responses", o.Components.Spec.Responses, opts)...)
		errs = append(errs, checkUnusedComponent("parameters", o.Components.Spec.Parameters, opts)...)
		errs = append(errs, checkUnusedComponent("examples", o.Components.Spec.Examples, opts)...)
		errs = append(errs, checkUnusedComponent("requestBodies", o.Components.Spec.RequestBodies, opts)...)
		errs = append(errs, checkUnusedComponent("headers", o.Components.Spec.Headers, opts)...)
		errs = append(errs, checkUnusedComponent("securitySchemes", o.Components.Spec.SecuritySchemes, opts)...)
		errs = append(errs, checkUnusedComponent("links", o.Components.Spec.Links, opts)...)
		errs = append(errs, checkUnusedComponent("callbacks", o.Components.Spec.Callbacks, opts)...)
		errs = append(errs, checkUnusedComponent("paths", o.Components.Spec.Paths, opts)...)
	}

	for k, v := range opts.linkToOperationID {
		if !opts.visited[joinLoc("operations", v)] {
			errs = append(errs, newValidationError(k, "'%s' not found", v))
		}
	}
	return errs
}
