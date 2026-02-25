package service

import "go.uber.org/fx"

type Registry struct {
	importers map[string]Importer
}

type RegistryParams struct {
	fx.In

	Importers []Importer `group:"importers"`
}

func NewRegistry(params RegistryParams) *Registry {
	imap := make(map[string]Importer, len(params.Importers))
	for _, imp := range params.Importers {
		imap[imp.SourceType()] = imp
	}
	return &Registry{importers: imap}
}

func (r *Registry) Get(sourceType string) Importer {
	if r == nil {
		return nil
	}
	return r.importers[sourceType]
}
