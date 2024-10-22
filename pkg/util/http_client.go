package util

import (
	"fmt"
	"card-game-server-prototype/pkg/config"
	"github.com/samber/lo"
	"net/url"
	"time"

	"github.com/imroc/req/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const maxBodySize = 2048 // 2KB
var shouldTruncateLongBody bool

// Note: We've been seeing AWS ALB sending GOAWAY back to http client occasionally.
// https://stackoverflow.com/questions/41592929/aws-alb-http-2-and-goaway

func ProvideHttpClient(loggerFactory *LoggerFactory, logCFG *config.LogConfig, cfg *config.Config) *req.Client {
	logger := loggerFactory.Create("HttpClient")
	isDebugLevelLog := logCFG.GetLevel(string(*cfg.GameType)) == zapcore.DebugLevel
	shouldTruncateLongBody = lo.Ternary(isDebugLevelLog, false, true)
	client := req.C(). // Use C() to create a client and set with chainable client settings.
		// Timeout of all requests.
		SetTimeout(5*time.Second).
		// Enable retry and set the maximum retry count.
		SetCommonRetryCount(5).
		// Set the retry sleep interval with a commonly used algorithm: capped exponential backoff with jitter (https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/).
		SetCommonRetryBackoffInterval(100*time.Millisecond, 3*time.Second).
		// Logging
		SetLogger(logger.Sugar()).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			req := resp.Request

			// resp.Err represents the underlying error, e.g. network error, or unmarshal error (SetResult or SetError was invoked before).
			// Neither an error response nor a success response (e.g. status code < 200)
			if resp.Err == nil && resp.IsErrorState() {
				resp.Err = fmt.Errorf("resp status[%s] body[%s]", resp.Status, resp.String())
			}

			if resp.Err != nil {
				logger.Error(
					fmt.Sprintf("<= error %s %s", req.Method, req.RawURL),
					zap.Error(resp.Err),
					zap.Object("req", &formatRequest{req}),
				)
				return nil // Skip the following logic if there is an underlying error.
			}

			logger.Info(
				fmt.Sprintf("<= %s %s", req.Method, req.RawURL),
				zap.Object("resp", &formatResponse{resp}),
			)
			return nil
		}).
		OnBeforeRequest(func(client *req.Client, req *req.Request) error {
			if req.RetryAttempt > 0 {
				logger.Warn(
					fmt.Sprintf("=> retry attempt %d %s %s", req.RetryAttempt, req.Method, req.RawURL),
					zap.Object("req", &formatRequest{req}),
				)
			} else {
				logger.Info(
					fmt.Sprintf("=> %s %s", req.Method, req.RawURL),
					zap.Object("req", &formatRequest{req}),
				)
			}
			return nil
		})

	return client
}

type formatRequest struct{ *req.Request }

func (f *formatRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Method", f.Method)
	enc.AddString("URL", f.RawURL)

	if string(f.Body) != "" {
		if shouldTruncateLongBody && len(f.Body) > maxBodySize {
			var truncatedBody []byte
			truncatedBody = append(truncatedBody, f.Body[:maxBodySize]...)

			enc.AddByteString("Body", append(truncatedBody, []byte("...")...))
		} else {
			enc.AddByteString("Body", f.Body)
		}
	}

	if len(f.PathParams) > 0 {
		_ = enc.AddObject("PathParams", formatStringMap(f.PathParams))
	}

	if len(f.QueryParams) > 0 {
		_ = enc.AddObject("QueryParams", formatStringArrayMap(f.QueryParams))
	}

	if len(f.FormData) > 0 {
		_ = enc.AddObject("FormData", formatStringArrayMap(f.FormData))
	}

	// if len(req.Headers) > 0 {
	// 	fields = append(fields, zap.Object("Headers", formatStringsMap(req.Headers)))
	// }
	return nil
}

type formatResponse struct{ *req.Response }

func (f *formatResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Method", f.Request.Method)
	enc.AddString("URL", f.Request.RawURL)
	enc.AddInt("StatusCode", f.StatusCode)

	body, err := f.ToString()
	if err != nil {
		return err
	}
	if body != "" {

		if shouldTruncateLongBody && len(body) > maxBodySize {
			enc.AddString("Body", body[:maxBodySize]+"...")
		} else {
			enc.AddString("Body", body)
		}
	}

	return nil

}

type formatStringMap map[string]string

func (f formatStringMap) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range f {
		enc.AddString(k, v)
	}
	return nil
}

type formatStringArrayMap url.Values

func (f formatStringArrayMap) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range f {
		_ = enc.AddArray(k, formatStringArray(v))
	}
	return nil
}

type formatStringArray []string

func (f formatStringArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range f {
		arr.AppendString(f[i])
	}
	return nil
}
