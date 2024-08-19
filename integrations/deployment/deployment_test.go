package deployment

import (
	"testing"

	"github.com/mheers/knoperator/config"
	"github.com/mheers/knoperator/k8sclient"
	"github.com/stretchr/testify/require"
)

func TestGetPods(t *testing.T) {
	c := config.GetFakeConfig()

	k8s, err := k8sclient.Init(c)
	require.NoError(t, err)

	api, err := NewAPI(c, k8s)
	require.NoError(t, err)

	pods, err := api.GetPods()
	require.NoError(t, err)
	require.NotNil(t, pods)
}
