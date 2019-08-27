package main

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/guozijn/nuntius"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_receiverConfByReceiver(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *ReceiverConf
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := receiverConfByReceiver(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("receiverConfByReceiver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_providerByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    nuntius.Provider
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := providerByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("providerByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("providerByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_errorHandler(t *testing.T) {
	type args struct {
		w        http.ResponseWriter
		status   int
		err      error
		provider string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorHandler(tt.args.w, tt.args.status, tt.args.err, tt.args.provider)
		})
	}
}
