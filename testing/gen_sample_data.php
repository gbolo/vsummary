<?php
ini_set('display_errors', 1);
ini_set('display_startup_errors', 1);
error_reporting(E_ALL);

$dc_pre_name = array(
	"Toronto",
	"Vancouver",
	"NewYork",
	"Florida",
	"Chicago",
	"LosAngeles",
	"Dublin",
	"Amsterdam",
	"London",
	"Paris",
	"Berlin",
	"SaoPaulo",
	"HongKong",
	"Shanghai",
	"Montevideo",
	"Sydney",
	"Budapest",
	"Dubai",
	"Sousse",
	"Johannesburg"
);

$api_endpoint = 'http://vsummary.midgar.dev/api/update.php';
$vc_id = 'TEST'.md5(mt_rand());
$vc_type = strtolower( $dc_pre_name[mt_rand(0, count($dc_pre_name) - 1)] );
$vc_fqdn = 'vca.'.$vc_type.'.sample.tld';

$vm_pre_name = array(
	"nginx", 
	"mysql", 
	"percona",
	"coreos",
	"docker",
	"web",
	"db",
	"cassandra",
	"elasticsearch",
	"kibana",
	"graylog",
	"logstash",
	"rethinkdb",
	"apache",
	"haproxy",
	"influxdb",
	"sensu",
	"opennms",
	"nagios",
	"cacti",
	"seafile",
	"ipfs",
	"syncthing",
	"nfs",
	"mail",
	"app",
	"puppet",
	"dns",
	"dhcp",
	"testvm",
	"ftp",
	"centos",
	"nas"
);

$pg_pre_name = array(
	"app_private",
	"app_public",
	"dmz",
	"backend",
	"mgmt",
	"storage",
	"isolated",
	"frontend",
	"monitoring",
	"oob",
	"backups",
	"vpn",
	"untrusted"
);

$guest_os = array(
	"rhel7_64Guest",
	"rhel6_64Guest",
	"rhel5_64Guest",
	"other26xLinux64Guest",
	"otherLinux64Guest",
	"centos64Guest",
	"centosGuest",
	"debian6_64Guest",
	"debian7_64Guest",
	"debian8_64Guest",
	"fedora64Guest",
	"sles12_64Guest",
	"sles11_64Guest",
	"opensuse64Guest"
);

$rp_name = array(
	"High",
	"Medium",
	"Low",
	"Verylow"
);

$folder_L1 = array(
	"DEVOPS",
	"NETWORKING",
	"DATABASE",
	"SECURITY",
	"WEB",
	"MAIL",
	"OTHER",
	"SHARED"
);

$folder_L2 = array(
	"testing",
	"production",
	"staging",
	"misc"
);

function vsummary_api_call($data){
                              
    global $api_endpoint;                                  
	$data_string = json_encode($data);                                                                                   
	                                                                                                                     
	$ch = curl_init($api_endpoint);                                                                      
	curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "POST");                                                                     
	curl_setopt($ch, CURLOPT_POSTFIELDS, $data_string);                                                                  
	curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);                                                                      
	curl_setopt($ch, CURLOPT_HTTPHEADER, array(                                                                          
	    'Content-Type: application/json',                                                                                
	    'Content-Length: ' . strlen($data_string))                                                                       
	);                                                                                                                   
	                                                                                                                     
	curl_exec($ch);

	// Check if any error occurred
	if (!curl_errno($ch)) {
	  	$info = curl_getinfo($ch);
	  	$status = ($info['http_code'] == 200 ? 'SUCCESS' : 'FAILED');
	  	$result = $status . '! RESPONSE: ' . $info['http_code'] . ' TIME: ' . $info['total_time'] . "s\n";
	  	return $result;
	} else {
		return "API REQUEST FAILED\n";
	}
	
}

function gen_mac(){
	$mac = implode(':',str_split(substr(md5(mt_rand()),0,12),2));
	return $mac;
}

