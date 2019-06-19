package cluster

import (
	"bytes"
	"encoding/json"
	"fmt"
	v3 "github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/pkg/controllers/user/cis"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a ActionHandler) runCISScan(actionName string, action *types.Action, apiContext *types.APIContext) error {
	cluster, err := a.ClusterClient.Get(apiContext.ID, v1.GetOptions{})
	if err != nil {
		return httperror.WrapAPIError(err, httperror.NotFound,
			fmt.Sprintf("cluster with id %v doesn't exist", apiContext.ID))
	}
	if cluster.DeletionTimestamp != nil {
		return httperror.NewAPIError(httperror.InvalidType,
			fmt.Sprintf("cluster with id %v is being deleted", apiContext.ID))
	}
	if !v3.ClusterConditionReady.IsTrue(cluster) {
		return httperror.WrapAPIError(err, httperror.ClusterUnavailable,
			fmt.Sprintf("cluster not ready"))
	}
	if _, ok := cluster.Annotations[cis.RunCISScanAnnotation]; ok {
		return httperror.WrapAPIError(err, httperror.Conflict,
			fmt.Sprintf("CIS scan already running on cluster"))
	}

	newCisScan := cis.NewCISScan(cluster)
	cisScan, err := a.ClusterScanClient.Create(newCisScan)
	if err != nil {
		return httperror.WrapAPIError(err, httperror.ServerError,
			fmt.Sprintf("failed to create cis scan object"))
	}

	updatedCluster := cluster.DeepCopy()
	updatedCluster.Annotations[cis.RunCISScanAnnotation] = cisScan.Name

	_, err = a.ClusterClient.Update(updatedCluster)
	if err != nil {
		return httperror.WrapAPIError(err, httperror.ServerError, "failed to update cluster annotation for cis scan")
	}

	cisScanJSON, err := json.Marshal(cisScan)
	if err != nil {
		return httperror.WrapAPIError(err, httperror.ServerError,
			fmt.Sprintf("failed to marshal cis scan object"))
	}

	logrus.Infof("CIS scan requested")
	apiContext.Response.Header().Set("Content-Type", "application/json")
	http.ServeContent(apiContext.Response, apiContext.Request, "clusterScan", time.Now(), bytes.NewReader(cisScanJSON))
	return nil
}
