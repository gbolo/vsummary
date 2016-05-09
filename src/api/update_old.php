<?php

/*

 Simple API (using the term API losely here) to receive the
 POST json data sent from the PowerCLI script and modify/insert it
 into the various mysql tables.

 **Disclaimer**
 This script currently does no error checking or validation. It expects to recieve 
 very specific post data.
 STILL UNDER HEAVY DEVELOPMENT

*/

$units = explode(' ', 'B KB MB GB TB PB');

function format_size($size) {
    global $units;

    $mod = 1024;

    for ($i = 0; $size > $mod; $i++) {
        $size /= $mod;
    }

    $endIndex = strpos($size, ".")+3;

    return substr( $size, 0, $endIndex).' '.$units[$i];
}


function InsertVM($data){

    $vcenter_id = $data[0]['vcenter_id'];
    $vm_query = "UPDATE vm SET present = 0 WHERE vcenter_id = '$vcenter_id' AND present = 1";
    mysql_query($vm_query) or die(mysql_error());

    foreach ($data as $vm) {

        $id = md5( $vm['vcenter_id'] . $vm['moref'] );
        $name = $vm['name'];
        $moref = $vm['moref'];
        $vmx_path = $vm['vmx_path'];
        $vcpu = $vm['vcpu'];
        $memory_mb = $vm['memory_mb'];
        $config_guest_os = $vm['config_guest_os'];
        $config_version = $vm['config_version'];
        $smbios_uuid = $vm['smbios_uuid'];
        $instance_uuid = $vm['instance_uuid'];
        $config_change_version = $vm['config_change_version'];
        $guest_tools_version = $vm['guest_tools_version'];
        $guest_tools_running = $vm['guest_tools_running'];
        $guest_hostname = $vm['guest_hostname'];
        $guest_ip = $vm['guest_ip'];
        $stat_cpu_usage = $vm['stat_cpu_usage'];
        $stat_host_memory_usage = $vm['stat_host_memory_usage'];
        $stat_guest_memory_usage = $vm['stat_guest_memory_usage'];
        $stat_uptime_sec = $vm['stat_uptime_sec'];
        $power_state = $vm['power_state'];
        $esxi_id = md5( $vm['vcenter_id'] . $vm['esxi_moref'] );
        $vcenter_id = $vm['vcenter_id'];

        if ( $guest_tools_running == 'guestToolsRunning' ) {
            $guest_tools_running = 'Yes';
        } else {
            $guest_tools_running = 'No';
            $guest_hostname = 'n/a';
            $guest_tools_version = 'n/a';
            $guest_ip = 'n/a';

        }


        $query = "INSERT INTO vm (id,name,moref,vmx_path,vcpu,memory_mb,config_guest_os,config_version,smbios_uuid,instance_uuid,config_change_version,guest_tools_version,guest_tools_running,guest_hostname,guest_ip,stat_cpu_usage,stat_host_memory_usage,stat_guest_memory_usage,stat_uptime_sec,power_state,esxi_id,vcenter_id) 
            VALUES ('$id', '$name', '$moref', '$vmx_path', '$vcpu', '$memory_mb', '$config_guest_os', '$config_version', '$smbios_uuid', '$instance_uuid', '$config_change_version', '$guest_tools_version', '$guest_tools_running', '$guest_hostname', '$guest_ip', '$stat_cpu_usage', '$stat_host_memory_usage', '$stat_guest_memory_usage', '$stat_uptime_sec', '$power_state', '$esxi_id', '$vcenter_id') 
            ON DUPLICATE KEY
            UPDATE name='$name', vmx_path='$vmx_path', vcpu='$vcpu', memory_mb='$memory_mb', config_guest_os='$config_guest_os', config_version='$config_version', smbios_uuid='$smbios_uuid', instance_uuid='$instance_uuid', config_change_version='$config_change_version', guest_tools_version='$guest_tools_version', guest_tools_running='$guest_tools_running', guest_hostname='$guest_hostname', guest_ip='$guest_ip', stat_cpu_usage='$stat_cpu_usage', stat_host_memory_usage='$stat_host_memory_usage', stat_guest_memory_usage='$stat_guest_memory_usage', stat_uptime_sec='$stat_uptime_sec', power_state='$power_state', esxi_id='$esxi_id', present=1";
        mysql_query($query) or die(mysql_error());


    }


}


