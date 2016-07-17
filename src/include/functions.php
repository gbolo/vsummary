<?php

if (isset($_GET['view'])){
	switch($_GET['view'])
	{
		case 'vm':
			$view = $_GET['view'];
			$view_title = 'Virtual Machine';
			break;
		case 'poller':
			$view = $_GET['view'];
			$view_title = 'Poller';
			break;
		case 'esxi':
			$view = $_GET['view'];
			$view_title = 'ESXi Host';
			break;
		case 'vdisk':
			$view = $_GET['view'];
			$view_title = 'Virtual Machine Disk';
			break;
		case 'vnic':
			$view = $_GET['view'];
			$view_title = 'Virtual Machine Network Adapter';
			break;
		case 'datastore':
			$view = $_GET['view'];
			$view_title = 'ESXi Datastore';
			break;
		case 'portgroup':
			$view = $_GET['view'];
			$view_title = 'Virtual Switch Portgroup';
			break;
		case 'vcenter':
			$view = $_GET['view'];
			$view_title = 'vCenter Environment';
			break;
		default;
			$view = 'vm';
			$view_title = 'Virtual Machine';
			break;
	}
}else{
	$view = 'vm';
	$view_title = 'Virtual Machine';
}

function datatables_html($view){

	include("include/${view}.html");

}

?>
