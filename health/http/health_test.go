package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "returns 204 status",
			args: args{
				opts: []Option{
					WithHealthCheckFuncs(func() bool { return true }),
				},
			},
			want: http.StatusNoContent,
		},
		{
			name: "returns 204 status with two healthcheckers",
			args: args{
				opts: []Option{
					WithHealthCheckFuncs(
						func() bool { return true },
						func() bool { return true },
					),
				},
			},
			want: http.StatusNoContent,
		},
		{
			name: "returns 204 status without healthcheckers",
			want: http.StatusNoContent,
		},
		{
			name: "returns 503 status",
			args: args{
				opts: []Option{
					WithHealthCheckFuncs(func() bool { return false }),
				},
			},
			want: http.StatusServiceUnavailable,
		},
		{
			name: "returns 503 status with two healthcheckers",
			args: args{
				opts: []Option{
					WithHealthCheckFuncs(
						func() bool { return false },
						func() bool { return false },
					),
				},
			},
			want: http.StatusServiceUnavailable,
		},
		{
			name: "returns 503 status with two healthcheckers and one returns false",
			args: args{
				opts: []Option{
					WithHealthCheckFuncs(
						func() bool { return true },
						func() bool { return false },
					),
				},
			},
			want: http.StatusServiceUnavailable,
		},
		{
			name: "returns 503 status global timeout",
			args: args{
				opts: []Option{
					WithTimeout(time.Millisecond),
					WithHealthCheckFuncs(
						func() bool {
							time.Sleep(time.Second)

							return true
						},
					),
				},
			},
			want: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "", nil)
			if err != nil {
				t.Errorf("Failed to create request")
			}

			HandlerFunc(tt.args.opts...)(ht, req)

			if ht.Code != tt.want {
				t.Errorf("Handler() = %d, want %d", ht.Code, tt.want)
			}
		})
	}
}
