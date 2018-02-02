<?php


/* ======================================= 
	START LOGIC
======================================= */

// Import Config
require_once('migration_config.php');

// Import functions
require_once('migration_functions.php');

$file_content = "Write-Host !!NOTICE!! PHASE 4 --> BEFORE RUNNING THIS SCRIPT:\n";
$file_content .= "Write-Host   - MAKE VSUMMARY INVENTORY HAS BEEN UPDATED BEFORE PROCEEDING\n";
$file_content .= "Write-Host   - ONCE THIS PHASE 4 SCRIPT IS COMPLETED, RE-ENABLE DRS AND RESTORE NETWORK REDUNDANCY\n";
$file_content .= "Write-Host Press any key to continue ...\n";
$file_content .= '$x = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")'."\n";

// STEP 0: RUN VSUMMARY BEFORE RUNNING THIS SCRIPT!!

// STEP 1: MOVE VMS BACK TO RESOURCEPOOLS AND VAPPS
$VM_ARRAY = gen_vm_array($destination_vcenter_id, $destination_datacenter_name, $destination_cluster_id);
$VM_SOURCE_ARRAY = json_decode(file_get_contents('PHASE1_VM_AND_TEMPLATES_FULL.json'), true);
$RESOURCEPOOL_ARRAY = gen_rpool_array($destination_cluster_id);
$VAPP_ARRAY = gen_vapp_array($destination_cluster_id);

$file_content = "Connect-vcenter $destination_vcenter_fqdn destination\n";
$file_content .= powercli_restore_vm_rpools($VM_ARRAY, $VM_SOURCE_ARRAY, $RESOURCEPOOL_ARRAY, $VAPP_ARRAY);

// STEP 3: CONVERT VMS BACK TO TEMPLATES
$file_content .= powercli_restore_vm_templates($VM_ARRAY, $VM_SOURCE_ARRAY);



echo $file_content;

/*  INSTRUCTIONS:

1 - run RESTORE-VM-POOLS-VAPPS.ps1
2 - run RESTORE-TEMPLATES.ps1
3 - restore redundancy to DVS and remove temp vswitch
4 - run vsummary again :)

*/