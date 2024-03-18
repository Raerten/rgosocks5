package slogger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestErrorf(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(nil)
	}()

	logger := Socks5Logger{}
	logger.Errorf("test %s", "message")

	want := "ERROR test message\\n"
	got := buf.String()

	if strings.HasSuffix(got, want) {
		t.Errorf("Errorf() = %q, want %q", got, want)
	}
}
