package model

type SecMinCounter struct {
	SecCounter TimeCounter
	MinCounter TimeCounter
}

func NewSecMinCounter() *SecMinCounter {
	return &SecMinCounter{
		SecCounter: &SecCounter{},
		MinCounter: &MinCounter{},
	}
}

type SecCounter struct {
	count   int
	curTime int64
}

func (p *SecCounter) Count(nowTime int64) (curCount int) {
	if p.curTime != nowTime {
		p.count = 1
		p.curTime = nowTime
		curCount = p.count
		return
	}
	p.count++
	curCount = p.count
	return
}

func (p *SecCounter) Check(nowTime int64) int {
	if p.curTime != nowTime {
		return 0
	}
	return p.count
}
