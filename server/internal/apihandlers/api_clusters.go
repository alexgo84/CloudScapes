package apihandlers

import (
	"CloudScapes/pkg/wire"
	"CloudScapes/server/internal/rqctx"
)

func ClustersGetHandler(c *rqctx.Context) rqctx.ResponseHandler {
	clusters, err := c.Clusters.GetClusters(c.Account.ID)
	if err != nil {
		return c.SendError(err)
	}
	return c.SendOK(clusters)
}

func ClustersPostHandler(c *rqctx.Context) rqctx.ResponseHandler {
	var newCluster wire.NewCluster
	if err := c.DecodeBody(&newCluster); err != nil {
		return c.SendError(err)
	}

	user, err := c.Clusters.CreateCluster(newCluster)
	if err != nil {
		return c.SendError(err)
	}

	return c.SendCreated(user)
}

func ClustersDeleteHandler(c *rqctx.Context) rqctx.ResponseHandler {
	clusterID, err := c.IdFromPath("clusterId")
	if err != nil {
		return c.SendError(err)
	}

	if err := c.Clusters.DeleteCluster(c.Account.ID, clusterID); err != nil {
		return c.SendError(convertToAPIIfNeeded("Cluster", clusterID, err))
	}
	return c.SendNothing()
}
