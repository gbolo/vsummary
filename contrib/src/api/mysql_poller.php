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
auto_poll,
last_poll
FROM vcenter
');

// Modify output
$dt->edit('id', function ($data){
    $GLOBALS['id'] = $data['id'];
    $edit = '<div class="btn-group btn-group-sm" role="group" aria-label="options">';
	  $edit .= '<a href="edit.php?id='.$data['id'].'" data-toggle="modal" data-target="#pollerModal" class="btn btn-primary btn-xs"><strong>EDIT</strong></a>';
    $edit .= '<a href="view.php?id='.$data['id'].'" class="btn btn-info btn-xs"><strong>VIEW</strong></a>';
    $edit .= '<a href="run.php?id='.$data['id'].'" data-toggle="modal" data-target="#pollerModal" class="btn btn-success btn-xs"><strong>RUN</strong></a>';
    $edit .= '</div>';
    return $edit;
});

$dt->edit('last_poll_status', function ($data){
    if ( $data['last_poll_status'] == 'Idle' ){
        $hr = '<span class="label label-pill label-default">IDLE</span>';
    } elseif ( $data['last_poll_status'] == 'running' ){
        $hr = '<span class="label label-pill label-success">RUNNING</span>';
    } else {
        $hr = '<span class="label label-pill label-warning">UNKNOWN</span>';
    }
    return $hr;
});

$dt->edit('last_poll_result', function ($data){
    if ( $data['last_poll_result'] == 0 ){
        $hr = '<span class="label label-pill label-default">UNKNOWN</span>';
    } elseif ( $data['last_poll_result'] == 1 ){
        $hr = '<span class="label label-pill label-success">SUCCESS</span>';
    } else {
        $hr = '<span class="label label-pill label-danger">FAILED</span>';
    }
    return $hr;
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
