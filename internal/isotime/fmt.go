package isotime

import (
	"database/sql"
	"time"
)

type String string

func Encode(t time.Time) String {
	t = t.UTC()
	buf := make([]byte, 0, len(time.RFC3339))
	buf = t.AppendFormat(buf, time.RFC3339)
	nsec := t.Nanosecond() / 1e6
	b, s := nsec/1e2, nsec%1e2
	s, g := s/1e1, s%1e1
	buf = append(buf[:len(buf)-1], '.', byte(b)+'0', byte(s)+'0', byte(g)+'0', 'Z')
	return String(buf)
}

func EncodeNow() String {
	return Encode(time.Now())
}

func OrNow(ns sql.Null[String]) String {
	if ns.Valid {
		return ns.V
	}
	return EncodeNow()
}
