package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
	"github.com/rackspace/gophercloud/pagination"
	log "github.com/sirupsen/logrus"
)

const (
	CREDENTIAL_FILE      = ".novassh"
	CREDENTIAL_FILE_MODE = 0600
)

type machine struct {
	Name   string
	Ipaddr string
	Uuid   string
}

func newMachine(s servers.Server, interfaceName string) (*machine, error) {
	m := &machine{
		Name: s.Name,
		Uuid: s.ID,
	}

	// For ConoHa
	for key, value := range s.Metadata {
		if key == "instance_name_tag" {
			m.Name = value.(string)
		}
	}

	//
	// Detect public IP address for connectiong SSH
	//

	// Try to detect any public IP addresses
	for name, addressSet := range s.Addresses {
		// Response of Rackspace API has the key either "public" or "private".
		// Response of ConoHa API has the prefix either "ext-" or "local-"
		if name == "public" || strings.HasPrefix(name, "ext-") || name == interfaceName {
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
					goto DETECTED
				}
			}
		}
	}

	// Try to detect any private floating IP address.
	for name, addressSet := range s.Addresses {
		if name == "private" {
			as, ok := addressSet.([]interface{})
			if !ok {
				return nil, fmt.Errorf("Invalid address set(type assertion failed).")
			}

			for _, v := range as {
				addr, ok := v.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("Invalid address set(type assertion failed).")
				}

				if addr["OS-EXT-IPS:type"] == "floating" {
					m.Ipaddr = addr["addr"].(string)
					goto DETECTED
				}
			}
		}
	}

	// Choose the first one when we can not detect.
	for _, addressSet := range s.Addresses {
		as, ok := addressSet.([]interface{})
		if !ok {
			return nil, fmt.Errorf("Invalid address set(type assertion failed).")
		}

		addr, ok := as[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Invalid address set(type assertion failed).")
		}
		m.Ipaddr = addr["addr"].(string)
		goto DETECTED
	}

DETECTED:

	return m, nil
}

type nova struct {
	machines         []*machine
	provider         *gophercloud.ServiceClient
	networkInterface string
}

type Credential struct {
	ComputeEndpoint string
	*tokens.Token
}

func NewNova(nicname string) *nova {
	nova := &nova{
		networkInterface: nicname,
		machines:         nil,
	}
	return nova
}

func (n *nova) Init(authcache bool) (err error) {
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

		if cred.ExpiresAt.After(time.Now()) {
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
	if n.provider == nil {
		return fmt.Errorf("Could not found the Compute endpoint")
	}

	// Store credential to cache file
	if authcache {
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
	}
	return nil
}

func (n *nova) Find(name string) (m *machine, err error) {
	if n.machines == nil {
		n.machines, err = n.List()
		if err != nil {
			return nil, err
		}
	}

	for _, machine := range n.machines {
		if strings.ToLower(machine.Name) == strings.ToLower(name) {
			return machine, nil
		}
	}

	return nil, nil
}

func (n *nova) List() ([]*machine, error) {
	pager := servers.List(n.provider, servers.ListOpts{})

	machines := []*machine{}
	pager.EachPage(func(page pagination.Page) (bool, error) {
		ss, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}

		for _, s := range ss {
			m, err := newMachine(s, n.networkInterface)
			if err != nil {
				return false, err
			}
			log.Debugf("Machine found: name=%s, ipaddr=%s", m.Name, m.Ipaddr)
			for name, _ := range s.Addresses {
				log.Debugf("InterfaceName: %s", name)
			}

			machines = append(machines, m)
		}

		return true, nil
	})

	return machines, nil
}

func (n *nova) GetConsoleUrl(m *machine) (string, error) {
	data, err := json.Marshal(map[string]interface{}{
		"os-getSerialConsole": map[string]string{
			"type": "serial",
		},
	})
	if err != nil {
		return "", err
	}

	url := n.provider.ServiceURL("servers", m.Uuid, "action")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range n.provider.AuthenticatedHeaders() {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response struct {
		Console struct {
			Type string `json:"type"`
			Url  string `json:"url"`
		} `json:"console"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	return response.Console.Url, nil
}

func (n *nova) credentialCachePath() string {
	d, err := homedir.Dir()
	if err == nil {
		return fmt.Sprintf("%s%c%s", d, filepath.Separator, CREDENTIAL_FILE)
	} else {
		return CREDENTIAL_FILE
	}
}

func (n *nova) RemoveCredentialCache() error {
	_, err := os.Stat(n.credentialCachePath())
	if err == nil {
		return os.Remove(n.credentialCachePath())
	} else {
		return nil
	}
}
