<?php

/* ======================================= 
	DEFINE VARIABLES  -- vcsa1 -> vcsa2
======================================= */
/*

// MIGRATION VARIABLES
$source_vcenter_fqdn = 'vcsa1.lab1.midgar.local';
$destination_vcenter_fqdn = 'vcsa2.lab1.midgar.local';

$source_vcenter_id = '0184679d-369a-4590-993a-5fbdf326a75a';
$destination_vcenter_id = '38ff814b-216f-42cc-a382-681253b26c94';

$source_datacenter_name = 'DC1';
$destination_datacenter_name = 'DC1';

$source_cluster_id = '34677d70de6db190a6a5458a553eb27f';
$destination_cluster_id = '7c52b05daa29e7eb223d47a69acb1ec3';

$source_cluster_name = 'CL1';
$destination_cluster_name = 'CL1';

// DVS PORTGROUP EXPORT
$source_dvs_id = '5fb6de4be73d4154d746ca485eec9dae';
$destination_dvs_id = 'UPDATE ME AFTER PHASE1'; //update this after phase1
$destination_dvs_name = 'DVS1';

// VSWITCH GENERATE
$vswitch_name = 'vSwitch1';
*/

/* ======================================= 
	DEFINE VARIABLES  -- vcsa2 -> vcsa1
======================================= */
/*
// MIGRATION VARIABLES
$source_vcenter_fqdn = 'vcsa2.lab1.midgar.local';
$destination_vcenter_fqdn = 'vcsa1.lab1.midgar.local';

$source_vcenter_id = '38ff814b-216f-42cc-a382-681253b26c94';
$destination_vcenter_id = '0184679d-369a-4590-993a-5fbdf326a75a';

$source_datacenter_name = 'DC1';
$destination_datacenter_name = 'DC1';

$source_cluster_id = '7c52b05daa29e7eb223d47a69acb1ec3';
$destination_cluster_id = '34677d70de6db190a6a5458a553eb27f';

$source_cluster_name = 'CL1';
$destination_cluster_name = 'CL1';

// DVS PORTGROUP EXPORT
$source_dvs_id = 'ae0a388e01b202241df9dd193c463125';
//$destination_dvs_id = '5fb6de4be73d4154d746ca485eec9dae'; //update this after phase1
$destination_dvs_id = 'UPDATE ME AFTER PHASE1'; //update this after phase1
$destination_dvs_name = 'DVS1';

// VSWITCH GENERATE
$vswitch_name = 'vSwitch1';
*/


/* ======================================= 
	DEFINE VARIABLES  -- vcsa1 CL2 -> vcsa2
======================================= */


// MIGRATION VARIABLES
$source_vcenter_fqdn = 'vcsa1.lab1.midgar.local';
$destination_vcenter_fqdn = 'vcsa2.lab1.midgar.local';

$source_vcenter_id = '0184679d-369a-4590-993a-5fbdf326a75a';
$destination_vcenter_id = '38ff814b-216f-42cc-a382-681253b26c94';

$source_datacenter_name = 'DC1';
$destination_datacenter_name = 'DC1';

$source_cluster_id = '0ab4bd682f2ff858f3f834cb3613ae70';
$destination_cluster_id = 'aac6f6768700e0a9a98cbcf24309b167';

$source_cluster_name = 'CL2';
$destination_cluster_name = 'CL2';

// DVS PORTGROUP EXPORT
$source_dvs_id = 'ebded42fed964e7c46d8ee5fb26106ae';
$destination_dvs_id = 'UPDATE ME AFTER PHASE1'; //update this after phase1
$destination_dvs_name = 'DVS2';

// VSWITCH GENERATE
$vswitch_name = 'vSwitch1';