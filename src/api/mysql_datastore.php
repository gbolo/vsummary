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
Editor::inst( $db, 'datastore' )
    ->fields(
        Field::inst( 'datastore.name' ),
        Field::inst( 'datastore.moref' ),
        Field::inst( 'datastore.status' ),
        Field::inst( 'datastore.capacity_bytes' ),
        Field::inst( 'datastore.free_bytes' ),
        Field::inst( 'datastore.uncommitted_bytes' ),
        Field::inst( 'datastore.type' ),
        Field::inst( 'datastore.vcenter_id' )
    )
    ->where( $key = 'datastore.present', $value = '1', $op = '=' )
    ->process( $_POST )
    ->json();
