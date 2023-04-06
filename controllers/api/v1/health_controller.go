package v1

import (
	"net/http"

	"github.com/Improwised/golang-api/constants"
	"github.com/Improwised/golang-api/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HealthController struct {
	db     *goqu.Database
	logger *zap.Logger
}

func NewHealthController(db *goqu.Database, logger *zap.Logger) (*HealthController, error) {
	return &HealthController{
		db:     db,
		logger: logger,
	}, nil
}

// Overall check overall health of application as well as dependencies health check
// swagger:route GET /healthz Healthcheck overallHealthCheck
//
//	Overall health check
//
//	Overall health check
//
// Produces:
// - application/json
//
// Responses:
//
//	200: GenericResOk
//	500: GenericResError
func (hc *HealthController) Overall(ctx *fiber.Ctx) error {
	err := healthDb(hc.db)
	if err != nil {
		hc.logger.Error("error while health checking of db", zap.Error(err))
		return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
	}

	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

func (hc *HealthController) Self(ctx *fiber.Ctx) error {
	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

// Database health check
// swagger:route GET /healthz/db Healthcheck dbHealthCheck
//
//	Database health check
//
//	Database health check
//
// Produces:
// - application/json
//
// Responses:
//
//	200: GenericResOk
//	500: GenericResError
func (hc *HealthController) Db(ctx *fiber.Ctx) error {
	err := healthDb(hc.db)
	if err != nil {
		hc.logger.Error("error while health checking of db", zap.Error(err))
		return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
	}

	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

///////////////////////
// HealthCheck CORE
//////////////////////

func healthDb(db *goqu.Database) error {
	_, err := db.Query("SELECT 1")
	if err != nil {
		return err
	}
	return nil
}
