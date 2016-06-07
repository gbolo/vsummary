<?php

// SQL server connection information
require_once('../src/api/lib/mysql_config.php');

// set up PDO
try {
    $dsn = "mysql:host={$sql_details['host']};dbname={$sql_details['db']};charset={$sql_details['charset']}";
    $opt = [
        PDO::ATTR_ERRMODE            => PDO::ERRMODE_EXCEPTION,
        PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
        PDO::ATTR_EMULATE_PREPARES   => false,
    ];
    $pdo = new PDO($dsn, $sql_details['user'], $sql_details['pass'], $opt);
}
catch (PDOException $e) {
    // return 500
    echo "error connecting to database: ".$e->getMessage();
    http_response_code(500);
}


// DEFINE FUNCTIONS

function gen_folder_array($vcenter_id, $dc_name){

	global $pdo;

	$query = "SELECT DISTINCT full_path, name, moref from folder 
	WHERE full_path LIKE '{$dc_name}/%'
	AND type = 'VirtualMachine' 
	AND parent != 'datacenter' 
	AND vcenter_id = '${vcenter_id}' 
	AND present = 1;
	";

	$sth = $pdo->prepare($query);
	$sth->execute();
	$result = $sth->fetchAll();

	return $result;

}

function gen_vm_array($vcenter_id, $dc_name, $cluster_id=null){

	global $pdo;

	$query = "
	select vm.name AS name, vm.moref AS moref, vm.instance_uuid AS instance_uuid, vm.vapp_id AS vapp_id, esxi.name AS esxi_name, vm.template AS template, folder.full_path AS folder, cluster.name AS cluster, resourcepool.name AS resourcepool, datacenter.name AS datacenter FROM vm 
	LEFT JOIN folder ON vm.folder_id = folder.id LEFT JOIN esxi ON vm.esxi_id=esxi.id 
	LEFT JOIN cluster ON esxi.cluster_id=cluster.id 
	LEFT JOIN resourcepool ON vm.resourcepool_id = resourcepool.id
	LEFT JOIN datacenter ON cluster.datacenter_id=datacenter.esxi_folder_id WHERE vm.vcenter_id='{$vcenter_id}' AND datacenter.name='{$dc_name}' AND vm.present=1;
	";

	if ( !is_null($cluster_id) ){
		$query = "
		select vm.name AS name, vm.moref AS moref, vm.instance_uuid AS instance_uuid, vm.vapp_id AS vapp_id, esxi.name AS esxi_name, vm.template AS template, folder.full_path AS folder, cluster.name AS cluster, resourcepool.name AS resourcepool, datacenter.name AS datacenter FROM vm 
		LEFT JOIN folder ON vm.folder_id = folder.id LEFT JOIN esxi ON vm.esxi_id=esxi.id 
		LEFT JOIN cluster ON esxi.cluster_id=cluster.id 
		LEFT JOIN resourcepool ON vm.resourcepool_id = resourcepool.id
		LEFT JOIN datacenter ON cluster.datacenter_id=datacenter.esxi_folder_id WHERE vm.vcenter_id='{$vcenter_id}' AND datacenter.name='{$dc_name}' AND vm.present=1 AND cluster.id='$cluster_id';
		";
	}

	$sth = $pdo->prepare($query);
	$sth->execute();
	$result = $sth->fetchAll();

	return $result;

}


function gen_dvs_pg_array($dvs_id){

	global $pdo;

	$query = "
	SELECT name,vlan,vlan_type,id FROM portgroup 
	WHERE vswitch_id='{$dvs_id}' AND present = 1
	";

	$sth = $pdo->prepare($query);
	$sth->execute();
	$result = $sth->fetchAll();

	return $result;

}

