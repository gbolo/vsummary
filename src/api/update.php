<?php
/*

 Simple API (using the term API losely here) to receive the
 POST json data sent from the PowerCLI script and modify/insert it
 into the various mysql tables.

 - each POST will affect a single table and executed in a single transaction
 - failed transactions or bad POST data will result in 500 response code

 **Disclaimer**
 This script currently does not do much error checking or validation. 
 It expects to recieve very specific post data.
 STILL UNDER HEAVY DEVELOPMENT

*/

ini_set('display_errors', 1);
ini_set('display_startup_errors', 1);
error_reporting(E_ALL);

// functions required

function update_vcenter($data){

    $id = $data['vc_uuid'];
    $fqdn = $data['vc_fqdn'];
    $short_name = $data['vc_shortname'];

    try {

        // grab the pdo object declared outside of this function
        global $pdo;

        // start transaction
        $pdo->beginTransaction();

        // prepare statement to avoid sql injections
        $stmt = $pdo->prepare('INSERT INTO vcenter (id, fqdn, short_name) ' . 
                'VALUES(:id, :fqdn, :short_name) ' .
                'ON DUPLICATE KEY UPDATE fqdn=VALUES(fqdn), short_name=VALUES(short_name)');

        $stmt->bindParam(':id', $id, PDO::PARAM_STR);
        $stmt->bindParam(':fqdn', $fqdn, PDO::PARAM_STR);
        $stmt->bindParam(':short_name', $short_name, PDO::PARAM_STR);

        // execute prepared statement
        $stmt->execute();

        // commit transaction
        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        http_response_code(500);
    }
}

function update_vnic($data){

    $vcenter_id = $data[0]['vcenter_id'];
    $type = 'vSwitch';

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE vnic SET present = 0 WHERE present = 1 AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO vnic (id, name, mac, type, connected, status, vm_id, portgroup_id, vcenter_id) ' . 
                    'VALUES(:id, :name, :mac, :type, :connected, :status, :vm_id, :portgroup_id, :vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE name=VALUES(name), mac=VALUES(mac), type=VALUES(type), connected=VALUES(connected), status=VALUES(status), vm_id=VALUES(vm_id), portgroup_id=VALUES(portgroup_id), present=1');

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

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':mac', $mac, PDO::PARAM_STR);
                $stmt->bindParam(':type', $type, PDO::PARAM_STR);
                $stmt->bindParam(':connected', $connected, PDO::PARAM_STR);
                $stmt->bindParam(':status', $status, PDO::PARAM_STR);
                $stmt->bindParam(':vm_id', $vm_id, PDO::PARAM_STR);
                $stmt->bindParam(':portgroup_id', $portgroup_id, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        http_response_code(500);
    }
}

