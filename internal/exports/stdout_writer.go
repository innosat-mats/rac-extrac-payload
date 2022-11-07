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
	errorStats := common.NewErrorStats()

	return func(pkg common.DataRecord) {
			errorStats.Register(pkg.Error)
			if writeTimeseries {
				if pkg.Error != nil {
					pkg.Error = fmt.Errorf(
						"%s %s",
						pkg.Error,
						common.MakePackageInfo(&pkg),
					)
				}
				fmt.Fprintf(out, "%+v\n", pkg)
			}
		}, func() {
			fmt.Fprint(out, errorStats.Summarize())
		}
}
