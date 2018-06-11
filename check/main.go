package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/dpb587/dynamic-metalink-resource/api"
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/repository"
	filter_and "github.com/dpb587/metalink/repository/filter/and"
	"github.com/dpb587/metalink/repository/sorter"
	sorter_fileversion "github.com/dpb587/metalink/repository/sorter/fileversion"
	"github.com/dpb587/metalink/repository/source"
)

func main() {
	var request Request

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		api.Fatal("check: bad stdin: parse error", err)
	}

	andFilter := filter_and.NewFilter()

	err = request.ApplyFilter(&andFilter)
	if err != nil {
		api.Fatal("check: bad stdin: filter error", err)
	}

	stdout, err := api.ExecuteScript(request.Source.VersionCheck, nil)
	if err != nil {
		api.Fatal("check: version check script", err)
	}

	var metalinks []repository.RepositoryMetalink

	for _, version := range strings.Split(strings.TrimSpace(string(stdout)), "\n") {
		if version == "" {
			continue
		}

		metalinks = append(
			metalinks,
			repository.RepositoryMetalink{
				Metalink: metalink.Metalink{
					Files: []metalink.File{
						{
							Version: version,
						},
					},
				},
			},
		)
	}

	metalinks, err = source.FilterInMemory(metalinks, andFilter)
	if err != nil {
		api.Fatal("check: bad filter", err)
	}

	sorter.Sort(metalinks, sorter_fileversion.Sorter{})

	response := Response{}

	for _, meta4 := range metalinks {
		response = append(
			response,
			api.Version{
				Version: meta4.Metalink.Files[0].Version,
			},
		)

		if request.Version == nil {
			break
		}
	}

	err = json.NewEncoder(os.Stdout).Encode(response)
	if err != nil {
		api.Fatal("check: bad stdout: json", err)
	}
}
