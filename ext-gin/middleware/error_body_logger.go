package middleware

import (
	"net/http"

	"github.com/Shyp/go-dberror"
	"github.com/gin-gonic/gin"
	"github.com/karsto/glew/common/types"
	uuid "github.com/satori/go.uuid"
)

// TODO: update incidentID to use an id that is referenced when logs are written
func ErrorBodyLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		// look for bind errors and if find them return with a 400 status code
		bindErrors := ctx.Errors.ByType(gin.ErrorTypeBind)
		if len(bindErrors) > 0 {
			errors := []string{}
			for _, v := range bindErrors {
				errors = append(errors, v.Error())
			}
			ctx.JSON(http.StatusBadRequest, types.WebError{
				Status:     http.StatusBadRequest,
				IncidentID: uuid.Must(uuid.NewV4(), nil).String(),
				Message:    "Bad Request: failed to bind due to validation or type issue. See errors for more detail.",
				MoreInfo:   "TODO: link to api swagger documentation goes here",
				Errors:     errors,
			})
			return
		}

		// attempt to parse db error codes out of errors
		internalErrors := []error{}
		badRequestErrors := []dberror.Error{}
		for _, v := range ctx.Errors {
			parsedError := dberror.GetError(v)
			switch e := parsedError.(type) {
			case *dberror.Error:
				if badRequstDbCodes[e.Code] {
					badRequestErrors = append(badRequestErrors, *e)
				} else {
					internalErrors = append(internalErrors, parsedError)
				}
			default:
				internalErrors = append(internalErrors, parsedError)
			}
		}

		// did we have any errors that should be a 400 status code
		if len(badRequestErrors) > 0 {
			ctx.JSON(http.StatusBadRequest, types.WebError{
				Status:     http.StatusBadRequest,
				IncidentID: uuid.Must(uuid.NewV4(), nil).String(),
				Message:    "Bad Request: failed due to issue with input. See errors for more detail.",
				MoreInfo:   "TODO: link to api swagger documentation goes here",
				Errors:     toStringSlice(badRequestErrors),
			})
			return
		}

		// default hide exception details, provide reporting concept
		ctx.JSON(http.StatusInternalServerError, types.WebError{
			Status:     http.StatusInternalServerError,
			IncidentID: uuid.Must(uuid.NewV4(), nil).String(),
			Message:    "Internal Server Error. Please create an issue and refernece this incident id.",
			MoreInfo:   "TODO: link to create a new issue goes here",
		})
	}
}

var badRequstDbCodes = map[string]bool{
	dberror.CodeForeignKeyViolation: true,
	dberror.CodeUniqueViolation:     true,
	dberror.CodeCheckViolation:      true,
}

func toStringSlice(errs []dberror.Error) []string {
	res := make([]string, 0, len(errs))
	for _, v := range errs {
		res = append(res, v.Error())
	}
	return res
}
