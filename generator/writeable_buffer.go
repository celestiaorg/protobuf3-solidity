package generator

import (
	"bytes"
	"fmt"
)

// WriteableBuffer is a writeable buffer wrapper
type WriteableBuffer struct {
	buf    bytes.Buffer
	indent string
}

// Indent increases indentation for subsequent writes
func (b *WriteableBuffer) Indent() {
	b.indent += "    "
}

// Unindent decreases indentation for subsequent writes
func (b *WriteableBuffer) Unindent() {
	if len(b.indent) > 0 {
		b.indent = b.indent[4:]
	}
}

// P writes an object to the buffer
func (b *WriteableBuffer) P(s ...interface{}) {
	if s != nil {
		b.buf.WriteString(b.indent)
		for _, v := range s {
			b.printAtom(v)
		}
	}
	b.buf.WriteByte('\n')
}

func (b *WriteableBuffer) printAtom(v interface{}) {
	switch v := v.(type) {
	case string:
		b.buf.WriteString(v)
	case bool, int, int32, int64, uint, uint32, uint64:
		b.buf.WriteString(fmt.Sprint(v))
	default:
		panic(fmt.Sprintf("unknown type in printer: %T", v))
	}
}

func (b *WriteableBuffer) String() string {
	return b.buf.String()
}
