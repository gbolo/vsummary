<?php
$uri_filename = basename($_SERVER['SCRIPT_NAME']);

$class_vm = "";
$class_esxi = "";

switch ($uri_filename) {
    case "vm.php":
        $class_vm='class="active"';
        break;
    case "esxi.php":
        $class_esxi='class="active"';
        break;
    case "vnic.php":
        $class_vnic='class="active"';
        break;
    case "vdisk.php":
        $class_vdisk='class="active"';
        break;
    case "datastore.php":
        $class_datastore='class="active"';
        break;
    case "vcenter.php":
        $class_vcenter='class="active"';
        break;
}

?>

          <h4 class="sub-header navbar-sub">Views</h4>
          <ul class="nav nav-sidebar">
            <li <?php echo $class_vm; ?> ><a href="vm.php">Virtual Machine</a></li>
            <li <?php echo $class_vnic; ?> ><a href="vnic.php">VM vNIC</a></li>
            <li <?php echo $class_vdisk; ?> ><a href="vdisk.php">VM vDisk</a></li>
            <li <?php echo $class_esxi; ?> ><a href="esxi.php">ESXi Host</a></li>
            <li <?php echo $class_datastore; ?> ><a href="datastore.php">Datastore</a></li>
          </ul>
          <br />
          <h4 class="sub-header navbar-sub">Analytics</h4>
          <ul class="nav nav-sidebar">
            <li <?php echo $class_vcpu_over; ?> ><a href="vcpu_over.php">vCPU Overprovision</a></li>
            <li <?php echo $class_vram_over; ?> ><a href="vram_over.php">vRAM Overprovision</a></li>
            <li <?php echo $class_vdisk_over; ?> ><a href="vdisk_over.php">vDisk Overprovision</a></li>
          </ul>
