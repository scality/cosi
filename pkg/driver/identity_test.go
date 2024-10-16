package driver

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cosiapi "sigs.k8s.io/container-object-storage-interface-spec"
)

func TestDriverGetInfo(t *testing.T) {
	longProvisioner := "scality-cosi-driver" + strings.Repeat("x", 1000)

	tests := []struct {
		name        string
		provisioner string
		request     *cosiapi.DriverGetInfoRequest
		want        *cosiapi.DriverGetInfoResponse
		wantErr     bool
		errCode     codes.Code
	}{
		{
			name:        "Valid provisioner name",
			provisioner: "scality-cosi-driver",
			request:     &cosiapi.DriverGetInfoRequest{},
			want:        &cosiapi.DriverGetInfoResponse{Name: "scality-cosi-driver"},
			wantErr:     false,
		},
		{
			name:        "Empty provisioner name",
			provisioner: "",
			request:     &cosiapi.DriverGetInfoRequest{},
			want:        nil,
			wantErr:     true,
			errCode:     codes.InvalidArgument,
		},
		{
			name:        "Empty request object",
			provisioner: "scality-cosi-driver",
			request:     nil, // Test for nil request to ensure function handles nil input gracefully.
			want:        &cosiapi.DriverGetInfoResponse{Name: "scality-cosi-driver"},
			wantErr:     false,
		},
		{
			name:        "Long provisioner name",
			provisioner: longProvisioner,
			request:     &cosiapi.DriverGetInfoRequest{},
			want:        &cosiapi.DriverGetInfoResponse{Name: longProvisioner},
			wantErr:     false,
		},
		{
			name:        "Provisioner name with special characters",
			provisioner: "scality-cosi-driver-ß∂ƒ©",
			request:     &cosiapi.DriverGetInfoRequest{},
			want:        &cosiapi.DriverGetInfoResponse{Name: "scality-cosi-driver-ß∂ƒ©"},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newIdentityServer(tt.provisioner)

			resp, err := server.DriverGetInfo(context.Background(), tt.request)

			if tt.wantErr {
				assertErrorWithCode(t, err, tt.errCode)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.want, resp)
			}
		})
	}
}

// Helper function to create identityServer
func newIdentityServer(provisioner string) *identityServer {
	return &identityServer{
		provisioner: provisioner,
	}
}

// Helper function to assert error codes
func assertErrorWithCode(t *testing.T, err error, expectedCode codes.Code) {
	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, expectedCode, st.Code())
}
