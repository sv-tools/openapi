package openapi

// OAuthFlows allows configuration of the supported OAuth Flows.
//
// https://spec.openapis.org/oas/v3.1.0#oauth-flows-object
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
//	  authorizationCode:
//	    authorizationUrl: https://example.com/api/oauth/dialog
//	    tokenUrl: https://example.com/api/oauth/token
//	    scopes:
//	      write:pets: modify pets in your account
//	      read:pets: read your pets
type OAuthFlows struct {
	// Configuration for the OAuth Implicit flow.
	Implicit *Extendable[OAuthFlow] `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	// Configuration for the OAuth Resource Owner Password flow.
	Password *Extendable[OAuthFlow] `json:"password,omitempty" yaml:"password,omitempty"`
	// Configuration for the OAuth Client Credentials flow.
	// Previously called application in OpenAPI 2.0.
	ClientCredentials *Extendable[OAuthFlow] `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	// Configuration for the OAuth Authorization Code flow.
	// Previously called accessCode in OpenAPI 2.0.
	AuthorizationCode *Extendable[OAuthFlow] `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

func (o *OAuthFlows) validateSpec(loc string, opts *specValidationOptions) []*validationError {
	var errs []*validationError
	if o.Implicit != nil {
		errs = append(errs, o.Implicit.validateSpec(joinLoc(loc, "implicit"), opts)...)
		if o.Implicit.Spec.AuthorizationURL == "" {
			errs = append(errs, newValidationError(joinLoc(loc, "implicit", "authorizationUrl"), ErrRequired))
		}
	}
	if o.Password != nil {
		errs = append(errs, o.Password.validateSpec(joinLoc(loc, "password"), opts)...)
		if o.Password.Spec.TokenURL == "" {
			errs = append(errs, newValidationError(joinLoc(loc, "password", "tokenUrl"), ErrRequired))
		}
	}
	if o.ClientCredentials != nil {
		errs = append(errs, o.ClientCredentials.validateSpec(joinLoc(loc, "clientCredentials"), opts)...)
		if o.ClientCredentials.Spec.TokenURL == "" {
			errs = append(errs, newValidationError(joinLoc(loc, "clientCredentials", "tokenUrl"), ErrRequired))
		}
	}
	if o.AuthorizationCode != nil {
		errs = append(errs, o.AuthorizationCode.validateSpec(joinLoc(loc, "authorizationCode"), opts)...)
		if o.AuthorizationCode.Spec.AuthorizationURL == "" {
			errs = append(errs, newValidationError(joinLoc(loc, "authorizationCode", "authorizationUrl"), ErrRequired))
		}
		if o.AuthorizationCode.Spec.TokenURL == "" {
			errs = append(errs, newValidationError(joinLoc(loc, "authorizationCode", "tokenUrl"), ErrRequired))
		}
	}

	return errs
}
