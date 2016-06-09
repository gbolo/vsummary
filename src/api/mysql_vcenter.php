<?php

// Load the library for datatables
require_once('lib/DB/DatabaseInterface.php');
require_once('lib/DB/MySQL.php');
require_once('lib/Datatables.php');

// Load some common configs
require_once('lib/mysql_config.php');
require_once('lib/common.php');

use Ozdemir\Datatables\Datatables;
use Ozdemir\Datatables\DB\MySQL;

// Create object
$dt = new Datatables(new MySQL($config));

// Query
$dt->query('SELECT 
fqdn,
short_name,
vms_vcpu_on,
vms_memory_on,
vms_on,
vms,
datacenters,
clusters,
esxi_hosts,
esxi_cpu,
esxi_memory,
vnics,
vdisks,
datastores,
portgroups,
vswitches,
resourcepools
FROM view_vcenter
');

// Modify output
$dt->edit('esxi_memory', function ($data){
	$hr = format_size($data['esxi_memory']);
    return $hr;
});

$dt->edit('vms_memory_on', function ($data){
    $hr = format_size(1000000 * $data['vms_memory_on']);
    return $hr;
});

// Respond with results
echo $dt->generate();


