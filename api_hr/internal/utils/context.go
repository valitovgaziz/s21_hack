package utils

import "context"

type contextKey string

const UserClaimsKey contextKey = "userClaims"

func GetUserIDFromContext(ctx context.Context) *uint {
	claims, ok := ctx.Value(UserClaimsKey).(*Claims)
	if !ok || claims == nil {
		return nil
	}
	return &claims.UserID
}

func GetUserEmailFromContext(ctx context.Context) string {
	claims, ok := ctx.Value(UserClaimsKey).(*Claims)
	if !ok || claims == nil {
		return ""
	}
	return claims.Email
}
