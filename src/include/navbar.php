
<?php
$uri_filename = basename($_SERVER['SCRIPT_NAME']);

$class_dashboard = "";
$class_admin = "";

if ( strcmp("admin.php", $uri_filename) == 0 ){
  $class_admin = 'class="active"';
} else {
  $class_dashboard = 'class="active"';
}


?>

    <nav class="navbar navbar-inverse navbar-fixed-top">
      <div class="container-fluid">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a class="navbar-brand" href="index.php"><span class="glyphicon glyphicon-eye-open" aria-hidden="true"></span> vSummary</a>
        </div>
        <div id="navbar" class="navbar-collapse collapse">
          <ul class="nav navbar-nav navbar-right">
            <li <?php echo $class_dashboard; ?> ><a href="index.php" ><span class="icon-hdtv"></span> DASHBOARD</a></li>
            <li <?php echo $class_admin; ?> ><a href="admin.php"><span class="icon-lock"></span> ADMIN LOGIN</a></li>
          </ul>
          <!--
          <form class="navbar-form navbar-right">
            <input type="text" class="form-control" placeholder="Search...">
          </form>
          -->
        </div>
      </div>
    </nav>


