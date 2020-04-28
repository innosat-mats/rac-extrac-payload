package exports

import (
	"fmt"
	"io"
)

// StdoutCallbackFactory returns a callback that writes to stdout
func StdoutCallbackFactory(
	out io.Writer,
	writeTimeseries bool,
) Callback {

	return func(pkg ExportablePackage) {
		if writeTimeseries {
			fmt.Fprintf(out, "%+v\n", pkg)
		}
	}
}
