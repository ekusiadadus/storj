// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package console_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"storj.io/common/testcontext"
	"storj.io/storj/private/testplanet"
	"storj.io/storj/satellite"
	"storj.io/storj/satellite/console"
)

const validResponseToken = "myResponseToken"

type mockRecaptcha struct{}

func (r mockRecaptcha) Verify(ctx context.Context, responseToken string, userIP string) (bool, error) {
	return responseToken == validResponseToken, nil
}

// TestRegistrationRecaptcha ensures that registration reCAPTCHA service is working properly.
func TestRegistrationRecaptcha(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount: 1,
		Reconfigure: testplanet.Reconfigure{
			Satellite: func(log *zap.Logger, index int, config *satellite.Config) {
				config.Console.Captcha.Registration.Recaptcha.Enabled = true
				config.Console.Captcha.Registration.Recaptcha.SecretKey = "mySecretKey"
				config.Console.Captcha.Registration.Recaptcha.SiteKey = "mySiteKey"
			},
		},
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		service := planet.Satellites[0].API.Console.Service
		require.NotNil(t, service)
		service.TestSwapCaptchaHandler(mockRecaptcha{})

		regToken1, err := service.CreateRegToken(ctx, 1)
		require.NoError(t, err)

		user, err := service.CreateUser(ctx, console.CreateUser{
			FullName:        "User",
			Email:           "u@mail.test",
			Password:        "password",
			CaptchaResponse: validResponseToken,
		}, regToken1.Secret)

		require.NotNil(t, user)
		require.NoError(t, err)

		regToken2, err := service.CreateRegToken(ctx, 1)
		require.NoError(t, err)

		user, err = service.CreateUser(ctx, console.CreateUser{
			FullName:        "User2",
			Email:           "u2@mail.test",
			Password:        "password",
			CaptchaResponse: "wrong",
		}, regToken2.Secret)

		require.Nil(t, user)
		require.True(t, console.ErrCaptcha.Has(err))
	})
}

// TestLoginRecaptcha ensures that login reCAPTCHA service is working properly.
func TestLoginRecaptcha(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount: 1,
		Reconfigure: testplanet.Reconfigure{
			Satellite: func(log *zap.Logger, index int, config *satellite.Config) {
				config.Console.Captcha.Login.Recaptcha.Enabled = true
				config.Console.Captcha.Login.Recaptcha.SecretKey = "mySecretKey"
				config.Console.Captcha.Login.Recaptcha.SiteKey = "mySiteKey"
			},
		},
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		service := planet.Satellites[0].API.Console.Service
		require.NotNil(t, service)
		service.TestSwapCaptchaHandler(mockRecaptcha{})

		regToken, err := service.CreateRegToken(ctx, 1)
		require.NoError(t, err)

		email := "user@mail.test"
		password := "password"

		user, err := service.CreateUser(ctx, console.CreateUser{
			FullName: "User",
			Email:    email,
			Password: password,
		}, regToken.Secret)

		require.NotNil(t, user)
		require.NoError(t, err)

		activationToken, err := service.GenerateActivationToken(ctx, user.ID, user.Email)
		require.NoError(t, err)

		user, err = service.ActivateAccount(ctx, activationToken)
		require.NotNil(t, user)
		require.NoError(t, err)

		token, err := service.Token(ctx, console.AuthUser{
			Email:           email,
			Password:        password,
			CaptchaResponse: validResponseToken,
		})

		require.NotEmpty(t, token)
		require.NoError(t, err)

		token, err = service.Token(ctx, console.AuthUser{
			Email:           email,
			Password:        password,
			CaptchaResponse: "wrong",
		})

		require.Empty(t, token)
		require.True(t, console.ErrCaptcha.Has(err))
	})
}
