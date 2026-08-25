package codec

import (
	"encoding/csv"
	"io"
	"strconv"
)

type Row struct {
	ID    string
	Value float64
}

func WriteRows(w io.Writer, rows []Row) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"id", "value"}); err != nil {
		return err
	}
	for _, row := range rows {
		if err := cw.Write([]string{row.ID, strconv.FormatFloat(row.Value, 'f', 3, 64)}); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
