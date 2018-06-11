package factory

import (
	"github.com/dpb587/metalink/cli/verification"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

var DynamicVerification = verification.NewDynamicVerifierImpl(boshsys.NewOsFileSystem(boshlog.NewLogger(boshlog.LevelError)))
