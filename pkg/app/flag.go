package cli

import (
	"github.com/spf13/pflag"
	"strings"
)

func normalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	// alias
	//switch name {
	//case "old-flag-name":
	//	name = "new-flag-name"
	//	break
	//}

	// --my-flag == --my_flag == --my.flag
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
	}
	return pflag.NormalizedName(name)
}
