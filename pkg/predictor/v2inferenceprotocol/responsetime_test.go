package v2inferenceprotocol

import (
	"net/http"
	"testing"

	"github.com/waok8s/wao-core/pkg/util"
)

func TestResponseTimePredictor_Endpoint(t *testing.T) {
	type fields struct {
		urlTemplate string
		client      *http.Client
		editorFns   []util.RequestEditorFn
	}
	type args struct {
		appName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "with template",
			fields: fields{
				urlTemplate: "http://example.com/{{.App}}",
			},
			args: args{
				appName: "foo",
			},
			want:    "http://example.com/foo",
			wantErr: false,
		},
		{
			name: "with undefined template field",
			fields: fields{
				urlTemplate: "http://example.com/{{.Foo}}",
			},
			args: args{
				appName: "foo",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "with invalid template",
			fields: fields{
				urlTemplate: "http://example.com/{{.App",
			},
			args: args{
				appName: "foo",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "without template",
			fields: fields{
				urlTemplate: "http://example.com",
			},
			args: args{
				appName: "foo",
			},
			want:    "http://example.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ResponseTimePredictor{
				urlTemplate: tt.fields.urlTemplate,
				client:      tt.fields.client,
				editorFns:   tt.fields.editorFns,
			}
			got, err := p.Endpoint(tt.args.appName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResponseTimePredictor.Endpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ResponseTimePredictor.Endpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