function gen_esxi_array($vcenter_id, $dc_name, $cluster_id=null){

	global $pdo;

	$query = "select esxi.name AS name, esxi.moref AS moref, esxi.id AS id, cluster.name AS cluster, datacenter.name AS datacenter FROM esxi 
	LEFT JOIN cluster ON esxi.cluster_id=cluster.id 
	LEFT JOIN datacenter ON cluster.datacenter_id=datacenter.esxi_folder_id WHERE esxi.vcenter_id = '$vcenter_id' AND datacenter.name='$dc_name' AND esxi.present=1;
	";

	if ( !is_null($cluster_id) ){
		$query = "select esxi.name AS name, esxi.moref AS moref, esxi.id AS id, cluster.name AS cluster, datacenter.name AS datacenter FROM esxi 
		LEFT JOIN cluster ON esxi.cluster_id=cluster.id 
		LEFT JOIN datacenter ON cluster.datacenter_id=datacenter.esxi_folder_id WHERE esxi.vcenter_id = '$vcenter_id' AND datacenter.name='$dc_name' AND esxi.cluster_id='$cluster_id' AND esxi.present=1;
		";
	}

	$sth = $pdo->prepare($query);
	$sth->execute();
	$result = $sth->fetchAll();

	return $result;

}


function gen_resourcepool_array($vcenter_id, $cluster_id){

	global $pdo;

	$query = "select name,type from resourcepool 
	WHERE parent!='cluster' AND type='ResourcePool' AND cluster_id='$cluster_id';
	";

	$sth = $pdo->prepare($query);
	$sth->execute();
	$result = $sth->fetchAll();

	return $result;

}



function export_vm_folder_csv($vm_array, $filename){

	$fp = fopen($filename, 'w');
	foreach($vm_array as $vm){
		fputcsv($fp, $vm);
	}
	fclose($fp);

	// create json file too
	file_put_contents($filename.'.json',json_encode($vm_array));

}


function powercli_import_vm_folders($folder_array, $filename){

	$file_content = "### VSUMMARY GENERATED -- IMPORT VM FOLDERS TO DESTINATION VCENTER DATACENTER\n";
	$file_content .= 'Import-Module .\vsummaryPowershellModule.psm1;'."\n";
	$folders = array();

	foreach($folder_array as $folder){
		$folder_depth = count( explode("/", $folder['full_path']) );
		$folders[] = array("depth" => $folder_depth, "name" => $folder['name'], "path" => $folder['full_path']);
	}
	
	// sort by folder depth
	usort($folders, function($a, $b) {
    	return $a['depth'] - $b['depth'];
	});

	$file_content .= '$folderArray = @()'."\n";
	foreach ($folders as $folder){

		$path = $folder['path'];
		$depth = $folder['depth'];
		$present = false;

		// try to remove some redundant folders from array 
		foreach ($folders as $fd){
			if ( $fd['depth'] > $depth && strpos($fd['path'], $path) === 0 ){
				$present = true;
			}
			
		}
		if (!$present){

			// array of objects
			/*
			echo '$folder = New-Object System.Object'."\n";
			echo "\$folder | Add-Member -type NoteProperty -name Name -Value '{$folder['name']}'\n";
			echo "\$folder | Add-Member -type NoteProperty -name Path -Value '{$folder['path']}'\n";
			echo '$folderArray += $folder'."\n";
			*/

			// array of strings

			// Powershell escapes single quote with a consecutive single quote!
			// http://technet.microsoft.com/en-us/library/dd315325.aspx
			$folder_path = str_replace("'", "''", $folder['path']);
			$file_content .= "\$folderArray += '$folder_path'\n";

		}
	}

	$file_content .= '$vc_fqdn = Read-Host "DESTINATION vCenter"'."\n";
	$file_content .= 'Connect-vcenter $vc_fqdn'."\n";
	$file_content .= 'Import-VM-Folders $folderArray'."\n";

	file_put_contents($filename, $file_content);

}


function powercli_import_dvs($dvs_array, $filename, $dvs_dst_name){

	$file_content = "### CREATING NEW PORTGROUPS ON DESTINATION DISTRIBUTED SWITCH ###\n";
	$file_content .= '$vc_fqdn = Read-Host "DESTINATION vCenter"'."\n";
	$file_content .= 'Connect-vcenter $vc_fqdn'."\n";
	foreach($dvs_array as $pg){
		if ( $pg['vlan_type'] === 'single' ){

			$file_content .= "Get-VDSwitch -Name '{$dvs_dst_name}' | New-VDPortgroup -Name '{$pg['name']}' -NumPorts 128 -VLanId {$pg['vlan']}\n";

		}
	}

	file_put_contents($filename, $file_content);

}


