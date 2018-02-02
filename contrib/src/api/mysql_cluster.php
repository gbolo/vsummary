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
name,
num_hosts,
avg_memory_per_host,
total_memory_bytes,
total_memory_used,
vms_on,
avg_vcpu_per_vm,
avg_memory_per_vm,
ratio_memory,
ratio_memory_80,
supported_failures_80,
vcenter_short_name
FROM view_cluster_capacity
');

// Modify output
$dt->edit('avg_memory_per_host', function ($data){
    $hr = format_size($data['avg_memory_per_host']);
    return $hr;
});

$dt->edit('avg_memory_per_vm', function ($data){
    $hr = format_size($data['avg_memory_per_vm']);
    return $hr;
});

$dt->edit('total_memory_bytes', function ($data){
    $hr = format_size($data['total_memory_bytes']);
    return $hr;
});

$dt->edit('total_memory_used', function ($data){
    $hr = format_size($data['total_memory_used']);
    return $hr;
});

// Respond with results
echo $dt->generate();
