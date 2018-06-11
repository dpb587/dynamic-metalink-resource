package factory

import (
	"github.com/dpb587/metalink/repository/source"
	source_factory "github.com/dpb587/metalink/repository/source/factory"
	source_fs "github.com/dpb587/metalink/repository/source/fs"
	source_git "github.com/dpb587/metalink/repository/source/git"
	source_http "github.com/dpb587/metalink/repository/source/http"
	source_s3 "github.com/dpb587/metalink/repository/source/s3"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

func getSourceFactory() source.Factory {
	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)
	cmdRunner := boshsys.NewExecCmdRunner(logger)

	sourceFactory := source_factory.NewFactory()
	sourceFactory.Add(source_fs.NewFactory(fs))
	sourceFactory.Add(source_http.NewFactory())
	sourceFactory.Add(source_git.NewFactory(fs, cmdRunner))
	sourceFactory.Add(source_s3.NewFactory())

	return sourceFactory
}

func GetSource(uri string, options map[string]interface{}) (source.Source, error) {
	return getSourceFactory().Create(uri, options)
}
