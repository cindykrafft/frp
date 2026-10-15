// Copyright 2026 The frp Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	v1 "github.com/fatedier/frp/pkg/config/v1"
)

func TestBoolFuncFlagSet(t *testing.T) {
	require := require.New(t)

	var trueCalls, falseCalls int
	f := &BoolFuncFlag{
		TrueFunc:  func() { trueCalls++ },
		FalseFunc: func() { falseCalls++ },
	}

	require.NoError(f.Set("true"))
	require.Equal(1, trueCalls)
	require.Equal(0, falseCalls)
	require.Equal("true", f.String())

	require.NoError(f.Set("false"))
	require.Equal(1, trueCalls)
	require.Equal(1, falseCalls)
	require.Equal("false", f.String())

	require.Error(f.Set("not-a-bool"))
}

func TestRegisterServerConfigFlagsDashboardTLSMode(t *testing.T) {
	for _, args := range [][]string{
		{"--dashboard_tls_mode=true"},
		{"--dashboard_tls_mode"},
	} {
		t.Run(args[0], func(t *testing.T) {
			require := require.New(t)

			cfg := v1.ServerConfig{}
			cmd := &cobra.Command{Use: "frps", RunE: func(*cobra.Command, []string) error { return nil }}
			RegisterServerConfigFlags(cmd, &cfg)
			cmd.SetArgs(append([]string{
				"--dashboard_port=7500",
				"--dashboard_tls_cert_file=server.crt",
				"--dashboard_tls_key_file=server.key",
			}, args...))
			require.NoError(cmd.Execute())

			require.NotNil(cfg.WebServer.TLS, "dashboard TLS should be enabled")
			require.Equal("server.crt", cfg.WebServer.TLS.CertFile)
			require.Equal("server.key", cfg.WebServer.TLS.KeyFile)
		})
	}
}