function csv_vm_list($vm_array, $filename){

	$fp = fopen($filename.'.csv', 'w');
	
	foreach($vm_array as $vm){
		fputcsv($fp, $vm);
	}

	fclose($fp);

	// create json file too
	file_put_contents($filename.'.json',json_encode($vm_array));

}



function powercli_templates_to_vms($esxi_array, $outut_dir){

	foreach($esxi_array as $esxi){

		$file_content = '$vc_fqdn = Read-Host "SOURCE vCenter"'."\n";
		$file_content .= 'Connect-vcenter $vc_fqdn'."\n";
		$file_content .= 'Get-Vmhost -Id HostSystem-'. $esxi['moref'] .' | get-template | %{ Set-Template -Template $_ -ToVM };';

		$file = $outut_dir . "convert-templates-to-vm_{$esxi['name']}.ps1";
		file_put_contents($file, $file_content);

	}

}



function powercli_move_vm_vnics($dvs_array, $esxi_array, $dst_svs_name, $outut_dir){

	global $pdo;

	$vnics_changed = array();
	foreach($esxi_array as $esxi){

		$file1 = $outut_dir . 'create_vswitch_pg_' . $esxi['name'] . '.ps1';
		$file2 = $outut_dir . 'change_vm_pg_' . $esxi['name'] . '.ps1';
		$file3 = $outut_dir . 'csv/VNICS_CHANGED.csv';

		$file1_content = "### CREATING NEW PORTGROUPS ON STANDARD VSWITCH FOR {$esxi['name']} ###\n";
		$file1_content .= '$vc_fqdn = Read-Host "SOURCE vCenter"'."\n";
		$file1_content .= 'Connect-vcenter $vc_fqdn'."\n";
		$file2_content = "### CHANGING PORTGROUPS FOR EACH VM ON {$esxi['name']} ###\n";
		$file2_content .= '$vc_fqdn = Read-Host "SOURCE vCenter"'."\n";
		$file2_content .= 'Connect-vcenter $vc_fqdn'."\n";

		foreach($dvs_array as $pg){
			if ( $pg['vlan_type'] === 'single' ){

				// create portgroup
				$file1_content .= "Get-VMHost {$esxi['name']} | Get-VirtualSwitch -Name '{$dst_svs_name}' | New-VirtualPortGroup -Name 'mig-{$pg['name']}' -VLanId {$pg['vlan']}\n";

			}
			// move every VM to it
			foreach($pdo->query("SELECT id, moref, name, instance_uuid FROM vm WHERE esxi_id = '{$esxi['id']}' AND present = 1") as $vm){

				// loop through vnics
				foreach($pdo->query("SELECT * FROM vnic WHERE vm_id = '{$vm['id']}' AND portgroup_id = '{$pg['id']}' AND present = 1") as $vnic){

					$vnics_changed[] = array( "vm_name" => $vm['name'], "instance_uuid" => $vm['instance_uuid'], "vnic" => $vnic['name'], "portgroup" => $pg['name'], "esxi" => $esxi['name'] );
					//$file2_content .= "Get-VMHost {$esxi['name']} | Get-VM -Id VirtualMachine-{$vm['moref']} | Get-NetworkAdapter -Name '{$vnic['name']}' | Set-NetworkAdapter -NetworkName 'mig-{$pg['name']}' -Confirm:\$false -RunAsync\n";
					$file2_content .= "Get-VM -Id VirtualMachine-{$vm['moref']} | Get-NetworkAdapter -Name '{$vnic['name']}' | Set-NetworkAdapter -NetworkName 'mig-{$pg['name']}' -Confirm:\$false -RunAsync\n";

				}		
			}

		}

		file_put_contents($file1, $file1_content);
		file_put_contents($file2, $file2_content);

		$fp = fopen($file3, 'w');
		foreach($vnics_changed as $vnic){
			fputcsv($fp, $vnic);
		}
		fclose($fp);

		// create json file too
		file_put_contents($outut_dir . 'csv/VNICS_CHANGED.json',json_encode($vnics_changed));

	}

}



