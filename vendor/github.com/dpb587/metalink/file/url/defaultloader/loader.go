package defaultloader

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/metalink/file/url"
	fileurl "github.com/dpb587/metalink/file/url/file"
	ftpurl "github.com/dpb587/metalink/file/url/ftp"
	httpurl "github.com/dpb587/metalink/file/url/http"
	s3url "github.com/dpb587/metalink/file/url/s3"
)

func New(fs boshsys.FileSystem) url.Loader {
	file := fileurl.NewLoader(fs)

	loader := url.NewLoaderFactory()
	loader.Add(file)
	loader.Add(ftpurl.Loader{})
	loader.Add(httpurl.Loader{})
	loader.Add(s3url.Loader{})
	loader.Add(fileurl.NewEmptyLoader(file))

	return loader
}
