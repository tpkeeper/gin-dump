package gindump

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

func Dump() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		//dump req header
		var strB strings.Builder
		strB.WriteString("[GIN-dump]:\nRequest-Header:\n")
		strB.WriteString(strHeader(ctx.Request.Header))

		//dump req body
		if ctx.Request.ContentLength > 0 {
			buf, err := ioutil.ReadAll(ctx.Request.Body)
			if err != nil {
				strB.WriteString(fmt.Sprintf("\nread bodyCache err \n %s", err.Error()))
				goto DumpRes
			}
			rdr := ioutil.NopCloser(bytes.NewBuffer(buf))
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

			switch strings.Split(ctx.Request.Header.Get("Content-Type"), ";")[0] {
			case gin.MIMEJSON:
				var mapReq map[string]interface{}
				bytes, err := ioutil.ReadAll(rdr)
				if err != nil {
					strB.WriteString(fmt.Sprintf("\nread rdr err \n %s", err.Error()))
					goto DumpRes
				}
				if err := json.Unmarshal(bytes, &mapReq); err != nil {
					strB.WriteString(fmt.Sprintf("\nparse bodyCache err \n" + err.Error()))
					goto DumpRes
				}

				strB.WriteString("\nRequest-Body:\n")
				strB.WriteString(strMap(mapReq))
			case gin.MIMEPOSTForm:
			case gin.MIMEMultipartPOSTForm:
			case gin.MIMEHTML:
			default:
			}
		}

	DumpRes:
		ctx.Writer = &bodyWriter{bodyCache: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Next()

		//dump res header
		strB.WriteString("\nResponse-Header:\n")
		strB.WriteString(strHeader(ctx.Writer.Header()))

		bw, ok := ctx.Writer.(*bodyWriter)
		if !ok {
			strB.WriteString("\nbodyWriter was override , can not read bodyCache")
			goto End
		}

		//dump res body
		if bodyAllowedForStatus(ctx.Writer.Status()) && bw.bodyCache.Len() > 0 {
			switch strings.Split(ctx.Writer.Header().Get("Content-Type"), ";")[0] {
			case gin.MIMEJSON:
				var mapRes map[string]interface{}
				if err := json.Unmarshal(bw.bodyCache.Bytes(), &mapRes); err != nil {
					strB.WriteString(fmt.Sprintf("\nparse bodyCache err \n" + err.Error()))
					goto End
				}
				strB.WriteString("\nResponse-Body:\n")
				strB.WriteString(strMap(mapRes))

			case gin.MIMEPOSTForm:
			case gin.MIMEMultipartPOSTForm:
			case gin.MIMEHTML:
			default:
			}
		}
	End:
		fmt.Print(strB.String())
	}
}

func strHeader(header http.Header) string {
	var strB strings.Builder
	strB.WriteString("	{\n")
	for key, value := range header {
		strB.WriteString(fmt.Sprintf("			%s : %s\n", key, value))
	}
	strB.WriteString("	}\n")
	return strB.String()
}

func strMap(m map[string]interface{}) string {
	var strB strings.Builder
	strB.WriteString("	{\n")
	for key, value := range m {
		strB.WriteString(fmt.Sprintf("			%s : %s\n", key, value))
	}
	strB.WriteString("	}\n")
	return strB.String()
}

type bodyWriter struct {
	gin.ResponseWriter
	bodyCache *bytes.Buffer
}

//rewrite Write()
func (w bodyWriter) Write(b []byte) (int, error) {
	w.bodyCache.Write(b)
	return w.ResponseWriter.Write(b)
}

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
