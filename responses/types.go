package responses

type Primes struct {
	IsPrime bool   `json:"isPrime"`
	Message string `json:"message"`
}

type Request struct {
	Number int `json:"number"`
	Count  int `json:"count"`
}

type History struct {
	Requests []Request `json:"requests"`
}

type Message struct {
	LowerLimit int    `json:"lowerLimit"`
	Message    string `json:"message"`
}

type MessageSlice []Message

func (ms MessageSlice) Len() int           { return len(ms) }
func (ms MessageSlice) Less(i, j int) bool { return ms[i].LowerLimit > ms[j].LowerLimit }
func (ms MessageSlice) Swap(i, j int)      { ms[i], ms[j] = ms[j], ms[i] }

type Messages struct {
	Messages MessageSlice `json:"messages"`
}
