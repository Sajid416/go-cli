package parser

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ---------- Errors ----------
var (
	ErrEmptyInput    = errors.New("empty input")
	ErrInvalidFormat = errors.New("invalid format")
)

// ---------- Types (slices, maps, methods) ----------
type Record map[string]string
type RecordSet []Record

func (rs RecordSet) Len() int { return len(rs) }

// ---------- Interface ----------
type Parser interface {
	Parse([]byte) (RecordSet, error)
	Name() string
}

// ---------- CSV implementation ----------
type CSVParser struct{}

// this is test branch
func (CSVParser) Name() string { return "csv" }

func (CSVParser) Parse(b []byte) (RecordSet, error) {
	if len(bytes.TrimSpace(b)) == 0 {
		return nil, ErrEmptyInput
	}
	r := csv.NewReader(strings.NewReader(string(b)))
	r.FieldsPerRecord = -1
	header, err := r.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("csv: %w", ErrEmptyInput)
		}
		return nil, fmt.Errorf("csv header: %w", err)
	}

	var out RecordSet
	for {
		row, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv row: %w", err)
		}
		rec := Record{}
		for i := range header {
			key := header[i]
			val := ""
			if i < len(row) {
				val = row[i]
			}
			rec[key] = val
		}
		out = append(out, rec)
	}
	return out, nil
}

// ---------- JSON implementation ----------
type JSONParser struct{}

func (JSONParser) Name() string { return "json" }

// accepts either an array of objects or {"data":[...]}
func (JSONParser) Parse(b []byte) (RecordSet, error) {
	if len(bytes.TrimSpace(b)) == 0 {
		return nil, ErrEmptyInput
	}

	type anyMap = map[string]any

	// Try array form first: []map[string]any
	var arr []anyMap
	if err := decodeJSON[[]anyMap](b, &arr); err == nil {
		return toRecordSet(arr), nil
	}

	// Try {"data":[...]}
	var obj struct {
		Data []anyMap `json:"data"`
	}
	if err := json.Unmarshal(b, &obj); err == nil && len(obj.Data) > 0 {
		return toRecordSet(obj.Data), nil
	}

	return nil, fmt.Errorf("json: %w", ErrInvalidFormat)
}

// ---------- Generics ----------
func decodeJSON[T any](b []byte, out *T) error {
	var t T
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	*out = t
	return nil
}

func toRecordSet(in []map[string]any) RecordSet {
	out := make(RecordSet, 0, len(in))
	for _, m := range in {
		rec := Record{}
		for k, v := range m {
			rec[k] = fmt.Sprint(v)
		}
		out = append(out, rec)
	}
	return out
}

// ---------- Factory ----------
func New(kind string) (Parser, error) {
	switch strings.ToLower(kind) {
	case "csv":
		return CSVParser{}, nil
	case "json":
		return JSONParser{}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidFormat, kind)
	}
}
