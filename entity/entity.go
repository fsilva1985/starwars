package entity

import (
	"bytes"
	"encoding/base32"

	"github.com/xtgo/uuid"
)

func GetNewId() string {
	var encoding = base32.NewEncoding("t8SaX57xiJpZ0iq3KO1EeNikfhrJvgEM")
	var b bytes.Buffer
	encoder := base32.NewEncoder(encoding, &b)
	w := uuid.NewRandom()
	encoder.Write(w.Bytes())
	encoder.Close()
	b.Truncate(26)

	return b.String()
}