function update_esxi($data){

    $vcenter_id = $data[0]['vcenter_id'];

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE esxi SET present = 0 WHERE present = 1 AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO esxi (id,name,cluster_id,moref,max_evc,current_evc,status,power_state,in_maintenance_mode,vendor,model,uuid,memory_bytes,cpu_model,cpu_mhz,cpu_sockets,cpu_cores,cpu_threads,nics,hbas,version,build,stat_cpu_usage,stat_memory_usage,stat_uptime_sec,vcenter_id) ' . 
                    'VALUES(:id,:name,:cluster_id,:moref,:max_evc,:current_evc,:status,:power_state,:in_maintenance_mode,:vendor,:model,:uuid,:memory_bytes,:cpu_model,:cpu_mhz,:cpu_sockets,:cpu_cores,:cpu_threads,:nics,:hbas,:version,:build,:stat_cpu_usage,:stat_memory_usage,:stat_uptime_sec,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE id=VALUES(id),name=VALUES(name),cluster_id=VALUES(cluster_id),moref=VALUES(moref),max_evc=VALUES(max_evc),current_evc=VALUES(current_evc),status=VALUES(status),power_state=VALUES(power_state),in_maintenance_mode=VALUES(in_maintenance_mode),vendor=VALUES(vendor),model=VALUES(model),uuid=VALUES(uuid),memory_bytes=VALUES(memory_bytes),cpu_model=VALUES(cpu_model),cpu_mhz=VALUES(cpu_mhz),cpu_sockets=VALUES(cpu_sockets),cpu_cores=VALUES(cpu_cores),cpu_threads=VALUES(cpu_threads),nics=VALUES(nics),hbas=VALUES(hbas),version=VALUES(version),build=VALUES(build),stat_cpu_usage=VALUES(stat_cpu_usage),stat_memory_usage=VALUES(stat_memory_usage),stat_uptime_sec=VALUES(stat_uptime_sec),vcenter_id=VALUES(vcenter_id),present=1');

            foreach ($data as $esxi) {

                $id = md5( $esxi['vcenter_id'] . $esxi['moref'] );
                $cluster_id = md5( $esxi['vcenter_id'] . $esxi['cluster_moref'] );
                $name = $esxi['name'];
                $moref = $esxi['moref'];
                $max_evc = $esxi['max_evc'];
                if ( is_null($esxi['current_evc']) ){
                    $current_evc = var_export( $esxi['current_evc'], true );
                } else {
                    $current_evc = $esxi['current_evc'];
                }
                $status = $esxi['status'];
                $power_state = $esxi['power_state'];
                if ( is_null($esxi['in_maintenance_mode']) ){
                    $in_maintenance_mode = var_export( $esxi['in_maintenance_mode'], true );
                } else {
                    $in_maintenance_mode = $esxi['in_maintenance_mode'];
                }
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

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':cluster_id', $cluster_id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':moref', $moref, PDO::PARAM_STR);
                $stmt->bindParam(':max_evc', $max_evc, PDO::PARAM_STR);
                $stmt->bindParam(':current_evc', $current_evc, PDO::PARAM_STR);
                $stmt->bindParam(':status', $status, PDO::PARAM_STR);
                $stmt->bindParam(':power_state', $power_state, PDO::PARAM_STR);
                $stmt->bindParam(':in_maintenance_mode', $in_maintenance_mode, PDO::PARAM_STR);
                $stmt->bindParam(':vendor', $vendor, PDO::PARAM_STR);
                $stmt->bindParam(':model', $model, PDO::PARAM_STR);
                $stmt->bindParam(':uuid', $uuid, PDO::PARAM_STR);
                $stmt->bindParam(':memory_bytes', $memory_bytes, PDO::PARAM_STR);
                $stmt->bindParam(':cpu_model', $cpu_model, PDO::PARAM_STR);
                $stmt->bindParam(':cpu_mhz', $cpu_mhz, PDO::PARAM_STR);
                $stmt->bindParam(':cpu_sockets', $cpu_sockets, PDO::PARAM_STR);
                $stmt->bindParam(':cpu_cores', $cpu_cores, PDO::PARAM_STR);
                $stmt->bindParam(':cpu_threads', $cpu_threads, PDO::PARAM_STR);
                $stmt->bindParam(':nics', $nics, PDO::PARAM_STR);
                $stmt->bindParam(':hbas', $hbas, PDO::PARAM_STR);
                $stmt->bindParam(':version', $version, PDO::PARAM_STR);
                $stmt->bindParam(':build', $build, PDO::PARAM_STR);
                $stmt->bindParam(':stat_cpu_usage', $stat_cpu_usage, PDO::PARAM_STR);
                $stmt->bindParam(':stat_memory_usage', $stat_memory_usage, PDO::PARAM_STR);
                $stmt->bindParam(':stat_uptime_sec', $stat_uptime_sec, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
}

function update_vm($data){

    $vcenter_id = $data[0]['vcenter_id'];

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE vm SET present = 0 WHERE present = 1 AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO vm (id,name,moref,vmx_path,vcpu,memory_mb,template,config_guest_os,config_version,smbios_uuid,instance_uuid,config_change_version,guest_tools_version,guest_tools_running,guest_hostname,guest_ip,stat_cpu_usage,stat_host_memory_usage,stat_guest_memory_usage,stat_uptime_sec,power_state,folder_id,vapp_id,resourcepool_id,esxi_id,vcenter_id) ' . 
                    'VALUES(:id,:name,:moref,:vmx_path,:vcpu,:memory_mb,:template,:config_guest_os,:config_version,:smbios_uuid,:instance_uuid,:config_change_version,:guest_tools_version,:guest_tools_running,:guest_hostname,:guest_ip,:stat_cpu_usage,:stat_host_memory_usage,:stat_guest_memory_usage,:stat_uptime_sec,:power_state,:folder_id,:vapp_id,:resourcepool_id,:esxi_id,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE id=VALUES(id),name=VALUES(name),moref=VALUES(moref),vmx_path=VALUES(vmx_path),vcpu=VALUES(vcpu),memory_mb=VALUES(memory_mb),template=VALUES(template),config_guest_os=VALUES(config_guest_os),config_version=VALUES(config_version),smbios_uuid=VALUES(smbios_uuid),instance_uuid=VALUES(instance_uuid),config_change_version=VALUES(config_change_version),guest_tools_version=VALUES(guest_tools_version),guest_tools_running=VALUES(guest_tools_running),guest_hostname=VALUES(guest_hostname),guest_ip=VALUES(guest_ip),stat_cpu_usage=VALUES(stat_cpu_usage),stat_host_memory_usage=VALUES(stat_host_memory_usage),stat_guest_memory_usage=VALUES(stat_guest_memory_usage),stat_uptime_sec=VALUES(stat_uptime_sec),power_state=VALUES(power_state),folder_id=VALUES(folder_id),vapp_id=VALUES(vapp_id),resourcepool_id=VALUES(resourcepool_id),esxi_id=VALUES(esxi_id),vcenter_id=VALUES(vcenter_id),present=1');

            foreach ($data as $vm) {

                $id = md5( $vm['vcenter_id'] . $vm['moref'] );
                if ( is_null($vm['vapp_moref']) ){
                    $vapp_id = 'none';
                    $folder_id = md5( $vm['vcenter_id'] . $vm['folder_moref'] );
                } else {
                    $vapp_id = md5( $vm['vcenter_id'] . $vm['vapp_moref'] );
                    $folder_id = 'vapp';
                }
                if ( is_null($vm['resourcepool_moref']) ){
                    $resourcepool_id = 'none';
                } else {
                    $resourcepool_id = md5( $vm['vcenter_id'] . $vm['resourcepool_moref'] );
                }
                $template = var_export( $vm['template'], true );
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

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':vapp_id', $vapp_id, PDO::PARAM_STR);
                $stmt->bindParam(':folder_id', $folder_id, PDO::PARAM_STR);
                $stmt->bindParam(':resourcepool_id', $resourcepool_id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':moref', $moref, PDO::PARAM_STR);
                $stmt->bindParam(':vmx_path', $vmx_path, PDO::PARAM_STR);
                $stmt->bindParam(':vcpu', $vcpu, PDO::PARAM_STR);
                $stmt->bindParam(':memory_mb', $memory_mb, PDO::PARAM_STR);
                $stmt->bindParam(':template', $template, PDO::PARAM_STR);
                $stmt->bindParam(':config_guest_os', $config_guest_os, PDO::PARAM_STR);
                $stmt->bindParam(':config_version', $config_version, PDO::PARAM_STR);
                $stmt->bindParam(':smbios_uuid', $smbios_uuid, PDO::PARAM_STR);
                $stmt->bindParam(':instance_uuid', $instance_uuid, PDO::PARAM_STR);
                $stmt->bindParam(':config_change_version', $config_change_version, PDO::PARAM_STR);
                $stmt->bindParam(':guest_tools_version', $guest_tools_version, PDO::PARAM_STR);
                $stmt->bindParam(':guest_tools_running', $guest_tools_running, PDO::PARAM_STR);
                $stmt->bindParam(':guest_hostname', $guest_hostname, PDO::PARAM_STR);
                $stmt->bindParam(':guest_ip', $guest_ip, PDO::PARAM_STR);
                $stmt->bindParam(':stat_cpu_usage', $stat_cpu_usage, PDO::PARAM_STR);
                $stmt->bindParam(':stat_host_memory_usage', $stat_host_memory_usage, PDO::PARAM_STR);
                $stmt->bindParam(':stat_guest_memory_usage', $stat_guest_memory_usage, PDO::PARAM_STR);
                $stmt->bindParam(':stat_uptime_sec', $stat_uptime_sec, PDO::PARAM_STR);
                $stmt->bindParam(':power_state', $power_state, PDO::PARAM_STR);
                $stmt->bindParam(':esxi_id', $esxi_id, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
}


