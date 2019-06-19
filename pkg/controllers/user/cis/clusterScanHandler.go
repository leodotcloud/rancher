package cis

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rancher/rancher/pkg/app/utils"
	"github.com/rancher/rancher/pkg/settings"
	"github.com/rancher/rancher/pkg/systemaccount"
	rcorev1 "github.com/rancher/types/apis/core/v1"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	mgmtv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	projv3 "github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type cisScanHandler struct {
	clusterLister                v3.ClusterLister
	mgmtCtxClusterClient         v3.ClusterInterface
	mgmtCtxProjClient            mgmtv3.ProjectInterface
	mgmtCtxAppClient             projv3.AppInterface
	mgmtCtxTemplateVersionClient mgmtv3.CatalogTemplateVersionInterface
	mgmtCtxClusterScanClient     mgmtv3.ClusterScanInterface
	userCtx                      *config.UserContext
	userCtxNSClient              rcorev1.NamespaceInterface
	clusterNamespace             string
	systemAccountManager         *systemaccount.Manager
	configMapsClient             rcorev1.ConfigMapInterface
}

func (csh *cisScanHandler) Create(cs *mgmtv3.ClusterScan) (runtime.Object, error) {
	logrus.Debugf("cisScanHandler: Create: %+v", cs)
	var err error
	cluster, err := csh.mgmtCtxClusterClient.Get(cs.ClusterName, metav1.GetOptions{})
	if err != nil {
		return cs, err
	}
	if !v3.ClusterConditionReady.IsTrue(cluster) {
		return cs, fmt.Errorf("cisScanHandler: Create: cluster not ready")
	}
	if !mgmtv3.ClusterScanConditionCreated.IsTrue(cs) {
		logrus.Infof("cisScanHandler: Create: deploying helm chart")
		// Deploy the system helm chart
		if err := csh.deployApp(cs.Spec.ClusterID, cs.Name); err != nil {
			logrus.Errorf("cisScanHandler: Create: error deploying app: %v", err)
			return cs, err
		}
		mgmtv3.ClusterScanConditionCreated.True(cs)
		mgmtv3.ClusterScanConditionCompleted.Unknown(cs)

		cs, err = csh.mgmtCtxClusterScanClient.Update(cs)
		if err != nil {
			return cs, err
		}
	}
	return cs, nil
}

func (csh *cisScanHandler) Remove(cs *mgmtv3.ClusterScan) (runtime.Object, error) {
	logrus.Debugf("cisScanHandler: Remove: %+v", cs)
	// Delete the configmap associated with this scan
	err := csh.configMapsClient.Delete(cs.Name, &metav1.DeleteOptions{})
	if err != nil {
		return cs, fmt.Errorf("cisScanHandler: Delete: error deleting cm=%v", cs.Name)
	}

	return cs, nil
}

func (csh *cisScanHandler) Updated(cs *mgmtv3.ClusterScan) (runtime.Object, error) {
	logrus.Debugf("cisScanHandler: Updated: %+v", cs)
	if !mgmtv3.ClusterScanConditionCompleted.IsUnknown(cs) &&
		mgmtv3.ClusterScanConditionCompleted.IsFalse(cs) {
		// Delete the system helm chart
		if err := csh.deleteApp(cs.Name); err != nil {
			logrus.Errorf("cisScanHandler: error deleting app: %v", err)
			return nil, err
		}

		cluster, err := csh.clusterLister.Get("", csh.clusterNamespace)
		if err != nil {
			return nil, fmt.Errorf("cisScanHandler: error getting cluster %v", err)
		}

		updatedCluster := cluster.DeepCopy()
		delete(updatedCluster.Annotations, RunCISScanAnnotation)
		if _, err := csh.mgmtCtxClusterClient.Update(updatedCluster); err != nil {
			return nil, fmt.Errorf("cisScanHandler: failed to update cluster about CIS scan completion")
		}

		v3.ClusterScanConditionCompleted.True(cs)
		_, err = csh.mgmtCtxClusterScanClient.Update(cs)
		if err != nil {
			return nil, fmt.Errorf("error updating condition of cluster scan object: %v", cs.Name)
		}
	}
	return cs, nil
}

func (csh *cisScanHandler) deployApp(clusterName, appName string) error {
	appCatalogID := settings.SystemCISBenchmarkCatalogID.Get()
	err := utils.DetectAppCatalogExistence(appCatalogID, csh.mgmtCtxTemplateVersionClient)
	if err != nil {
		return errors.Wrapf(err, "failed to find cis system catalog %q", appCatalogID)
	}
	appDeployProjectID, err := utils.GetSystemProjectID(csh.mgmtCtxProjClient)
	if err != nil {
		return err
	}

	appProjectName, err := utils.EnsureAppProjectName(csh.userCtxNSClient, appDeployProjectID, clusterName, DefaultNamespaceForCis)
	if err != nil {
		return err
	}

	creator, err := csh.systemAccountManager.GetSystemUser(clusterName)
	if err != nil {
		return err
	}

	appAnswers := map[string]string{
		"owner": appName,
	}

	app := &projv3.App{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{creatorIDAnno: creator.Name},
			Name:        appName,
			Namespace:   appDeployProjectID,
		},
		Spec: projv3.AppSpec{
			Answers:         appAnswers,
			Description:     "Rancher CIS Benchmark",
			ExternalID:      appCatalogID,
			ProjectName:     appProjectName,
			TargetNamespace: DefaultNamespaceForCis,
		},
	}

	_, err = utils.DeployApp(csh.mgmtCtxAppClient, appDeployProjectID, app, false)
	if err != nil {
		return err
	}

	return nil
}

func (csh *cisScanHandler) deleteApp(appName string) error {
	appDeployProjectID, err := utils.GetSystemProjectID(csh.mgmtCtxProjClient)
	if err != nil {
		return err
	}

	if err := utils.DeleteApp(csh.mgmtCtxAppClient, appDeployProjectID, appName); err != nil {
		return err
	}

	return nil
}
