package queries

import "strings"

const (
	getById = `
query Q($id: string!) {
  ` + namePlaceHolder + `(func: uid($id), first: 1) {
    ` + fragmentPlaceHolder + `
  }
}
`

	fragmentPlaceHolder = "%%FRAGMENT%%"
	namePlaceHolder     = "%%NAME%%"
)

func GetById(name string, fragment string) string {
	return strings.ReplaceAll(strings.ReplaceAll(getById, namePlaceHolder, name), fragmentPlaceHolder, fragment)
}
