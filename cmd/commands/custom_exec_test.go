package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecTarget(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		wantSelector  string
		wantContainer string
		wantErr       string
	}{
		{
			name:          "default toolbox",
			wantSelector:  "app=rook-ceph-tools",
			wantContainer: "rook-ceph-tools",
		},
		{
			name:          "toolbox",
			args:          []string{"toolbox"},
			wantSelector:  "app=rook-ceph-tools",
			wantContainer: "rook-ceph-tools",
		},
		{
			name:          "operator",
			args:          []string{"operator"},
			wantSelector:  "app=rook-ceph-operator",
			wantContainer: "rook-ceph-operator",
		},
		{
			name:          "mon",
			args:          []string{"mon-a"},
			wantSelector:  "app=rook-ceph-mon,mon=a,mon_canary!=true",
			wantContainer: "mon",
		},
		{
			name:          "mgr",
			args:          []string{"mgr-b"},
			wantSelector:  "app=rook-ceph-mgr,mgr=b",
			wantContainer: "mgr",
		},
		{
			name:          "osd",
			args:          []string{"osd-42"},
			wantSelector:  "app=rook-ceph-osd,osd=42",
			wantContainer: "osd",
		},
		{
			name:    "missing daemon ID",
			args:    []string{"mon-"},
			wantErr: `invalid component "mon-": expected toolbox, operator, mon-<id>, mgr-<id>, osd-<id>, mds-<id>, or rgw-<id>`,
		},
		{
			name:    "unsupported component",
			args:    []string{"nfs-a"},
			wantErr: `invalid component "nfs-a": expected toolbox, operator, mon-<id>, mgr-<id>, osd-<id>, mds-<id>, or rgw-<id>`,
		},
		{
			name:    "invalid daemon ID",
			args:    []string{"osd-bad/value"},
			wantErr: `invalid component "osd-bad/value": expected toolbox, operator, mon-<id>, mgr-<id>, osd-<id>, mds-<id>, or rgw-<id>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector, container, err := execTarget(tt.args)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				assert.Empty(t, selector)
				assert.Empty(t, container)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantSelector, selector)
			assert.Equal(t, tt.wantContainer, container)
		})
	}
}
