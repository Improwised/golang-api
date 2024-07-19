package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Improwised/golang-api/constants"
	helpers "github.com/Improwised/golang-api/helpers/flipt"
	"github.com/Improwised/golang-api/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"go.flipt.io/flipt-grpc"
	"go.uber.org/zap"
)

type HealthController struct {
	db     *goqu.Database
	logger *zap.Logger
	fc     *flipt.FliptClient
}

func NewHealthController(db *goqu.Database, logger *zap.Logger, fc *flipt.FliptClient) (*HealthController, error) {
	return &HealthController{
		db:     db,
		logger: logger,
		fc:     fc,
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

	if hc.fc != nil {
		fc := *hc.fc

		countryFlagResp, err := helpers.GetVarientFlag(fc, "country_key", "1234", map[string]string{"country": "ind"})
		if err != nil {
			hc.logger.Error("error while health checking of flipt", zap.Error(err))
			return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
		}

		if countryFlagResp.Match {
			fmt.Println("======================================")
			fmt.Println("country is enabled")
			switch countryFlagResp.Value {
			case "orange":
				fmt.Println("======================================")
				fmt.Println("country is enabled for india and color is orange")
				fmt.Println("======================================")
			case "red":
				fmt.Println("======================================")
				fmt.Println("country is enabled for india and color is red")
				fmt.Println("======================================")
			default:
				fmt.Println("======================================")
				fmt.Println("country is enabled for default is white")
				fmt.Println("======================================")
			}

		} else {
			fmt.Println("======================================")
			fmt.Println("country is disabled")
			fmt.Println("======================================")
		}
		fmt.Println()
		fmt.Println("======================================")

		advertisementFlag, err := helpers.GetBooleanFlag(fc, "advertisement")
		if err != nil {
			hc.logger.Error("error while health checking of flipt", zap.Error(err))
			return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
		}

		fmt.Println("======================================")
		fmt.Println()
		fmt.Println("advertisement response", advertisementFlag.Enabled)
		fmt.Println()
		fmt.Println("======================================")

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
	_, err := db.ExecContext(context.TODO(), "SELECT 1")
	if err != nil {
		return err
	}
	return nil
}
