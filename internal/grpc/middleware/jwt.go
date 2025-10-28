package middleware

import (
	"context"
	"time"
	pb "vado-client/api/pb/auth"
	"vado-client/internal/app"

	"google.golang.org/grpc/metadata"
)

const TokenAliveMinutes = 15

func WithAuth(appCtx *app.Context, ctx context.Context) context.Context {
	access := appCtx.Prefs.AccessToken()
	refresh := appCtx.Prefs.RefreshToken()
	exp := time.Unix(appCtx.Prefs.ExpiresAt(), 0)

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

		appCtx.Prefs.SetAccessToken(resp.Token)
		appCtx.Prefs.SetExpiresAt(time.Now().Add(TokenAliveMinutes * time.Minute).Unix())
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+resp.Token)
	} else {
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+access)
	}
}
