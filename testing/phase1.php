<?php

/* ======================================= 
	DEFINE VARIABLES
======================================= */

$outut_dir = 'output/phase1/';

// MIGRATION VARIABLES
$source_vcenter_id = '0184679d-369a-4590-993a-5fbdf326a75a';
$source_datacenter_name = 'DC1';
$source_cluster_id = '34677d70de6db190a6a5458a553eb27f';

// DVS PORTGROUP EXPORT
$source_dvs_id = '5fb6de4be73d4154d746ca485eec9dae';
$destination_dvs_name = 'DVS2';

/* ======================================= 
	START LOGIC
======================================= */

// Import functions
require_once('migration_functions.php');


// STEP 0: STOP DRS

// STEP 1: EXPORT ALL VM FOLDERS IN GIVEN DC AND CREATE IMPORT PS1 SCRIPT
$FOLDER_ARRAY = gen_folder_array($source_vcenter_id, $source_datacenter_name);
$filename = $outut_dir . 'IMPORT_VM_FOLDERS.ps1';
powercli_import_vm_folders($FOLDER_ARRAY , $filename);

// STEP 2: EXPORT DVS PORTGROUPS
$DVS_PG_ARRAY = gen_dvs_pg_array($source_dvs_id);
$filename = $outut_dir . 'IMPORT_DVS.ps1';
powercli_import_dvs($DVS_PG_ARRAY, $filename, $destination_dvs_name);

// STEP 3: EXPORT FULL LIST OF VMS AND TEMPLATES IN CSV FORMAT
$VM_ARRAY = gen_vm_array($source_vcenter_id, $source_datacenter_name, $source_cluster_id);
$filename = $outut_dir . 'csv/VM_AND_TEMPLATES_FULL.csv';
csv_vm_list($VM_ARRAY, $filename);

// STEP 4: CONVERT TEMPLATES TO VMS ON EACH ESXI HOST
$ESXI_ARRAY = gen_esxi_array($source_vcenter_id, $source_datacenter_name, $source_cluster_id);
powercli_templates_to_vms($ESXI_ARRAY, $outut_dir);
