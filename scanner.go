package log2csv

import (
	"bufio"
	"io"
	"time"
)

type Scanner struct {
	sc      *bufio.Scanner
	formats []*Format
	err     error

	lastMatched *Format
}

func (s *Scanner) Scan() *Log {
	if s.err != nil {
		return nil
	}

	var log *Log
	for s.sc.Scan() {
		now := time.Now()
		line := s.sc.Text()

		if format, fields := s.match(line); format != nil {
			log = new(Log)
			log.Timestamp = now
			log.Format = format
			log.Fields = fields
			break
		}
		// invalid data
	}
	if s.sc.Err() != nil {
		s.err = s.sc.Err()
	}

	return log
}

func (s *Scanner) match(line string) (*Format, []string) {
	if s.lastMatched != nil {
		if fields := s.lastMatched.Pattern.FindStringSubmatch(line); fields != nil {
			return s.lastMatched, fields[1:]
		}
	}

	for _, f := range s.formats {
		if fields := f.Pattern.FindStringSubmatch(line); fields != nil {
			s.lastMatched = f
			return f, fields[1:]
		}
	}

	return nil, nil
}

func (s *Scanner) Err() error {
	return s.err
}

func NewScanner(r io.Reader, formats []*Format) *Scanner {
	s := new(Scanner)
	s.sc = bufio.NewScanner(r)
	s.formats = formats

	return s
}
