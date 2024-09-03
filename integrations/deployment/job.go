package deployment

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"path"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"
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

func (di *DeploymentIntegration) WatchJobs(nc *nats.Conn) error {
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		return di.k8s.BatchV1().Jobs(di.config.K8sNamespace).Watch(context.Background(), metav1.ListOptions{})
	}

	watcher, err := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})
	if err != nil {
		return err
	}

	for event := range watcher.ResultChan() {
		item := event.Object.(*batchv1.Job)

		switch event.Type {
		case watch.Modified:
			logrus.Infof("job modified: %s", item.GetName())
			if err := nc.Publish("knoperator.jobs.modified", []byte(item.GetName())); err != nil {
				return fmt.Errorf("error publishing modified job: %s", err)
			}

			// only if not deleted
			if item.Status.Active == 0 && item.Status.Succeeded == 0 && item.Status.Failed == 0 {
				if err := nc.Publish("knoperator.jobs.completed", []byte(item.GetName())); err != nil {
					return fmt.Errorf("error publishing modified job: %s", err)
				}
				if err := nc.Publish(fmt.Sprintf("knoperator.jobs.completed.%s", item.GetName()), []byte(item.GetName())); err != nil {
					return fmt.Errorf("error publishing modified job: %s", err)
				}
			}

			// checks if job is failed
			if item.Status.Failed == 1 {
				if err := nc.Publish("knoperator.jobs.failed", []byte(item.GetName())); err != nil {
					return fmt.Errorf("error publishing modified job: %s", err)
				}
			}

		case watch.Bookmark:
			logrus.Infof("job bookmark: %s", item.GetName())
			if err := nc.Publish("knoperator.jobs.bookmark", []byte(item.GetName())); err != nil {
				return fmt.Errorf("error publishing modified job: %s", err)
			}
		case watch.Error:
			logrus.Infof("job error: %s", item.GetName())
			if err := nc.Publish("knoperator.jobs.error", []byte(item.GetName())); err != nil {
				return fmt.Errorf("error publishing modified job: %s", err)
			}
		case watch.Deleted:
			logrus.Infof("job deleted: %s", item.GetName())
			if err := nc.Publish("knoperator.jobs.deleted", []byte(item.GetName())); err != nil {
				return fmt.Errorf("error publishing modified job: %s", err)
			}
		case watch.Added:
			logrus.Infof("job added: %s", item.GetName())
			if err := nc.Publish("knoperator.jobs.added", []byte(item.GetName())); err != nil {
				return fmt.Errorf("error publishing modified job: %s", err)
			}
		}

		type WatchMsg struct {
			JobName string
			Type    watch.EventType
		}

		msg := WatchMsg{
			JobName: item.GetName(),
			Type:    event.Type,
		}

		msgJSON, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("error marshalling watch message: %s", err)
		}

		if err := nc.Publish("knoperator.jobs.watch", msgJSON); err != nil {
			return fmt.Errorf("error publishing watch message: %s", err)
		}
	}

	return nil
}

func (di *DeploymentIntegration) CreateJob(name, image string, command, args []string, env map[string]string, mountpoints map[string]string) error {
	volumes := []corev1.Volume{
		{
			Name: "default-env",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: "default-env",
				},
			},
		},
		// {
		// 	Name: "ssh",
		// 	VolumeSource: corev1.VolumeSource{
		// 		Secret: &corev1.SecretVolumeSource{
		// 			SecretName: "ssh",
		// 		},
		// 	},
		// },
	}

	volumeMounts := []corev1.VolumeMount{
		// {
		// 	Name:      "ssh",
		// 	MountPath: "/root/.ssh",
		// 	ReadOnly:  true,
		// 	SubPath:   "id_rsa",
		// },
	}

	for k, v := range mountpoints {
		kmd5 := fmt.Sprintf("%x", md5.Sum([]byte(k)))

		path := path.Join(di.config.BaseHostPath, k)

		// // create the host path if it does not exist (would only work when running in a pod)
		// if err := os.MkdirAll(path, 0777); err != nil {
		// 	if !os.IsExist(err) {
		// 		return err
		// 	}
		// }

		hostPathType := corev1.HostPathDirectoryOrCreate
		// recursiveReadOnly := corev1.RecursiveReadOnlyDisabled // only for newer k8s versions
		mountPropagation := corev1.MountPropagationHostToContainer

		volumes = append(volumes, corev1.Volume{
			Name: kmd5,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: path,
					Type: &hostPathType,
				},
			},
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      kmd5,
			MountPath: v,
			ReadOnly:  false,
			// RecursiveReadOnly: &recursiveReadOnly, // only for newer k8s versions
			MountPropagation: &mountPropagation,
		})
	}

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
							VolumeMounts: volumeMounts,
						},
					},
					Volumes: volumes,
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