function update_resourcepool($data){

    $vcenter_id = $data[0]['vcenter_id'];

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE resourcepool SET present = 0 WHERE present = 1 AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO resourcepool (id,moref,name,type,status,vapp_state,parent,parent_moref,cluster_id,configured_memory_mb,cpu_reservation,cpu_limit,mem_reservation,mem_limit,vcenter_id) ' . 
                    'VALUES(:id,:moref,:name,:type,:status,:vapp_state,:parent,:parent_moref,:cluster_id,:configured_memory_mb,:cpu_reservation,:cpu_limit,:mem_reservation,:mem_limit,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE name=VALUES(name),type=VALUES(type),status=VALUES(status),vapp_state=VALUES(vapp_state),parent=VALUES(parent),parent_moref=VALUES(parent_moref),cluster_id=VALUES(cluster_id),configured_memory_mb=VALUES(configured_memory_mb),cpu_reservation=VALUES(cpu_reservation),cpu_limit=VALUES(cpu_limit),mem_reservation=VALUES(mem_reservation),mem_limit=VALUES(mem_limit),present=1');

            foreach ($data as $res) {

                $id = md5( $res['vcenter_id'] . $res['moref'] );
                $moref = $res['moref'];
                $name = $res['name'];
                $type = $res['type'];
                $status = $res['status'];
                $vapp_state = $res['vapp_state'];
                if ( strpos($res['parent_moref'], 'domain-') === false ){
                    # this is the cluster's default root resource pool
                    $parent = md5( $res['vcenter_id'] . $res['parent_moref'] );
                    $parent_moref = $res['parent_moref'];
                } else {
                    $parent = 'cluster';
                    $parent_moref = $res['parent_moref'];
                }
                $cluster_id = md5( $res['vcenter_id'] . $res['cluster_moref'] );
                $configured_memory_mb = $res['configured_memory_mb'];
                $cpu_reservation = $res['cpu_reservation'];
                $cpu_limit = $res['cpu_limit'];
                $mem_reservation = $res['mem_reservation'];
                $mem_limit = $res['mem_limit'];
                $vcenter_id = $res['vcenter_id'];

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':moref', $moref, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':type', $type, PDO::PARAM_STR);
                $stmt->bindParam(':status', $status, PDO::PARAM_STR);
                $stmt->bindParam(':vapp_state', $vapp_state, PDO::PARAM_STR);
                $stmt->bindParam(':parent', $parent, PDO::PARAM_STR);
                $stmt->bindParam(':parent_moref', $parent_moref, PDO::PARAM_STR);
                $stmt->bindParam(':cluster_id', $cluster_id, PDO::PARAM_STR);
                $stmt->bindParam(':configured_memory_mb', $configured_memory_mb);
                $stmt->bindParam(':cpu_reservation', $cpu_reservation);
                $stmt->bindParam(':cpu_limit', $cpu_limit);
                $stmt->bindParam(':mem_reservation', $mem_reservation);
                $stmt->bindParam(':mem_limit', $mem_limit);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

        // Update resourcepool full path if successful
        update_resourcepool_full_path($vcenter_id);

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }

}



