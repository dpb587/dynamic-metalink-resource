package factory

import (
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink/file"
	"github.com/dpb587/metalink/file/metaurl"
)

func GetMetaURLLoaderFactory() *metaurl.LoaderFactory {
	return metaurl.NewLoaderFactory()
}

func GetOriginURL(ref metalink.MetaURL) (file.Reference, error) {
	return GetMetaURLLoaderFactory().Load(ref)
}
