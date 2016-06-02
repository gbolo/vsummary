<?php

/* ======================================= 
	DEFINE VARIABLES
======================================= */

$outut_dir = 'output/phase3/';

// MIGRATION VARIABLES
$source_vcenter_id = '0184679d-369a-4590-993a-5fbdf326a75a';
$source_datacenter_name = 'DC1';
$source_cluster_id = '34677d70de6db190a6a5458a553eb27f';

// VSWITCH GENERATE
$vswitch_name = 'vSwitch1';

// DVS PORTGROUP EXPORT
$source_dvs_id = '5fb6de4be73d4154d746ca485eec9dae';
$destination_dvs_name = 'DVS2';

/* ======================================= 
	START LOGIC
======================================= */

// Import functions
require_once('migration_functions.php');


// STEP 0: RUN VSUMMARY BEFORE RUNNING THIS SCRIPT!!

// STEP 1: EXPORT ALL VM FOLDERS IN GIVEN DC AND CREATE IMPORT PS1 SCRIPT
$VM_ARRAY = gen_vm_array($source_vcenter_id, $source_datacenter_name, $source_cluster_id);
$filename = $outut_dir . 'csv/VM_AND_TEMPLATES_FULL.csv';
csv_vm_list($VM_ARRAY, $filename);

// STEP 2 & 3: CREATE PORTGROUPS ON STANDARD VSWITCH ON EACH ESXI HOST & MOVE ALL VM VNICS TO NEW STANDARTD VSWITCH PGS
$ESXI_ARRAY = gen_esxi_array($source_vcenter_id, $source_datacenter_name, $source_cluster_id);
$DVS_ARRAY = gen_dvs_pg_array($source_dvs_id);

powercli_move_vm_vnics($DVS_ARRAY, $ESXI_ARRAY, $vswitch_name, $outut_dir);