function update_datastore($data){

    $vcenter_id = $data[0]['vcenter_id'];

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE datastore SET present = 0 WHERE present = 1 AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO datastore (id,name,moref,status,capacity_bytes,free_bytes,uncommitted_bytes,type,vcenter_id) ' . 
                    'VALUES(:id,:name,:moref,:status,:capacity_bytes,:free_bytes,:uncommitted_bytes,:type,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE name=VALUES(name),status=VALUES(status),capacity_bytes=VALUES(capacity_bytes),free_bytes=VALUES(free_bytes),uncommitted_bytes=VALUES(uncommitted_bytes),type=VALUES(type),present=1');

            foreach ($data as $datastore) {

                $id = md5( $datastore['vcenter_id'] . $datastore['moref'] );
                $name = $datastore['name'];
                $moref = $datastore['moref'];
                $status = $datastore['status'];
                $capacity_bytes = $datastore['capacity_bytes'];
                $free_bytes = $datastore['free_bytes'];
                $uncommitted_bytes = $datastore['uncommitted_bytes'];
                $type = $datastore['type'];
                $vcenter_id = $datastore['vcenter_id'];

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':moref', $moref, PDO::PARAM_STR);
                $stmt->bindParam(':status', $status, PDO::PARAM_INT);
                $stmt->bindParam(':capacity_bytes', $capacity_bytes);
                $stmt->bindParam(':free_bytes', $free_bytes);
                $stmt->bindParam(':uncommitted_bytes', $uncommitted_bytes);
                $stmt->bindParam(':type', $type, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
}

function update_vdisk($data){

    $vcenter_id = $data[0]['vcenter_id'];

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE vdisk SET present = 0 WHERE present = 1 AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO vdisk (id,name,capacity_bytes,path,thin_provisioned,datastore_id,uuid,disk_object_id,vm_id,esxi_id,vcenter_id) ' . 
                    'VALUES(:id,:name,:capacity_bytes,:path,:thin_provisioned,:datastore_id,:uuid,:disk_object_id,:vm_id,:esxi_id,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE name=VALUES(name),capacity_bytes=VALUES(capacity_bytes),path=VALUES(path),thin_provisioned=VALUES(thin_provisioned),datastore_id=VALUES(datastore_id),uuid=VALUES(uuid),vm_id=VALUES(vm_id),esxi_id=VALUES(esxi_id),present=1');

            foreach ($data as $vdisk) {

                // older vm versions will not have disk_object_id. add path to the hash
                $id = md5( $vdisk['vcenter_id'] . $vdisk['disk_object_id'] . $vdisk['path'] );
                $name = $vdisk['name'];
                if ( is_null($vdisk['capacity_bytes']) ){
                    $capacity_bytes = 1024 * $vdisk['capacity_kb'];
                } else {
                    $capacity_bytes = $vdisk['capacity_bytes'];
                }
                $path = $vdisk['path'];
                $thin_provisioned = $vdisk['thin_provisioned'];
                $datastore_id = md5( $vdisk['vcenter_id'] . $vdisk['datastore_moref'] );
                $uuid = $vdisk['uuid'];
                $disk_object_id = $vdisk['disk_object_id'];
                $vm_id = md5( $vdisk['vcenter_id'] . $vdisk['vm_moref'] );
                $esxi_id= md5( $vdisk['vcenter_id'] . $vdisk['esxi_moref'] );
                $vcenter_id = $vdisk['vcenter_id'];

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':capacity_bytes', $capacity_bytes);
                $stmt->bindParam(':path', $path, PDO::PARAM_STR);
                $stmt->bindParam(':thin_provisioned', $thin_provisioned, PDO::PARAM_STR);
                $stmt->bindParam(':datastore_id', $datastore_id, PDO::PARAM_STR);
                $stmt->bindParam(':uuid', $uuid, PDO::PARAM_STR);
                $stmt->bindParam(':disk_object_id', $disk_object_id, PDO::PARAM_STR);
                $stmt->bindParam(':vm_id', $vm_id, PDO::PARAM_STR);
                $stmt->bindParam(':esxi_id', $esxi_id, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
}

function update_pnic($data){
    
    $vcenter_id = $data[0]['vcenter_id'];

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE pnic SET present = 0 WHERE present = 1 AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO pnic (id,name,mac,link_speed,driver,esxi_id,vcenter_id) ' . 
                    'VALUES(:id,:name,:mac,:link_speed,:driver,:esxi_id,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE name=VALUES(name),mac=VALUES(mac),link_speed=VALUES(link_speed),driver=VALUES(driver),present=1');

            foreach ($data as $pnic) {

                $name = $pnic['name'];
                $esxi_id = md5( $pnic['vcenter_id'] . $pnic['esxi_moref'] );
                $id = md5( $esxi_id . $name );
                $mac = $pnic['mac'];
                $link_speed = $pnic['link_speed'];
                $driver = $pnic['driver'];
                $vcenter_id = $pnic['vcenter_id'];

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':mac', $mac, PDO::PARAM_STR);
                $stmt->bindParam(':link_speed', $link_speed, PDO::PARAM_INT);
                $stmt->bindParam(':driver', $driver, PDO::PARAM_STR);
                $stmt->bindParam(':esxi_id', $esxi_id, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
}

