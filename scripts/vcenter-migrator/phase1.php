<?php
/* ======================================= 
	START LOGIC
======================================= */

// Import functions
require_once('migration_functions.php');
require_once('migration_config.php');

$file_content = "Write-Host !!NOTICE!! PHASE 1 --> BEFORE RUNNING THIS SCRIPT:\n";
$file_content .= "Write-Host   - SET DRS IN THE SOURCE CLUSTER TO MANUAL\n";
$file_content .= "Write-Host   - MAKE SURE AN EMPTY STANDARD VSWITCH IS CREATED ON EACH HOST IN THE CLUSTER\n";
$file_content .= "Write-Host   - MAKE VSUMMARY INVENTORY HAS BEEN UPDATED AFTER COMPLETING ABOVE STEPS\n";
$file_content .= "Write-Host Press any key to continue ...\n";
$file_content .= '$x = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")'."\n";


$file_content .= "
Write-Host 'Importing vsummary migration module'
Import-Module .\\vSummaryMigratorModule.psm1;

Write-Host 'Please ENTER Credentials for SOURCE vCenter: $source_vcenter_fqdn'
Get-Credential | Export-Clixml .\source_vcenter_creds.xml
Write-Host 'Please ENTER Credentials for DESTINATION vCenter: $destination_vcenter_fqdn'
Get-Credential | Export-Clixml .\destination_vcenter_creds.xml



";



// STEP 0: STOP DRS

// STEP 1: EXPORT ALL VM FOLDERS IN GIVEN DC AND CREATE IMPORT PS1 SCRIPT
$FOLDER_ARRAY = gen_folder_array($source_vcenter_id, $source_datacenter_name);
$file_content .= '
Write-Host "IMPORTING FOLDERS `n
FROM: '.$source_vcenter_fqdn.' `n
TO: '.$destination_vcenter_fqdn.' `n"

Connect-vcenter '.$destination_vcenter_fqdn.' "destination"
';
$file_content .= powercli_import_vm_folders($FOLDER_ARRAY);



// STEP 2: EXPORT DVS PORTGROUPS
$DVS_PG_ARRAY = gen_dvs_pg_array($source_dvs_id);
$file_content .= '
Write-Host "IMPORTING DVS `n
FROM: '.$source_vcenter_fqdn.' `n
TO: '.$destination_vcenter_fqdn.' `n"

Connect-vcenter '.$destination_vcenter_fqdn.' "destination"
';

$file_content .= powercli_import_dvs($DVS_PG_ARRAY, $destination_dvs_name);


// STEP 3: EXPORT FULL LIST OF VMS AND TEMPLATES IN CSV FORMAT
$VM_ARRAY = gen_vm_array($source_vcenter_id, $source_datacenter_name, $source_cluster_id);
$filename = 'PHASE1_VM_AND_TEMPLATES_FULL';
csv_vm_list($VM_ARRAY, $filename);

// STEP 4: CONVERT TEMPLATES TO VMS ON EACH ESXI HOST
$ESXI_ARRAY = gen_esxi_array($source_vcenter_id, $source_datacenter_name, $source_cluster_id);
$file_content .= '
Write-Host "CONVERTING VM TEMPLATES TO VMS `n
FROM: '.$source_vcenter_fqdn.' `n"

Connect-vcenter '.$source_vcenter_fqdn.' "source"
';
$file_content .= powercli_templates_to_vms($ESXI_ARRAY);

echo $file_content;

/*  INSTRUCTIONS:

1 - run convert-vm-to-template*.ps1 scripts
2 - RUN BOTH IMPORT* scripts in destination vcenter




*/


