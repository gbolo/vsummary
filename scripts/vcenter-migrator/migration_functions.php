<?php

// SQL server connection information
require_once('../api/lib/mysql_config.php');

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
	select vm.name AS name, vm.moref AS moref, vm.instance_uuid AS instance_uuid, vm.vapp_id AS vapp_id, esxi.name AS esxi_name, vm.template AS template, folder.full_path AS folder, cluster.name AS cluster, resourcepool.name AS resourcepool, resourcepool.full_path AS rpool_full_path, datacenter.name AS datacenter FROM vm 
	LEFT JOIN folder ON vm.folder_id = folder.id LEFT JOIN esxi ON vm.esxi_id=esxi.id 
	LEFT JOIN cluster ON esxi.cluster_id=cluster.id 
	LEFT JOIN resourcepool ON vm.resourcepool_id = resourcepool.id
	LEFT JOIN datacenter ON cluster.datacenter_id=datacenter.esxi_folder_id WHERE vm.vcenter_id='{$vcenter_id}' AND datacenter.name='{$dc_name}' AND vm.present=1;
	";

	if ( !is_null($cluster_id) ){
		$query = "
		select vm.name AS name, vm.moref AS moref, vm.instance_uuid AS instance_uuid, vm.vapp_id AS vapp_id, esxi.name AS esxi_name, vm.template AS template, folder.full_path AS folder, cluster.name AS cluster, resourcepool.name AS resourcepool, resourcepool.full_path AS rpool_full_path, datacenter.name AS datacenter FROM vm 
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

function gen_rpool_array($cluster_id){

	global $pdo;

	$query = "
	SELECT moref, full_path FROM resourcepool 
	WHERE present=1 AND type='ResourcePool' AND parent != 'cluster'
	AND cluster_id='{$cluster_id}'
	";

	$sth = $pdo->prepare($query);
	$sth->execute();
	$result = $sth->fetchAll();

	return $result;

}

function gen_vapp_array($cluster_id){

	global $pdo;

	$query = "
	SELECT moref, full_path FROM resourcepool 
	WHERE present=1 AND type='VirtualApp'
	AND cluster_id='{$cluster_id}'
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


function powercli_import_vm_folders($folder_array){

	$file_content = "### VSUMMARY GENERATED -- IMPORT VM FOLDERS TO DESTINATION VCENTER DATACENTER\n";
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

	$file_content .= 'Import-VM-Folders $folderArray'."\n";

	return $file_content;

}


function powercli_import_dvs($dvs_array, $dvs_dst_name){

	$file_content = "### CREATING NEW PORTGROUPS ON DESTINATION DISTRIBUTED SWITCH ###\n";
	foreach($dvs_array as $pg){
		if ( $pg['vlan_type'] === 'single' ){

			$file_content .= "Get-VDSwitch -Name '{$dvs_dst_name}' | New-VDPortgroup -Name '{$pg['name']}' -NumPorts 128 -VLanId {$pg['vlan']}\n";

		}
	}

	return $file_content;

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



function powercli_templates_to_vms($esxi_array){

	$file_content = "### CONVERTING TEMPLATES TO VMS \n";

	foreach($esxi_array as $esxi){

		$file_content .= "Write-Host 'TEMPLATE VMS ON ESXI {$esxi['name']} WILL BE CONVERTED TO VMS.' \n";
		$file_content .= '$confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "'."\n";
    	$file_content .= 'if ($confirmation -eq \'y\') {'."\n";
		$file_content .= '  Get-Vmhost -Id HostSystem-'. $esxi['moref'] .' | get-template | %{ Set-Template -Template $_ -ToVM };'."\n";
		$file_content .= '}'."\n";
	}

	return $file_content;

}



function powercli_move_vm_vnics($dvs_array, $esxi_array, $dst_svs_name, $file_pre){

	global $pdo;

	$vnics_changed = array();

	$file1_content = '';
	$file2_content = '';

	foreach($esxi_array as $esxi){

		$file3 = $file_pre . '_VNICS_CHANGED.csv';

		$file1_content .= "Write-Host CREATING NEW PORTGROUPS ON STANDARD VSWITCH FOR {$esxi['name']}\n";
		$file2_content .= "write-Host CHANGING PORTGROUPS FOR EACH VM ON {$esxi['name']} \n";
		$file1_content .= '$confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "'."\n";
      	$file2_content .= '$confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "'."\n";
		$file1_content .= 'if ($confirmation -eq \'y\') {'."\n";
        $file2_content .= 'if ($confirmation -eq \'y\') {'."\n";



		foreach($dvs_array as $pg){
			if ( $pg['vlan_type'] === 'single' ){

				// create portgroup
				$file1_content .= "   Get-VMHost {$esxi['name']} | Get-VirtualSwitch -Name '{$dst_svs_name}' | New-VirtualPortGroup -Name 'mig-{$pg['name']}' -VLanId {$pg['vlan']}\n";

			}
			// move every VM to it
			foreach($pdo->query("SELECT id, moref, name, instance_uuid FROM vm WHERE esxi_id = '{$esxi['id']}' AND present = 1") as $vm){

				// loop through vnics
				foreach($pdo->query("SELECT * FROM vnic WHERE vm_id = '{$vm['id']}' AND portgroup_id = '{$pg['id']}' AND present = 1") as $vnic){

					$vnics_changed[] = array( "vm_name" => $vm['name'], "instance_uuid" => $vm['instance_uuid'], "vnic" => $vnic['name'], "portgroup" => $pg['name'], "esxi" => $esxi['name'] );
					//$file2_content .= "Get-VMHost {$esxi['name']} | Get-VM -Id VirtualMachine-{$vm['moref']} | Get-NetworkAdapter -Name '{$vnic['name']}' | Set-NetworkAdapter -NetworkName 'mig-{$pg['name']}' -Confirm:\$false -RunAsync\n";
					$file2_content .= "   Get-VM -Id VirtualMachine-{$vm['moref']} | Get-NetworkAdapter -Name '{$vnic['name']}' | Set-NetworkAdapter -NetworkName 'mig-{$pg['name']}' -Confirm:\$false -RunAsync\n";

				}		
			}

		}

		$file1_content .= "}\n";
		$file2_content .= "}\n";

		$output = array();
		$output['vswitch'] = $file1_content;
		$output['vnic'] = $file2_content;

		$fp = fopen($file3, 'w');
		foreach($vnics_changed as $vnic){
			fputcsv($fp, $vnic);
		}
		fclose($fp);

		// create json file too
		file_put_contents($file_pre . '_VNICS_CHANGED.json',json_encode($vnics_changed));

	}

	return $output;

}



function powercli_restore_vm_vnics($vm_array, $vnics_changed){

	// GENERATE POWERCLI
	$file_content = "### MOVING VMS BACK TO DVS\n";
	$file_content .= "Write-Host VM vNICS THAT WERE CHANGED IN PREVIOUS PHASE WILL NOW BE CHANGED BACK TO DVS\n";
	$file_content .= '$confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "'."\n";
	$file_content .= 'if ($confirmation -eq \'y\') {'."\n";

	foreach ( $vnics_changed as $vnic ){

		$key = array_search($vnic['instance_uuid'], array_column($vm_array, 'instance_uuid'));

		if ( $key !== false ){
			$vm_moref = $vm_array[$key]['moref'];
			$vm_name = $pool_name = str_replace("'", "''", $vm_array[$key]['name']);
			$file_content .= "   Write-Host 'Changing {$vm_name} {$vnic['vnic']}'\n";
			$file_content .= "   Get-VM -Id VirtualMachine-{$vm_moref} | Get-NetworkAdapter -Name '{$vnic['vnic']}' | Set-NetworkAdapter -NetworkName '{$vnic['portgroup']}' -Confirm:\$false -RunAsync\n";
		}

	}

	$file_content .= "}\n";
	return $file_content;

}

function powercli_restore_vm_templates($vm_array, $vm_source_array){

	// GENERATE POWERCLI
	$file_content = "### RESTORING VM TEMPLATES\n";
	$file_content .= "Write-Host VMS THAT WERE TEMPLATES WILL NOW BE CONVERTED BACK TO TEMPLATES\n";
	$file_content .= '$confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "'."\n";
	$file_content .= 'if ($confirmation -eq \'y\') {'."\n";

	foreach ( $vm_source_array as $vm_source ){

		if ($vm_source['template'] === 'true'){
			
			$key = array_search($vm_source['instance_uuid'], array_column($vm_array, 'instance_uuid'));
			if ( $key !== false ){
				$vm_moref = $vm_array[$key]['moref'];
				$file_content .= "   Get-VM -Id VirtualMachine-{$vm_moref} | Set-VM -ToTemplate\n";
			}
			

		}

	}

	$file_content .= "}\n";
	return $file_content;

}


function powercli_export_vapps($vcenter_id, $cluster_id){

	global $pdo;

	$query = "select name,moref from resourcepool 
	WHERE present=1 AND vapp_in_path = 0 AND type='VirtualApp' 
	AND vcenter_id='$vcenter_id' AND cluster_id='$cluster_id'";

	$sth = $pdo->prepare($query);
	$sth->execute();
	$result = $sth->fetchAll();

	$file_content = "### EXPORT VAPPS FROM SOURCE VCENTER\n";
	$file_content .= "Write-Host EXPORT VAPPS FROM SOURCE VCENTER TO LOCAL DISK. MAKE SURE SOURCE VAPPS ARE EMPTY!\n";
	$file_content .= '$confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "'."\n";
	$file_content .= 'if ($confirmation -eq \'y\') {'."\n";

	foreach ($result as $vapp){

		$file_content .= "   Get-VApp -Id 'VirtualApp-{$vapp['moref']}' | Stop-VApp -force \n";
		$file_content .= "   Get-VApp -Id 'VirtualApp-{$vapp['moref']}' | Export-VApp -destination '.\\vapps' \n";

	}

	$file_content .= "}\n";
	return $file_content;

}


function powercli_import_resourcepools($vcenter_id, $cluster_id){

	global $pdo;

	$query = "SELECT name,moref,full_path FROM resourcepool 
	WHERE present=1 AND type='ResourcePool' AND parent != 'cluster'
	AND vcenter_id='$vcenter_id' AND cluster_id='$cluster_id'
	ORDER BY full_path";

	$sth = $pdo->prepare($query);
	$sth->execute();
	$result = $sth->fetchAll();

	$file_content = "### IMPORT RESOURCEPOOL STRUCTURE FROM SOURCE VCENTER\n";
	$file_content .= "Write-Host IMPORTING RESOURCEPOOL STRUCTURE FROM SOURCE VCENTER\n";
	$file_content .= '$confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "'."\n";
	$file_content .= 'if ($confirmation -eq \'y\') {'."\n";


	$rpool_array = array();
	foreach ($result as $rpool){

		# create the 1st level resourcepools in cluster
		$paths = explode('/',$rpool['full_path']);
		if ( count( $paths ) == 3 ){
			$pool_name = str_replace("'", "''", $rpool['name']);
			$file_content .= "   \$pool_name = [System.Web.HttpUtility]::UrlDecode('$pool_name')\n";
			$file_content .= "   \$cluster = Get-Cluster -Name {$paths[1]}\n";
			$file_content .= "   New-ResourcePool -Location \$cluster -Name \$pool_name\n";
		} else {
			$rpool_array[] = $rpool;
		}
	}

	$file_content .= '   $resourcePoolArray = @()'."\n";

	foreach ($rpool_array as $rpool){

			# get parent path
			$full_path_array = explode('/',$rpool['full_path']);
			$parent_path_array = array_slice($full_path_array, 0, -1);
			$parent_path = str_replace( "'", "''", implode('/', $parent_path_array) );
			$rpool_parent_path = str_replace("'", "''", $parent_path);
			$pool_name = str_replace( "'", "''", $rpool['name'] );

			$file_content .= '   $rpool = New-Object System.Object'."\n";
			$file_content .= "   \$rpool | Add-Member -type NoteProperty -name Name -Value '$pool_name'\n";
			$file_content .= "   \$rpool | Add-Member -type NoteProperty -name ParentPath -Value '$parent_path'\n";
			$file_content .= '   $resourcePoolArray += $rpool'."\n";

	}

	$file_content .= '   Import-Cluster-ResoucePools $resourcePoolArray'."\n";
	$file_content .= "}\n";

	return $file_content;

}


function powercli_restore_vm_folders($vm_source_array, $vm_array){

	$vm_folder_move = array();

	foreach($vm_source_array as $vm){

		$folder = $vm['folder'];
		if ( !is_null($folder) ){
			$key = array_search($vm['instance_uuid'], array_column($vm_array, 'instance_uuid'));
			if ( $key !== false ){
				$vm_moref = $vm_array[$key]['moref'];
				$vm_name = $vm_array[$key]['name'];
				$vm_folder_move[$folder][] = array('moref' => $vm_moref, 'name' => $vm_name);
			}
		}

	}

	// GENERATE POWERCLI
	$file_content = "### MOVING VMS BACK TO ORIGINAL FOLDERS\n";
	$file_content .= "Write-Host VMs WILL NOW BE MOVED BACK TO ORIGINAL FOLDERS\n";
	$file_content .= '$confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "'."\n";
	$file_content .= 'if ($confirmation -eq \'y\') {'."\n";
	

	foreach ( $vm_folder_move as $folder => $vms ){

		$file_content .= '   $vmsByFolderArray = @()'."\n";

		foreach ($vms as $vm){
			$vm_moref = 'VirtualMachine-'.$vm['moref'];
			
			$vm_name = str_replace( "'", "''", $vm['name']);
			$file_content .= '   $vm = New-Object System.Object'."\n";
			$file_content .= "   \$vm | Add-Member -type NoteProperty -name Id -Value '$vm_moref'\n";
			$file_content .= "   \$vm | Add-Member -type NoteProperty -name Name -Value '$vm_name'\n";
			$file_content .= '   $vmsByFolderArray += $vm'."\n";

		}

		$vm_folder_path = str_replace( "'", "''", $folder);
		$file_content .= "   Restore-VMs-ByFolder '$vm_folder_path' \$vmsByFolderArray \n";

	}

	$file_content .= "}\n";
	return $file_content;

}

function powercli_import_vapps($cluster_name){

	$folder = 'vapps';
	$file_content = "### IMPORTING VAPPS FROM SOURCE VCENTER\n";
	$file_content .= "Write-Host IMPORTING VAPPS FROM SOURCE VCENTER TO DESTINATION VCENTER\n";
	$file_content .= '$confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "'."\n";
	$file_content .= 'if ($confirmation -eq \'y\') {'."\n";
	$file_content .= "  Import-Cluster-Vapps '$folder' '$cluster_name'\n";

	$file_content .= "}\n";
	return $file_content;
}


function powercli_restore_vm_rpools($VM_ARRAY, $VM_SOURCE_ARRAY, $RESOURCEPOOL_ARRAY, $VAPP_ARRAY){


	// GENERATE POWERCLI
	$file_content = "### MOVING VMS BACK TO ORIGINAL RESOURCEPOOLS AND VAPPS\n";
	$file_content .= "Write-Host MOVING VMS BACK TO ORIGINAL RESOURCEPOOLS AND VAPPS\n";
	$file_content .= '$confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "'."\n";
	$file_content .= 'if ($confirmation -eq \'y\') {'."\n";

	foreach( $VM_ARRAY as $vm ){

		// match migrated vms by instance uuid
		$key = array_search($vm['instance_uuid'], array_column($VM_SOURCE_ARRAY, 'instance_uuid'));

		// found VM
		if ( $key !== false  ){

			// check if its not a vapp
			if ( $VM_SOURCE_ARRAY[$key]['vapp_id'] == 'none' ){

				// move only if vm was not in root resourcepool before migration
				if ( $vm['rpool_full_path'] != $VM_SOURCE_ARRAY[$key]['rpool_full_path'] ){

					// find resourcepool new moref
					$key2 = array_search($VM_SOURCE_ARRAY[$key]['rpool_full_path'], array_column($RESOURCEPOOL_ARRAY, 'full_path'));
					$file_content .= "   \$pool = Get-ResourcePool -Id ResourcePool-{$RESOURCEPOOL_ARRAY[$key2]['moref']}\n";
					$file_content .= "   Get-VM -Id VirtualMachine-{$vm['moref']} | Move-VM -Destination \$pool \n";
				}

			} else {
				// find vapp new moref
				$key2 = array_search($VM_SOURCE_ARRAY[$key]['rpool_full_path'], array_column($VAPP_ARRAY, 'full_path'));
				if ( $key2 !== false  ){
					$file_content .= "   \$vapp = Get-vApp -Id VirtualApp-{$VAPP_ARRAY[$key2]['moref']}\n";
					$file_content .= "   Get-VM -Id VirtualMachine-{$vm['moref']} | Move-VM -Destination \$vapp \n";
				}
			}

		} else {
			$file_content .= "   Write-Host VM {$vm['name']} NOT FOUND!! {$vm['instance_uuid']}\n";
		}


	}

	$file_content .= "}\n";
	return $file_content;


}


/**
 * This file is part of the array_column library
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 *
 * @copyright Copyright (c) Ben Ramsey (http://benramsey.com)
 * @license http://opensource.org/licenses/MIT MIT
 */

if (!function_exists('array_column')) {
    /**
     * Returns the values from a single column of the input array, identified by
     * the $columnKey.
     *
     * Optionally, you may provide an $indexKey to index the values in the returned
     * array by the values from the $indexKey column in the input array.
     *
     * @param array $input A multi-dimensional array (record set) from which to pull
     *                     a column of values.
     * @param mixed $columnKey The column of values to return. This value may be the
     *                         integer key of the column you wish to retrieve, or it
     *                         may be the string key name for an associative array.
     * @param mixed $indexKey (Optional.) The column to use as the index/keys for
     *                        the returned array. This value may be the integer key
     *                        of the column, or it may be the string key name.
     * @return array
     */
    function array_column($input = null, $columnKey = null, $indexKey = null)
    {
        // Using func_get_args() in order to check for proper number of
        // parameters and trigger errors exactly as the built-in array_column()
        // does in PHP 5.5.
        $argc = func_num_args();
        $params = func_get_args();

        if ($argc < 2) {
            trigger_error("array_column() expects at least 2 parameters, {$argc} given", E_USER_WARNING);
            return null;
        }

        if (!is_array($params[0])) {
            trigger_error(
                'array_column() expects parameter 1 to be array, ' . gettype($params[0]) . ' given',
                E_USER_WARNING
            );
            return null;
        }

        if (!is_int($params[1])
            && !is_float($params[1])
            && !is_string($params[1])
            && $params[1] !== null
            && !(is_object($params[1]) && method_exists($params[1], '__toString'))
        ) {
            trigger_error('array_column(): The column key should be either a string or an integer', E_USER_WARNING);
            return false;
        }

        if (isset($params[2])
            && !is_int($params[2])
            && !is_float($params[2])
            && !is_string($params[2])
            && !(is_object($params[2]) && method_exists($params[2], '__toString'))
        ) {
            trigger_error('array_column(): The index key should be either a string or an integer', E_USER_WARNING);
            return false;
        }

        $paramsInput = $params[0];
        $paramsColumnKey = ($params[1] !== null) ? (string) $params[1] : null;

        $paramsIndexKey = null;
        if (isset($params[2])) {
            if (is_float($params[2]) || is_int($params[2])) {
                $paramsIndexKey = (int) $params[2];
            } else {
                $paramsIndexKey = (string) $params[2];
            }
        }

        $resultArray = array();

        foreach ($paramsInput as $row) {
            $key = $value = null;
            $keySet = $valueSet = false;

            if ($paramsIndexKey !== null && array_key_exists($paramsIndexKey, $row)) {
                $keySet = true;
                $key = (string) $row[$paramsIndexKey];
            }

            if ($paramsColumnKey === null) {
                $valueSet = true;
                $value = $row;
            } elseif (is_array($row) && array_key_exists($paramsColumnKey, $row)) {
                $valueSet = true;
                $value = $row[$paramsColumnKey];
            }

            if ($valueSet) {
                if ($keySet) {
                    $resultArray[$key] = $value;
                } else {
                    $resultArray[] = $value;
                }
            }

        }

        return $resultArray;
    }

}