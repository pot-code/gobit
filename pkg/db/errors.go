package db

import (
	"go.uber.org/zap/zapcore"
)

type SqlDBError struct {
	Sql  string
	Type string
	Args []interface{}
	Err  error
}

func (se SqlDBError) Error() string {
	if se.Err != nil {
		return se.Err.Error()
	}
	return "sql db error"
}

func (se SqlDBError) Unwrap() error {
	return se.Err
}

func (se SqlDBError) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("sql", se.Sql)
	return enc.AddReflected("args", se.Args)
}
