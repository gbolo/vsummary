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
Editor::inst( $db, 'esxi' )
    ->fields(
        Field::inst( 'esxi.name' ),
        Field::inst( 'esxi.max_evc' ),
        Field::inst( 'esxi.current_evc' ),
        Field::inst( 'esxi.status' ),
        Field::inst( 'esxi.power_state' ),
        Field::inst( 'esxi.in_maintenance_mode' ),
        Field::inst( 'esxi.vendor' ),
        Field::inst( 'esxi.model' ),
        Field::inst( 'esxi.memory_bytes' ),
        Field::inst( 'esxi.cpu_model' ),
        Field::inst( 'esxi.cpu_mhz' ),
        Field::inst( 'esxi.cpu_sockets' ),
        Field::inst( 'esxi.cpu_cores' ),
        Field::inst( 'esxi.cpu_threads' ),
        Field::inst( 'esxi.nics' ),
        Field::inst( 'esxi.hbas' ),
        Field::inst( 'esxi.version' ),
        Field::inst( 'esxi.build' ),
        Field::inst( 'esxi.stat_cpu_usage' ),
        Field::inst( 'esxi.stat_memory_usage' ),
        Field::inst( 'esxi.stat_uptime_sec' ),
        Field::inst( 'esxi.vcenter_id' )
    )
    ->where( $key = 'esxi.present', $value = '1', $op = '=' )
    ->process( $_POST )
    ->json();