function gen_datacenter($num){

	global $vc_id;
	global $dc_pre_name;
	$arr;

	for ($i = 1; $i <= $num; $i++) {

		$name = $dc_pre_name[mt_rand(0, count($dc_pre_name) - 1)].'-DC'.rand(1, 9);

		$json = '
	    {
	        "vcenter_id":  "'.$vc_id.'",
	        "name":  "'.$name.'",
	        "esxi_folder_moref":  "group-h'.$i.'",
	        "objecttype":  "DC",
	        "vm_folder_moref":  "group-v'.$i.'",
	        "moref":  "datacenter-'.$i.'"
	    }';
		$arr[] = json_decode($json, true);
	}
	return $arr;
}

function gen_folder($dc){

	global $vc_id;
	global $folder_L1;
	global $folder_L2;

	$folders = [];

	// dc root folder must be present
	$json = '
    {
        "vcenter_id":  "'.$vc_id.'",
        "name":  "vm",
        "objecttype":  "FOLDER",
        "parent_moref":  "'.$dc['moref'].'",
        "type":  "Folder VirtualMachine VirtualApp",
        "moref":  "'.$dc['vm_folder_moref'].'"
    }';
    $root_folder[] = json_decode($json, true);

	$i = 100;
	foreach ($folder_L1 as $name){
		$json = '
	    {
	        "vcenter_id":  "'.$vc_id.'",
	        "name":  "'.$name.'",
	        "objecttype":  "FOLDER",
	        "parent_moref":  "'.$dc['vm_folder_moref'].'",
	        "type":  "Folder VirtualMachine VirtualApp",
	        "moref":  "group-v'.$i.'"
	    }';
	    $arrL1[] = json_decode($json, true);
	    $i++;
	}
	
	foreach ($arrL1 as $folder){

		$name = $folder_L2[mt_rand(0, count($folder_L2) - 1)];
		$json = '
	    {
	        "vcenter_id":  "'.$vc_id.'",
	        "name":  "'.$name.'",
	        "objecttype":  "FOLDER",
	        "parent_moref":  "'.$folder['moref'].'",
	        "type":  "Folder VirtualMachine VirtualApp",
	        "moref":  "group-v'.$i.'"
	    }';
	    $arrL2[] = json_decode($json, true);
	    $i++;
	}

	$folders = array_merge($root_folder, $arrL1, $arrL2);
	return $folders;
}

function gen_cluster($num, $dc){

	global $vc_id;

	for ($i = 1; $i <= $num; $i++) {

		$json = '
	    {
	        "ha_enabled":  "True",
	        "target_balance":  200,
	        "total_cpu_mhz":  10784,
	        "moref":  "domain-c'.$i.'",
	        "drs_enabled":  "True",
	        "name":  "CLUSTER-'.$i.'",
	        "total_vmotions":  '.rand(50, 5000).',
	        "total_memory_bytes":  17178796032,
	        "status":  "green",
	        "num_hosts":  2,
	        "datacenter_moref":  "'.$dc['esxi_folder_moref'].'",
	        "current_balance":  '.rand(10, 180).',
	        "drs_behaviour":  2,
	        "total_cpu_threads":  4,
	        "objecttype":  "CLUSTER",
	        "vcenter_id":  "'.$vc_id.'"
	    }';
		$arr[] = json_decode($json, true);
	}

	return $arr;
}

function gen_resourcepool($num, $cl){

	global $vc_id;
	global $rp_name;

	for ($i = 1; $i <= $num; $i++) {

		$name = $rp_name[mt_rand(0, count($rp_name) - 1)].'-'.$i;

		$json = '
	    {
	        "mem_reservation":  0,
	        "vapp_state":  "n/a",
	        "objecttype":  "RES",
	        "cpu_limit":  -1,
	        "mem_limit":  -1,
	        "moref":  "resgroup-'.$i.'",
	        "parent_moref":  "resgroup-1",
	        "name":  "'.$name.'",
	        "type":  "ResourcePool",
	        "status":  1,
	        "cluster_moref":  "'.$cl['moref'].'",
	        "vcenter_id":  "'.$i.'",
	        "configured_memory_mb":  2048,
	        "cpu_reservation":  0
	    }';

		$arr[] = json_decode($json, true);
	}
	return $arr;

}

