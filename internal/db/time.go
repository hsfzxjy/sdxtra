package db

// type ISO8601TimeString string

// const _ISO8601_FORMAT = "2006-01-02T15:04:05.000Z"

// func EncodeTime(t time.Time) ISO8601TimeString {
// 	return ISO8601TimeString(t.UTC().Format(_ISO8601_FORMAT))
// }

// func EncodeTimeOrNow(t sql.Null[ISO8601TimeString]) ISO8601TimeString {
// 	if t.Valid {
// 		return t.V
// 	}
// 	return EncodeTime(time.Now())
// }
