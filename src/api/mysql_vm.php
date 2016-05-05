<?php

/*
 * Example PHP implementation used for the index.html example
 */

// DataTables PHP library
include( "lib/DataTables.php" );

// Alias Editor classes so they are easy to use
use
	DataTables\Editor,
	DataTables\Editor\Field,
	DataTables\Editor\Format,
	DataTables\Editor\Join,
	DataTables\Editor\Upload,
	DataTables\Editor\Validate;

if ( isset($_POST['action']) && $_POST['action'] === 'remove' ) {
	header("HTTP/1.0 204 No Response");
	exit;
}

// Build our Editor instance and process the data coming from _POST
Editor::inst( $db, 'vm' )
    ->fields(
        Field::inst( 'vm.name' ),
        Field::inst( 'vm.vmx_path' ),
        Field::inst( 'vm.vcpu' ),
        Field::inst( 'vm.memory_mb' ),
        Field::inst( 'vm.config_guest_os' ),
        Field::inst( 'vm.config_version' ),
        Field::inst( 'vm.config_change_version' ),
        Field::inst( 'vm.guest_tools_version' ),
        Field::inst( 'vm.guest_tools_running' ),
        Field::inst( 'vm.guest_hostname' ),
        Field::inst( 'vm.guest_ip' ),
        Field::inst( 'vm.stat_cpu_usage' ),
        Field::inst( 'vm.stat_host_memory_usage' ),
        Field::inst( 'vm.stat_guest_memory_usage' ),
        Field::inst( 'vm.stat_uptime_sec' ),
        Field::inst( 'vm.power_state' )
    )
    ->where( $key = 'vm.present', $value = '1', $op = '=' )
    ->process( $_POST )
    ->json();
