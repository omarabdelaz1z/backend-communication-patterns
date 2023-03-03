package main

import (
	"context"
	"sync"
	"time"
)

type Service struct {
	changeTicker *time.Ticker
	buffer       chan *NetStat
	events       chan []byte
}

func (s *Service) Capture(ctx context.Context) error {
	var stat *NetStat

	for {
		select {
		case <-ctx.Done():
			s.changeTicker.Stop()
			close(s.buffer)
			return nil
		case <-s.changeTicker.C:
			if nil == stat {
				stat, _ = BriefStat()
			}

			var next *NetStat
			next, _ = BriefStat()

			delta := Diff(next, stat)

			s.buffer <- delta
			stat = next
		}
	}
}

func (s *Service) Transform(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case stat := <-s.buffer:
				s.events <- stat.Bytes()
			}
		}
	}()

	wg.Wait()

	return nil
}
