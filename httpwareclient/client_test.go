package httpwareclient

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
)

func helperHTTPClient(in RequestFunc, fixture []byte) {
	httpClient = &HttpClientMock{
		Func: func(req *http.Request) (*http.Response, error) {
			if in != nil {
				in(req)
			}

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBuffer(fixture)),
			}, nil
		},
	}
}

func TestSend_ReturnsValidBody(t *testing.T) {
	senddata := &sendData{}
	recvdata := &recvData{}
	httpClient = &HttpClientMock{
		Func: func(*http.Request) (*http.Response, error) {
			buffer := bytes.NewBuffer([]byte(`{"value":"test"}`))

			return &http.Response{
				Body: ioutil.NopCloser(buffer),
			}, nil
		},
	}

	in := &SendIn{
		Method:   http.MethodPost,
		URL:      "",
		Coder:    GetCoder(JSON),
		BodySend: senddata,
		BodyRecv: recvdata,
	}

	if err := Send(context.Background(), in); err != nil {
		t.Error(err)
	}

	if in.BodyRecv.(*recvData).Value != "test" {
		t.Error("values not equals")
	}
}

func TestSend_WithHeaders_ReturnsValidHeader(t *testing.T) {
	type (
		sendData struct {
			Field string
		}
		recvData struct {
			Value string
		}
	)

	senddata := &sendData{}
	recvdata := &recvData{}
	httpClient = &HttpClientMock{
		Func: func(req *http.Request) (*http.Response, error) {
			value := req.Header.Get("header")
			if value != "myheader" {
				t.Error("header is not equals")
			}

			buffer := bytes.NewBuffer([]byte(`{"value":"test"}`))
			body := ioutil.NopCloser(buffer)

			return &http.Response{
				Body: body,
			}, nil
		},
	}

	in := &SendIn{
		Method:   http.MethodPost,
		URL:      "",
		Coder:    GetCoder(JSON),
		BodySend: senddata,
		BodyRecv: recvdata,
		Headers: map[string]string{
			"header": "myheader",
		},
	}

	if err := Send(context.Background(), in); err != nil {
		t.Error(err)
	}
}

type (
	sendData struct {
		Field string
	}
	recvData struct {
		Value string
	}
)

func TestSend(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *SendIn
	}
	type hc struct {
		in      RequestFunc
		fixture []byte
	}
	tests := []struct {
		name    string
		args    args
		hc      hc
		wantErr bool
	}{
		{
			name: "empty SendIn returns no error",
			args: args{
				ctx: context.Background(),
				in:  &SendIn{},
			},
			wantErr: false,
		},
		{
			name: "BodyRecv not set returns no error",
			args: args{
				ctx: context.Background(),
				in: &SendIn{
					Method: http.MethodPost,
					URL:    "",
					Coder:  GetCoder(JSON),
					BodySend: &sendData{
						Field: "value",
					},
				},
			},
			hc: hc{
				in:      func(*http.Request) {},
				fixture: []byte{},
			},
			wantErr: false,
		},
		{
			name: "BodySend not set returns no error",
			args: args{
				ctx: context.Background(),
				in: &SendIn{
					Method:   http.MethodPost,
					URL:      "",
					Coder:    GetCoder(JSON),
					BodyRecv: &recvData{},
				},
			},
			hc: hc{
				in:      func(*http.Request) {},
				fixture: []byte(`{"value":"test"}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helperHTTPClient(tt.hc.in, tt.hc.fixture)
			if err := Send(tt.args.ctx, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
