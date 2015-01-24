package log2csv

import "io"

type Converter struct {
	r io.Reader
	w Writer
}

func (c *Converter) Convert() error {
	sc := NewScanner(c.r, formats)

	for {
		log := sc.Scan()
		if log == nil {
			break
		}

		if err := c.w.Write(log); err != nil {
			return err
		}
	}

	if sc.Err() != nil {
		return sc.Err()
	}
	if f, ok := c.w.(Flusher); ok {
		return f.Flush()
	}
	return nil
}

func NewConverter(r io.Reader, w Writer) *Converter {
	c := new(Converter)
	c.r = r
	c.w = w

	return c
}
