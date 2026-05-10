package source

type ErrorResolver struct {
	Err error
}

func NewErrorResolver(err error) *ErrorResolver {
	return &ErrorResolver{Err: err}
}

func (r *ErrorResolver) ListContracts(pluginVersion string) ([]Contract, error) {
	return nil, r.Err
}

func (r *ErrorResolver) GetContract(pluginVersion string, name string) (*Contract, error) {
	return nil, r.Err
}
