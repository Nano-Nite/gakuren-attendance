package v1

import (
	"strings"

	"gakuren-system.com/pkg/db"
	"gakuren-system.com/pkg/helper"
	"gakuren-system.com/pkg/model/attendance"
	"github.com/gofiber/fiber/v3"
)

func SetupLocationRoutes(app *fiber.App, version string) {
	baseURL := "/" + strings.Trim(version, "/") + "/attendance/location"

	app.Post(baseURL+"/create", func(c fiber.Ctx) error {
		payload := new(attendance.CreateLocationRequest)

		//* validate header
		schoolUUID, tenantUUID, requesterUUID, err := helper.ValidateRequestHeader(c)
		if err != nil {
			return helper.ReturnResponse(c, fiber.StatusUnauthorized, "Missing or invalid authentication data", nil, err)
		}

		//* validate payload
		if err := c.Bind().Body(payload); err != nil {
			return helper.ReturnResponse(c, fiber.StatusBadRequest, "Missing or invalid body", nil, err)
		}

		//* permission check
		if ok, permissionErr := helper.GetUserPermission(requesterUUID.String(), helper.LOCATION_CREATE_PERMISSION); permissionErr != nil || !ok {
			return helper.ReturnResponse(c, fiber.StatusUnauthorized, "Access Denied", nil, permissionErr)
		}

		//* validate existing data in db
		if err = helper.ValidateRequestLocation(*payload); err != nil {
			return helper.ReturnResponse(c, fiber.StatusConflict, "Data already exist", nil, err)
		}

		createPayload := new(attendance.LocationModel{
			TenantUUID:  tenantUUID,
			SchoolUUID:  schoolUUID,
			Code:        payload.Code,
			Name:        payload.Name,
			Latitude:    payload.Latitude,
			Longitude:   payload.Longitude,
			RadiusMeter: payload.RadiusMeter,
			CreatedBy:   &requesterUUID,
		})

		//* check if user can bypass workflow
		canBypass, err := helper.ValidateApprovalBypass(requesterUUID)
		if err != nil {
			return helper.ReturnResponse(c, fiber.StatusInternalServerError, "Failed to check approval bypass", nil, err)
		}
		if canBypass {
			tx, err := db.Conn.Begin(c.Context())
			if err != nil {
				return err
			}
			defer tx.Rollback(c.Context())

			if err = helper.CreateLocation(*createPayload, tx, c.Context()); err != nil {
				return helper.ReturnResponse(c, fiber.StatusInternalServerError, "Failed to create location", nil, err)
			}

			tx.Commit(c.Context())
			return helper.ReturnResponse(c, fiber.StatusOK, "success", map[string]any{"payload": payload}, nil)
		}

		//* using worklow approval
		workflow, err := helper.DetermineWorkflow(schoolUUID, tenantUUID, requesterUUID, helper.LOCATION_CREATE_PERMISSION, helper.ACTION_CODE_CREATE)
		if err != nil {
			return helper.ReturnResponse(c, fiber.StatusInternalServerError, "Failed to determine approval workflow", nil, err)
		}

		//* fallback when workflow didn't exist, shall we continue or not execute at all
		if workflow == nil {
			err = helper.ExecuteWorkflowFallback(func() error {
				tx, err := db.Conn.Begin(c.Context())
				if err != nil {
					return err
				}
				defer tx.Rollback(c.Context())

				if err = helper.CreateLocation(*createPayload, tx, c.Context()); err != nil {
					return helper.ReturnResponse(c, fiber.StatusInternalServerError, "Failed to create location", nil, err)
				}

				return tx.Commit(c.Context())
			})
			if err != nil {
				return helper.ReturnResponse(c, fiber.StatusInternalServerError, "Create rejected by workflow configuration", nil, err)
			}
			return helper.ReturnResponse(c, fiber.StatusOK, "success", map[string]any{"payload": payload}, nil)
		}

		//* create approval instance
		approvalUUID, err := helper.CreateApproval(*workflow, schoolUUID, tenantUUID, requesterUUID, nil, helper.ACTION_CODE_CREATE, helper.LOCATION_ENTITY_TYPE, helper.LOCATION_MODULE_CODE, payload)
		if err != nil {
			return helper.ReturnResponse(c, fiber.StatusInternalServerError, "Failed to create user approval", nil, err)
		}

		return helper.ReturnResponse(c, fiber.StatusOK, "success", map[string]any{"approval_uuid": approvalUUID}, nil)
	})

	// get detail
	app.Get(baseURL+"/get-all", func(c fiber.Ctx) error {
		schoolUUID, tenantUUID, _, err := helper.ValidateRequestHeader(c)
		if err != nil {
			return helper.ReturnResponse(c, fiber.StatusUnauthorized, "Missing or invalid authentication data", nil, err)
		}

		result, err := helper.GetAllLocation(schoolUUID, tenantUUID)
		if err != nil {
			return helper.ReturnResponse(c, fiber.StatusBadRequest, "Failed to search users", nil, err)
		}
		return helper.ReturnResponse(c, fiber.StatusOK, "success", result, nil)
	})

}