function update_dvs($data){
    
    $vcenter_id = $data[0]['vcenter_id'];
    $type = 'DVS';

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }

    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE vswitch SET present = 0 WHERE present = 1 AND type = "'.$type.'" AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO vswitch (id,name,type,version,max_mtu,ports,vcenter_id) ' . 
                    'VALUES(:id,:name,:type,:version,:max_mtu,:ports,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE name=VALUES(name),type=VALUES(type),version=VALUES(version),max_mtu=VALUES(max_mtu),ports=VALUES(ports),present=1');

            foreach ($data as $dvs) {

                $id = md5( $dvs['vcenter_id'] . $dvs['moref'] );
                $name = $dvs['name'];
                $version = $dvs['version'];
                $max_mtu = $dvs['max_mtu'];
                $ports = $dvs['ports'];
                $vcenter_id = $dvs['vcenter_id'];

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':type', $type, PDO::PARAM_STR);
                $stmt->bindParam(':version', $version, PDO::PARAM_STR);
                $stmt->bindParam(':max_mtu', $max_mtu, PDO::PARAM_INT);
                $stmt->bindParam(':ports', $ports, PDO::PARAM_INT);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
}

function update_svs($data){
    
    $vcenter_id = $data[0]['vcenter_id'];
    $type = 'vSwitch';

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }

    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE vswitch SET present = 0 WHERE present = 1 AND type = "'.$type.'" AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO vswitch (id,name,type,esxi_id,max_mtu,ports,vcenter_id) ' . 
                    'VALUES(:id,:name,:type,:esxi_id,:max_mtu,:ports,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE max_mtu=VALUES(max_mtu),ports=VALUES(ports),present=1');

            foreach ($data as $svs) {

                $id = md5( $svs['vcenter_id'] . $svs['esxi_moref'] . $svs['name'] );
                $name = $svs['name'];
                $max_mtu = $svs['max_mtu'];
                $ports = $svs['ports'];
                $esxi_id = md5( $svs['vcenter_id'] . $svs['esxi_moref'] );
                $vcenter_id = $svs['vcenter_id'];

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':type', $type, PDO::PARAM_STR);
                $stmt->bindParam(':esxi_id', $esxi_id, PDO::PARAM_STR);
                $stmt->bindParam(':max_mtu', $max_mtu, PDO::PARAM_STR);
                $stmt->bindParam(':ports', $ports, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
}

function update_dvspg($data){
    
    $vcenter_id = $data[0]['vcenter_id'];
    $type = 'DVS';

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE portgroup SET present = 0 WHERE present = 1 AND type = "'.$type.'" AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO portgroup (id,name,type,vlan,vlan_type,vswitch_id,vcenter_id) ' . 
                    'VALUES(:id,:name,:type,:vlan,:vlan_type,:vswitch_id,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE name=VALUES(name),type=VALUES(type),vlan=VALUES(vlan),vlan_type=VALUES(vlan_type),vswitch_id=VALUES(vswitch_id),vcenter_id=VALUES(vcenter_id),present=1');

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

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':type', $type, PDO::PARAM_STR);
                $stmt->bindParam(':vlan', $vlan, PDO::PARAM_STR);
                $stmt->bindParam(':vlan_type', $vlan_type, PDO::PARAM_STR);
                $stmt->bindParam(':vswitch_id', $vswitch_id, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
}

function update_svspg($data){
    
    $vcenter_id = $data[0]['vcenter_id'];
    $type = 'vSwitch';

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE portgroup SET present = 0 WHERE present = 1 AND type = "'.$type.'" AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO portgroup (id,name,type,vlan,vlan_type,vswitch_id,vcenter_id) ' . 
                    'VALUES(:id,:name,:type,:vlan,:vlan_type,:vswitch_id,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE name=VALUES(name),type=VALUES(type),vlan=VALUES(vlan),vlan_type=VALUES(vlan_type),vswitch_id=VALUES(vswitch_id),vcenter_id=VALUES(vcenter_id),present=1');

            foreach ($data as $pg) {

                $esxi_id = md5( $pg['vcenter_id'] . $pg['esxi_moref'] );
                $id = md5( $pg['vcenter_id'] . $pg['esxi_moref'] . $pg['name'] );
                $vswitch_id = md5( $pg['vcenter_id'] . $pg['esxi_moref'] . $pg['vswitch_name'] );
                $vlan = $pg['vlan'];
                $vlan_type = 'single';
                $name = $pg['name'];
                $vcenter_id = $pg['vcenter_id'];

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':type', $type, PDO::PARAM_STR);
                $stmt->bindParam(':vlan', $vlan, PDO::PARAM_STR);
                $stmt->bindParam(':vlan_type', $vlan_type, PDO::PARAM_STR);
                $stmt->bindParam(':vswitch_id', $vswitch_id, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
}

