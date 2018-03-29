# This configuration file is used in conjunction with:
# https://github.com/gbolo/vSummary/blob/master/collectors/powershell/vsummary_collect.ps1


# ADD YOUR vSUMMARY API ENDPOINT HERE:
$vsummary_url = 'http://vsummary.linuxctl.com/api/update.php'

# ADD YOUR VCENTER SERVERS LIKE THIS:
$vcenters = @{
    LAB = @{ fqdn = 'vcsa1.lab.linuxctl.com'; readonly_user = 'readonly@vsphere.local'; password = 'changeme'; };
    VDI = @{ fqdn = 'vcsa1.vdi.linuxctl.com'; readonly_user = 'ro@vsphere.local'; password = 'changeme'; }; 
    PROD = @{ fqdn = 'vcsa1.prod.linuxctl.com'; readonly_user = 'ro@vsphere.local'; password = 'changeme'; }; 
    DR = @{ fqdn = 'vcsa1.dr.linuxctl.com'; readonly_user = 'ro@vsphere.local'; password = 'changeme'; }; 
}