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
folder,
vcpu,
memory_mb,
power_state,
config_guest_os,
config_version,
config_change_version,
guest_tools_version,
guest_tools_running,
guest_hostname,
guest_ip,
cluster,
pool,
datacenter,
stat_cpu_usage,
stat_host_memory_usage,
stat_guest_memory_usage,
stat_uptime_sec,
esxi_name,
esxi_current_evc,
esxi_status,
esxi_cpu_model,
vdisks,
vnics,
vmx_path,
vcenter_fqdn,
vcenter_short_name
FROM view_vm
');

// Modify output
$dt->edit('folder', function ($data){
    if ( is_null($data['folder']) ){
        return 'vApp not supported yet';
    } else {
        return $data['folder'];
    }
});

$dt->edit('memory_mb', function ($data){
    $hr = $data['memory_mb'] . ' MB';
    return $hr;
});

$dt->edit('stat_uptime_sec', function ($data){
    $hr = uptime_human_readable($data['stat_uptime_sec']);
    return $hr;
});

// Respond with results
echo $dt->generate();