function update_datacenter($data){
    
    $vcenter_id = $data[0]['vcenter_id'];

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE datacenter SET present = 0 WHERE present = 1 AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO datacenter (id,vm_folder_id,esxi_folder_id,name,vcenter_id) ' . 
                    'VALUES(:id,:vm_folder_id,:esxi_folder_id,:name,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE vm_folder_id=VALUES(vm_folder_id),esxi_folder_id=VALUES(esxi_folder_id),name=VALUES(name),present=1');

            foreach ($data as $dc) {

                $id = md5( $dc['vcenter_id'] . $dc['moref'] );
                $vm_folder_id = $dc['vm_folder_moref'];
                $esxi_folder_id = md5( $dc['vcenter_id'] . $dc['esxi_folder_moref'] );
                $name = $dc['name'];
                $vcenter_id = $dc['vcenter_id'];

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':vm_folder_id', $vm_folder_id, PDO::PARAM_STR);
                $stmt->bindParam(':esxi_folder_id', $esxi_folder_id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
}

function update_cluster($data){
    
    $vcenter_id = $data[0]['vcenter_id'];

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE cluster SET present = 0 WHERE present = 1 AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO cluster (id,name,datacenter_id,current_balance,target_balance,total_cpu_threads,total_cpu_mhz,total_memory_bytes,total_vmotions,num_hosts,drs_enabled,drs_behaviour,ha_enabled,status,vcenter_id) ' . 
                    'VALUES(:id,:name,:datacenter_id,:current_balance,:target_balance,:total_cpu_threads,:total_cpu_mhz,:total_memory_bytes,:total_vmotions,:num_hosts,:drs_enabled,:drs_behaviour,:ha_enabled,:status,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE name=VALUES(name),current_balance=VALUES(current_balance),target_balance=VALUES(target_balance),datacenter_id=VALUES(datacenter_id),total_cpu_threads=VALUES(total_cpu_threads),total_cpu_mhz=VALUES(total_cpu_mhz),total_memory_bytes=VALUES(total_memory_bytes),total_vmotions=VALUES(total_vmotions),num_hosts=VALUES(num_hosts),drs_enabled=VALUES(drs_enabled),drs_behaviour=VALUES(drs_behaviour),ha_enabled=VALUES(ha_enabled),status=VALUES(status),present=1');

            foreach ($data as $cl) {

                $id = md5( $cl['vcenter_id'] . $cl['moref'] );
                $datacenter_id = md5( $cl['vcenter_id'] . $cl['datacenter_moref'] );
                $current_balance = $cl['current_balance'];
                $target_balance = $cl['target_balance'];
                $total_cpu_threads = $cl['total_cpu_threads'];
                $total_cpu_mhz = $cl['total_cpu_mhz'];
                $total_memory_bytes = $cl['total_memory_bytes'];
                $total_vmotions = $cl['total_vmotions'];
                $num_hosts = $cl['num_hosts'];
                $drs_enabled = $cl['drs_enabled'];
                $drs_behaviour = $cl['drs_behaviour'];
                $ha_enabled = $cl['ha_enabled'];
                $status = $cl['status'];
                $name = $cl['name'];
                $vcenter_id = $cl['vcenter_id'];

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':datacenter_id', $datacenter_id, PDO::PARAM_STR);
                $stmt->bindParam(':current_balance', $current_balance, PDO::PARAM_INT);
                $stmt->bindParam(':target_balance', $target_balance, PDO::PARAM_INT);
                $stmt->bindParam(':total_cpu_threads', $total_cpu_threads, PDO::PARAM_INT);
                $stmt->bindParam(':total_cpu_mhz', $total_cpu_mhz);
                $stmt->bindParam(':total_memory_bytes', $total_memory_bytes);
                $stmt->bindParam(':total_vmotions', $total_vmotions, PDO::PARAM_INT);
                $stmt->bindParam(':num_hosts', $num_hosts, PDO::PARAM_INT);
                $stmt->bindParam(':drs_enabled', $drs_enabled, PDO::PARAM_STR);
                $stmt->bindParam(':drs_behaviour', $drs_behaviour, PDO::PARAM_STR);
                $stmt->bindParam(':ha_enabled', $ha_enabled, PDO::PARAM_STR);
                $stmt->bindParam(':status', $status, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
}

