package k8s

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClient(t *testing.T) {
	_, err := grpc.NewClient("172.18.0.2:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
}

func Test_isEqualIPs(t *testing.T) {
	type args struct {
		s1 []string
		s2 []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{s1: []string{"1", "2"}, s2: []string{"2", "1"}},
			want: true,
		},
		{
			name: "",
			args: args{s1: []string{"1", "2"}, s2: []string{"1", "2"}},
			want: true,
		},
		{
			name: "",
			args: args{s1: []string{"1", "2"}, s2: []string{"1"}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEqualIPs(tt.args.s1, tt.args.s2); got != tt.want {
				t.Errorf("isEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
