# This configuration file is used in conjunction with:
# https://github.com/gbolo/vsummary/blob/master/contrib/collectors/powershell/vsummary_collect.ps1


# ADD YOUR VSUMMARY-SERVER URL HERE:
$vsummary_url = 'http://127.0.0.1:8080'

# ADD YOUR VCENTER SERVER(S) LIKE THIS:
$vcenters = @{
  LAB = @{ fqdn = 'vcsa1.lab.linuxctl.com'; readonly_user = 'readonly@vsphere.local'; password = 'changeme'; };
  VDI = @{ fqdn = 'vcsa1.vdi.linuxctl.com'; readonly_user = 'ro@vsphere.local'; password = 'changeme'; };
  PROD = @{ fqdn = 'vcsa1.prod.linuxctl.com'; readonly_user = 'ro@vsphere.local'; password = 'changeme'; };
  DR = @{ fqdn = 'vcsa1.dr.linuxctl.com'; readonly_user = 'ro@vsphere.local'; password = 'changeme'; };
  VCSIM = @{ fqdn = '127.0.0.1'; port = 8989; readonly_user = 'user'; password = 'pass'; };
}
