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
capacity_bytes,
path,
thin_provisioned,
vm_power_state,
datastore_name,
datastore_type,
esxi_name,
vcenter_fqdn,
vcenter_short_name
FROM view_vdisk
');

// Modify output
$dt->edit('capacity_bytes', function ($data){
    $hr = format_size($data['capacity_bytes']);
    return $hr;
});

// Respond with results
echo $dt->generate();