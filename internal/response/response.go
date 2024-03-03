package response

var ErrorMsgsByCode = map[int]string{
	500: "internal server error",
	404: "item is not found",
	422: "unprocessable entity",
	400: "bad request",
}

type Response struct {
	Body interface{} `json:"body"`
}
