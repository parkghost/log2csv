package log2csv

const defaultBufSize = 100

type Converter struct {
	sc *Scanner
	w  Writer

	logCh chan *Log
	errCh chan error

	quit chan struct{}
}

func (c *Converter) Convert() error {
	if c.sc.Err() != nil {
		return c.sc.Err()
	}

	go c.scanLoop()
	go c.writeLoop()

	<-c.quit
	select {
	case err := <-c.errCh:
		return err
	default:
		return nil
	}
}

func (c *Converter) scanLoop() {
loop:
	for {
		select {
		case <-c.quit:
			break loop
		default:
			log := c.sc.Scan()
			if log != nil {
				c.logCh <- log
				continue
			}

			if c.sc.Err() != nil {
				c.errCh <- c.sc.Err()
			}

			break loop
		}
	}

	close(c.logCh)
}

func (c *Converter) writeLoop() {
	for log := range c.logCh {
		if err := c.w.Write(log); err != nil {
			c.errCh <- err
			break
		}

		if f, ok := c.w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				c.errCh <- err
				break
			}
		}
	}

	close(c.quit)
}

func NewConverter(sc *Scanner, w Writer) *Converter {
	return NewConverterSize(sc, w, defaultBufSize)
}

func NewConverterSize(sc *Scanner, w Writer, bufSize int) *Converter {
	c := new(Converter)
	c.sc = sc
	c.w = w

	c.logCh = make(chan *Log, bufSize)
	c.errCh = make(chan error, 1)

	c.quit = make(chan struct{})

	return c
}
