package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelmChartUpdater_Files(t *testing.T) {
	assert.Equal(t, []string{"Chart.yaml"}, HelmChart().Files())
}

func TestHelmChartUpdater_CreateNewFiles(t *testing.T) {
	assert.False(t, HelmChart().CreateNewFiles())
}

func TestHelmChartUpdater_Update(t *testing.T) {
	tests := []updaterTestCase{
		{
			name:    "simple Chart.yaml",
			content: "apiVersion: v2\nname: test-chart\nversion: 1.0.0",
			info: ReleaseInfo{
				Version: "v2.0.5",
			},
			want:    "apiVersion: v2\nname: test-chart\nversion: 2.0.5",
			wantErr: assert.NoError,
		},
		{
			name:    "Chart.yaml with description and appVersion",
			content: "apiVersion: v2\nname: test-chart\ndescription: A Helm chart for Kubernetes\ntype: application\nversion: 1.0.0\nappVersion: \"1.16.0\"",
			info: ReleaseInfo{
				Version: "v2.1.3",
			},
			want:    "apiVersion: v2\nname: test-chart\ndescription: A Helm chart for Kubernetes\ntype: application\nversion: 2.1.3\nappVersion: \"1.16.0\"",
			wantErr: assert.NoError,
		},
		{
			name:    "Chart.yaml with spaces around colon",
			content: "apiVersion: v2\nname: test-chart\nversion  :  1.0.0\nappVersion: \"1.16.0\"",
			info: ReleaseInfo{
				Version: "v3.2.1",
			},
			want:    "apiVersion: v2\nname: test-chart\nversion  :  3.2.1\nappVersion: \"1.16.0\"",
			wantErr: assert.NoError,
		},
		{
			name:    "Chart.yaml with inline comment",
			content: "apiVersion: v2\nname: test-chart\nversion: 1.0.0  # This is the chart version",
			info: ReleaseInfo{
				Version: "v2.0.0",
			},
			want:    "apiVersion: v2\nname: test-chart\nversion: 2.0.0  # This is the chart version",
			wantErr: assert.NoError,
		},
		{
			name:    "Chart.yaml with indented version",
			content: "apiVersion: v2\nname: test-chart\n  version: 1.0.0\nappVersion: \"1.16.0\"",
			info: ReleaseInfo{
				Version: "v1.5.2",
			},
			want:    "apiVersion: v2\nname: test-chart\n  version: 1.5.2\nappVersion: \"1.16.0\"",
			wantErr: assert.NoError,
		},
		{
			name:    "Chart.yaml with quoted version",
			content: "apiVersion: v2\nname: test-chart\nversion: \"1.0.0\"\nappVersion: \"1.16.0\"",
			info: ReleaseInfo{
				Version: "v2.3.4",
			},
			want:    "apiVersion: v2\nname: test-chart\nversion: 2.3.4\nappVersion: \"1.16.0\"",
			wantErr: assert.NoError,
		},
		{
			name:    "invalid yaml",
			content: `not yaml`,
			info: ReleaseInfo{
				Version: "v2.0.0",
			},
			want:    `not yaml`,
			wantErr: assert.NoError,
		},
		{
			name:    "yaml without version",
			content: "apiVersion: v2\nname: test-chart",
			info: ReleaseInfo{
				Version: "v2.0.0",
			},
			want:    "apiVersion: v2\nname: test-chart",
			wantErr: assert.NoError,
		},
		{
			name:    "empty version value",
			content: "apiVersion: v2\nname: test-chart\nversion:\nappVersion: \"1.16.0\"",
			info: ReleaseInfo{
				Version: "v1.0.0",
			},
			want:    "apiVersion: v2\nname: test-chart\nversion: 1.0.0\nappVersion: \"1.16.0\"",
			wantErr: assert.NoError,
		},
		{
			name:    "version without v prefix",
			content: "apiVersion: v2\nname: test-chart\nversion: 1.0.0",
			info: ReleaseInfo{
				Version: "2.0.0", // no "v" prefix
			},
			want:    "apiVersion: v2\nname: test-chart\nversion: 2.0.0",
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runUpdaterTest(t, HelmChart(), tt)
		})
	}
}