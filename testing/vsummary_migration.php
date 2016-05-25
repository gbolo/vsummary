<?php

ini_set('display_errors', 1);
ini_set('display_startup_errors', 1);
error_reporting(E_ALL);

// MIGRATION VARIABLES
$source_vcenter_id = '0184679d-369a-4590-993a-5fbdf326a75a';
$source_datacenter_name = 'DC1';

// VSWITCH GENERATE
$esxi_hostname = '';
$vswitch_name = 'vSwitch1';

// DVS PORTGROUP EXPORT
$gen_dvs_portgroup = true;
$source_dvs_id = '5fb6de4be73d4154d746ca485eec9dae';
$destination_dvs_name = 'DVS2';

// FOLDER STRUCTURE EXPORT
$gen_folder_structure = true;



// SQL server connection information
require_once('src/api/lib/mysql_config.php');

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

function gen_vswitch_pg($src_dvs_id, $dst_svs_name){
	global $pdo;

	$query = "SELECT id, name, max_mtu, vcenter_id FROM vswitch WHERE id='$src_dvs_id'";
	$stmt = $pdo->query($query); 
	$dvs = $stmt->fetch(PDO::FETCH_ASSOC);
	//print_r($dvs);

	$esxi_pg = "";
	$esxi_vm = "";

	$query = "SELECT id, name FROM esxi WHERE present=1 AND vcenter_id='{$dvs['vcenter_id']}'";
	foreach($pdo->query($query) as $esxi){

		$query = "SELECT * FROM portgroup WHERE vswitch_id='{$dvs['id']}' AND present = 1";

		$esxi_pg .= "\n\n### CREATING NEW PORTGROUPS ON STANDARD VSWITCH FOR {$esxi['name']} ###\n";
		$esxi_vm .= "\n\n### CHANGING PORTGROUPS FOR EACH VM ON {$esxi['name']} ###\n";
		foreach($pdo->query($query) as $pg){
			if ( $pg['vlan_type'] === 'single' ){

				# create portgroup
				$esxi_pg .= "Get-VMHost {$esxi['name']} | Get-VirtualSwitch -Name '{$dst_svs_name}' | New-VirtualPortGroup -Name 'mig-{$pg['name']}' -VLanId {$pg['vlan']}\n";

			}
			# move every VM to it
			foreach($pdo->query("SELECT id, name FROM vm WHERE esxi_id = '{$esxi['id']}' AND present = 1") as $vm){

				# loop through vnics
				foreach($pdo->query("SELECT * FROM vnic WHERE vm_id = '{$vm['id']}' AND portgroup_id = '{$pg['id']}' AND present = 1") as $vnic){

					$esxi_vm .= "Get-VMHost {$esxi['name']} | Get-VM '{$vm['name']}' | Get-NetworkAdapter -Name '{$vnic['name']}' | Set-NetworkAdapter -NetworkName 'mig-{$pg['name']}' -Confirm:\$false -RunAsync\n";

				}		
			}

		}



	}

	echo $esxi_pg;
	echo $esxi_vm;

}

function gen_dvs_pg($dvs_id, $dvs_dst_name){
	global $pdo;
	$query = "SELECT * FROM portgroup WHERE vswitch_id='{$dvs_id}' AND present = 1";
	echo "\n\n### CREATING NEW PORTGROUPS ON DESTINATION DISTRIBUTED SWITCH ###\n";
	foreach($pdo->query($query) as $pg){
		if ( $pg['vlan_type'] === 'single' ){

			echo "Get-VDSwitch -Name '{$dvs_dst_name}' | New-VDPortgroup -Name '{$pg['name']}' -NumPorts 128 -VLanId {$pg['vlan']}\n";

		}
	}

}


function export_vm_folders($vcenter_id, $dc_name){
	
	global $pdo;

	$query = "select DISTINCT full_path, name from folder 
	WHERE full_path LIKE '{$dc_name}/%'
	AND type = 'VirtualMachine' 
	AND parent != 'datacenter' 
	AND vcenter_id = '${vcenter_id}' 
	AND present = 1;";

	$folders = array();
	foreach($pdo->query($query) as $folder){
		$folder_depth = count( explode("/", $folder['full_path']) );
		$folders[] = array("depth" => $folder_depth, "name" => $folder['name'], "path" => $folder['full_path']);
	}
	
	// sort by folder depth
	usort($folders, function($a, $b) {
    	return $a['depth'] - $b['depth'];
	});

	// generate powercli commands
	echo "\n\n### VIRTUAL MACHINE FOLDER IMPORT ###\n";
	echo '$folderArray = @()'."\n";
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
			echo "\$folderArray += '{$folder['path']}'\n";

		}

	}

	// echo powercli logic
	echo '
$folderArray | % {
 $startFolder = Get-Datacenter -Name $datacenter | Get-Folder -Name \'vm\' -NoRecursion
    $path = $_
 
    $location = $startFolder
    echo $location
    $path.Split(\'/\') | Select -skip 1 | %{
        $folder=$_
        Try {
            echo "GET: $folder LOC: $location"
            $location = Get-Folder -Name $folder -Location $location -NoRecursion -ErrorAction Stop
        }
        Catch{
            echo "NEW: $folder LOC: $location"
            $location = New-Folder -Name $folder -Location $location
        }
    } 
    echo "======="
}

';


}

// PREPERATION ON DESTINATION VCENTER
echo "###################################################\n";
echo "#    PREPARATION ON DESTINATION VCENTER SERVER    #\n";
echo "###################################################\n";
echo '
Add-PSSnapin VMware.VimAutomation.Core
If ($globale:DefaultVIServers) {
	Disconnect-VIServer -Server $global:DefaultVIServers -Force
}


$destVI = Read-Host "DESTINATION vCenter"
$datacenter = Read-Host "DESTINATION DataCenter Name"
$creds = get-credential
connect-viserver -server $destVI -Credential $creds
';

if ($gen_folder_structure){
	export_vm_folders($source_vcenter_id, $source_datacenter_name);
}

if ($gen_dvs_portgroup){
	gen_dvs_pg($source_dvs_id, $destination_dvs_name);
}

echo '
### DISCONNECT FROM DESTINATION VCENTER
Disconnect-VIServer "*" -confirm:$false


';


echo "###################################################\n";
echo "#        CHANGES ON SOURCE VCENTER SERVER         #\n";
echo "###################################################\n";
gen_vswitch_pg($source_dvs_id, $vswitch_name);