function powercli_restore_vm_vnics($vm_array, $vnics_changed, $outut_dir){

	// GENERATE POWERCLI
	$file_content = "### MOVING VMS BACK TO DVS\n";
	$file_content .= 'Import-Module .\vsummaryPowershellModule.psm1;'."\n";
	$file_content .= '$vc_fqdn = Read-Host "DESTINATION vCenter"'."\n";
	$file_content .= 'Connect-vcenter $vc_fqdn'."\n";

	foreach ( $vnics_changed as $vnic ){

		$key = array_search($vnic['instance_uuid'], array_column($vm_array, 'instance_uuid'));
		$vm_moref = $vm_array[$key]['moref'];

		$file_content .= "Get-VM -Id VirtualMachine-{$vm_moref} | Get-NetworkAdapter -Name '{$vnic['vnic']}' | Set-NetworkAdapter -NetworkName '{$vnic['portgroup']}' -Confirm:\$false -RunAsync\n";

	}

	file_put_contents($outut_dir.'RESTORE-VM-PORTGROUPS.ps1', $file_content);

}


function powercli_restore_vm_folders($vm_source_array, $vm_array, $outut_dir){

	$vm_folder_move = array();

	foreach($vm_source_array as $vm){

		$folder = $vm['folder'];
		if ( !is_null($folder) ){
			$key = array_search($vm['instance_uuid'], array_column($vm_array, 'instance_uuid'));
			$vm_moref = $vm_array[$key]['moref'];
			$vm_folder_move[$folder][] = $vm_moref;
		}

	}

	// GENERATE POWERCLI
	$file_content = "### MOVING VMS BACK TO ORIGINAL FOLDERS\n";
	$file_content .= 'Import-Module .\vsummaryPowershellModule.psm1;'."\n";
	$file_content .= '$vc_fqdn = Read-Host "DESTINATION vCenter"'."\n";
	$file_content .= 'Connect-vcenter $vc_fqdn'."\n";

	foreach ( $vm_folder_move as $folder => $vm ){

		$file_content .= '$folder = Get-FolderByPath -Path "'.$folder."\" -Separator '/'\n";
		foreach ($vm as $moref){
			$file_content .= "Get-VM -Id 'VirtualMachine-{$moref}' | Move-VM -Destination \$folder\n";
		}

	}

	file_put_contents($outut_dir.'RESTORE-VM-FOLDERS.ps1', $file_content);

}

function powercli_export_vapps($vcenter_id, $cluster_id, $outut_dir){

	global $pdo;

	$query = "select name,moref from resourcepool 
	WHERE present=1 AND vapp_in_path = 0 AND type='VirtualApp' 
	AND vcenter_id='$vcenter_id' AND cluster_id='$cluster_id'";

	$sth = $pdo->prepare($query);
	$sth->execute();
	$result = $sth->fetchAll();

	$file_content = "### EXPORT VAPPS FROM SOURCE VCENTER\n";
	$file_content .= 'Import-Module .\vsummaryPowershellModule.psm1;'."\n";
	$file_content .= '$vc_fqdn = Read-Host "SOURCE vCenter"'."\n";
	$file_content .= 'Connect-vcenter $vc_fqdn'."\n";

	foreach ($result as $vapp){

		$file_content .= "Get-VApp -Id 'VirtualApp-{$vapp['moref']}' | Stop-VApp -force \n";
		$file_content .= "Get-VApp -Id 'VirtualApp-{$vapp['moref']}' | Export-VApp -destination '.\\vapps' \n";

	}

	file_put_contents($outut_dir.'EXPORT-VAPPS.ps1', $file_content);

}