function gen_esxi($num, $cl){
	global $vc_type;
	global $vc_id;
	$arr;
	for ($i = 1; $i <= $num; $i++) {

		$json = '
		{
		    "max_evc":  "intel-sandybridge",
		    "cpu_threads":  12,
		    "name":  "esxi-'.$i.'.'.$vc_type.'.linuxctl.com",
		    "vendor":  "Supermicro",
		    "stat_memory_usage":  29971,
		    "version":  "6.0.0",
		    "model":  "X9SRA/X9SRA-3",
		    "vcenter_id":  "'.$vc_id.'",
		    "objecttype":  "ESXI",
		    "build":  "3073146",
		    "cpu_sockets":  1,
		    "current_evc":  "intel-sandybridge",
		    "cpu_model":  "Intel(R) Xeon(R) CPU E5-2630L 0 @ 2.00GHz",
		    "nics":  4,
		    "power_state":  0,
		    "in_maintenance_mode":  "false",
		    "cpu_mhz":  1999,
		    "hbas":  3,
		    "stat_uptime_sec":  '.rand(2017352, 99017352).',
		    "cpu_cores":  6,
		    "status":  2,
		    "stat_cpu_usage":  '.rand(500, 5000).',
		    "cluster_moref":  "'.$cl['moref'].'",
		    "moref":  "host-'.rand(1, 1000).'",
		    "memory_bytes":  103043387392,
		    "uuid":  "'.md5(uniqid()).'"
		}';

		$arr[] = json_decode($json, true);
	}
	return $arr;
}

function gen_dvs($num){
	global $vc_id;
	global $vc_type;
	$arr;
	for ($i = 1; $i <= $num; $i++) {

		$json = '
	    {
	        "version":  "6.0.0",
	        "name":  "DVS-'.$vc_type.'-'.$i.'",
	        "objecttype":  "DVS",
	        "max_mtu":  1500,
	        "vcenter_id":  "'.$vc_id.'",
	        "ports":  '.rand(64, 2048).',
	        "moref":  "dvs'.rand(1, 100).'"
	    }';
		$arr[] = json_decode($json, true);
	}
	return $arr;
}

function gen_ds($num){
	global $vc_id;
	$arr;
	for ($i = 1; $i <= $num; $i++) {

		$json = '
		{
		    "moref":  "datastore-'.rand(1, 1000).'",
		    "type":  "VMFS",
		    "objecttype":  "DS",
		    "vcenter_id":  "'.$vc_id.'",
		    "capacity_bytes":  999922073600,
		    "name":  "DS'.$i.'-1TB",
		    "uncommitted_bytes":  '.rand(102327252475, 302327252475).',
		    "free_bytes":  '.rand(102327252475, 602327252475).',
		    "status":  1
		}';

		$arr[] = json_decode($json, true);
	}
	return $arr;
}

function gen_pg($num, $dvs){
	global $pg_pre_name;
	$arr;
	for ($i = 1; $i <= $num; $i++) {
		$name = $pg_pre_name[array_rand($pg_pre_name)]."-".$i;
		$json = '
        {
	        "moref":  "dvportgroup-'.$i.'",
	        "vlan_type":  "VmwareDistributedVirtualSwitchVlanIdSpec",
	        "vlan":  '.rand(5,4000).',
	        "vlan_end":  "na",
	        "objecttype":  "DVSPG",
	        "vlan_start":  "na",
	        "name":  "'.$name.'",
	        "vcenter_id":  "'.$dvs['vcenter_id'].'",
	        "dvs_moref":  "'.$dvs['moref'].'"
	    }';
		$arr[] = json_decode($json, true);
	}
	return $arr;
}

