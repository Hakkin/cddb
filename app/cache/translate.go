package cache

import (
	"fmt"
	"strconv"
	"strings"
)

func translate(id string) (string, error) {
	idNumber := strings.Split(id, "-")[0]
	idInt32, err := strconv.ParseUint(idNumber, 10, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%08x", idInt32), nil
}
