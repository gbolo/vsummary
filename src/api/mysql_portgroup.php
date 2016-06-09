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
type,
vlan,
vlan_type,
vswitch_name,
vswitch_type,
vswitch_max_mtu,
vcenter_fqdn,
vcenter_short_name
FROM view_portgroup
');

// Respond with results
echo $dt->generate();