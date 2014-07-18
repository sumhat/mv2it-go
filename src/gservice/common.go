package gservice

var (
	clientKeys = [...]string{"201407168391"}
)

func IsValidClientKey(ck string) bool {
	for _, v := range clientKeys {
		if v == ck {
			return true
		}
	}
	return false
}