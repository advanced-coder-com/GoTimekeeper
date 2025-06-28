package logs

import (
	"bytes"
	"encoding/json"
	"strings"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type prettyEncoder struct {
	zapcore.Encoder
}

func (e *prettyEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	// Распарсим JSON в map
	var raw map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		return buf, nil // fallback если невалидный JSON
	}

	if val, ok := raw["stacktrace"].(string); ok {
		raw["stacktrace"] = strings.Split(val, "\n")
	}

	var indented bytes.Buffer
	enc := json.NewEncoder(&indented)
	enc.SetIndent("", "  ")
	if err := enc.Encode(raw); err != nil {
		return buf, nil
	}

	final := buffer.NewPool().Get()
	final.Write(indented.Bytes())
	return final, nil
}

func WrapEncoderAsPretty(enc zapcore.Encoder) zapcore.Encoder {
	return &prettyEncoder{enc}
}
