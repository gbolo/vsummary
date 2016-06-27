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
vm_name,
name,
mac,
type,
connected,
status,
esxi_name,
portgroup_name,
vlan,
vswitch_name,
vswitch_type,
vswitch_max_mtu,
vcenter_fqdn,
vcenter_short_name
FROM view_vnic
');

// Modify output
$dt->edit('name', function ($data){
    $hr = str_replace("Network adapter", "vNIC #", $data['name']);
    return $hr;
});

$dt->edit('portgroup_name', function ($data){
    if ($data['portgroup_name'] === 'ORPHANED'){
        return "NULL";
    }else{
        return $data['portgroup_name'];
    }
});

$dt->edit('vlan', function ($data){
    if ($data['vlan'] === '4095'){
        return "ALL";
    }elseif ($data['vlan'] === '0'){
        return "None";
    }else {
        return $data['vlan'];
    }
});

// Respond with results
echo $dt->generate();