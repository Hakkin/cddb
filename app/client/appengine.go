// +build appengine

package appenginemain

import (
	"net/http"

	"github.com/Hakkin/cddb/app/handler"
)

func init() {
	http.HandleFunc("/cddb", handler.CDDB)
	http.HandleFunc("/cddb/", handler.CDDB)
}