function InsertESXi($data){

    $vcenter_id = $data[0]['vcenter_id'];
    $vm_query = "UPDATE esxi SET present = 0 WHERE vcenter_id = '$vcenter_id' AND present = 1";
    mysql_query($vm_query) or die(mysql_error());

    foreach ($data as $esxi) {

        $id = md5( $esxi['vcenter_id'] . $esxi['moref'] );
        $name = $esxi['name'];
        $moref = $esxi['moref'];
        $max_evc = $esxi['max_evc'];
        $current_evc = $esxi['current_evc'];
        $status = $esxi['status'];
        $power_state = $esxi['power_state'];
        $in_maintenance_mode = $esxi['in_maintenance_mode'];
        $vendor = $esxi['vendor'];
        $model = $esxi['model'];
        $uuid = $esxi['uuid'];
        $memory_bytes = $esxi['memory_bytes'];
        $cpu_model = $esxi['cpu_model'];
        $cpu_mhz = $esxi['cpu_mhz'];
        $cpu_sockets = $esxi['cpu_sockets'];
        $cpu_cores = $esxi['cpu_cores'];
        $cpu_threads = $esxi['cpu_threads'];
        $nics = $esxi['nics'];
        $hbas = $esxi['hbas'];
        $version = $esxi['version'];
        $build = $esxi['build'];
        $stat_cpu_usage = $esxi['stat_cpu_usage'];
        $stat_memory_usage = $esxi['stat_memory_usage'];
        $stat_uptime_sec = $esxi['stat_uptime_sec'];
        $vcenter_id = $esxi['vcenter_id'];

        $query = "INSERT INTO esxi (id,name,moref,max_evc,current_evc,status,power_state,in_maintenance_mode,vendor,model,uuid,memory_bytes,cpu_model,cpu_mhz,cpu_sockets,cpu_cores,cpu_threads,nics,hbas,version,build,stat_cpu_usage,stat_memory_usage,stat_uptime_sec,vcenter_id) 
            VALUES ('$id','$name','$moref','$max_evc','$current_evc','$status','$power_state','$in_maintenance_mode','$vendor','$model','$uuid','$memory_bytes','$cpu_model','$cpu_mhz','$cpu_sockets','$cpu_cores','$cpu_threads','$nics','$hbas','$version','$build','$stat_cpu_usage','$stat_memory_usage','$stat_uptime_sec','$vcenter_id') 
            ON DUPLICATE KEY
            UPDATE name='$name', max_evc='$max_evc', current_evc='$current_evc', status='$status', power_state='$power_state', in_maintenance_mode='$in_maintenance_mode', vendor='$vendor', model='$model', uuid='$uuid', memory_bytes='$memory_bytes', cpu_model='$cpu_model', cpu_mhz='$cpu_mhz', cpu_sockets='$cpu_sockets', cpu_cores='$cpu_cores', cpu_threads='$cpu_threads', nics='$nics', hbas='$hbas', version='$version', build='$build', stat_cpu_usage='$stat_cpu_usage', stat_memory_usage='$stat_memory_usage', stat_uptime_sec='$stat_uptime_sec', present=1";
        mysql_query($query) or die(mysql_error());


    }


}

function InsertDatastore($data){

    $vcenter_id = $data[0]['vcenter_id'];
    $vm_query = "UPDATE datastore SET present = 0 WHERE vcenter_id = '$vcenter_id' AND present = 1";
    mysql_query($vm_query) or die(mysql_error());

    foreach ($data as $datastore) {

        $id = md5( $datastore['vcenter_id'] . $datastore['moref'] );
        $name = $datastore['name'];
        $moref = $datastore['moref'];
        $status = $datastore['status'];
        $capacity_bytes = format_size( $datastore['capacity_bytes'] );
        $free_bytes = $datastore['free_bytes'];
        $uncommitted_bytes = $datastore['uncommitted_bytes'];
        $type = $datastore['type'];
        $vcenter_id = $datastore['vcenter_id'];

        $query = "INSERT INTO datastore (id,name,moref,status,capacity_bytes,free_bytes,uncommitted_bytes,type,vcenter_id) 
            VALUES ('$id','$name','$moref', '$status', '$capacity_bytes', '$free_bytes', '$uncommitted_bytes', '$type', '$vcenter_id') 
            ON DUPLICATE KEY
            UPDATE name='$name', moref='$moref', status='$status', capacity_bytes='$capacity_bytes', free_bytes='$free_bytes', uncommitted_bytes='$uncommitted_bytes', type='$type', present=1";
        mysql_query($query) or die(mysql_error());


    }


}


