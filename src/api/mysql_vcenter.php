<?php
<<<<<<< HEAD
=======

// Load the library for datatables
>>>>>>> origin/master
require_once('lib/DB/DatabaseInterface.php');
require_once('lib/DB/MySQL.php');
require_once('lib/Datatables.php');

<<<<<<< HEAD
require_once('lib/mysql_config.php');
=======
// Load some common configs
require_once('lib/mysql_config.php');
require_once('lib/common.php');
>>>>>>> origin/master

use Ozdemir\Datatables\Datatables;
use Ozdemir\Datatables\DB\MySQL;

<<<<<<< HEAD
// used to convert bytes
$units = explode(' ', 'B KB MB GB TB PB');
function format_size($size) {
    global $units;
    $mod = 1024;
    for ($i = 0; $size > $mod; $i++) {
        $size /= $mod;
    }
    $endIndex = strpos($size, ".")+3;
    return substr( $size, 0, $endIndex).' '.$units[$i];
}

=======
>>>>>>> origin/master
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


