package helper

const (
	MetricInvalidParams = "invalid_params"
	MetricNoHealth = "no_health"
	MetricQueryError = "query_error"
	MetricMarshalReqError = "marshal_request_body_error"
	MetricInitReqError = "init_request_body_error"
	MetricPostReqError = "post_request_error"
	MetricGetReqError = "get_request_error"
	MetricNon200Error = "response_non200_error"
	MetricReadBodyError = "read_response_body_error"
	MetricUnmarshalBodyError = "unmarshal_response_body_error"
)

const (
	MetricSvcNameOrder = "order"
	MetricSvcNameCustomer = "customer"
	MetricSvcNameItem = "item"
)

const (
	MetricDBErrNoRow = "no_row"
	MetricDBErrInternal = "internal"
)

type SuccessResponse struct {
	Payload any `json:"payload"`
}

type ErrorResponse struct {
	ErrCode string `json:"error_code"`
	ErrMsg string `json:"error_message"`
}

var DataNotFoundResponse = ErrorResponse{
	ErrCode: "404",
	ErrMsg: "NOT_FOUND",
}

var BadRequestResponse = ErrorResponse{
	ErrCode: "400",
	ErrMsg: "BAD_REQUEST",
}

var InternalErrorResponse = ErrorResponse{
	ErrCode: "500",
	ErrMsg: "INTERNAL_ERROR",
}
