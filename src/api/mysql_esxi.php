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
max_evc,
current_evc,
status,
power_state,
in_maintenance_mode,
vendor,
model,
memory_bytes,
cpu_model,
cpu_mhz,
cpu_sockets,
cpu_cores,
cpu_threads,
nics,
hbas,
version,
build,
stat_cpu_usage,
stat_memory_usage,
stat_uptime_sec,
vms_powered_on,
vcpus_powered_on,
vmemory_mb_powered_on,
pnics,
cluster,
datacenter,
vcenter_fqdn,
vcenter_short_name
FROM view_esxi
');

// Modify output
$dt->edit('memory_bytes', function ($data){
    $hr = format_size($data['memory_bytes']);
    return $hr;
});

$dt->edit('cpu_mhz', function ($data){
    $hr = $data['cpu_mhz'] . ' MHZ';
    return $hr;
});

$dt->edit('stat_memory_usage', function ($data){
    $hr = format_size($data['stat_memory_usage']*1000*1000);
    return $hr;
});

$dt->edit('stat_cpu_usage', function ($data){
    $hr = $data['stat_cpu_usage'] . ' MHZ';
    return $hr;
});

$dt->edit('stat_uptime_sec', function ($data){
    $hr = uptime_human_readable($data['stat_uptime_sec']);
    return $hr;
});

$dt->edit('status', function ($data){
    if ($data['status'] === 'green'){
        return '<span class="label label-pill label-success">green</span>';
    }elseif ($data['status'] === 'red'){
        return '<span class="label label-pill label-danger">red</span>';
    }elseif ($data['status'] === 'yellow'){
        return '<span class="label label-pill label-warning">yellow</span>';
    }else{
        return $data['status'];
    }
});

$dt->edit('power_state', function ($data){
    if ($data['power_state'] === '1'){
        return '<span class="label label-pill label-success">1 - ON</span>';
    }elseif ($data['power_state'] === '0'){
        return '<span class="label label-pill label-danger">0 - OFF</span>';
    }else{
        return $data['power_state'];
    }
});

// Respond with results
echo $dt->generate();
