package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zufardhiyaulhaq/ngrok-operator/pkg/builder"

	ngrokcomv1alpha1 "github.com/zufardhiyaulhaq/ngrok-operator/api/v1alpha1"
)

type configurationBuilderTestCase struct {
	name                  string
	spec                  *ngrokcomv1alpha1.NgrokSpec
	expectedConfiguration string
	expectedError         bool
}

var testGrid = []configurationBuilderTestCase{
	{
		name: "empty service",
		spec: &ngrokcomv1alpha1.NgrokSpec{
			Service:       "",
			Port:          80,
			AuthTokenType: "plain",
			Protocol:      "http",
			Inspect:       false,
		},
		expectedError:         true,
		expectedConfiguration: "",
	},
	{
		name: "simple HTTP service",
		spec: &ngrokcomv1alpha1.NgrokSpec{
			Service:       "foo.bar",
			Port:          80,
			AuthTokenType: "plain",
			Protocol:      "http",
			Inspect:       false,
		},
		expectedError: false,
		expectedConfiguration: `web_addr: 0.0.0.0:4040
tunnels:
  app:
    proto: http
    addr: foo.bar:80
    inspect: false`,
	},
	{
		name: "full HTTP service",
		spec: &ngrokcomv1alpha1.NgrokSpec{
			Service:       "foo.bar",
			Port:          80,
			Region:        "us",
			Auth:          "foo:bar",
			AuthToken:     "foo",
			AuthTokenType: "plain",
			Protocol:      "http",
			BindTLS:       "true",
			Inspect:       false,
			HostHeader:    "foo.dev",
			Hostname:      "foo.bar",
		},
		expectedError: false,
		expectedConfiguration: `web_addr: 0.0.0.0:4040
authtoken: foo
region: us
tunnels:
  app:
    proto: http
    addr: foo.bar:80
    inspect: false
    auth: foo:bar
    host_header: foo.dev
    bind_tls: true
    hostname: foo.bar`,
	},
	{
		name: "simple TCP service",
		spec: &ngrokcomv1alpha1.NgrokSpec{
			Service:    "foo.bar",
			Port:       80,
			Protocol:   "tcp",
			RemoteAddr: "1.tcp.ngrok.io:12345",
		},
		expectedError: false,
		expectedConfiguration: `web_addr: 0.0.0.0:4040
tunnels:
  app:
    proto: tcp
    addr: foo.bar:80
    remote_addr: 1.tcp.ngrok.io:12345`,
	},
}

func TestConfigurationBuilder(t *testing.T) {
	for _, test := range testGrid {
		t.Run(test.name, func(t *testing.T) {
			configuration, err := builder.NewNgrokConfigurationBuilder().
				SetSpec(test.spec).
				Build()

			assert.Equal(t, configuration, test.expectedConfiguration)

			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
