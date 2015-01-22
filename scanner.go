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
	sc := s.sc
	for sc.Scan() {
		now := time.Now()
		line := sc.Text()

		if f := s.match(line); f != nil {
			log := new(Log)
			log.Timestamp = now
			log.Format = f
			log.Fields = f.Pattern.FindStringSubmatch(line)[1:]
			return log
		}
	}

	if sc.Err() != nil {
		s.err = sc.Err()
	}

	return nil
}

func (s *Scanner) match(line string) *Format {
	if s.lastMatched != nil && s.lastMatched.Pattern.MatchString(line) {
		return s.lastMatched
	}

	for _, f := range s.formats {
		if f.Pattern.MatchString(line) {
			s.lastMatched = f
			return f
		}
	}

	return nil
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
