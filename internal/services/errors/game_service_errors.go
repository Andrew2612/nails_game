package errors

type NotFoundError struct {
	error
}

type UnauthorizedError struct {
	error
}

type InvalidOperationError struct {
	error
}
