<?php
 
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
$table = 'view_datastore';
 
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
    array( 'db' => 'name', 'dt' => 0 ),
    array( 'db' => 'status', 'dt' => 1 ),
    array( 
        'db' => 'capacity_bytes', 
        'dt' => 2,
        'formatter' => function( $d, $row ) {
            return format_size($d);
        }
    ),
    array( 
        'db' => 'free_bytes', 
        'dt' => 3,
        'formatter' => function( $d, $row ) {
            return format_size($d);
        }
    ),
    array( 
        'db' => 'uncommitted_bytes', 
        'dt' => 4,
        'formatter' => function( $d, $row ) {
            return format_size($d);
        }
    ),
    array( 'db' => 'type', 'dt' => 5 ),
    array( 'db' => 'vcenter_fqdn', 'dt' => 6 ),
    array( 'db' => 'vcenter_short_name', 'dt' => 7 ),
);
 
// Load MYSQL connection details
require_once( 'lib/mysql_config.php' );
  
// Load ssp class
require_once( 'lib/ssp.class.php' );
 
echo json_encode(
    SSP::simple( $_GET, $sql_details, $table, $primaryKey, $columns )
);