package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/shirou/gopsutil/v3/net"
)

const allInterface = false

type NetStat struct {
	BytesSent  uint64
	BytesRecv  uint64
	BytesTotal uint64
}

func (n *NetStat) Set(new NetStat) {
	n.BytesSent = new.BytesSent
	n.BytesRecv = new.BytesRecv
	n.BytesTotal = new.BytesTotal
}

func BriefStat() (*NetStat, error) {
	stats, err := net.IOCounters(allInterface)

	if err != nil {
		return nil, fmt.Errorf("failed to capture network stat: %v", err)
	}

	return &NetStat{
		BytesSent:  stats[0].BytesSent,
		BytesRecv:  stats[0].BytesRecv,
		BytesTotal: stats[0].BytesSent + stats[0].BytesRecv,
	}, nil
}

// Diff return the delta between two Netstat structs.
func Diff(current, previous *NetStat) *NetStat {
	return &NetStat{
		BytesSent:  current.BytesSent - previous.BytesSent,
		BytesRecv:  current.BytesRecv - previous.BytesRecv,
		BytesTotal: current.BytesTotal - previous.BytesTotal,
	}
}

// IncrBy increases "current" counters by "new" then return a new Netstat struct.
func IncrBy(current, new *NetStat) NetStat {
	return NetStat{
		BytesSent:  current.BytesSent + new.BytesSent,
		BytesRecv:  current.BytesRecv + new.BytesRecv,
		BytesTotal: current.BytesTotal + new.BytesTotal,
	}
}

func (s *NetStat) Bytes() []byte {
	var wg sync.WaitGroup

	formatted := make([]string, 3)

	for i, v := range [3]uint64{s.BytesSent, s.BytesRecv, s.BytesTotal} {
		wg.Add(1)

		go func(i int, v uint64) {
			defer wg.Done()
			formatted[i] = FormatBytes(v)
		}(i, v)
	}

	wg.Wait()

	data, _ := json.Marshal(map[string]string{
		"sent":     formatted[0],
		"received": formatted[1],
		"total":    formatted[2],
	})

	return data
}
