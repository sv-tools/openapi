package openapi

const (
	TypeApiKey        = "apiKey"
	TypeHTTP          = "http"
	TypeMutualTLS     = "mutualTLS"
	TypeOAuth2        = "oauth2"
	TypeOpenIDConnect = "openIdConnect"
)

// SecurityScheme defines a security scheme that can be used by the operations.
// Supported schemes are HTTP authentication, an API key (either as a header, a cookie parameter or as a query parameter),
// mutual TLS (use of a client certificate), OAuth2’s common flows (implicit, password, client credentials and authorization code)
// as defined in [RFC6749], and OpenID Connect Discovery.
// Please note that as of 2020, the implicit flow is about to be deprecated by OAuth 2.0 Security Best Current Practice.
// Recommended for most use case is Authorization Code Grant flow with PKCE.
//
// https://spec.openapis.org/oas/v3.1.0#security-scheme-object
//
// Example:
//
//	type: oauth2
//	flows:
//	  implicit:
//	    authorizationUrl: https://example.com/api/oauth/dialog
//	    scopes:
//	      write:pets: modify pets in your account
//	      read:pets: read your pets
type SecurityScheme struct {
	// REQUIRED.
	// The type of the security scheme.
	// Valid values are "apiKey", "http", "mutualTLS", "oauth2", "openIdConnect".
	//
	// Applies To: any
	Type string `json:"type" yaml:"type"`
	// A description for security scheme.
	// CommonMark syntax MAY be used for rich text representation.
	//
	// Applies To: any
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// REQUIRED.
	// The name of the header, query or cookie parameter to be used.
	//
	// Applies To: apiKey
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// REQUIRED.
	// The location of the API key.
	// Valid values are "query", "header" or "cookie".
	//
	// Applies To: apiKey
	In string `json:"in,omitempty" yaml:"in,omitempty"`
	// REQUIRED.
	// The name of the HTTP Authorization scheme to be used in the Authorization header as defined in [RFC7235].
	// The values used SHOULD be registered in the IANA Authentication Scheme registry.
	//
	// Applies To: http
	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	// A hint to the client to identify how the bearer token is formatted.
	// Bearer tokens are usually generated by an authorization server, so this information is primarily for documentation purposes.
	//
	// Applies To: http ("bearer")
	BearerFormat string `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	// REQUIRED.
	// An object containing configuration information for the flow types supported.
	//
	// Applies To: oauth2
	Flows *Extendable[OAuthFlows] `json:"flows,omitempty" yaml:"flows,omitempty"`
	// REQUIRED.
	// OpenId Connect URL to discover OAuth2 configuration values.
	// This MUST be in the form of a URL.
	// The OpenID Connect standard requires the use of TLS.
	//
	// Applies To: openIdConnect
	OpenIDConnectURL string `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

func (o *SecurityScheme) validateSpec(location string, opts *specValidationOptions) []*validationError {
	var errs []*validationError
	if o.Type == "" {
		errs = append(errs, newValidationError(joinLoc(location, "type"), ErrRequired))
	} else {
		switch o.Type {
		case TypeApiKey:
			if o.Name == "" {
				errs = append(errs, newValidationError(joinLoc(location, "name"), ErrRequired))
			}
			if o.In == "" {
				errs = append(errs, newValidationError(joinLoc(location, "in"), ErrRequired))
			} else {
				switch o.In {
				case InQuery, InHeader, InCookie:
				default:
					errs = append(errs, newValidationError(joinLoc(location, "in"), "must be one of [%s, %s, %s], but got '%s'", InQuery, InHeader, InCookie, o.In))
				}
			}
		case TypeHTTP:
			if o.Scheme == "" {
				errs = append(errs, newValidationError(joinLoc(location, "scheme"), ErrRequired))
			}
		case TypeOAuth2:
			if o.Flows == nil {
				errs = append(errs, newValidationError(joinLoc(location, "flows"), ErrRequired))
			} else {
				errs = o.Flows.validateSpec(joinLoc(location, "flows"), opts)
			}
		case TypeOpenIDConnect:
			if o.OpenIDConnectURL == "" {
				errs = append(errs, newValidationError(joinLoc(location, "openIdConnectUrl"), ErrRequired))
			}
		case TypeMutualTLS:
		default:
			errs = append(errs, newValidationError(joinLoc(location, "type"), "must be one of [%s, %s, %s, %s, %s], but got '%s'", TypeApiKey, TypeHTTP, TypeMutualTLS, TypeOAuth2, TypeOpenIDConnect, o.Type))
		}
	}
	return errs
}
