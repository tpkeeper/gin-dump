package gindump

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func Dump() gin.HandlerFunc {
	return DumpWithOptions(true, true, true, true, true, nil)
}

func DumpWithOptions(showReq bool, showResp bool, showBody bool, showHeaders bool, showCookies bool, cb func(dumpStr string)) gin.HandlerFunc {
	headerHiddenFields := make([]string, 0)
	bodyHiddenFields := make([]string, 0)

	if !showCookies {
		headerHiddenFields = append(headerHiddenFields, "cookie")
	}

	return func(ctx *gin.Context) {
		var strB strings.Builder

		if showReq && showHeaders {
			//dump req header
			s, err := FormatToBeautifulJson(ctx.Request.Header, headerHiddenFields)

			if err != nil {
				strB.WriteString(fmt.Sprintf("\nparse req header err \n" + err.Error()))
			} else {
				strB.WriteString("Request-Header:\n")
				strB.WriteString(string(s))
			}
		}

		if showReq && showBody {
			//dump req body
			if ctx.Request.ContentLength > 0 {
				buf, err := ioutil.ReadAll(ctx.Request.Body)
				if err != nil {
					strB.WriteString(fmt.Sprintf("\nread bodyCache err \n %s", err.Error()))
					goto DumpRes
				}
				rdr := ioutil.NopCloser(bytes.NewBuffer(buf))
				ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
				ctGet := ctx.Request.Header.Get("Content-Type")
				ct, _, err := mime.ParseMediaType(ctGet)
				if err != nil {
					strB.WriteString(fmt.Sprintf("\ncontent_type: %s parse err \n %s", ctGet, err.Error()))
					goto DumpRes
				}

				switch ct {
				case gin.MIMEJSON:
					bts, err := ioutil.ReadAll(rdr)
					if err != nil {
						strB.WriteString(fmt.Sprintf("\nread rdr err \n %s", err.Error()))
						goto DumpRes
					}

					s, err := BeautifyJsonBytes(bts, bodyHiddenFields)
					if err != nil {
						strB.WriteString(fmt.Sprintf("\nparse req body err \n" + err.Error()))
						goto DumpRes
					}

					strB.WriteString("\nRequest-Body:\n")
					strB.WriteString(string(s))
				case gin.MIMEPOSTForm:
					bts, err := ioutil.ReadAll(rdr)
					if err != nil {
						strB.WriteString(fmt.Sprintf("\nread rdr err \n %s", err.Error()))
						goto DumpRes
					}
					val, err := url.ParseQuery(string(bts))

					s, err := FormatToBeautifulJson(val, bodyHiddenFields)
					if err != nil {
						strB.WriteString(fmt.Sprintf("\nparse req body err \n" + err.Error()))
						goto DumpRes
					}
					strB.WriteString("\nRequest-Body:\n")
					strB.WriteString(string(s))

				case gin.MIMEMultipartPOSTForm:
				default:
				}
			}

		DumpRes:
			ctx.Writer = &bodyWriter{bodyCache: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
			ctx.Next()
		}

		if showResp && showHeaders {
			//dump res header
			sHeader, err := FormatToBeautifulJson(ctx.Writer.Header(), headerHiddenFields)
			if err != nil {
				strB.WriteString(fmt.Sprintf("\nparse res header err \n" + err.Error()))
			} else {
				strB.WriteString("\nResponse-Header:\n")
				strB.WriteString(string(sHeader))
			}
		}

		if showResp && showBody {
			bw, ok := ctx.Writer.(*bodyWriter)
			if !ok {
				strB.WriteString("\nbodyWriter was override , can not read bodyCache")
				goto End
			}

			//dump res body
			if bodyAllowedForStatus(ctx.Writer.Status()) && bw.bodyCache.Len() > 0 {
				ctGet := ctx.Writer.Header().Get("Content-Type")
				ct, _, err := mime.ParseMediaType(ctGet)
				if err != nil {
					strB.WriteString(fmt.Sprintf("\ncontent-type: %s parse  err \n %s", ctGet, err.Error()))
					goto End
				}
				switch ct {
				case gin.MIMEJSON:

					s, err := BeautifyJsonBytes(bw.bodyCache.Bytes(), bodyHiddenFields)
					if err != nil {
						strB.WriteString(fmt.Sprintf("\nparse bodyCache err \n" + err.Error()))
						goto End
					}
					strB.WriteString("\nResponse-Body:\n")

					strB.WriteString(string(s))
				case gin.MIMEHTML:
				default:
				}
			}
		}

	End:
		if cb != nil {
			cb(strB.String())
		} else {
			fmt.Println(strB.String())
		}
	}
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
