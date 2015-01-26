package log2csv

import (
	"regexp"
	"time"
)

type Log struct {
	Timestamp time.Time
	Format    *Format
	Fields    []string
}

type Format struct {
	Name    string
	Header  string
	Pattern *regexp.Regexp
}

var GCTraceFormats = []*Format{
	{
		"Go 1.0",
		"numgc,nproc,mark,sweep,cleanup,heap0,heap1,obj0,obj1,nmalloc,nfree,nhandoff",
		regexp.MustCompile(`gc(\d+)\((\d+)\): (\d+)\+(\d+)\+(\d+) \w+ (\d+) -> (\d+) \w+ (\d+) -> (\d+) \((\d+)-(\d+)\) objects (\d+) handoff`),
	},
	{
		"Go 1.1",
		"numgc,nproc,mark,sweep,cleanup,heap0,heap1,obj0,obj1,nmalloc,nfree,nhandoff,nhandoffcnt,nsteal,nstealcnt,nprocyield,nosyield,nsleep",
		regexp.MustCompile(`gc(\d+)\((\d+)\): (\d+)\+(\d+)\+(\d+) \w+, (\d+) -> (\d+) \w+ (\d+) -> (\d+) \((\d+)-(\d+)\) objects, (\d+)\((\d+)\) handoff, (\d+)\((\d+)\) steal, (\d+)\/(\d+)\/(\d+) yields`),
	},
	{
		"Go 1.3",
		"numgc,nproc,seq,sweep,mark,wait,heap0,heap1,obj,nmalloc,nfree,nspan,nbgsweep,npausesweep,nhandoff,nhandoffcnt,nsteal,nstealcnt,nprocyield,nosyield,nsleep",
		regexp.MustCompile(`gc(\d+)\((\d+)\): (\d+)\+(\d+)\+(\d+)\+(\d+) \w+, (\d+) -> (\d+) \w+, (\d+) \((\d+)-(\d+)\) objects, (\d+)\/(\d+)\/(\d+) sweeps, (\d+)\((\d+)\) handoff, (\d+)\((\d+)\) steal, (\d+)\/(\d+)\/(\d+) yields`),
	},
	{
		"Go 1.4",
		"numgc,nproc,seq,sweep,mark,wait,heap0,heap1,obj,nmalloc,nfree,goroutines,nspan,nbgsweep,npausesweep,nhandoff,nhandoffcnt,nsteal,nstealcnt,nprocyield,nosyield,nsleep",
		regexp.MustCompile(`gc(\d+)\((\d+)\): (\d+)\+(\d+)\+(\d+)\+(\d+) \w+, (\d+) -> (\d+) \w+, (\d+) \((\d+)-(\d+)\) objects, (\d+) goroutines, (\d+)\/(\d+)\/(\d+) sweeps, (\d+)\((\d+)\) handoff, (\d+)\((\d+)\) steal, (\d+)\/(\d+)\/(\d+) yields`),
	},
}
