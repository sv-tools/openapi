package openapi

type Webhooks = map[string]*RefOrSpec[Extendable[PathItem]]

func NewWebhooks() Webhooks {
	return make(Webhooks)
}
