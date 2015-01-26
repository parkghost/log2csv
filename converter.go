package log2csv

type Converter struct {
	sc *Scanner
	w  Writer
}

func (c *Converter) Convert() error {
	for {
		log := c.sc.Scan()
		if log == nil {
			break
		}

		if err := c.w.Write(log); err != nil {
			return err
		}
	}

	if c.sc.Err() != nil {
		return c.sc.Err()
	}
	if f, ok := c.w.(Flusher); ok {
		return f.Flush()
	}
	return nil
}

func NewConverter(sc *Scanner, w Writer) *Converter {
	c := new(Converter)
	c.sc = sc
	c.w = w

	return c
}
