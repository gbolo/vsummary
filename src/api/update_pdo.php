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
        $conn->rollback();
        // return 500
        http_response_code(500);
    }


}

function update_vnic($data){

    $vcenter_id = $data[0]['vcenter_id'];
    $type = 'vSwitch';
   
    try {

        global $pdo;

        $pdo->beginTransaction();
        $pdo->query( 'UPDATE vnic SET present = 0 WHERE present = 1 AND vcenter_id = ' . $pdo->quote($vcenter_id) );

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
        $pdo->commit();

    } catch (PDOException $e) {
        // rollback transaction on error
        $conn->rollback();
        // return 500
        http_response_code(500);
    }
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

// Pass the POST data to the correct function
if ( strcasecmp($data['objecttype'],"VCENTER")==0 ){
    update_vcenter($data);
}
elseif ( strcasecmp($data[0]['objecttype'],"VNIC")==0 ){
    update_vnic($data);
} else{
    echo "Invalid data";
    http_response_code(500);
}










?>

