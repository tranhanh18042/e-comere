package helper

const (
	MetricInvalidParams = "invalid_params"
	MetricNoHealth = "no_health"
	MetricQueryError = "query_error"
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
	Payload any
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
