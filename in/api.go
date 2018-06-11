package main

import (
	"github.com/dpb587/dynamic-metalink-resource/api"
	"github.com/dpb587/metalink/repository/filter/and"
	"github.com/dpb587/metalink/repository/filter/fileversion"
)

type Request struct {
	Source  api.Source  `json:"source"`
	Version api.Version `json:"version"`
	Params  Params      `json:"params"`
}

func (r Request) ApplyFilter(filter *and.Filter) error {
	err := r.Source.ApplyFilter(filter)
	if err != nil {
		return err
	}

	if r.Version.Version != "" {
		addFilter, err := fileversion.CreateFilter(r.Version.Version)
		if err != nil {
			return err
		}

		filter.Add(addFilter)
	}

	return nil
}

type Params struct {
	SkipDownload bool `json:"skip_download"`
}

type Response struct {
	Version  api.Version    `json:"version"`
	Metadata []api.Metadata `json:"metadata,omitempty"`
}
