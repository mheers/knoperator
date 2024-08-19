package deployment

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (di *DeploymentIntegration) GetJobs() ([]batchv1.Job, error) {
	list, err := di.k8s.BatchV1().Jobs(di.config.K8sNamespace).List(context.TODO(), metav1.ListOptions{
		Limit: 100,
	})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (di *DeploymentIntegration) DeleteJob(name string) error {
	err := di.k8s.BatchV1().Jobs(di.config.K8sNamespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func (di *DeploymentIntegration) CreateJob(name, image string, command, args []string, env map[string]string) error {
	e := []corev1.EnvVar{}
	for k, v := range env {
		e = append(e, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	d := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: di.config.K8sNamespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    name,
							Image:   image,
							Command: command,
							Args:    args,
							Env:     e,
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "default-env",
										},
									},
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "default-env",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "default-env",
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := di.k8s.BatchV1().Jobs(di.config.K8sNamespace).Create(context.TODO(), d, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
