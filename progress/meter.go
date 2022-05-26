package progress

// "github.com/vbauerster/mpb/v7"
// "github.com/vbauerster/mpb/v7/decor"

import (
	"context"
	"log"
	"sync"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

type Meter interface {
	AddBar(barId string, title string, total int64)
	IncrBar(barId string, amount int)
	Wait()
	SetTotal(barId string, to int64, done bool)
	SetBar(barId string, to int64)
	AbortBar(barId string)

	AddUnboundedBar(barId string, title string, progressFn ProgressFunc)
	IncrUnboundedProgress(barId string, count int64)
}

var WrapCount int64 = 1000000

type ProgressFunc func() string

type MbpMeter struct {
	progress             *mpb.Progress
	bars                 map[string]*mpb.Bar
	getOverallProgressFn ProgressFunc
	mutex                sync.RWMutex
}

func Setup(ctx context.Context) *MbpMeter {
	m := MbpMeter{
		progress: mpb.NewWithContext(ctx, mpb.WithWidth(70)),
		bars:     map[string]*mpb.Bar{},
	}
	return &m
}

func (m *MbpMeter) getOverallProgress(stats decor.Statistics) string {
	if m.getOverallProgressFn != nil {
		progress := m.getOverallProgressFn()
		return progress
	}
	return ""
}

func (m *MbpMeter) AddUnboundedBar(barId string, title string, progressFn ProgressFunc) {
	m.getOverallProgressFn = progressFn
	b := m.progress.New(
		0,
		mpb.BarStyle().Lbound("[").Filler("-").Tip("|").Padding("-").Rbound("]"),
		mpb.PrependDecorators(
			decor.Name(title, decor.WCSyncSpaceR),
			decor.Any(m.getOverallProgress),
		),
		mpb.AppendDecorators(
			decor.Elapsed(decor.ET_STYLE_GO),
		),
	)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.bars[barId] = b
}

func (m *MbpMeter) AddBar(barId string, title string, totalExpected int64) {
	b := m.progress.New(
		totalExpected,
		mpb.BarStyle().Lbound("[").Filler("=").Tip("|").Padding("-").Rbound("]"),
		mpb.PrependDecorators(
			decor.Name(title, decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d ", decor.WCSyncWidth),
			decor.Percentage(),
			// decor.Elapsed(decor.ET_STYLE_GO),
		),
		mpb.AppendDecorators(
			decor.Elapsed(decor.ET_STYLE_GO),
		),
	)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.bars[barId] = b
}

func (m *MbpMeter) Wait() {
	m.progress.Wait()
}

func (m *MbpMeter) getBar(barId string) *mpb.Bar {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.bars[barId]
}

func (m *MbpMeter) AbortBar(barId string) {
	bar := m.getBar(barId)
	if bar != nil {
		log.Println("bar aborted")
		bar.Abort(true)
	}
}

func (m *MbpMeter) IncrBar(barId string, amount int) {
	m.getBar(barId).IncrBy(amount)
}

func (m *MbpMeter) SetBar(barId string, to int64) {
	m.getBar(barId).SetCurrent(to)
}

func (m *MbpMeter) SetTotal(barId string, to int64, done bool) {
	m.getBar(barId).SetTotal(to, done)
}

func (m *MbpMeter) IncrUnboundedProgress(barId string, currentCount int64) {
	if m != nil {
		if currentCount == 1 {
			m.SetTotal(barId, WrapCount, false)
		}

		if currentCount%WrapCount == 0 {
			m.SetBar(barId, 1)
		} else {
			m.IncrBar(barId, 1)
		}
	}
}
