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

function is_valid_domain($url){

    $validation = FALSE;
    /*Parse URL*/
    $urlparts = parse_url(filter_var($url, FILTER_SANITIZE_URL));
    /*Check host exist else path assign to host*/
    if(!isset($urlparts['host'])){
        $urlparts['host'] = $urlparts['path'];
    }

    if($urlparts['host']!=''){
       /*Add scheme if not found*/
        if (!isset($urlparts['scheme'])){
            $urlparts['scheme'] = 'http';
        }
        /*Validation*/
        if(checkdnsrr($urlparts['host'], 'A') && in_array($urlparts['scheme'],array('http','https')) && ip2long($urlparts['host']) === FALSE){
            $urlparts['host'] = preg_replace('/^www\./', '', $urlparts['host']);
            $url = $urlparts['scheme'].'://'.$urlparts['host']. "/";

            if (filter_var($url, FILTER_VALIDATE_URL) !== false && @get_headers($url)) {
                $validation = TRUE;
            }
        }
    }

  if(!$validation){
      return false;
  }else{
      return true;
  }

}


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
  );

  curl_setopt_array($curl, $opt);
  $output = curl_exec($curl);

  $result['status'] = curl_getinfo($curl, CURLINFO_HTTP_CODE);
  $result['output'] = $output;

  return $result;

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

}
