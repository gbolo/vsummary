<?php

$POLLER_ENABLED = true;


$dt_dom = "<'row'<'col-sm-6'l><'col-sm-6 text-right'B>><'row'<'col-sm-12'tr>><'row'<'col-sm-5'i><'col-sm-7'p>>";
$dt_select = "true";
$dt_buttons = "
'copy',
'csv',
{ extend: 'colvis',
	className: 'colvis',
	text: 'Custom View'
}
";

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
			$dt_select = "false";
			if ( $POLLER_ENABLED == true ){
					$dt_buttons = "
					{
							text: '<i class=\"fa fa-plus-square fa-fw\"></i> <b>Add vCenter Server</b>',
							className: 'btn-success',
							action: function ( e, dt, node, config ) {
									$(\"#pollerModal\").find(\".modal-content\").load(\"add.php\");
									$(\"#pollerModal\").modal() ;
							}
					}
					";
			} else {
				$view_title = 'Poller is Disabled for Demo';
				$dt_buttons = "
				{
						text: '<i class=\"fa fa-ban fa-fw\"></i> <b>Add vCenter Server (disabled for demo)</b>',
						className: 'btn-danger'
				}
				";
			}

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
