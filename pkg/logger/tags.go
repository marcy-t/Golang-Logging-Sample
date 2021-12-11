package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Tag struct {
	Field string
	Value interface{}
}

func (t *Tag) Log() zapcore.Field {
	d := t
	d.Prepare()
	return zap.Any(d.Field, d.Value)
}

func (t *Tag) Prepare() {
	switch d := t.Value.(type) {
	case []*Tag:
		var data []interface{}
		for _, v := range d {
			v.Prepare()
			data = append(data, map[string]interface{}{v.Field: v.Value})
		}
		log.Println(data)
		t.Value = data
	case *Tag:
		d.Prepare()
		t.Value = map[string]interface{}{d.Field: d.Value}
	default:
		break
	}
	return
}

func NewTag(f string, v interface{}) (tag *Tag) {
	tag = &Tag{
		Field: f,
		Value: v,
	}
	return
}
