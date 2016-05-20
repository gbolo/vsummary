<?php

// SQL server connection information
$sql_details = array(
    'user'    => 'vsummary',
    'pass'    => 'changeme',
    'db'      => 'vsummary',
    'host'    => 'localhost',
    'charset' => 'utf8'
);


// Format for new library
$config = array(
	'host'     => $sql_details['host'],
	'port'     => '3306',
	'username' => $sql_details['user'],
	'password' => $sql_details['pass'],
	'database' => $sql_details['db']
);


?>