function update_folder($data){
    
    $vcenter_id = $data[0]['vcenter_id'];

    if ( isset($data[0]['empty']) ){
        $EMPTY = true;
    } else {
        $EMPTY = false;
    }
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE folder SET present = 0 WHERE present = 1 AND vcenter_id = ' . $pdo->quote($vcenter_id) );

        if (!$EMPTY){

            $stmt = $pdo->prepare('INSERT INTO folder (id,moref,name,type,parent,parent_datacenter_id,vcenter_id) ' . 
                    'VALUES(:id,:moref,:name,:type,:parent,:parent_datacenter_id,:vcenter_id) ' .
                    'ON DUPLICATE KEY UPDATE name=VALUES(name),type=VALUES(type),parent=VALUES(parent),parent_datacenter_id=VALUES(parent_datacenter_id),present=1');

            foreach ($data as $folder) {

                $id = md5( $folder['vcenter_id'] . $folder['moref'] );
                $moref = $folder['moref'];
                if ( strpos($folder['type'], 'VirtualMachine') === false ) {
                    # dont handle non-vm folders for now
                    $type = 'not_vm';
                }else {
                    $type = 'VirtualMachine';
                }
                if ( strpos($folder['parent_moref'], 'datacenter-') === false ) {
                    $parent = md5( $folder['vcenter_id'] . $folder['parent_moref'] );
                    $parent_datacenter_id = 'n/a';
                }else {
                    # This is a root folder DC
                    $parent = 'datacenter';
                    $parent_datacenter_id = md5( $folder['vcenter_id'] . $folder['parent_moref'] );
                }
                $name = $folder['name'];
                $vcenter_id = $folder['vcenter_id'];

                $stmt->bindParam(':id', $id, PDO::PARAM_STR);
                $stmt->bindParam(':moref', $moref, PDO::PARAM_STR);
                $stmt->bindParam(':type', $type, PDO::PARAM_STR);
                $stmt->bindParam(':parent', $parent, PDO::PARAM_STR);
                $stmt->bindParam(':parent_datacenter_id', $parent_datacenter_id, PDO::PARAM_STR);
                $stmt->bindParam(':name', $name, PDO::PARAM_STR);
                $stmt->bindParam(':vcenter_id', $vcenter_id, PDO::PARAM_STR);

                $stmt->execute();

            }
        }

        $pdo->commit();

        // Update resourcepool full path if successful
        update_folder_full_path($vcenter_id);

    } catch (PDOException $e) {
        // rollback transaction on error
        $pdo->rollback();
        // return 500
        echo "Error in transaction: ".$e->getMessage();
        http_response_code(500);
    }
    
}

function update_folder_full_path($vcenter_id){

    // import array_column function to support php versions older than 5.5
    require_once('lib/array_column.php');

    global $pdo;

    // make arrays to save sql queries and make function faster
    $dcs = array();
    $folders = array();
    $full_path_results = array();

    // populate folders
    $query = "select id,name,parent FROM folder WHERE vcenter_id = '$vcenter_id' AND parent != 'datacenter' AND type = 'VirtualMachine' AND present=1";
    foreach($pdo->query($query) as $folder){
        $folders[] = $folder;
    }

    // populate dcs
    $query = "select folder.id, folder.name, folder.parent, datacenter.name from folder LEFT JOIN datacenter ON folder.parent_datacenter_id = datacenter.id WHERE folder.vcenter_id = '$vcenter_id' AND folder.parent='datacenter' AND folder.type='VirtualMachine' AND folder.present=1 AND datacenter.present=1";
    foreach($pdo->query($query) as $dc){
        $dcs[ $dc['id'] ] = $dc['name'];
    }


    // Loop through folders and determine full path
    foreach ( $folders as $id => $info ){
        $continue = true;
        $search_id = $info['parent'];

        while ( $continue ) {
            $key = array_search( $search_id, array_column($folders,'id') );
            if ( $key === false || ( array_search($search_id, $dcs) != false ) ){
                # this parent folder is a dc
                $folders[ $id ]['dc_id'] = $search_id;
                $folders[ $id ]['dc_name'] = $dcs[ $search_id ];
                $continue = false;
            } else {
                $folders[ $id ]['path'][] = $folders[ $key ]['name'];
                $folders[ $id ]['path_id'][] = $search_id;
                if ( isset( $folders[ $key ]['parent'] ) && ( array_search($folders[ $key ]['parent'], $dcs) == false ) ){
                    $search_id = $folders[ $key ]['parent'];
                    $continue = true;
                } else {
                    $continue = false;
                }
            }
            
        }
    }


    // add full path to root folders
    // populate folders
    $query = "select null AS name, folder.id AS id,datacenter.name AS dc_name FROM folder LEFT JOIN datacenter on folder.parent_datacenter_id=datacenter.id WHERE folder.parent = 'datacenter' AND folder.vcenter_id = '$vcenter_id' AND folder.type = 'VirtualMachine' AND folder.present=1";
    foreach($pdo->query($query) as $folder){
        $folders[] = $folder;
    }


    foreach ( $folders as $folder ){
        $full_path = $folder['dc_name'];
        if ( isset($folder['path']) ){
            $ordered_path = array_reverse($folder['path']);
            foreach ( $ordered_path as $path ){
                $full_path .= "/$path";
            }
        }

        $full_path .= "/".$folder['name'];
        

        try {

            // start transaction
            $pdo->beginTransaction();

            // prepare statement to avoid sql injections
            $stmt = $pdo->prepare('UPDATE folder ' . 
                    'SET full_path=:full_path ' .
                    'WHERE id=' . $pdo->quote($folder['id']) );

            $stmt->bindParam(':full_path', $full_path, PDO::PARAM_STR);

            // execute prepared statement
            $stmt->execute();

            // commit transaction
            $pdo->commit();

        } catch (PDOException $e) {
            // rollback transaction on error
            $pdo->rollback();
            // return 500
            http_response_code(500);
        }

    }

}


