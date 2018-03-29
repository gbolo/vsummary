<?php

/* ======================================= 
	START LOGIC
======================================= */

// Import functions
require_once('migration_functions.php');
require_once('migration_config.php');

// STEP 0: RUN VSUMMARY BEFORE RUNNING THIS SCRIPT!!
$file_content = "Write-Host !!NOTICE!! PHASE 2 --> BEFORE RUNNING THIS SCRIPT:\n";
$file_content .= "Write-Host   - MAKE VSUMMARY INVENTORY HAS BEEN UPDATED BEFORE PROCEEDING\n";
$file_content .= "Write-Host   - ONCE THIS PHASE 2 SCRIPT IS COMPLETED, YOU MUST DISCONNECT EACH ESXI HOST FROM SOURCE VCENTER\n";
$file_content .= "Write-Host Press any key to continue ...\n";
$file_content .= '$x = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")'."\n";


// STEP 1: EXPORT ALL VM FOLDERS IN GIVEN DC AND CREATE IMPORT PS1 SCRIPT
$VM_ARRAY = gen_vm_array($source_vcenter_id, $source_datacenter_name, $source_cluster_id);
$filename = 'PHASE2_VM_AND_TEMPLATES_FULL';
csv_vm_list($VM_ARRAY, $filename);

// STEP 2 & 3: CREATE PORTGROUPS ON STANDARD VSWITCH ON EACH ESXI HOST & MOVE ALL VM VNICS TO NEW STANDARTD VSWITCH PGS
$ESXI_ARRAY = gen_esxi_array($source_vcenter_id, $source_datacenter_name, $source_cluster_id);
$DVS_ARRAY = gen_dvs_pg_array($source_dvs_id);

$file_content .= '
Write-Host "CREATE TEMP STANDARD VSWITCH AND MOVE ALL VMS TO IT `n
FROM: '.$source_vcenter_fqdn.' `n"

Connect-vcenter '.$source_vcenter_fqdn.' "source"
';

$output = powercli_move_vm_vnics($DVS_ARRAY, $ESXI_ARRAY, $vswitch_name, 'PHASE2');

$file_content .= $output['vswitch'];
$file_content .= $output['vnic'];

echo $file_content;
