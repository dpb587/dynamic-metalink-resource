package factory

import (
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file"
	"github.com/dpb587/metalink/file/url"
	"github.com/dpb587/metalink/file/url/defaultloader"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

func GetURLLoaderFactory() url.Loader {
	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)

	return defaultloader.New(fs)
}

func GetOrigin(ref metalink.URL) (file.Reference, error) {
	return GetURLLoaderFactory().Load(ref)
}
