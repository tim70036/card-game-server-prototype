package util

import "strings"

// JoinStrings to replace strings.Join. strings.Join has very bad performance.
// BenchmarkStringConcat-4               19          61431933 ns/op        632845167 B/op     10005 allocs/op
// BenchmarkStringSprintf-4              10         109283838 ns/op        1075688336 B/op    29688 allocs/op
// BenchmarkStringJoin-4                 15          75854431 ns/op        632844905 B/op     10003 allocs/op
// BenchmarkStringBuilder-4           15961             74010 ns/op          522224 B/op         23 allocs/op

func JoinStrings(strs []string, sep ...string) string {
	var s strings.Builder
	for i, v := range strs {
		if len(sep) > 0 && i > 0 {
			s.WriteString(sep[0])
		}
		s.WriteString(v)
	}
	return s.String()
}
