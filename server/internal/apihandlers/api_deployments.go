package apihandlers

import (
	"CloudScapes/pkg/wire"
	"CloudScapes/server/internal/rqctx"
)

func DeploymentsGetHandler(c *rqctx.Context) rqctx.ResponseHandler {
	Deployments, err := c.Deployments.GetDeployments()
	if err != nil {
		return c.SendError(err)
	}
	return c.SendOK(Deployments)
}

func DeploymentsPostHandler(c *rqctx.Context) rqctx.ResponseHandler {
	var newDeployment wire.NewDeployment
	if err := c.DecodeBody(&newDeployment); err != nil {
		return c.SendError(err)
	}

	_, err := c.Plans.GetPlan(newDeployment.PlanID)
	if err != nil {
		return c.SendError(convetErrIfNeeded("Plan", newDeployment.PlanID, err))
	}

	user, err := c.Deployments.CreateDeployment(newDeployment, c.User.ID)
	if err != nil {
		return c.SendError(err)
	}

	return c.SendCreated(user)
}

func DeploymentsPutHandler(c *rqctx.Context) rqctx.ResponseHandler {
	deploymentID, err := c.IdFromPath("deploymentId")
	if err != nil {
		return c.SendError(err)
	}

	var newDeployment wire.NewDeployment
	if err := c.DecodeBody(&newDeployment); err != nil {
		return c.SendError(err)
	}
	if _, err := c.Plans.GetPlan(newDeployment.PlanID); err != nil {
		return c.SendError(convetErrIfNeeded("Plan", newDeployment.PlanID, err))
	}

	Deployment, err := c.Deployments.UpdateDeployment(deploymentID, c.User.ID, newDeployment)
	if err != nil {
		return c.SendError(err)
	}

	return c.SendOK(Deployment)
}

func DeploymentsDeleteHandler(c *rqctx.Context) rqctx.ResponseHandler {
	deploymentID, err := c.IdFromPath("deploymentId")
	if err != nil {
		return c.SendError(err)
	}

	if err := c.Deployments.DeleteDeployment(deploymentID, c.User.ID); err != nil {
		return c.SendError(convetErrIfNeeded("Deployment", deploymentID, err))
	}
	return c.SendNothing()
}
