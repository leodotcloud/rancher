package ingresshostgen

import (
	"context"
	"fmt"
	"github.com/rancher/norman/types/convert"
	"github.com/sirupsen/logrus"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rancher/rancher/pkg/controllers/user/approuter"
	"github.com/rancher/rancher/pkg/controllers/user/ingress"
	"github.com/rancher/rancher/pkg/settings"
	v1beta12 "github.com/rancher/types/apis/extensions/v1beta1"
	"github.com/rancher/types/config"
	"k8s.io/api/extensions/v1beta1"
)

const (
	ingressStateAnnotation = "field.cattle.io/ingressState"
)

type IngressHostGen struct {
	ingress v1beta12.IngressInterface
}

func Register(ctx context.Context, userOnlyContext *config.UserOnlyContext) {
	c := &IngressHostGen{
		ingress: userOnlyContext.Extensions.Ingresses(""),
	}
	userOnlyContext.Extensions.Ingresses("").AddHandler(ctx, "ingress-host-gen", c.sync)
}

func isGeneratedDomain(obj *v1beta1.Ingress, host, domain string) bool {
	parts := strings.Split(host, ".")
	is := strings.HasSuffix(host, "."+domain) && len(parts) == 8 && parts[1] == obj.Namespace
	logrus.Infof("MP: ingresshostgen: isGeneratedDomain: host: %v is: %v", host, is)
	return is
}

func (i *IngressHostGen) sync(key string, obj *v1beta1.Ingress) (runtime.Object, error) {
	if obj == nil {
		return nil, nil
	}

	ipDomain := settings.IngressIPDomain.Get()
	if ipDomain == "" {
		return nil, nil
	}

	var xipHost string
	for _, status := range obj.Status.LoadBalancer.Ingress {
		if status.IP != "" {
			xipHost = fmt.Sprintf("%s.%s.%s.%s", obj.Name, obj.Namespace, status.IP, ipDomain)
			break
		}
	}

	logrus.Infof("MP: ingresshostgen: xipHost: %v", xipHost)

	if xipHost == "" {
		return nil, nil
	}

	changed := false
	for _, rule := range obj.Spec.Rules {
		if (isGeneratedDomain(obj, rule.Host, ipDomain) || rule.Host == ipDomain) && rule.Host != xipHost && ipDomain != approuter.RdnsIPDomain {
			changed = true
			break
		}
	}

	if !changed {
		return nil, nil
	}

	obj = obj.DeepCopy()
	ingressState := ingress.GetIngressState(obj)
	for i, rule := range obj.Spec.Rules {
		if strings.HasSuffix(rule.Host, ipDomain) {

			for _, path := range rule.HTTP.Paths {
				//oldStateKey := fmt.Sprintf("%s/%s/%s/%s/%s", obj.Name, obj.Namespace, obj.Spec.Rules[i].Host, path.Path, convert.ToString(path.Backend.ServicePort.IntVal))
				//newStateKey := fmt.Sprintf("%s/%s/%s/%s/%s", obj.Name, obj.Namespace, xipHost, path.Path, convert.ToString(path.Backend.ServicePort.IntVal))
				oldStateKey := ingress.GetStateKey(obj.Name, obj.Namespace, obj.Spec.Rules[i].Host, path.Path, convert.ToString(path.Backend.ServicePort.IntVal))
				newStateKey := ingress.GetStateKey(obj.Name, obj.Namespace, xipHost, path.Path, convert.ToString(path.Backend.ServicePort.IntVal))
				oldStateKeyValue := ingressState[oldStateKey]
				delete(ingressState, oldStateKey)
				logrus.Infof("replacing oldStateKey: %v with newStateKey: %v and value: %v", oldStateKey, newStateKey, oldStateKeyValue)
				ingressState[newStateKey] = oldStateKeyValue
			}
			obj.Spec.Rules[i].Host = xipHost
		}
	}
	if err := ingress.SetIngressState(obj, ingressState); err != nil {
		return nil, err
	}
	return i.ingress.Update(obj)
}
