package openapi

// OAuthFlow configuration details for a supported OAuth Flow
//
// https://spec.openapis.org/oas/v3.1.1#oauth-flow-object
//
// Example:
//
//	implicit:
//	  authorizationUrl: https://example.com/api/oauth/dialog
//	  scopes:
//	    write:pets: modify pets in your account
//	    read:pets: read your pets
//	authorizationCode
//	  authorizationUrl: https://example.com/api/oauth/dialog
//	  scopes:
//	    write:pets: modify pets in your account
//	    read:pets: read your pets
type OAuthFlow struct {
	// REQUIRED.
	// The available scopes for the OAuth2 security scheme.
	// A map between the scope name and a short description for it.
	// The map MAY be empty.
	//
	// Applies To: oauth2
	Scopes map[string]string `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	// REQUIRED.
	// The authorization URL to be used for this flow.
	// This MUST be in the form of a URL.
	// The OAuth2 standard requires the use of TLS.
	//
	// Applies To:oauth2 ("implicit", "authorizationCode")
	AuthorizationURL string `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	// REQUIRED.
	// The token URL to be used for this flow.
	// This MUST be in the form of a URL.
	// The OAuth2 standard requires the use of TLS.
	//
	// Applies To: oauth2 ("password", "clientCredentials", "authorizationCode")
	TokenURL string `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	// The URL to be used for obtaining refresh tokens.
	// This MUST be in the form of a URL.
	// The OAuth2 standard requires the use of TLS.
	//
	// Applies To: oauth2
	RefreshURL string `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
}

func (o *OAuthFlow) validateSpec(path string, opts *specValidationOptions) []*validationError {
	// all the validations are done in the parent object
	return nil
}

type OAuthFlowBuilder struct {
	spec *Extendable[OAuthFlow]
}

func NewOAuthFlowBuilder() *OAuthFlowBuilder {
	return &OAuthFlowBuilder{
		spec: NewExtendable[OAuthFlow](&OAuthFlow{}),
	}
}

func (b *OAuthFlowBuilder) Build() *Extendable[OAuthFlow] {
	return b.spec
}

func (b *OAuthFlowBuilder) Extensions(v map[string]any) *OAuthFlowBuilder {
	b.spec.Extensions = v
	return b
}

func (b *OAuthFlowBuilder) AddExt(name string, value any) *OAuthFlowBuilder {
	b.spec.AddExt(name, value)
	return b
}

func (b *OAuthFlowBuilder) Scopes(v map[string]string) *OAuthFlowBuilder {
	b.spec.Spec.Scopes = v
	return b
}

func (b *OAuthFlowBuilder) AddScope(name, value string) *OAuthFlowBuilder {
	if b.spec.Spec.Scopes == nil {
		b.spec.Spec.Scopes = make(map[string]string, 1)
	}
	b.spec.Spec.Scopes[name] = value
	return b
}

func (b *OAuthFlowBuilder) AuthorizationURL(v string) *OAuthFlowBuilder {
	b.spec.Spec.AuthorizationURL = v
	return b
}

func (b *OAuthFlowBuilder) TokenURL(v string) *OAuthFlowBuilder {
	b.spec.Spec.TokenURL = v
	return b
}

func (b *OAuthFlowBuilder) RefreshURL(v string) *OAuthFlowBuilder {
	b.spec.Spec.RefreshURL = v
	return b
}