function update_resourcepool_full_path($vcenter_id){

    // import array_column function to support php versions older than 5.5
    require_once('lib/array_column.php');

    global $pdo;

    // make arrays to save sql queries and make function faster
    $resourcepools = array();

    // populate folders
    $query = "select r.name, r.moref, r.parent_moref, c.name AS cluster, d.name AS datacenter FROM resourcepool r LEFT JOIN cluster c ON r.cluster_id = c.id LEFT JOIN datacenter d ON c.datacenter_id = d.esxi_folder_id WHERE r.vcenter_id = '$vcenter_id' AND r.present=1";
    foreach($pdo->query($query) as $respool){
        $resourcepools[] = $respool;
    }

    // Loop through resourcepools and determine full path
    foreach ( $resourcepools as $res ){
        $continue = true;
        $search_id = $res['parent_moref'];
        $parent_type = explode("-", $search_id)[0];

        $res['full_path'] = '';

        // initialize and wipe array
        $full_path_array = array();

        // stop when parent is a datacenter
        if ($parent_type == 'domain'){
            $continue = false;
        }

        // for migration purposes, keep track of any vapps in path
        $res['vapp_in_path'] = 0;

        while ( $continue ) {
            $key = array_search( $search_id, array_column($resourcepools,'moref') );
            if ( $key === false  ){
                // stop if parent not found
                $continue = false;
            } else {
                // since we are here, check if parent is vapp for migration purposes.
                // we can do this by checking for a v after - 
                if ( explode("-", $search_id)[1][0] == 'v' ){
                    $res['vapp_in_path'] = 1;
                }

                // continue with path logic
                $full_path_array[] = $resourcepools[$key]['name'];
                $search_id = $resourcepools[$key]['parent_moref'];
            }
            
        }

        // add the cluster and datacenter name
        $full_path_array[] = $res['cluster'];
        $full_path_array[] = $res['datacenter'];

        // reverse order to start from cluster
        $ordered_path = array_reverse($full_path_array);

        // remove the hidden root pool of the cluster that's always called Resources
        unset($ordered_path[2]);

        // create a string from the paths
        foreach ( $ordered_path as $part ){
            $res['full_path'] .= $part.'/';
        }

        // add the name of this resource to the end
        $res['full_path'] .= $res['name'];

        // calculate id of this resource
        $rp_id = md5( $vcenter_id . $res['moref'] );

        try {

            // start transaction
            $pdo->beginTransaction();

            // prepare statement to avoid sql injections
            $stmt = $pdo->prepare('UPDATE resourcepool ' . 
                    'SET full_path=:full_path, vapp_in_path=:vapp_in_path ' .
                    'WHERE id=' . $pdo->quote($rp_id) );

            $stmt->bindParam(':full_path', $res['full_path'], PDO::PARAM_STR);
            $stmt->bindParam(':vapp_in_path', $res['vapp_in_path'], PDO::PARAM_INT);

            // execute prepared statement
            $stmt->execute();

            // commit transaction
            $pdo->commit();

        } catch (PDOException $e) {
            // rollback transaction on error
            $pdo->rollback();
            // return 500
            http_response_code(500);
        }


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


// Load MYSQL connection details
require_once( 'lib/mysql_config.php' );

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

// Get POST data
$data = json_decode(file_get_contents('php://input'), true);

// quick data validation and manipulation to support single arrays 
if ( isset($data['objecttype']) && strcasecmp($data['objecttype'],"VCENTER") == 0 ){
    update_vcenter($data);
    exit();
} elseif ( isset($data['objecttype']) && strcasecmp($data['objecttype'],"VCENTER") != 0 ){
    $post_data[0] = $data;
    $object_type = $data['objecttype'];
} elseif ( isset($data[0]['objecttype']) ){
    $post_data = $data;
    $object_type = $data[0]['objecttype'];
} else {
    echo "Invalid data";
    http_response_code(500);
    exit();
}

// pass post_data to correct function based on object_type
switch ($object_type) {
    case "ESXI":
        update_esxi($post_data);
        break;
    case "VM":
        update_vm($post_data);
        break;
    case "VNIC":
        update_vnic($post_data);
        break;
    case "DS":
        update_datastore($post_data);
        break;
    case "VDISK":
        update_vdisk($post_data);
        break;
    case "PNIC":
        update_pnic($post_data);
        break;
    case "DVS":
        update_dvs($post_data);
        break;
    case "SVS":
        update_svs($post_data);
        break;
    case "DVSPG":
        update_dvspg($post_data);
        break;
    case "SVSPG":
        update_svspg($post_data);
        break;
    case "RES":
        update_resourcepool($post_data);
        break;
    case "DC":
        update_datacenter($post_data);
        break;
    case "FOLDER":
        update_folder($post_data);
        break;
    case "CLUSTER":
        update_cluster($post_data);
        break;
    default:
        echo "Invalid data";
        http_response_code(500);
}














?>

