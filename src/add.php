<form role="edit" method="post">
<div class="modal-header">
  <button type="button" class="close" data-dismiss="modal" aria-label="Close">
    <span aria-hidden="true">&times;</span>
  </button>
  <h4 class="modal-title" id="myModalLabel">Add vCenter Credentials</h4>
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

function add_vcenter($pdo, $vc_uuid){

  $stmt = $pdo->prepare("SELECT * FROM vcenter WHERE id =?");
  $stmt->bindValue(1, $vc_uuid, PDO::PARAM_STR);
  $stmt->execute();
  return $stmt->fetchAll(PDO::FETCH_ASSOC);

}

if ($_SERVER['REQUEST_METHOD'] === 'GET') {
  echo '
  <table class="table table-striped">
    <tbody>
      <tr>
        <td><label for="fqdn">FQDN or IP:</label></td>
        <td><input name="host" type="text" style="width:100%"></td>
      </tr>
      <tr>
        <td><label for="short_name">Environment:</label></td>
        <td><input name="short_name" type="text" style="width:100%"></td>
      </tr>
      <tr>
        <td><label for="user_name">Username:</label></td>
        <td><input name="user" type="text" style="width:100%"></td>
      </tr>
      <tr>
        <td><label for="password">Password:</label></td>
        <td><input name="pass" type="password" style="width:100%"></td>
      </tr>
      <tr>
        <td><label for="autopoll">Enable Auto Poll</label>:</td>
        <td><input type="checkbox" name="auto_poll"></td>
      </tr>
    </tbody>
  </table>
  <div id="message" class="alert alert-danger hidden" role="alert">
  </div>
  ';
}

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
  echo "<br /><hr />POSTED!!";
}



?>
</div>
<div class="modal-footer">
  <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
  <button type="submit" id="submit" class="btn btn-primary">Add vCenter</button>
</div>
</form>

<script>
function hideModal(){
  $("#pollerModal").modal('hide');
}

$("#submit").click(function(e){
    e.preventDefault();
    //make ajax call
    $.ajax({
        type: "POST",
        url: "bridge.php",
        data: $('form').serialize(),
        beforeSend: function(){
            $("#message").html('<img src="img/ripple.gif" /> Executing Request');
            $("#message").removeClass("hidden alert-danger");
            $("#message").addClass("alert-info");
        },
        success: function(msg){
            $("#message").html("<i class='fa fa-check fa-fw'></i> " + msg);
            $("#message").removeClass("hidden alert-info alert-danger");
            $("#message").addClass("alert-success");
            setTimeout(hideModal, 1000);
        },
        error: function(jqXHR){
            $("#message").html("<i class='fa fa-times fa-fw'></i> " + jqXHR.responseText);
            $("#message").removeClass("hidden alert-info");
            $("#message").addClass("alert-danger");
        }
    });

});
</script>
