package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/mdp/qrterminal/v3"
)

func PrintQR(uri string, dst io.Writer) {
	if dst == nil {
		dst = os.Stdout
	}

	_, _ = fmt.Fprintln(dst)

	qrterminal.GenerateHalfBlock(uri, qrterminal.L, dst)

	_, _ = fmt.Fprintln(dst)
	_, _ = fmt.Fprintln(dst, uri)
	_, _ = fmt.Fprintln(dst)
}
