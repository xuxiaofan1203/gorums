package dev_test

import (
	"testing"
	"time"

	rpc "github.com/relab/gorums/dev"
	"github.com/relab/gorums/idutil"
)

func TestEqualGlobalConfigurationIDs(t *testing.T) {
	// Equal set of addresses, but different order.
	addrsOne := []string{"localhost:8080", "localhost:8081", "localhost:8082"}
	addrsTwo := []string{"localhost:8081", "localhost:8082", "localhost:8080"}

	mgrOne, err := rpc.NewManager(addrsOne, rpc.WithNoConnect())
	if err != nil {
		t.Fatal(err)
	}

	mgrTwo, err := rpc.NewManager(addrsTwo, rpc.WithNoConnect())
	if err != nil {
		t.Fatal(err)
	}

	// Create a configuration in each manager using all nodes.
	// Global ids should be equal.

	configOne, err := mgrOne.NewConfiguration(mgrOne.NodeIDs(), 1, 1)
	if err != nil {
		t.Fatalf("error creating config one: %v", err)
	}

	configTwo, err := mgrTwo.NewConfiguration(mgrTwo.NodeIDs(), 1, 1)
	if err != nil {
		t.Fatalf("error creating config two: %v", err)
	}

	if configOne.GlobalID() != configTwo.GlobalID() {
		t.Errorf("global configuration ids differ, %d != %d",
			configOne.GlobalID(), configTwo.GlobalID())
	}
}

func TestWithSelfAddrOption(t *testing.T) {
	addrs := []string{"localhost:8080", "localhost:8081", "localhost:8082"}
	selfAddr := "localhost:8081"
	wantSize := len(addrs)

	mgr, err := rpc.NewManager(addrs, rpc.WithNoConnect(), rpc.WithSelfAddr(selfAddr))
	if err != nil {
		t.Fatal(err)
	}

	gotSize, _ := mgr.Size()
	if gotSize != wantSize {
		t.Errorf("got manager node size %d, want %d", gotSize, wantSize)
	}

	ids := mgr.NodeIDs()
	if len(ids) != wantSize {
		t.Errorf("got %d node ids from manager, want %d", len(ids), wantSize)
	}

	gids := mgr.NodeGlobalIDs()
	if len(gids) != wantSize {
		t.Errorf("got %d node ids from manager, want %d", len(ids), wantSize)
	}

	nodes := mgr.Nodes(false)
	if len(nodes) != wantSize {
		t.Errorf("got %d nodes from manager, want %d", len(nodes), wantSize)
	}

	nodes = mgr.Nodes(true)
	if len(nodes) != wantSize-1 {
		t.Errorf("got %d nodes from manager, want %d", len(nodes), wantSize-1)
	}

	notPresentAddr := "localhost:8083"
	_, err = rpc.NewManager(addrs, rpc.WithNoConnect(), rpc.WithSelfAddr(notPresentAddr))
	if err == nil {
		t.Error("got no manager creation error, want error due to invaild WithSelf option")
	}
}

func TestWithSelfGidOption(t *testing.T) {
	addrs := []string{"localhost:8080", "localhost:8081", "localhost:8082"}
	selfAddr := "localhost:8081"
	selfGid, err := idutil.IDFromAddress(selfAddr)
	if err != nil {
		t.Fatal(err)
	}
	wantSize := len(addrs)

	mgr, err := rpc.NewManager(addrs, rpc.WithNoConnect(), rpc.WithSelfGid(selfGid))
	if err != nil {
		t.Fatal(err)
	}

	gotSize, _ := mgr.Size()
	if gotSize != wantSize {
		t.Errorf("got manager node size %d, want %d", gotSize, wantSize)
	}

	ids := mgr.NodeIDs()
	if len(ids) != wantSize {
		t.Errorf("got %d node ids from manager, want %d", len(ids), wantSize)
	}

	gids := mgr.NodeGlobalIDs()
	if len(gids) != wantSize {
		t.Errorf("got %d node gids from manager, want %d", len(gids), wantSize)
	}

	nodes := mgr.Nodes(false)
	if len(nodes) != wantSize {
		t.Errorf("got %d nodes from manager, want %d", len(nodes), wantSize)
	}

	nodes = mgr.Nodes(true)
	if len(nodes) != wantSize-1 {
		t.Errorf("got %d nodes from manager, want %d", len(nodes), wantSize-1)
	}

	var notPresentGid uint32 = 42
	_, err = rpc.NewManager(addrs, rpc.WithNoConnect(), rpc.WithSelfGid(notPresentGid))
	if err == nil {
		t.Error("got no manager creation error, want error due to invaild WithSelfGlobalID option")
	}
}

func TestCreateConfiguration(t *testing.T) {
	addrs := []string{"localhost:8080", "localhost:8081", "localhost:8082"}
	mgr, err := rpc.NewManager(addrs, rpc.WithNoConnect())
	if err != nil {
		t.Fatal(err)
	}
	ids := mgr.NodeIDs()
	config, err := mgr.NewConfiguration(ids, len(ids), time.Second)
	if err != nil {
		t.Errorf("got error creating configuration, want none (%v)", err)
	}

	cids := config.Nodes()
	if !equal(cids, ids) {
		t.Errorf("ids from Manager (got %v) and ids from configuration containing all nodes (got %v) should be equal",
			ids, cids)
	}

	_, size := mgr.Size()
	if size != 1 {
		t.Errorf("got #%d configurations from Manager, want %d", size, 1)
	}
}

func TestCreateConfiguratonWithSelfOption(t *testing.T) {
	addrs := []string{"localhost:8080", "localhost:8081", "localhost:8082"}
	selfAddr := "localhost:8081"
	selfGid, err := idutil.IDFromAddress(selfAddr)
	if err != nil {
		t.Fatal(err)
	}

	mgr, err := rpc.NewManager(addrs, rpc.WithNoConnect(), rpc.WithSelfGid(selfGid))
	if err != nil {
		t.Fatal(err)
	}

	ids := mgr.NodeIDs()
	_, err = mgr.NewConfiguration(ids, len(ids), time.Second)
	if err == nil {
		t.Error("expected error creating configuration with self, got none")
	}

	nodes := mgr.Nodes(true)
	var nids []int
	for _, node := range nodes {
		nids = append(nids, node.ID())
	}

	config, err := mgr.NewConfiguration(nids, len(nids), time.Second)
	if err != nil {
		t.Errorf("got error creating configuration, want none (%v)", err)
	}

	cids := config.Nodes()
	if !equal(cids, nids) {
		t.Errorf("ids from Manager (got %v) and ids from configuration containing all nodes (got %v) should be equal",
			ids, cids)
	}

	_, size := mgr.Size()
	if size != 1 {
		t.Errorf("got #%d configurations from Manager, want %d", size, 1)
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}
