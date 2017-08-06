// +build appengine

package appenginemain

import (
	"github.com/Hakkin/cddb"
	"net/http"
)

func init() {
	http.HandleFunc("/cddb", cddb.CddbHttp)
	http.HandleFunc("/cddb/", cddb.CddbHttp)
}
