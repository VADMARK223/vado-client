package middleware

import (
	"context"
	"time"
	pb "vado-client/api/pb/auth"
	"vado-client/internal/app"
	"vado-client/internal/constants/code"

	"google.golang.org/grpc/metadata"
)

const TokenAliveMinutes = 15

func WithAuth(appCtx *app.Context, ctx context.Context) context.Context {
	prefs := appCtx.App.Preferences()
	access := prefs.String(code.AccessToken)
	refresh := prefs.String(code.RefreshToken)
	exp := time.Unix(int64(prefs.Int("expires_at")), 0)

	if time.Now().Before(exp) {
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+access)
	}

	authClient := pb.NewAuthServiceClient(appCtx.GRPC)

	appCtx.Log.Debugw("Refresh token", "refresh", refresh)

	if refresh != "" {
		resp, err := authClient.Refresh(ctx, &pb.RefreshRequest{RefreshToken: refresh})
		if err != nil {
			appCtx.Log.Warnw("Error refresh token", "error", err)
			return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+access)
		}

		prefs.SetString(code.AccessToken, resp.Token)
		prefs.SetInt(code.ExpiresAt, int(time.Now().Add(TokenAliveMinutes*time.Minute).Unix()))
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+resp.Token)
	} else {
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+access)
	}
}
