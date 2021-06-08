package dirsync

import (
	"sync/atomic"
)

type Progress struct {
	finished int64
	total    int64

	subscribers []func(*Progress)
}

// NewProgress creates a new initialized progress
func NewProgress(total int64) *Progress {
	return &Progress{
		finished: 0,
		total:    total,

		subscribers: make([]func(*Progress), 0),
	}
}

// Percent gets a percent from 0-100 of the progress completed
func (p *Progress) Percent() float64 {
	return float64(p.finished) / float64(p.total) * 100
}

// Total gets the total items in the progress
func (p *Progress) Total() int64 {
	return p.total
}

// Finished gets the number of finished items in the progress
func (p *Progress) Finished() int64 {
	return p.finished
}

// Add increases the progress by a number
func (p *Progress) Add(num int64) {
	atomic.AddInt64(&p.finished, num)

	for _, listener := range p.subscribers {
		listener(p)
	}
}

// SetTotal sets the progress total
func (p *Progress) SetTotal(num int64) {
	atomic.StoreInt64(&p.total, num)

	for _, listener := range p.subscribers {
		listener(p)
	}
}

// AddTotal adds to the progress total
func (p *Progress) AddTotal(num int64) {
	atomic.AddInt64(&p.total, num)

	for _, listener := range p.subscribers {
		listener(p)
	}
}

// Subscribe subscribes a function to the progress
func (p *Progress) Subscribe(sub func(*Progress)) {
	p.subscribers = append(p.subscribers, sub)
}

//
