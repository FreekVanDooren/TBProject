package responses

type Prime struct {
	IsPrime bool   `json:"is_prime"`
	Message string `json:"message"`
}
