package apihandlers

import (
	"CloudScapes/pkg/wire"
	"CloudScapes/server/internal/rqctx"
)

func PlansGetHandler(c *rqctx.Context) rqctx.ResponseHandler {
	plans, err := c.Plans.GetPlans(c.Account.ID)
	if err != nil {
		return c.SendError(err)
	}
	return c.SendOK(plans)
}

func PlansPostHandler(c *rqctx.Context) rqctx.ResponseHandler {
	var newPlan wire.NewPlan
	if err := c.DecodeBody(&newPlan); err != nil {
		return c.SendError(err)
	}

	user, err := c.Plans.CreatePlan(newPlan)
	if err != nil {
		return c.SendError(err)
	}

	return c.SendCreated(user)
}

func PlansPutHandler(c *rqctx.Context) rqctx.ResponseHandler {
	var newPlan wire.NewPlan
	if err := c.DecodeBody(&newPlan); err != nil {
		return c.SendError(err)
	}

	planID, err := c.IdFromPath("plans")
	if err != nil {
		return c.SendError(err)
	}

	user, err := c.Plans.UpdatePlan(planID, newPlan)
	if err != nil {
		return c.SendError(err)
	}

	return c.SendCreated(user)
}
