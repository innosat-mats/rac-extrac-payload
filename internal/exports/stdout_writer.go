package exports

import (
	"fmt"
	"io"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

// StdoutCallbackFactory returns a callback that writes to stdout
func StdoutCallbackFactory(
	out io.Writer,
	writeTimeseries bool,
) (common.Callback, common.CallbackTeardown) {

	return func(pkg common.ExportablePackage) {
		if writeTimeseries {
			fmt.Fprintf(out, "%+v\n", pkg)
		}
	}, func() {}
}
