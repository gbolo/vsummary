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
Editor::inst( $db, 'vdisk' )
    ->fields(
        Field::inst( 'vdisk.name' ),
        Field::inst( 'vdisk.capacity_bytes' ),
        Field::inst( 'vdisk.path' ),
        Field::inst( 'vdisk.thin_provisioned' ),
        Field::inst( 'vdisk.datastore_id' ),
        Field::inst( 'vdisk.uuid' ),
        Field::inst( 'vdisk.disk_object_id' ),
        Field::inst( 'vdisk.vm_id' ),
        Field::inst( 'vdisk.esxi_id' ),
        Field::inst( 'vdisk.vcenter_id' )
    )
    ->where( $key = 'vdisk.present', $value = '1', $op = '=' )
    ->process( $_POST )
    ->json();
