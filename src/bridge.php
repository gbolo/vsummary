<?php

ini_set('display_errors', 1);
ini_set('display_startup_errors', 1);
error_reporting(E_ALL);

// Load MYSQL connection details
require_once( 'api/lib/mysql_config.php' );

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

// functions required

function validate_post_vars(){

  if ( !empty($_POST['host'])
    && !empty($_POST['short_name'])
    && !empty($_POST['user'])
    && !empty($_POST['pass'])
   ){
     $bool_host = filter_var($_POST['host'] ,FILTER_VALIDATE_IP);

  } else {
    echo '<b>Error</b> - Some Parameters are Missing!';
    http_response_code(500);
    exit();
  }

}

function update_vcenter($vc_uuid){

    try {

        // grab the pdo object declared outside of this function
        global $pdo;

        // start transaction
        $pdo->beginTransaction();

        // prepare statement to avoid sql injections
        $stmt = $pdo->prepare('INSERT INTO vcenter (id, fqdn, short_name, user_name, password) ' .
                'VALUES(:id, :fqdn, :short_name, :user_name, :password) ' .
                'ON DUPLICATE KEY UPDATE fqdn=VALUES(fqdn), short_name=VALUES(short_name), user_name=VALUES(user_name), password=VALUES(password)');

        $stmt->bindParam(':id', $vc_uuid, PDO::PARAM_STR);
        $stmt->bindParam(':fqdn', $_POST['host'], PDO::PARAM_STR);
        $stmt->bindParam(':short_name', $_POST['short_name'], PDO::PARAM_STR);
        $stmt->bindParam(':user_name', $_POST['user'], PDO::PARAM_STR);
        $stmt->bindParam(':password', $_POST['pass'], PDO::PARAM_STR);

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

function remove_vcenter($vc_uuid){

    try {

        // grab the pdo object declared outside of this function
        global $pdo;

        // start transaction
        $pdo->beginTransaction();

        // prepare statement to avoid sql injections
        $stmt = $pdo->prepare('DELETE FROM vcenter WHERE id = :id');

        $stmt->bindParam(':id', $vc_uuid, PDO::PARAM_STR);

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

function test_vcenter_creds($host, $user, $pass){
  $curl = curl_init();

  $fields = array(
    'user' => urlencode($user),
    'pwd' => urlencode($pass),
    'host' => urlencode($host),
  );

  $fields_string = "";
  foreach($fields as $key=>$value){
    $fields_string .= $key .'='. $value .'&';
  }
  rtrim($fields_string, '&');

  $opt = array(
    CURLOPT_URL            => "http://10.0.77.77:5000/vc_uuid",
    CURLOPT_USERAGENT      => "LinuxAPI",
    CURLOPT_CUSTOMREQUEST  => "POST",
    CURLOPT_POST           => count($fields),
    CURLOPT_POSTFIELDS     => $fields_string,
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_SSL_VERIFYHOST => 0,
    CURLOPT_SSL_VERIFYPEER => false,
    CURLOPT_FOLLOWLOCATION => true,
    CURLOPT_CONNECTTIMEOUT => 5,
    CURLOPT_TIMEOUT        => 7,
  );

  curl_setopt_array($curl, $opt);
  $output = curl_exec($curl);
  $curl_errno = curl_errno($curl);
  $curl_error = curl_error($curl);

  if ($curl_errno > 0) {
      $result['output'] = "Error ($curl_errno): $curl_error";
  } else {
      $result['output'] = $output;
  }

  $result['status'] = curl_getinfo($curl, CURLINFO_HTTP_CODE);

  return $result;

}

function poll_vcenter($host, $user, $pass){
  $curl = curl_init();

  $fields = array(
    'user' => urlencode($user),
    'pwd' => urlencode($pass),
    'host' => urlencode($host),
  );

  $fields_string = "";
  foreach($fields as $key=>$value){
    $fields_string .= $key .'='. $value .'&';
  }
  rtrim($fields_string, '&');

  $opt = array(
    CURLOPT_URL            => "http://10.0.77.77:5000/poll",
    CURLOPT_USERAGENT      => "LinuxAPI",
    CURLOPT_CUSTOMREQUEST  => "POST",
    CURLOPT_POST           => count($fields),
    CURLOPT_POSTFIELDS     => $fields_string,
    CURLOPT_RETURNTRANSFER => true,
    CURLOPT_SSL_VERIFYHOST => 0,
    CURLOPT_SSL_VERIFYPEER => false,
    CURLOPT_FOLLOWLOCATION => true,
    CURLOPT_CONNECTTIMEOUT => 5,
  );

  curl_setopt_array($curl, $opt);
  $output = curl_exec($curl);
  $curl_errno = curl_errno($curl);
  $curl_error = curl_error($curl);

  if ($curl_errno > 0) {
      $result['output'] = "Error ($curl_errno): $curl_error";
  } else {
      $result['output'] = $output;
  }

  $result['status'] = curl_getinfo($curl, CURLINFO_HTTP_CODE);

  return $result;

}

function get_vcenter_creds($vc_uuid){

  global $pdo;
  $stmt = $pdo->prepare("SELECT * FROM vcenter WHERE id =?");
  $stmt->bindValue(1, $vc_uuid, PDO::PARAM_STR);
  $stmt->execute();
  return $stmt->fetchAll(PDO::FETCH_ASSOC);

}

// start logic

if ($_SERVER['REQUEST_METHOD'] === 'POST') {

  validate_post_vars();

  if ( !empty($_GET['action']) && $_GET['action'] == 'remove' && !empty($_POST['vc_uuid']) ){
    remove_vcenter($_POST['vc_uuid']);
    echo "REMOVED ".$_POST['vc_uuid'];
    exit();
  }

  $result = test_vcenter_creds($_POST['host'], $_POST['user'], $_POST['pass']);

  if ( $result['status'] == 200 ){
    $vc_uuid = $result['output'];
    update_vcenter($vc_uuid);
    echo "Login Successful";
  } else {
    echo $result['output'];
    http_response_code(500);
  }

} else {

  if ( !empty($_GET['action']) && $_GET['action'] == 'poll' && !empty($_GET['vc_uuid']) ){
    echo "<pre>";
    $vcenter_creds = get_vcenter_creds($_GET['vc_uuid']);
    //print_r($vcenter_creds);
    $result = poll_vcenter($vcenter_creds[0]['fqdn'], $vcenter_creds[0]['user_name'], $vcenter_creds[0]['password']);
    print_r( json_decode($result['output']) );
  } else {
    echo "unknown request";
    http_response_code(500);
  }


}
