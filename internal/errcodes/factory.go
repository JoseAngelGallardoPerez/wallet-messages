package errcodes

import (
	"github.com/Confialink/wallet-pkg-errors"
	"github.com/gin-gonic/gin"
)

func CreatePublicError(code string, title ...string) *errors.PublicError {
	err := &errors.PublicError{Code: code}
	if errDetails, ok := details[code]; ok {
		err.Details = errDetails
	}

	if len(title) > 0 {
		err.Title = title[0]
	}

	err.HttpStatus = HttpStatusCodeByErrCode(code)

	return err
}

func AddError(c *gin.Context, code string) {
	publicErr := &errors.PublicError{
		Code:       code,
		HttpStatus: HttpStatusCodeByErrCode(code),
	}
	if errDetails, ok := details[code]; ok {
		publicErr.Details = errDetails
	}

	errors.AddErrors(c, publicErr)
}