function gen_vm($num, $esxi){
	global $vc_id;
	global $vc_type;
	global $vm_pre_name;
	global $guest_os;
	global $fd_total;
	$arr;
	for ($i = 0; $i < $num; $i++) {
		$name = $vm_pre_name[mt_rand(0, count($vm_pre_name) - 1)].'-'.rand(1, 9);
		$ram = 1024 * rand(1, 16);
		$power_state = rand(0, 1);
		$guest = $guest_os[array_rand($guest_os)];
	 	$date = new DateTime(date('Y-m-d', strtotime('-'.rand(1, 365).' days')));
	 	$folder = $fd_total[mt_rand(0, count($fd_total) - 1)];
	 	while ( $folder['name'] === 'vm' ){
	 		$folder = $fd_total[mt_rand(0, count($fd_total) - 1)];
	 	}
		$config_date = date_format($date, 'Y-m-d\TH:i:s.').rand(100000, 400000).'Z';
		if ($power_state == 1){
			$vmtools = 'guestToolsRunning';
		} else {
			$vmtools = 'guestToolsNotRunning';
		}
		$json = '
		    {
		        "name":  "'.$name.'",
		        "moref":  "vm-'.rand(1, 10000).'",
		        "vmx_path":  "[SAMPLE] '.$name.'/'.$name.'.vmx",
		        "vcpu":  '.rand(1, 4).',
		        "memory_mb":  '.$ram.',
		        "config_guest_os":  "'.$guest.'",
		        "config_version":  "vmx-0'.rand(7, 9).'",
		        "resourcepool_moref": "resgroup-1",
		        "folder_moref":  "'.$folder['moref'].'",
		        "vapp_moref":  null,
		        "template":  false,
		        "smbios_uuid":  "'.md5(uniqid()).'",
		        "instance_uuid":  "'.md5(uniqid()).'",
		        "config_change_version":  "'.$config_date.'",
		        "guest_tools_version":  "'.rand(9000, 10000).'",
		        "guest_tools_running":  "'.$vmtools.'",
		        "guest_hostname":  "'.$name.'",
		        "guest_ip":  "10.'.rand(1, 254).'.'.rand(1, 254).'.'.rand(1, 254).'",
		        "stat_cpu_usage":  '.rand(1, 200).',
		        "stat_host_memory_usage":  '.rand(1, 1000).',
		        "stat_guest_memory_usage":  '.rand(1, 300).',
		        "stat_uptime_sec":  '.rand(10000, 31557600).',
		        "power_state":  '.$power_state.',
		        "esxi_moref":  "'.$esxi['moref'].'",
		        "vcenter_id":  "'.$vc_id.'",
		        "objecttype":  "VM"
		    }';
		$arr[] = json_decode($json, true);
	}
	return $arr;
}

function gen_vnic($vm){
	global $pg_total;
	global $vc_id;
	$arr = [];
	$num = rand(1, 2);
	for ($i = 1; $i <= $num; $i++) {
		$pg = $pg_total[mt_rand(0, count($pg_total) - 1)];
		$json = '
	    {
	        "portgroup_name":  "'.$pg['name'].'",
	        "esxi_moref":  "'.$vm['esxi_moref'].'",
	        "mac":  "'.gen_mac().'",
	        "objecttype":  "VNIC",
	        "vm_moref":  "'.$vm['moref'].'",
	        "name":  "Network adapter '.$i.'",
	        "type":  "VirtualVmxnet3",
	        "vswitch_type":  "VmwareDistributedVirtualSwitch",
	        "status":  "ok",
	        "vcenter_id":  "'.$vc_id.'",
	        "portgroup_moref":  "'.$pg['moref'].'",
	        "connected":  true,
	        "vswitch_name":  "SAMPLE"
	    }';
		$arr[] = json_decode($json, true);
	}
	return $arr;
}

function gen_vdisk($vm){
	global $ds_total;
	global $vc_id;
	$arr = [];
	$num = rand(1, 2);
	for ($i = 1; $i <= $num; $i++) {
		$ds = $ds_total[mt_rand(0, count($ds_total) - 1)];
		$json = '
	    {
	        "path":  "['.$ds['name'].'] '.$vm['name'].'/'.$vm['name'].'.vmdk",
	        "datastore_moref":  "'.$ds['moref'].'",
	        "objecttype":  "VDISK",
	        "thin_provisioned":  true,
	        "esxi_moref":  "'.$vm['esxi_moref'].'",
	        "capacity_bytes":  '.rand(10000000000, 50000000000).',
	        "name":  "Hard disk '.$i.'",
	        "vm_moref":  "'.$vm['moref'].'",
	        "vcenter_id":  "'.$vc_id.'",
	        "uuid":  "'.md5(uniqid()).'",
	        "disk_object_id":  "'.rand(10, 90).'-'.rand(1000, 5000).'"
	    }';
		$arr[] = json_decode($json, true);
	}
	return $arr;
}