function InsertvDisk($data){

    $vcenter_id = $data[0]['vcenter_id'];
    $vm_query = "UPDATE vdisk SET present = 0 WHERE vcenter_id = '$vcenter_id' AND present = 1";
    mysql_query($vm_query) or die(mysql_error());

    foreach ($data as $vdisk) {

        $id = md5( $vdisk['vcenter_id'] . $vdisk['disk_object_id'] );
        $name = $vdisk['name'];
        $capacity_bytes = $vdisk['capacity_bytes'];
        $path = $vdisk['path'];
        $thin_provisioned = $vdisk['thin_provisioned'];
        $datastore_id = md5( $vdisk['vcenter_id'] . $vdisk['datastore_moref'] );
        $uuid = $vdisk['uuid'];
        $disk_object_id = $vdisk['disk_object_id'];
        $vm_id = md5( $vdisk['vcenter_id'] . $vdisk['vm_moref'] );
        $esxi_id= md5( $vdisk['vcenter_id'] . $vdisk['esxi_moref'] );
        $vcenter_id = $vdisk['vcenter_id'];

        $query = "INSERT INTO vdisk (id,name,capacity_bytes,path,thin_provisioned,datastore_id,uuid,disk_object_id,vm_id,esxi_id,vcenter_id) 
            VALUES ('$id', '$name', '$capacity_bytes', '$path', '$thin_provisioned', '$datastore_id', '$uuid', '$disk_object_id', '$vm_id', '$esxi_id', '$vcenter_id') 
            ON DUPLICATE KEY
            UPDATE name='$name',capacity_bytes='$capacity_bytes',path='$path',thin_provisioned='$thin_provisioned',datastore_id='$datastore_id',uuid='$uuid',vm_id='$vm_id',esxi_id='$esxi_id', present=1";
        mysql_query($query) or die(mysql_error());


    }


}


function InsertPnic($data){

    $vcenter_id = $data[0]['vcenter_id'];
    $vm_query = "UPDATE pnic SET present = 0 WHERE vcenter_id = '$vcenter_id' AND present = 1";
    mysql_query($vm_query) or die(mysql_error());

    foreach ($data as $pnic) {
        $name = $pnic['name'];
        $esxi_id = md5( $pnic['vcenter_id'] . $pnic['esxi_moref'] );
        $id = md5( $esxi_id . $name );
        $mac = $pnic['mac'];
        $link_speed = $pnic['link_speed'];
        $driver = $pnic['driver'];
        $vcenter_id = $pnic['vcenter_id'];

        $query = "INSERT INTO pnic (id,name,mac,link_speed,driver,esxi_id,vcenter_id) 
            VALUES ('$id','$name','$mac','$link_speed','$driver','$esxi_id','$vcenter_id') 
            ON DUPLICATE KEY
            UPDATE name='$name', mac='$mac', link_speed='$link_speed', driver='$driver', present=1";
        mysql_query($query) or die(mysql_error());


    }


}

function InsertDVS($data){

    $vcenter_id = $data[0]['vcenter_id'];
    $type = 'DVS';
    $vm_query = "UPDATE vswitch SET present = 0 WHERE vcenter_id = '$vcenter_id' AND type = '$type' AND present = 1";
    mysql_query($vm_query) or die(mysql_error());

    foreach ($data as $dvs) {
        $id = md5( $dvs['vcenter_id'] . $dvs['moref'] );
        $name = $dvs['name'];
        $version = $dvs['version'];
        $max_mtu = $dvs['max_mtu'];
        $ports = $dvs['ports'];
        $vcenter_id = $dvs['vcenter_id'];

        $query = "INSERT INTO vswitch (id, name, type, version, max_mtu, ports, vcenter_id) 
            VALUES ('$id', '$name', '$type', '$version', '$max_mtu', '$ports', '$vcenter_id') 
            ON DUPLICATE KEY
            UPDATE name='$name', type='$type', version='$version', max_mtu='$max_mtu', ports='$ports', present=1";
        mysql_query($query) or die(mysql_error());

    }

}


