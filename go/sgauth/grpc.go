//
// Copyright 2018 Google Inc.
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
package sgauth

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"crypto/x509"
	"google.golang.org/grpc/credentials"

	"fmt"
	"github.com/google/oauth2l/go/sgauth/internal"
)

func NewGrpcConn(ctx context.Context, settings *Settings, host string, port string) (*grpc.ClientConn, error) {
	if settings == nil {
		settings = &Settings {
			Scope: DefaultScope,
		}
	}

	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := credentials.NewClientTLSFromCert(pool, "")
	var perRPC credentials.PerRPCCredentials

	if settings.APIKey != "" {
		perRPC = internal.GrpcApiKey{Value: settings.APIKey}
	} else if settings.Scope != "" {
		perRPC, _ = NewGrpcApplicationDefault(ctx, settings)
	} else {
		perRPC, _ = NewGrpcJWT(ctx, settings.Audience)
	}
	return grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
	)
}
