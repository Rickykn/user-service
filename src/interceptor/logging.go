package interceptor

import (
	"context"
	"github.com/Rickykn/user-service/src/logger"

	"google.golang.org/grpc"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Tambahkan method ke logger
		log := logger.Get().With().Str("method", info.FullMethod).Logger()

		// Inject ke context
		ctx = logger.InjectContext(ctx, log)

		// Eksekusi handler
		resp, err := handler(ctx, req)
		if err != nil {
			logger.WithContext(ctx).Error().Err(err).Msg("Request failed")
		} else {
			logger.WithContext(ctx).Info().Msg("Request success")
		}

		return resp, err
	}
}
