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
$table = 'view_vdisk';
 
// Table's primary key
$primaryKey = 'id';
 

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

// Array of database columns which should be read and sent back to DataTables.
// The `db` parameter represents the column name in the database, while the `dt`
// parameter represents the DataTables column identifier. In this case simple
// indexes
$columns = array(
    array( 'db' => 'vm_name', 'dt' => 0 ),
    array( 'db' => 'name', 'dt' => 1 ),
    array( 
        'db' => 'capacity_bytes', 
        'dt' => 2,
        'formatter' => function( $d, $row ) {
            return format_size($d);
        }
    ),
    array( 'db' => 'path', 'dt' => 3 ),
    array( 'db' => 'thin_provisioned', 'dt' => 4 ),
    array( 'db' => 'vm_power_state', 'dt' => 5 ),
    array( 'db' => 'datastore_name', 'dt' => 6 ),
    array( 'db' => 'datastore_type', 'dt' => 7 ),
    array( 'db' => 'esxi_name', 'dt' => 8 ),
    array( 'db' => 'vcenter_fqdn', 'dt' => 9 ),
    array( 'db' => 'vcenter_short_name', 'dt' => 10 )
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
>>>>>>> origin/master
