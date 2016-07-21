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

function add_vcenter($pdo, $vc_uuid){

  $stmt = $pdo->prepare("SELECT * FROM vcenter WHERE id =?");
  $stmt->bindValue(1, $vc_uuid, PDO::PARAM_STR);
  $stmt->execute();
  return $stmt->fetchAll(PDO::FETCH_ASSOC);

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
  );

  curl_setopt_array($curl, $opt);
  $output = curl_exec($curl);

  $result['status'] = curl_getinfo($curl, CURLINFO_HTTP_CODE);
  $result['output'] = $output;

  return $result;

}

// start logic

if ($_SERVER['REQUEST_METHOD'] === 'POST') {

  $result = test_vcenter_creds($_POST['host'], $_POST['user'], $_POST['pass']);

  if ( $result['status'] == 200 ){
    echo "SUCCESS";
  } else {
    echo $result['output'];
    http_response_code(500);
  }

}
