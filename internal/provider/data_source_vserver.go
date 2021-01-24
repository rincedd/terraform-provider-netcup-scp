package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rincedd/terraform-provider-netcup-scp/internal/client"
	"strings"
)

func dataSourceVServer() *schema.Resource {
	return &schema.Resource{
		Read: resourceServerRead,

		Schema: map[string]*schema.Schema{
			"server_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// computed properties
			"ipv4_addrs": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ipv6_addrs": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nickname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	loginName := m.(ProviderConfig).LoginName
	password := m.(ProviderConfig).Password
	serverName := d.Get("server_name").(string)

	wsClient := client.Client{LoginName: loginName, Password: password}
	serverInfo, err := wsClient.GetVServerInformation(serverName)
	if err != nil {
		return err
	}
	var ipv4s, ipv6s []string
	for _, ip := range serverInfo.IPs {
		if strings.ContainsRune(ip, ':') {
			ipv6s = append(ipv6s, ip)
		} else {
			ipv4s = append(ipv4s, ip)
		}
	}
	d.SetId(serverName)
	d.Set("ipv4_addrs", ipv4s)
	d.Set("ipv6_addrs", ipv6s)
	d.Set("state", serverInfo.Status)
	d.Set("nickname", serverInfo.Nickname)

	return nil
}
