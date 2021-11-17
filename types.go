package ksqldb

import (
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
	index := hdr.Index(column)

	if index == -1 {
		return nil
	}

	//value := this[index].(int) // seems to be float64 even when defined as int
	value := int(this[index].(float64))
	return &value
}

func (this Row) GetFloat(column string, hdr Header) *float64 {
	index := hdr.Index(column)

	if index == -1 {
		return nil
	}

	value := this[index].(float64)
	return &value
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
