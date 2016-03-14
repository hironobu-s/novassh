package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
	"github.com/rackspace/gophercloud/pagination"
)

const (
	CREDENTIAL_FILE      = ".novassh"
	CREDENTIAL_FILE_MODE = 0600
)

type machine struct {
	Name   string
	Ipaddr string
}

func newMachine(s servers.Server) (*machine, error) {
	m := &machine{
		Name: s.Name,
	}

	// For ConoHa
	for key, value := range s.Metadata {
		if key == "instance_name_tag" {
			m.Name = value.(string)
		}
	}

	// Detecting public IP address for connectiong SSH
	for name, addressSet := range s.Addresses {
		if name == "public" {
			// Response of rackspace API has the key either "public" or "private"

			// TODO: Rackspace API and other OpenStack Systems

		} else if strings.HasPrefix(name, "ext-") {
			// Response of ConoHa API has the prefix either "ext-" or "local-"
			as, ok := addressSet.([]interface{})
			if !ok {
				return nil, fmt.Errorf("Invalid address set(type assertion failed).")
			}

			for _, v := range as {
				addr, ok := v.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("Invalid address set(type assertion failed).")
				}

				// TODO: support IPv6 address
				if addr["version"] == 4.0 {
					m.Ipaddr = addr["addr"].(string)
				}
			}
		}
	}

	return m, nil
}

type nova struct {
	machines []*machine
	provider *gophercloud.ServiceClient
}

type Credential struct {
	ComputeEndpoint string
	*tokens.Token
}

func NewNova() *nova {
	nova := &nova{
		machines: nil,
	}
	return nova
}

func (n *nova) Init() (err error) {
	// Credentials from env
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return err
	}

	// Endpoint options
	eo := gophercloud.EndpointOpts{
		Type:   "compute",
		Region: os.Getenv("OS_REGION_NAME"),
	}

	// Try to use cached credential if file exists.
	_, err = os.Stat(n.credentialCachePath())
	if err == nil {
		// Use cache file
		strdata, err := ioutil.ReadFile(n.credentialCachePath())
		if err != nil {
			log.Warnf("Failed to load cache file: %v", err)
			goto AUTH
		}

		cred := &Credential{Token: &tokens.Token{}}
		if err = json.Unmarshal(strdata, cred); err != nil {
			log.Warnf("Failed to unmarchal the cache data: %v", err)
			goto AUTH
		}

		if time.Now().Before(cred.ExpiresAt) {
			log.Debugf("Token has expired. Try to reauth.")
			goto AUTH
		}

		client, err := openstack.NewClient(opts.IdentityEndpoint)
		client.EndpointLocator = func(o gophercloud.EndpointOpts) (string, error) {
			return cred.ComputeEndpoint, nil
		}
		client.TokenID = cred.ID

		// Set service client
		n.provider, err = openstack.NewComputeV2(client, eo)
		return nil
	}

AUTH:

	client, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return err
	}

	// Set service client
	n.provider, err = openstack.NewComputeV2(client, eo)

	// Store credential to cache file
	cred := &Credential{
		ComputeEndpoint: n.provider.Endpoint,
		Token: &tokens.Token{
			ID:        client.TokenID,
			ExpiresAt: time.Now(),
		},
	}
	strdata, err := json.Marshal(cred)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(n.credentialCachePath(), strdata, CREDENTIAL_FILE_MODE); err != nil {
		log.Warnf("Can not write the credential cache file: %s", n.credentialCachePath())
	}

	return nil
}

func (n *nova) Find(name string) (m *machine, err error) {
	if n.machines == nil {
		n.machines, err = n.listMachines()
		if err != nil {
			return nil, err
		}
	}

	for _, machine := range n.machines {
		if machine.Name == name {
			return machine, nil
		}
	}

	return nil, nil
}

func (n *nova) listMachines() ([]*machine, error) {
	pager := servers.List(n.provider, servers.ListOpts{})

	machines := []*machine{}
	pager.EachPage(func(page pagination.Page) (bool, error) {
		ss, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}

		for _, s := range ss {
			m, err := newMachine(s)
			if err != nil {
				return false, err
			}
			log.Debugf("Machine is found: name=%s, ipaddr=%s", m.Name, m.Ipaddr)

			machines = append(machines, m)
		}

		return true, nil
	})

	return machines, nil
}

func (n *nova) credentialCachePath() string {
	d, err := homedir.Dir()
	if err == nil {
		return fmt.Sprintf("%s%c%s", d, filepath.Separator, CREDENTIAL_FILE)
	} else {
		return CREDENTIAL_FILE
	}
}