function gen_pnic($esxi){
	global $vc_id;
	$arr = [];
	$num = rand(2, 6);
	for ($i = 1; $i <= $num; $i++) {
		$json = '
	    {
	        "vcenter_id":  "'.$vc_id.'",
	        "name":  "vmnic'.$i.'",
	        "esxi_moref":  "'.$esxi['moref'].'",
	        "objecttype":  "PNIC",
	        "mac":  "'.gen_mac().'",
	        "link_speed":  10000,
	        "driver":  "ixgbe"
	    }';
		$arr[] = json_decode($json, true);
	}
	return $arr;
}


//========================//
// START GENERATION LOGIC //
//========================//

// !!support only 1 dc per vcenter for random generation!!
$dc_total = gen_datacenter(1);

$cl_total = [];
$fd_total = [];
foreach ($dc_total as $dc){
	$n = rand(1, 3);
	$cl = gen_cluster($n,$dc);
	$cl_total = array_merge($cl_total, $cl);

	$fd = gen_folder($dc);
	$fd_total = array_merge($fd_total, $fd);
}

$esxi_total = [];
$rp_total = [];
foreach ($cl_total as $cl){
	$n = rand(2, 12);
	$esxi = gen_esxi($n,$cl);
	$esxi_total = array_merge($esxi_total, $esxi);

	$n = rand(1, 3);
	$rp = gen_resourcepool($n, $cl);
	$rp_total = array_merge($rp_total, $rp);
}

$dvs_total = gen_dvs(2);
$ds_total = gen_ds(rand(4, 12));

$pg_total = [];
foreach ($dvs_total as $dvs){
	$n = rand(5, 20);
	$pg = gen_pg($n, $dvs);
	$pg_total = array_merge($pg_total, $pg);
}

$vm_total = [];
$pnic_total = [];
foreach ($esxi_total as $esxi){
	$n = rand(5, 20);
	$vm = gen_vm($n, $esxi);
	$vm_total = array_merge($vm_total, $vm);
	$pnic = gen_pnic($esxi);
	$pnic_total = array_merge($pnic_total, $pnic);
}

$vnic_total = [];
$vdisk_total = [];
foreach ($vm_total as $vm){
	$vnic = gen_vnic($vm);
	$vnic_total = array_merge($vnic_total, $vnic);
	$vdisk = gen_vdisk($vm);
	$vdisk_total = array_merge($vdisk_total, $vdisk);
}

$vcenter_arr = array(
    "vc_uuid" => $vc_id,
    "objecttype" =>  "VCENTER",
    "vc_shortname" => strtoupper($vc_type),
    "vc_fqdn" =>  $vc_fqdn
);


echo "POSTING RANDOM SAMPLE DATA FOR VSUMMARY API: $api_endpoint\n---\n";
echo '[vcenter] ' . vsummary_api_call($vcenter_arr);
echo '[datacenter] ' . vsummary_api_call($dc_total);
echo '[cluster] ' . vsummary_api_call($cl_total);
echo '[resourcepool] ' . vsummary_api_call($rp_total);
echo '[esxi] ' . vsummary_api_call($esxi_total);
echo '[dvs] ' . vsummary_api_call($dvs_total);
echo '[datastore] ' . vsummary_api_call($ds_total);
echo '[vm] ' . vsummary_api_call($vm_total);
echo '[portgroup] ' . vsummary_api_call($pg_total);
echo '[pnic] ' . vsummary_api_call($pnic_total);
echo '[vnic] ' . vsummary_api_call($vnic_total);
echo '[vdisk] ' . vsummary_api_call($vdisk_total);
echo '[folder] ' . vsummary_api_call($fd_total);

