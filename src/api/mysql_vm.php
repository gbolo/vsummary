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
$table = 'view_vm';
 
// Table's primary key
$primaryKey = 'id';
 
// Array of database columns which should be read and sent back to DataTables.
// The `db` parameter represents the column name in the database, while the `dt`
// parameter represents the DataTables column identifier. In this case simple
// indexes
$columns = array(
    array( 'db' => 'name', 'dt' => 0 ),
    array( 
        'db' => 'folder', 
        'dt' => 1,
        'formatter' => function( $d, $row ) {
            if ( is_null($d) ){
                return 'vApp not supported yet';
            } else {
                return $d;
            }
        }
    ),
    array( 'db' => 'vcpu', 'dt' => 2 ),
    array( 'db' => 'memory_mb', 'dt' => 3 ),
    array( 
        'db' => 'memory_mb', 
        'dt' => 3,
        'formatter' => function( $d, $row ) {
            return $d . " MB";
        }
    ),
    array( 'db' => 'power_state', 'dt' => 4 ),
    array( 'db' => 'config_guest_os', 'dt' => 5 ),
    array( 'db' => 'config_version', 'dt' => 6 ),
    array( 'db' => 'config_change_version', 'dt' => 7 ),
    array( 'db' => 'guest_tools_version', 'dt' => 8 ),
    array( 'db' => 'guest_tools_running', 'dt' => 9 ),
    array( 'db' => 'guest_hostname', 'dt' => 10 ),
    array( 'db' => 'guest_ip', 'dt' => 11 ),
    array( 'db' => 'stat_cpu_usage', 'dt' => 12 ),
    array( 'db' => 'stat_host_memory_usage', 'dt' => 13 ),
    array( 'db' => 'stat_guest_memory_usage', 'dt' => 14 ),
    array( 'db' => 'stat_uptime_sec', 'dt' => 15 ),
    array( 'db' => 'esxi_name', 'dt' => 16 ),
    array( 'db' => 'esxi_current_evc', 'dt' => 17 ),
    array( 'db' => 'esxi_status', 'dt' => 18 ),
    array( 'db' => 'esxi_cpu_model', 'dt' => 19 ),
    array( 'db' => 'vdisks', 'dt' => 20 ),
    array( 'db' => 'vnics', 'dt' => 21 ),
    array( 'db' => 'vmx_path', 'dt' => 22 ),
    array( 'db' => 'vcenter_fqdn', 'dt' => 23 ),
    array( 'db' => 'vcenter_short_name', 'dt' => 24 )
);
 
// Load MYSQL connection details
require_once( 'lib/mysql_config.php' );
  
// Load ssp class
require_once( 'lib/ssp.class.php' );
 
echo json_encode(
    SSP::simple( $_GET, $sql_details, $table, $primaryKey, $columns )
);