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

	_, err := c.Clusters.GetCluster(c.Account.ID, newPlan.ClusterID)
	if err != nil {
		return c.SendError(convertToAPIIfNeeded("Cluster", newPlan.ClusterID, err))
	}

	user, err := c.Plans.CreatePlan(c.Account.ID, newPlan)
	if err != nil {
		return c.SendError(err)
	}

	return c.SendCreated(user)
}

func PlansPutHandler(c *rqctx.Context) rqctx.ResponseHandler {
	planID, err := c.IdFromPath("planId")
	if err != nil {
		return c.SendError(err)
	}

	var newPlan wire.NewPlan
	if err := c.DecodeBody(&newPlan); err != nil {
		return c.SendError(err)
	}

	if _, err := c.Clusters.GetCluster(c.Account.ID, newPlan.ClusterID); err != nil {
		return c.SendError(convertToAPIIfNeeded("Cluster", newPlan.ClusterID, err))
	}

	user, err := c.Plans.UpdatePlan(planID, newPlan)
	if err != nil {
		return c.SendError(err)
	}

	return c.SendCreated(user)
}

func PlansDeleteHandler(c *rqctx.Context) rqctx.ResponseHandler {
	planID, err := c.IdFromPath("clusterId")
	if err != nil {
		return c.SendError(err)
	}

	if err := c.Plans.DeletePlan(c.Account.ID, planID); err != nil {
		return c.SendError(convertToAPIIfNeeded("Plan", planID, err))
	}
	return c.SendNothing()
}