function InsertvSwitch($data){

    $vcenter_id = $data[0]['vcenter_id'];
    $type = 'vSwitch';
    $vm_query = "UPDATE vswitch SET present = 0 WHERE vcenter_id = '$vcenter_id' AND type = '$type' AND present = 1";
    mysql_query($vm_query) or die(mysql_error());

    foreach ($data as $vs) {
        $id = md5( $vs['vcenter_id'] . $vs['esxi_moref'] . $vs['name'] );
        $name = $vs['name'];
        $max_mtu = $vs['max_mtu'];
        $ports = $vs['ports'];
        $esxi_id = md5( $vs['vcenter_id'] . $vs['esxi_moref'] );
        $vcenter_id = $vs['vcenter_id'];

        $query = "INSERT INTO vswitch (id, name, type, esxi_id, max_mtu, ports, vcenter_id) 
            VALUES ('$id', '$name', '$type', '$esxi_id', '$max_mtu', '$ports', '$vcenter_id') 
            ON DUPLICATE KEY
            UPDATE name='$name', type='$type', esxi_id='$esxi_id', max_mtu='$max_mtu', ports='$ports', present=1";
        mysql_query($query) or die(mysql_error());

    }

}

function vlan_trunk_to_string($vlan_start, $vlan_end){

    $v_start = explode( " ", $vlan_start );
    $v_end = explode( " ", $vlan_end );

    $vlan = "";

    for ($i = 0; $i < count($v_start); $i++) {
        if ( $v_start[$i] == $v_end[$i] ){
            $vlan .=  "{$v_start[$i]}, ";
        } else {
            $vlan .=  "{$v_start[$i]} - {$v_end[$i]}, ";
        }
    }

    $vlan = rtrim($vlan, ", ");
    return $vlan;

}


function InsertDVSPortgroup($data){

##vswitch pg id can hash of vcenter uuid and pg name
    $vcenter_id = $data[0]['vcenter_id'];
    $type = 'DVS';
    $vm_query = "UPDATE portgroup SET present = 0 WHERE vcenter_id = '$vcenter_id' AND type = '$type' AND present = 1";
    mysql_query($vm_query) or die(mysql_error());
   

    foreach ($data as $pg) {
        $id = md5( $pg['vcenter_id'] . $pg['moref'] );
        $vswitch_id = md5( $pg['vcenter_id'] . $pg['dvs_moref'] );
        $name = $pg['name'];
        $vcenter_id = $pg['vcenter_id'];

        if ( $pg['vlan_type'] == "VmwareDistributedVirtualSwitchTrunkVlanSpec" ) {
            $vlan = vlan_trunk_to_string( $pg['vlan_start'], $pg['vlan_end'] );
            $vlan_type = "trunk";
        } else {
            $vlan = $pg['vlan'];
            $vlan_type = "single";
        }

        $query = "INSERT INTO portgroup (id, name, type, vlan, vlan_type, vswitch_id, vcenter_id) 
            VALUES ('$id', '$name', '$type', '$vlan', '$vlan_type', '$vswitch_id', '$vcenter_id') 
            ON DUPLICATE KEY
            UPDATE name='$name', type='$type', vlan='$vlan', vlan_type='$vlan_type', vswitch_id='$vswitch_id', present=1";
        mysql_query($query) or die(mysql_error());


    }


}


function InsertSVSPortgroup($data){

    $vcenter_id = $data[0]['vcenter_id'];
    $type = 'vSwitch';
    $vm_query = "UPDATE portgroup SET present = 0 WHERE vcenter_id = '$vcenter_id' AND type = '$type' AND present = 1";
    mysql_query($vm_query) or die(mysql_error());
   

    foreach ($data as $pg) {
        $esxi_id = md5( $pg['vcenter_id'] . $pg['esxi_moref'] );
        $id = md5( $pg['vcenter_id'] . $pg['esxi_moref'] . $pg['name'] );
        $vswitch_id = md5( $pg['vcenter_id'] . $pg['esxi_moref'] . $pg['vswitch_name'] );
        $vlan = $pg['vlan'];
        $vlan_type = 'single';
        $name = $pg['name'];
        $vcenter_id = $pg['vcenter_id'];


        $query = "INSERT INTO portgroup (id, name, type, vlan, vlan_type, vswitch_id, vcenter_id) 
            VALUES ('$id', '$name', '$type', '$vlan', '$vlan_type', '$vswitch_id', '$vcenter_id') 
            ON DUPLICATE KEY
            UPDATE name='$name', type='$type', vlan='$vlan', vlan_type='$vlan_type', vswitch_id='$vswitch_id', present=1";
        mysql_query($query) or die(mysql_error());


    }


}



