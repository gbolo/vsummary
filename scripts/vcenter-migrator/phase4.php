<?php

$outut_dir = 'output/phase4/';

/* ======================================= 
	START LOGIC
======================================= */

// Import Config
require_once('migration_config.php');

// Import functions
require_once('migration_functions.php');


// STEP 0: RUN VSUMMARY BEFORE RUNNING THIS SCRIPT!!

// STEP 1: MOVE VMS BACK TO RESOURCEPOOLS AND VAPPS
$VM_ARRAY = gen_vm_array($destination_vcenter_id, $destination_datacenter_name, $destination_cluster_id);
$VM_SOURCE_ARRAY = json_decode(file_get_contents('output/phase2/csv/VM_AND_TEMPLATES_FULL.json'), true);
$RESOURCEPOOL_ARRAY = gen_rpool_array($destination_cluster_id);
$VAPP_ARRAY = gen_vapp_array($destination_cluster_id);
powercli_restore_vm_rpools($VM_ARRAY, $VM_SOURCE_ARRAY, $RESOURCEPOOL_ARRAY, $VAPP_ARRAY, $outut_dir);

// STEP 3: CONVERT VMS BACK TO TEMPLATES


/*  INSTRUCTIONS:

1 - run RESTORE-VM-POOLS-VAPPS.ps1
2 - run RESTORE-TEMPLATES.ps1
3 - restore redundancy to DVS and remove temp vswitch
4 - run vsummary again :)

*/