package middleware

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// source originally from https://stackoverflow.com/questions/38501325/how-to-log-response-body-in-gin
func ResponseLogger(onlyErrors bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		statusCode := c.Writer.Status()

		if !onlyErrors || statusCode >= 400 && statusCode != 404 {
			//ok this is an request with error, let's make a record for it
			// now print body (or log in your preferred way)
			fmt.Println("Response body: " + blw.body.String())
		}
	}
}
