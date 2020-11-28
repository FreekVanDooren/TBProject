package responses

type Primes struct {
	IsPrime bool   `json:"is_prime"`
	Message string `json:"message"`
}

type Request struct {
	Number int `json:"number"`
	Count  int `json:"count"`
}

type History struct {
	Requests []Request `json:"requests"`
}
