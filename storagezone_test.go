//go:build integrationtest
// +build integrationtest

package bunny_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nrdcg/bunny-go"
)

func TestStorageZoneCRUD(t *testing.T) {
	clt := newClient(t)

	szOrigin := "http://bunny.net"

	szAddopts := bunny.StorageZoneAddOptions{
		Name:               pointer(randomResourceName("storagezone")),
		OriginURL:          pointer(szOrigin),
		Region:             pointer("NY"),
		ReplicationRegions: []string{"DE"},
	}

	listSzBefore, err := clt.StorageZone.List(context.Background(), nil)
	require.NoError(t, err, "storage zone list failed before add")

	sz := createStorageZone(t, clt, &szAddopts)

	// get the newly created storage zone
	getSz, err := clt.StorageZone.Get(context.Background(), deref(sz.ID))
	require.NoError(t, err, "storage zone get failed after adding")
	assert.NotNil(t, getSz.ID)
	assert.Equal(t, getSz.ReplicationRegions[0], "DE", "storage zone replication region should be set correctly")

	// update the storage zone
	updateOpts := bunny.StorageZoneUpdateOptions{
		OriginURL:          pointer(szOrigin + "/updated"),
		Rewrite404To200:    pointer(true),
		ReplicationRegions: []string{"LA"},
	}
	updateErr := clt.StorageZone.Update(context.Background(), deref(sz.ID), &updateOpts)
	assert.Nil(t, updateErr)

	// get the updated storage zone and validate updated properties
	getUpdatedSz, err := clt.StorageZone.Get(context.Background(), deref(sz.ID))
	assert.NotNil(t, getUpdatedSz.ID)
	assert.Equal(t, "LA", getUpdatedSz.ReplicationRegions[len(getUpdatedSz.ReplicationRegions)-1], "storage zone replication region should be updated correctly")

	// check the total number of storage zones is the expected amount
	listSzAfter, err := clt.StorageZone.List(context.Background(), nil)
	require.NoError(t, err, "storage zone list failed after add")
	assert.Equal(t, deref(listSzBefore.TotalItems+1), deref(listSzAfter.TotalItems), "storage zones total items should increase by exactly 1")
}
