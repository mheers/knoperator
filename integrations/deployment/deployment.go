package deployment

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	config "github.com/mheers/knoperator/config"
)

type DeploymentIntegration struct {
	config *config.Config
	k8s    *kubernetes.Clientset
}

func NewAPI(c *config.Config, k8s *kubernetes.Clientset) (*DeploymentIntegration, error) {
	si := &DeploymentIntegration{
		config: c,
		k8s:    k8s,
	}

	return si, nil
}

func (di *DeploymentIntegration) GetPods() ([]corev1.Pod, error) {
	list, err := di.k8s.CoreV1().Pods(di.config.K8sNamespace).List(context.TODO(), metav1.ListOptions{
		Limit: 100,
	})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (di *DeploymentIntegration) DeletePod(name string) error {
	err := di.k8s.CoreV1().Pods(di.config.K8sNamespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func (di *DeploymentIntegration) GetDeployments() ([]appsv1.Deployment, error) {
	list, err := di.k8s.AppsV1().Deployments(di.config.K8sNamespace).List(context.TODO(), metav1.ListOptions{
		Limit: 100,
	})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (di *DeploymentIntegration) DeleteDeployment(name string) error {
	err := di.k8s.AppsV1().Deployments(di.config.K8sNamespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func (di *DeploymentIntegration) UpdateDeployment(name, image string, command, args []string, env map[string]string) error {
	d, err := di.k8s.AppsV1().Deployments(di.config.K8sNamespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	r := int32(1)
	e := []corev1.EnvVar{}
	for k, v := range env {
		e = append(e, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	d.ObjectMeta.Name = name
	d.ObjectMeta.Namespace = di.config.K8sNamespace
	d.Spec.Selector.MatchLabels["app"] = name
	d.Spec.Replicas = &r
	d.Spec.Template.ObjectMeta.Labels["app"] = name
	d.Spec.Template.Spec.Containers[0].Name = name
	d.Spec.Template.Spec.Containers[0].Image = image
	d.Spec.Template.Spec.Containers[0].Command = command
	d.Spec.Template.Spec.Containers[0].Args = args
	d.Spec.Template.Spec.Containers[0].Env = e

	_, err = di.k8s.AppsV1().Deployments(di.config.K8sNamespace).Update(context.TODO(), d, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (di *DeploymentIntegration) CreateDeployment(name, image string, command, args []string, env map[string]string) error {
	r := int32(1)
	e := []corev1.EnvVar{}
	for k, v := range env {
		e = append(e, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: di.config.K8sNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Replicas: &r,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    name,
							Image:   image,
							Command: command,
							Args:    args,
							Env:     e,
						},
					},
				},
			},
		},
	}

	_, err := di.k8s.AppsV1().Deployments(di.config.K8sNamespace).Create(context.TODO(), d, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (di *DeploymentIntegration) ScaleDeployment(name string, nReplicas int) error {
	s, err := di.k8s.AppsV1().Deployments(di.config.K8sNamespace).GetScale(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	s.Spec.Replicas = int32(nReplicas)

	_, err = di.k8s.AppsV1().Deployments(di.config.K8sNamespace).UpdateScale(context.TODO(), name, s, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
