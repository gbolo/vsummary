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

// Keep track of ID
$id = null;

// Create object
$dt = new Datatables(new MySQL($config));

// Query
$dt->query('SELECT
id,
fqdn,
short_name,
last_poll_status,
last_poll_result,
last_poll_output,
auto_poll
FROM vcenter
');

// Modify output
$dt->edit('id', function ($data){
    $GLOBALS['id'] = $data['id'];
	  $edit = '<a href="edit.php?id='.$data['id'].'" class="btn btn-info btn-sm"><strong>EDIT</strong></a>';
    return $edit;
});

$dt->edit('last_poll_status', function ($data){
	  $hr = '<span class="label label-pill label-default">Idle</span>';
    return $hr;
});

$dt->edit('last_poll_result', function ($data){
	  $hr = '<span class="label label-pill label-success">Success</span>';
    return $hr;
});

$dt->edit('last_poll_output', function ($data){
    global $id;
	  $edit = '<a href="view.php?id='.$GLOBALS['id'].'" class="btn btn-info btn-sm"><strong>VIEW</strong></a>';
    return $edit;
});

$dt->edit('auto_poll', function ($data){
    if ( $data['auto_poll'] == 1 ){
        $hr = '<span class="label label-pill label-success">ENABLED</span>';
    } elseif ( $data['auto_poll'] == 0 ){
        $hr = '<span class="label label-danger label-">DISABLED</span>';
    } else {
        $hr = '<span class="label label-warning label-">'.$data['auto_poll'].'</span>';
    }
    return $hr;
});


// Respond with results
echo $dt->generate();
