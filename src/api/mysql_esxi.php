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
$table = 'view_esxi';
 
// Table's primary key
$primaryKey = 'id';
 
// Array of database columns which should be read and sent back to DataTables.
// The `db` parameter represents the column name in the database, while the `dt`
// parameter represents the DataTables column identifier. In this case simple
// indexes
$columns = array(
    array( 'db' => 'name', 'dt' => 0 ),
    array( 'db' => 'max_evc', 'dt' => 1 ),
    array( 'db' => 'current_evc', 'dt' => 2 ),
    array( 'db' => 'status', 'dt' => 3 ),
    array( 'db' => 'power_state', 'dt' => 4 ),
    array( 'db' => 'in_maintenance_mode', 'dt' => 5 ),
    array( 'db' => 'vendor', 'dt' => 6 ),
    array( 'db' => 'model', 'dt' => 7 ),
    array( 'db' => 'memory_bytes', 'dt' => 8 ),
    array( 'db' => 'cpu_model', 'dt' => 9 ),
    array( 'db' => 'cpu_mhz', 'dt' => 10 ),
    array( 'db' => 'cpu_sockets', 'dt' => 11 ),
    array( 'db' => 'cpu_cores', 'dt' => 12 ),
    array( 'db' => 'cpu_threads', 'dt' => 13 ),
    array( 'db' => 'nics', 'dt' => 14 ),
    array( 'db' => 'hbas', 'dt' => 15 ),
    array( 'db' => 'version', 'dt' => 16 ),
    array( 'db' => 'build', 'dt' => 17 ),
    array( 'db' => 'stat_cpu_usage', 'dt' => 18 ),
    array( 'db' => 'stat_memory_usage', 'dt' => 19 ),
    array( 'db' => 'stat_uptime_sec', 'dt' => 20 ),
    array( 'db' => 'vms_powered_on', 'dt' => 21 ),
    array( 'db' => 'vcpus_powered_on', 'dt' => 22 ),
    array( 'db' => 'vmemory_mb_powered_on', 'dt' => 23 ),
    array( 'db' => 'pnics', 'dt' => 24 ),
    array( 'db' => 'vcenter_fqdn', 'dt' => 25 ),
    array( 'db' => 'vcenter_short_name', 'dt' => 26 ),
);
 
// Load MYSQL connection details
require_once( 'lib/mysql_config.php' );
  
// Load ssp class
require_once( 'lib/ssp.class.php' );
 
echo json_encode(
    SSP::simple( $_GET, $sql_details, $table, $primaryKey, $columns )
);