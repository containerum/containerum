package builder

import (
	"os"
	"testing"

	"github.com/containerum/containerum/embark/pkg/models/requirements"
)

func TestClient_DownloadRequirements(t *testing.T) {
	type args struct {
		dir  string
		reqs requirements.Requirements
	}

	var client, _ = NewCLient()
	os.Remove("data")

	var req requirements.Requirements
	LoadYAML("testdata/requirements.yaml", &req)
	req.Dependencies = req.Dependencies[:1]
	tests := []struct {
		name    string
		fields  *Client
		args    args
		wantErr bool
	}{
		{
			name:   "download requirements",
			fields: client,
			args: args{
				dir:  "data",
				reqs: req,
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.fields
			if err := client.DownloadRequirements(tt.args.dir, tt.args.reqs); (err != nil) != tt.wantErr {
				t.Errorf("Client.DownloadRequirements() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
