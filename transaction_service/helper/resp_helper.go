package helper

type respHelper[T interface{}] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func RespHelper[T interface{}](msg string, d T) respHelper[T] {
	return respHelper[T]{Message: msg, Data: d}
}
