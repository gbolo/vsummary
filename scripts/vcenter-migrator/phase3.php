<?php

/* ======================================= 
	START LOGIC
======================================= */

// Import Config
require_once('migration_config.php');

// Import functions
require_once('migration_functions.php');


// STEP 0: RUN VSUMMARY BEFORE RUNNING THIS SCRIPT!!
$file_content = "Write-Host !!NOTICE!! PHASE 3 --> BEFORE RUNNING THIS SCRIPT:\n";
$file_content .= "Write-Host   - DISCONNECT EACH ESXI HOST FROM THE SOURCE VCENTER CLUSTER\n";
$file_content .= "Write-Host   - CONNECT EACH ESXI HOST TO THE DESTINATION VCENTER CLUSTER\n";
$file_content .= "Write-Host   - CONNECT EACH ESXI HOST TO THE DVS IN THE DESTINATION VCENTER\n";
$file_content .= "Write-Host   - RUN VSUMMARY COLLECTOR AFTER ALL THE ABOVE IS DONE SUCCESSFULLY\n";
$file_content .= "Write-Host Press any key to continue ...\n";
$file_content .= '$x = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")'."\n";

// STEP 1: MOVE VMS BACK TO DVS
$VM_ARRAY = gen_vm_array($destination_vcenter_id, $destination_datacenter_name, $destination_cluster_id);
$VM_SOURCE_ARRAY = json_decode(file_get_contents('PHASE2_VM_AND_TEMPLATES_FULL.json'), true);
$VNICS_CHANGED = json_decode(file_get_contents('PHASE2_VNICS_CHANGED.json'), true);

$file_content .= "Connect-vcenter $destination_vcenter_fqdn destination\n";
$file_content .= powercli_restore_vm_vnics($VM_ARRAY, $VNICS_CHANGED);


// STEP 2: MOVE VMS BACK TO FOLDERS
$file_content .= powercli_restore_vm_folders($VM_SOURCE_ARRAY, $VM_ARRAY);

// STEP 3: EXPORT AND IMPORT VAPPS
$file_content .= "New-Item -ItemType Directory -Force vapps\n";
$file_content .= "Connect-vcenter $source_vcenter_fqdn source\n";
$file_content .= powercli_export_vapps($source_vcenter_id, $source_cluster_id);

$file_content .= "Connect-vcenter $destination_vcenter_fqdn destination\n";
$file_content .= powercli_import_vapps($source_cluster_name);

// STEP 4: IMPORT RESOURCE POOL
$file_content .= powercli_import_resourcepools($source_vcenter_id, $source_cluster_id);




echo $file_content;
/*  INSTRUCTIONS:

1 - run EXPORT_VAPPS
2 - run IMPORT_VAPPS
3 - run RESTORE-VM-FOLDERS.ps1
4 - run RESTORE-VM-PORTGROUPS.ps1
5 - run IMPORT-RESOURCEPOOLS.ps1
6 - run vsummary

*/