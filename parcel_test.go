package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (

	randSource = rand.NewSource(time.Now().UnixNano())

	randRange = rand.New(randSource)
)

type ParcelWithoutNumber struct {
	Client    int
	Status    string
	Address   string
	CreatedAt string
}

func toParcelWithoutNumber(p Parcel) ParcelWithoutNumber {
	return ParcelWithoutNumber{
		Client:    p.Client,
		Status:    p.Status,
		Address:   p.Address,
		CreatedAt: p.CreatedAt,
	}
}


func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()


	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	retrievedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, toParcelWithoutNumber(parcel), toParcelWithoutNumber(retrievedParcel))

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err)
}


func TestSetAddress(t *testing.T) {

	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	retrievedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, retrievedParcel.Address)
}


func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	
	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	
	retrievedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newStatus, retrievedParcel.Status)
}


func TestGetByClient(t *testing.T) {
	
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		
		parcels[i].Number = id

		
		parcelMap[id] = parcels[i]
	}

	
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))

	
	for _, parcel := range storedParcels {
		
		expectedParcel, exists := parcelMap[parcel.Number]
		require.True(t, exists)

		
		require.Equal(t, toParcelWithoutNumber(expectedParcel), toParcelWithoutNumber(parcel))

	}
	for _, parcel := range parcels {
		err = store.Delete(parcel.Number)
		require.NoError(t, err)
	}

}