function InsertvNIC($data){

    $vcenter_id = $data[0]['vcenter_id'];
    $type = 'vSwitch';
    $vm_query = "UPDATE vnic SET present = 0 WHERE vcenter_id = '$vcenter_id' AND present = 1";
    mysql_query($vm_query) or die(mysql_error());
   

    foreach ($data as $vnic) {


        if ( $vnic['vswitch_type'] === 'HostVirtualSwitch' ){
            # standard vswitch
            $portgroup_id = md5( $vnic['vcenter_id'] . $vnic['esxi_moref'] . $vnic['portgroup_name'] );
        } elseif ( $vnic['vswitch_type'] === 'VmwareDistributedVirtualSwitch' ) {
            # DVS
            $portgroup_id = md5( $vnic['vcenter_id'] . $vnic['portgroup_moref'] );
        } else {
            $portgroup_id = 'ORPHANED';
        }

        $id = md5( $vnic['vcenter_id'] . $vnic['vm_moref'] . $vnic['name'] );
        $vm_id = md5( $vnic['vcenter_id'] . $vnic['vm_moref'] );
        $name = $vnic['name'];
        $mac = $vnic['mac'];
        $type = $vnic['type'];
        $connected = var_export( $vnic['connected'], true );
        $status = $vnic['status'];
        $vcenter_id = $vnic['vcenter_id'];


        $query = "INSERT INTO vnic (id, name, mac, type, connected, status, vm_id, portgroup_id, vcenter_id) 
            VALUES ('$id', '$name', '$mac', '$type', '$connected', '$status', '$vm_id', '$portgroup_id', '$vcenter_id') 
            ON DUPLICATE KEY
            UPDATE name='$name', mac='$mac', type='$type', connected='$connected', status='$status', vm_id='$vm_id', portgroup_id='$portgroup_id', present=1";
        mysql_query($query) or die(mysql_error());


    }


}


function update_vcenter($data){

    $id = $data['vc_uuid'];
    $fqdn = $data['vc_fqdn'];
    $short_name = $data['vc_shortname'];

    $query = "INSERT INTO vcenter (id, fqdn, short_name) 
        VALUES ('$id', '$fqdn', '$short_name') 
        ON DUPLICATE KEY
        UPDATE fqdn='$fqdn', short_name='$short_name', present=1";
    mysql_query($query) or die(mysql_error());

}


$data = json_decode(file_get_contents('php://input'), true);

$sql_user = 'vsummary';
$sql_pass = 'changeme';
$sql_db   = 'vsummary';
$sql_host = 'localhost';

$dbhandle = mysql_connect($sql_host, $sql_user, $sql_pass)
  or die("Couldn't connect to SQL Server on $sql_host");

$selected = mysql_select_db($sql_db, $dbhandle)
  or die("Couldn't open database $sql_db");



if ( strcasecmp($data[0]['objecttype'],"VM")==0 ){
    InsertVM($data);
}
elseif ( strcasecmp($data[0]['objecttype'],"ESXI")==0 ){
    InsertESXi($data);
}
elseif ( strcasecmp($data[0]['objecttype'],"DS")==0 ){
    InsertDatastore($data);
}
elseif ( strcasecmp($data[0]['objecttype'],"VDISK")==0 ){
    InsertvDisk($data);
}
elseif ( strcasecmp($data[0]['objecttype'],"VNIC")==0 ){
    InsertvNIC($data);
}
elseif ( strcasecmp($data[0]['objecttype'],"PNIC")==0 ){
    InsertPnic($data);
}
elseif ( strcasecmp($data[0]['objecttype'],"DVS")==0 ){
    InsertDVS($data);
}
elseif ( strcasecmp($data[0]['objecttype'],"SVS")==0 ){
    InsertvSwitch($data);
}
elseif ( strcasecmp($data[0]['objecttype'],"DVSPG")==0 ){
    InsertDVSPortgroup($data);
}
elseif ( strcasecmp($data[0]['objecttype'],"SVSPG")==0 ){
    InsertSVSPortgroup($data);
}
elseif ( strcasecmp($data['objecttype'],"VCENTER")==0 ){
    update_vcenter($data);
}

mysql_close($dbhandle);







?>

