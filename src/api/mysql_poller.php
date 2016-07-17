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
	  $edit = '<button>'.$data['id'].'</button>';
    return $edit;
});

// Respond with results
echo $dt->generate();
