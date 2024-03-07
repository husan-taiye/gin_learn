package logger

func Error(err error) Field {
	return Field{
		Key:   "error",
		Value: err,
	}
}

func String(key, val string) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}
