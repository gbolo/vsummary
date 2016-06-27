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
status,
capacity_bytes,
free_bytes,
uncommitted_bytes,
type,
vcenter_fqdn,
vcenter_short_name
FROM view_datastore
');

// Modify output
$dt->edit('capacity_bytes', function ($data){
    $hr = format_size($data['capacity_bytes']);
    return $hr;
});

$dt->edit('free_bytes', function ($data){
    $hr = format_size($data['free_bytes']);
    return $hr;
});

$dt->edit('uncommitted_bytes', function ($data){
    $hr = format_size($data['uncommitted_bytes']);
    return $hr;
});

// Respond with results
echo $dt->generate();