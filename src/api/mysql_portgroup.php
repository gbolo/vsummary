<?php
 
<<<<<<< HEAD
/*
 * DataTables example server-side processing script.
 *
 * Please note that this script is intentionally extremely simply to show how
 * server-side processing can be implemented, and probably shouldn't be used as
 * the basis for a large complex system. It is suitable for simple use cases as
 * for learning.
 *
 * See http://datatables.net/usage/server-side for full details on the server-
 * side processing requirements of DataTables.
 *
 * @license MIT - http://datatables.net/license_mit
 */
 

// DB table to use
$table = 'view_portgroup';
 
// Table's primary key
$primaryKey = 'name';
 
// Array of database columns which should be read and sent back to DataTables.
// The `db` parameter represents the column name in the database, while the `dt`
// parameter represents the DataTables column identifier. In this case simple
// indexes
$columns = array(
    array( 'db' => 'name', 'dt' => 0 ),
    array( 'db' => 'type', 'dt' => 1 ),
    array( 'db' => 'vlan', 'dt' => 2 ),
    array( 'db' => 'vlan_type', 'dt' => 3 ),
    array( 'db' => 'vswitch_name', 'dt' => 4 ),
    array( 'db' => 'vswitch_type', 'dt' => 5 ),
    array( 'db' => 'vswitch_max_mtu', 'dt' => 6 ),
    array( 'db' => 'vcenter_fqdn', 'dt' => 7 ),
    array( 'db' => 'vcenter_short_name', 'dt' => 8 )
);
 
// Load MYSQL connection details
require_once( 'lib/mysql_config.php' );
  
// Load ssp class
require_once( 'lib/ssp.class.php' );
 
echo json_encode(
    SSP::simple( $_GET, $sql_details, $table, $primaryKey, $columns )
);
=======
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
vnics,
vcenter_fqdn,
vcenter_short_name
FROM view_portgroup
');

// Respond with results
echo $dt->generate();
>>>>>>> origin/master
