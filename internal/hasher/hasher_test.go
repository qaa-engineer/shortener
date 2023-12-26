package hasher

import (
	"testing"
)

func TestGetShortLink(t *testing.T) {
	fullURL := "https://practicum.yandex.ru/learn/go-advanced-self-paced/courses/8bca0296-484d-45dc-b9ab-f01e0f44f9f4/sprints/145736/topics/63027ac1-f19b-405d-bad5-49e3bbddf30b/lessons/572d89a8-1713-457a-927a-90c2280757bc/"
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive test",
			args: args{
				input: fullURL,
			},
			want:    "KFDaxKze",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetShortLink(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShortLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetShortLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}
