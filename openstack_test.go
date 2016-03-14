package main

import "testing"

func TestNewNova(t *testing.T) {
	n := NewNova()
	if n.servers != nil {
		t.Errorf("'servers' attribute is not nil")
	}
}

// func TestInitAndCache(t *testing.T) {
// 	var err error
// 	n := NewNova()

// 	// Remove credential cache file
// 	os.Remove(n.credentialCachePath())

// 	// Init
// 	if err = n.Init(); err != nil {
// 		t.Errorf("%v", err)
// 	}

// 	// Verify credential cache file
// 	_, err = os.Stat(n.credentialCachePath())
// 	if err != nil {
// 		t.Errorf("%v", err)
// 	}

// 	strdata, err := ioutil.ReadFile(n.credentialCachePath())
// 	cred := &Credential{Token: &tokens.Token{}}
// 	if err = json.Unmarshal(strdata, cred); err != nil {
// 		t.Errorf("%v", err)
// 	}
// }

// func TestInitAndCache2(t *testing.T) {
// 	// NOTE: Credential cache file has already created by previous test
// 	n := NewNova()

// 	// Init() uses it instead of authenticating
// 	if err := n.Init(); err != nil {
// 		t.Errorf("%v", err)
// 	}
// }

func TestFind(t *testing.T) {
	n := NewNova()
	n.Init()

	servers, err := n.listServers()
	if err != nil {
		t.Errorf("%v", err)

	} else if len(servers) == 0 {
		t.Skipf("Skip beause no servers found")
	}

	s, err := newMachine(servers[0])
	if err != nil {
		t.Errorf("%v", err)
	}

	ss, err := n.Find(s.Name)
	if err != nil {
		t.Errorf("%v", err)

	} else if ss == nil || ss.Name != s.Name {
		t.Errorf("Find() did not return the server: name=%s", s.Name)
	}
}
