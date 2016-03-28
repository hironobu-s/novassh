package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
)

func TestNewNova(t *testing.T) {
	n := NewNova()
	if n.machines != nil {
		t.Errorf("'servers' attribute is not nil")
	}
}

func TestInitAndCache(t *testing.T) {
	var err error
	n := NewNova()

	// Remove credential cache file
	os.Remove(n.credentialCachePath())

	// Init
	if err = n.Init(true); err != nil {
		t.Errorf("%v", err)
	}

	// Verify credential cache file
	_, err = os.Stat(n.credentialCachePath())
	if err != nil {
		t.Errorf("%v", err)
	}

	strdata, err := ioutil.ReadFile(n.credentialCachePath())
	cred := &Credential{Token: &tokens.Token{}}
	if err = json.Unmarshal(strdata, cred); err != nil {
		t.Errorf("%v", err)
	}
}

func TestInitAndCache2(t *testing.T) {
	// NOTE: Credential cache file has already created by previous test
	n := NewNova()

	// Init() uses it instead of authenticating
	if err := n.Init(false); err != nil {
		t.Errorf("%v", err)
	}
}

func TestFind(t *testing.T) {
	n := NewNova()
	n.Init(false)

	machines, err := n.List()
	if err != nil {
		t.Errorf("%v", err)

	} else if len(machines) == 0 {
		t.Skipf("Skip beause no servers found")
	}

	ss, err := n.Find(machines[0].Name)
	if err != nil {
		t.Errorf("%v", err)

	} else if ss == nil || ss.Name != machines[0].Name {
		t.Errorf("Find() did not return the server: name=%s", machines[0].Name)
	}
}

func TestFind2(t *testing.T) {
	n := NewNova()
	n.Init(false)

	ss, err := n.Find("undefinded-instance-name")
	if err != nil {
		t.Errorf("%v", err)

	} else if ss != nil {
		t.Errorf("Find() should return 'nil'")
	}
}
