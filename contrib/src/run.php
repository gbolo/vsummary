<?php
require_once('include/functions.php');
if ( $POLLER_ENABLED != true ){
  exit('POLLER IS DISABLED FOR DEMO');
}
?>

<div class="modal-header">
  <button type="button" class="close" data-dismiss="modal" aria-label="Close">
    <span aria-hidden="true">&times;</span>
  </button>
  <h4 class="modal-title" id="myModalLabel">Run vCenter Poll</h4>
</div>
<div class="modal-body">


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
<h3>Begin Polling of vCenter: '.$rows[0]['fqdn'].'</h3>
<hr />
<strong>UUID:</strong> '.$rows[0]['id'].'<br />
<strong>Username:</strong> '.$rows[0]['user_name'].'<br />
<hr />
<div id="message" class="alert alert-danger hidden" role="alert">
</div>
';




?>
</div>
<div class="modal-footer">
  <button type="button" class="btn btn-danger" data-dismiss="modal"><i class='fa fa-times fa-fw'></i> Close</button>
  <button type="button" id="run" class="btn btn-success"><i class='fa fa-play fa-fw'></i> RUN NOW</button>
</div>

<script>
function hideModal(){
  $("#pollerModal").modal('hide');
}

$("#run").click(function(e){
    e.preventDefault();
    //make ajax call
    $.ajax({
        type: "GET",
        url: "bridge.php?action=poll&vc_uuid=<?php echo $_GET['id']; ?>",
        beforeSend: function(){
            $("#message").html('<img src="img/ripple.gif" /> Executing Request');
            $("#message").removeClass("hidden alert-danger");
            $("#message").addClass("alert-info");
        },
        success: function(msg){
            $("#message").html("<i class='fa fa-check fa-fw'></i> " + msg);
            $("#message").removeClass("hidden alert-info alert-danger");
            $("#message").addClass("alert-success");
            //setTimeout(hideModal, 1000);
        },
        error: function(jqXHR){
            $("#message").html("<i class='fa fa-times fa-fw'></i> " + jqXHR.responseText);
            $("#message").removeClass("hidden alert-info");
            $("#message").addClass("alert-danger");
        }
    });

});
</script>
