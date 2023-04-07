package constants

import "time"

const DatabaseTimeout = 3 * time.Second
const AuthTokenExpirySeconds = 60 * 15              // 15 mins
const RefreshTokenExpirySeconds = 60 * 60 * 24 * 30 // 30 days
