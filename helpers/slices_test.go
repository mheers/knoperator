package helpers

import (
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	type args struct {
		slice []string
		val   string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name: "found",
			args: args{
				slice: []string{"apple", "orange", "banana"},
				val:   "orange",
			},
			want:  1,
			want1: true,
		},
		{
			name: "not found",
			args: args{
				slice: []string{"apple", "orange", "banana"},
				val:   "lime",
			},
			want:  -1,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Find(tt.args.slice, tt.args.val)
			if got != tt.want {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Find() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetKeysOfMap(t *testing.T) {
	type args struct {
		m map[string]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty map",
			args: args{
				m: map[string]string{},
			},
			want: []string{},
		},
		{
			name: "one element map",
			args: args{
				m: map[string]string{
					"key": "value",
				},
			},
			want: []string{"key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKeysOfMap(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKeysOfMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
