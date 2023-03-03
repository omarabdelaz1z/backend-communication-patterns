package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"
)

const (
	measures = "kMGTPE"
	unit     = 1024
)

type Event struct {
	id,
	event,
	data string
	retry int
}

func NewEvent(id, event, data string) *Event {
	return &Event{
		id:    id,
		event: event,
		data:  data,
		retry: 0,
	}
}

func (m *Event) Buffer() *bytes.Buffer {
	var buffer bytes.Buffer

	if len(m.id) > 0 {
		buffer.WriteString(fmt.Sprintf("id: %s\n", m.id))
	}

	if m.retry > 0 {
		buffer.WriteString(fmt.Sprintf("retry: %d\n", m.retry))
	}

	if len(m.event) > 0 {
		buffer.WriteString(fmt.Sprintf("event: %s\n", m.event))
	}

	if len(m.data) > 0 {
		for _, line := range strings.Split(m.data, "\n") {
			buffer.WriteString(fmt.Sprintf("data: %s\n", line))
		}
	}

	buffer.WriteString("\n")

	return &buffer
}

func FormatBytes(bytes uint64) string {
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	exp := int(math.Log(float64(bytes)) / math.Log(unit))

	return fmt.Sprintf("%.1f %cB", float64(bytes)/math.Pow(float64(unit), float64(exp)), measures[int(exp)-1])
}
