package times

import (
	"fmt"
	"testing"
	"time"
)

func TestIsZero(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "good1",
			args: args{t: time.Time{}},
			want: true,
		},
		{
			name: "good2",
			args: args{t: time.Unix(0, 0)},
			want: true,
		},
		{
			name: "bad",
			args: args{t: time.Now()},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsZero(tt.args.t); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTimeToStr(t *testing.T) {
	dataStr := parseTimeToStr(time.Now(), LayoutDateTime)
	fmt.Println(dataStr) //2023-01-08 16:56:42
}
