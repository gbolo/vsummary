<form role="edit" method="post">
<div class="modal-header">
  <button type="button" class="close" data-dismiss="modal" aria-label="Close">
    <span aria-hidden="true">&times;</span>
  </button>
  <h4 class="modal-title" id="myModalLabel">Modify vCenter Credentials</h4>
</div>
<div class="modal-body">
<div class="alert alert-warning" role="alert">
  <strong>Warning!</strong>
  <p>Please provide a user with <strong>read-only</strong> access to vcenter
    since no greater permissions are required.<br />
    Also keep in mind that passwords are not stored securely for auto-polling purposes.</p>
</div>

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

function get_vcenter_info($pdo, $vc_uuid){

  $stmt = $pdo->prepare("SELECT * FROM vcenter WHERE id =?");
  $stmt->bindValue(1, $vc_uuid, PDO::PARAM_STR);
  $stmt->execute();
  return $stmt->fetchAll(PDO::FETCH_ASSOC);

}

if ( isset($_GET['id']) ){
  $rows = get_vcenter_info($pdo, $_GET['id']);
} else {
  exit('Error in GET VAR');
}


echo '
<table class="table table-striped">
  <tbody>
    <tr>
      <td>ID:</td>
      <td>'.$rows[0]['id'].'</td>
    </tr>
    <tr>
      <td>FQDN:</td>
      <td><input name="fqdn" type="text" value="'.$rows[0]['fqdn'].'"></td>
    </tr>
    <tr>
      <td>ENV:</td>
      <td><input name="short_name" type="text" value="'.$rows[0]['short_name'].'"></td>
    </tr>
    <tr>
      <td>UserName:</td>
      <td><input name="user_name" type="text" value="'.$rows[0]['user_name'].'"></td>
    </tr>
    <tr>
      <td>Password:</td>
      <td><input name="password" type="password" value="'.$rows[0]['password'].'"></td>
    </tr>
    <tr>
      <td>Enable Auto Poll:</td>
      <td><input type="checkbox" name="auto_poll"></td>
    </tr>
  </tbody>
</table>

';





?>
</div>
<div class="modal-footer">
  <button type="button" class="btn btn-danger" data-dismiss="modal">Cancel</button>
  <button type="submit" class="btn btn-primary">Save changes</button>
</div>
</form>
