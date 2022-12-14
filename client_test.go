package helvargo

import (
	"fmt"
	"path"
	"sync"
	"testing"

	"github.com/nuqz/chanfan"
	"github.com/nuqz/helvar-go/message"
	ht "github.com/nuqz/helvar-go/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	fakeSrvPort = 49876
)

func TestClient(t *testing.T) {
	defaultNet := ht.MustNetFromYAMLFile(path.Join("testing", "test_net.yml"))
	fakeSrv := ht.NewRouter(fmt.Sprintf(":%d", fakeSrvPort), defaultNet)
	require.NoError(t, fakeSrv.Listen())

	client := NewClient("localhost", fakeSrvPort)
	_, err := client.Transceive(message.NewQueryTime())
	assert.Error(t, err)

	errs, err := client.Connect(4, 4)
	require.NoError(t, err)

	go func() {
		for err := range chanfan.Merge(errs) {
			require.NoError(t, err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		clusters, err := client.GetClusters()
		require.NoError(t, err)

		for _, expected := range defaultNet.Clusters {
			found := false
			for _, actual := range clusters {
				found = actual.ID == expected.ID
				if found {
					break
				}
			}
			assert.True(t, found)
		}
		wg.Done()
	}()

	go func() {
		routers, err := client.GetRouters()
		require.NoError(t, err)

		for _, expected := range defaultNet.Routers {
			found := false
			for _, actual := range routers {
				found = actual.ID == expected.ID
				if found {
					break
				}
			}
			assert.True(t, found)
		}
		wg.Done()
	}()

	go func() {
		groups, err := client.GetGroups()
		require.NoError(t, err)

		for _, group := range groups {
			name, err := client.GetGroupName(group)
			require.NoError(t, err)

			group.Name = name

			devs, err := client.GetDevices(group)
			require.NoError(t, err)

			t.Logf("%+v", devs)

			for _, dev := range devs {
				name, err := client.GetDeviceName(dev)
				require.NoError(t, err)

				dev.Name = name

				t.Logf("GID: %d Name: %s / Address: %s Name %s",
					group.ID, group.Name, dev.Address, dev.Name)
			}
		}
		wg.Done()
	}()

	go func() {
		netTime, err := client.GetTime()
		require.NoError(t, err)

		t.Log(netTime)
		wg.Done()
	}()

	wg.Wait()
	client.Disconnect()

	_, err = client.Transceive(message.NewQueryTime())
	assert.Error(t, err)
}
