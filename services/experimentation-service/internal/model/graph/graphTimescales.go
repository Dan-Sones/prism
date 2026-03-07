package graph

type GraphTimeScale string

const (
	Minute     GraphTimeScale = "minute"
	TenMinute  GraphTimeScale = "ten_minute"
	HalfHour   GraphTimeScale = "half_hour"
	ScaleHour  GraphTimeScale = "hour"
	ScaleDay   GraphTimeScale = "day"
	ScaleWeek  GraphTimeScale = "week"
	ScaleMonth GraphTimeScale = "month"
)

func (g GraphTimeScale) String() string {
	return string(g)
}

func GetListOfGraphTimeScales() string {
	return string(Minute) + ", " + string(TenMinute) + ", " + string(HalfHour) + ", " + string(ScaleHour) + ", " + string(ScaleDay) + ", " + string(ScaleWeek) + ", " + string(ScaleMonth)
}

func ParseGraphTimeScale(scale string) (GraphTimeScale, bool) {
	switch GraphTimeScale(scale) {
	case Minute, TenMinute, HalfHour, ScaleHour, ScaleDay, ScaleWeek, ScaleMonth:
		return GraphTimeScale(scale), true
	default:
		return "", false
	}
}
