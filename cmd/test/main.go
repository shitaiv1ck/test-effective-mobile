package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/shitaiv1ck/test-effective-mobile/internal/core/logger"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/repository/postgres"
	httpserver "github.com/shitaiv1ck/test-effective-mobile/internal/core/transport/http/server"
	statsrep "github.com/shitaiv1ck/test-effective-mobile/internal/features/statistics/repository"
	statssrvc "github.com/shitaiv1ck/test-effective-mobile/internal/features/statistics/service"
	statshttp "github.com/shitaiv1ck/test-effective-mobile/internal/features/statistics/transport"
	subsrep "github.com/shitaiv1ck/test-effective-mobile/internal/features/subscriptions/repository"
	subssrvc "github.com/shitaiv1ck/test-effective-mobile/internal/features/subscriptions/service"
	subshttp "github.com/shitaiv1ck/test-effective-mobile/internal/features/subscriptions/transport/http"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logger, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		panic(err)
	}

	logger.Debug("init postgres conn pool...")
	connPool, err := postgres.NewConnPool(ctx, postgres.NewConfigMust())
	if err != nil {
		panic(err)
	}
	defer connPool.Close()

	logger.Debug("init feature: subscriptions...")
	subsRep := subsrep.NewSubsRepository(connPool)
	subsService := subssrvc.NewSubsService(subsRep)
	subsHTTP := subshttp.NewSubHTTP(subsService)

	logger.Debug("init feature: statistics...")
	statsRep := statsrep.NewStatsRepository(connPool)
	statsService := statssrvc.NewStatsService(statsRep)
	statsHTTP := statshttp.NewStatsHTTP(statsService)

	router := http.NewServeMux()
	router.Handle("POST /api/subscriptions", subsHTTP.CreateSubHandler())
	router.Handle("GET /api/subscriptions", subsHTTP.GetSubsHandler())
	router.Handle("GET /api/subscriptions/{sub_id}", subsHTTP.GetSubHandler())
	router.Handle("PATCH /api/subscriptions/{sub_id}", subsHTTP.PatchSubHandler())
	router.Handle("DELETE /api/subscriptions/{sub_id}", subsHTTP.DeleteSubHandler())
	router.Handle("GET /api/subscriptions/statistics", statsHTTP.GetStatisticsHandler())

	logger.Debug("init http server...")
	httpServer := httpserver.NewHTTPServer(router, httpserver.NewConfigMust(), logger)
	if err := httpServer.Run(ctx); err != nil {
		panic(err)
	}
}
