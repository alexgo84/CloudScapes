package apihandlers

import (
	"CloudScapes/pkg/wire"
	"CloudScapes/server/internal/rqctx"
)

func PlansGetHandler(c *rqctx.Context) rqctx.ResponseHandler {
	plans, err := c.Plans.GetPlans()
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

	_, err := c.Clusters.GetCluster(newPlan.ClusterID)
	if err != nil {
		return c.SendError(convetErrIfNeeded("Cluster", newPlan.ClusterID, err))
	}

	user, err := c.Plans.CreatePlan(newPlan)
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

	if _, err := c.Clusters.GetCluster(newPlan.ClusterID); err != nil {
		return c.SendError(convetErrIfNeeded("Cluster", newPlan.ClusterID, err))
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

	if err := c.Plans.DeletePlan(planID); err != nil {
		return c.SendError(convetErrIfNeeded("Plan", planID, err))
	}
	return c.SendNothing()
}
