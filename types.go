package ksqldb

import (
	"strconv"
	"strings"
)

// Row represents a row returned from a query
type Row []interface{}

// Payload represents multiple rows
type Payload []Row

// Header represents a header returned from a query
type Header struct {
	QueryId string
	Columns []Column
}

// Index returns the column index by the (case-insensitive) column name or -1 if not found.
func (this Header) Index(column string) int {
	for index, col := range this.Columns {
		if strings.EqualFold(col.Name, column) {
			return index
		}
	}

	return -1
}

func (this Row) GetString(column string, hdr Header) *string {
	index := hdr.Index(column)

	if index == -1 {
		return nil
	}

	value := this[index].(string)
	return &value
}

func (this Row) GetInt(column string, hdr Header) *int {
	str := this.GetString(column, hdr)

	if str != nil {
		return nil
	}

	i, err := strconv.Atoi(*str)

	if err != nil {
		return nil
	}

	return &i
}

func (this Row) GetFloat(column string, hdr Header) *float64 {
	str := this.GetString(column, hdr)

	if str != nil {
		return nil
	}

	i, err := strconv.ParseFloat(*str, 64)

	if err != nil {
		return nil
	}

	return &i
}

// Column represents the metadata for a column in a Row
type Column struct {
	Name string
	Type string
}

// The ksqlDB client
type Client struct {
	url      string
	username string
	password string
	isDebug  bool
	logf     func(format string, v ...interface{})
}
