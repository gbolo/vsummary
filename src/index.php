<?php require_once('include/functions.php'); ?>
<!DOCTYPE html>
<html lang="en">

<head>

    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="">
    <meta name="author" content="">

    <link rel="icon" href="img/favicon.ico">
    <title>vSummary</title>

    <!-- Bootstrap core CSS -->
    <link href="bootstrap/css/bootstrap.min.css" rel="stylesheet">

    <!-- MetisMenu CSS -->
    <link href="css/metisMenu.min.css" rel="stylesheet">

    <!-- DataTables with extensions CSS -->
    <link rel="stylesheet" type="text/css" href="datatables/DataTables-1.10.11/css/dataTables.bootstrap.min.css"/>
    <link rel="stylesheet" type="text/css" href="datatables/Buttons-1.1.2/css/buttons.bootstrap.min.css"/>
    <link rel="stylesheet" type="text/css" href="datatables/FixedColumns-3.2.1/css/fixedColumns.bootstrap.min.css"/>
    <link rel="stylesheet" type="text/css" href="datatables/FixedHeader-3.1.1/css/fixedHeader.bootstrap.min.css"/>
    <link rel="stylesheet" type="text/css" href="datatables/Responsive-2.0.2/css/responsive.bootstrap.min.css"/>
    <link rel="stylesheet" type="text/css" href="datatables/Scroller-1.4.1/css/scroller.bootstrap.min.css"/>
    <link rel="stylesheet" type="text/css" href="datatables/Select-1.1.2/css/select.bootstrap.min.css"/>

    <!-- Custom CSS -->
    <link href="css/sb-admin-2.css" rel="stylesheet">

    <!-- Custom Fonts -->
    <link href="font-awesome/css/font-awesome.min.css" rel="stylesheet" type="text/css">

    <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
        <script src="https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
        <script src="https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->

</head>

<body>

    <div id="wrapper">

        <!-- Navigation -->
        <nav class="navbar navbar-default navbar-static-top" role="navigation" style="margin-bottom: 0">
            <div class="navbar-header">
                <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
                    <span class="sr-only">Toggle navigation</span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                </button>
                <a class="navbar-brand" href="index.html">
                    <img src="img/vsummary_logo.png" alt="vSummary Logo">
                </a>
            </div>
            <!-- /.navbar-header -->

            <ul class="nav navbar-top-links navbar-right">
                <li class="nav-button">
                    <!--menu toggle button -->
                    <button id="menu-toggle" type="button" data-toggle="button" class="btn btn-primary">
                        Full Width
                    </button>
                </li>
            </ul>
            <!-- /.navbar-top-links -->

            <!-- Sidebar wrapper over SB Admin 2 sidebar -->
            <div id="sidebar-wrapper">
                <div class="navbar-default sidebar" role="navigation">
                    <div class="sidebar-nav navbar-collapse">
                        <ul class="nav" id="side-menu">
                            <li class="sidebar-search">
                                <div class="input-group custom-search-form">
                                    <input type="text" class="form-control" placeholder="Search...">
                                    <span class="input-group-btn">
                                    <button class="btn btn-default" type="button">
                                        <i class="fa fa-search"></i>
                                    </button>
                                </span>
                                </div>
                                <!-- /input-group -->
                            </li>
                            <li>
                                <a href="index.php?view=vm"><i class="fa fa-desktop fa-fw"></i> Virtual Machines</a>
                            </li>
                            <li>
                                <a href="index.php?view=esxi"><i class="fa fa-server fa-fw"></i> ESXi Hosts</a>
                            </li>
                            <li>
                                <a href="index.php?view=portgroup"><i class="fa fa-cloud fa-fw"></i> Network Portgroups</a>
                            </li>
                            <li>
                                <a href="index.php?view=vdisk"><i class="fa fa-hdd-o fa-fw"></i> Virtual Disks</a>
                            </li>
                            <li>
                                <a href="index.php?view=datastore"><i class="fa fa-hdd-o fa-fw"></i> Datastores</a>
                            </li>
                            <li>
                                <a href="index.php?view=vnic"><i class="fa fa-plug fa-fw"></i> Virtual NICs</a>
                            </li>
                            <li>
                                <a href="#"><i class="fa fa-external-link fa-fw"></i> External Links<span class="fa arrow"></span></a>
                                <ul class="nav nav-second-level">
                                    <li>
                                        <a href="#">Log Insight</a>
                                    </li>
                                    <li>
                                        <a href="http://pubs.vmware.com/vsphere-60/topic/com.vmware.powercli.cmdletref.doc/Overview.html">PowerCLI Reference</a>
                                    </li>
                                </ul>
                                <!-- /.nav-second-level -->
                            </li>
                        </ul>
                    </div>
                    <!-- /.sidebar-collapse -->
                </div>
                <!-- /.navbar-static-side -->
            </div>
        </nav>

        <div id="page-wrapper" >
            <div class="row">
                <div class="col-lg-12">
                    <h1 class="page-header"><?php echo $view_title; ?> Summary</h1>
                </div>
                <!-- /.col-lg-12 -->
            </div>
            <!-- /.row -->
            <div class="row">
                <div class="col-lg-12">
                    <!-- /.dataTable_wrapper -->
                    <div class="dataTable_wrapper">
                        <table width="100%" class="table table-striped table-bordered table-hover" id="dt-vsummary-<?php echo $view; ?>">
						  <?php datatables_html($view); ?>
                        </table>
                    </div>
                </div>
                <!-- /.col-lg-12 -->
            </div>
        </div>
        <!-- /#page-wrapper -->

    </div>
    <!-- /#wrapper -->

    <!-- jQuery JavaScript -->
    <script src="js/jquery-2.2.0.min.js"></script>

    <!-- Bootstrap Core JavaScript -->
    <script src="bootstrap/js/bootstrap.min.js"></script>

    <!-- Metis Menu Plugin JavaScript -->
    <script src="js/metisMenu.min.js"></script>

    <!-- DataTables with extensions JavaScript -->
    <script type="text/javascript" src="datatables/JSZip-2.5.0/jszip.min.js"></script>
    <script type="text/javascript" src="datatables/pdfmake-0.1.18/build/pdfmake.min.js"></script>
    <script type="text/javascript" src="datatables/pdfmake-0.1.18/build/vfs_fonts.js"></script>
    <script type="text/javascript" src="datatables/DataTables-1.10.11/js/jquery.dataTables.min.js"></script>
    <script type="text/javascript" src="datatables/DataTables-1.10.11/js/dataTables.bootstrap.min.js"></script>
    <script type="text/javascript" src="datatables/Buttons-1.1.2/js/dataTables.buttons.min.js"></script>
    <script type="text/javascript" src="datatables/Buttons-1.1.2/js/buttons.bootstrap.min.js"></script>
    <script type="text/javascript" src="datatables/Buttons-1.1.2/js/buttons.colVis.min.js"></script>
    <script type="text/javascript" src="datatables/Buttons-1.1.2/js/buttons.html5.min.js"></script>
    <script type="text/javascript" src="datatables/Buttons-1.1.2/js/buttons.print.min.js"></script>
    <script type="text/javascript" src="datatables/FixedColumns-3.2.1/js/dataTables.fixedColumns.min.js"></script>
    <script type="text/javascript" src="datatables/FixedHeader-3.1.1/js/dataTables.fixedHeader.min.js"></script>
    <script type="text/javascript" src="datatables/Responsive-2.0.2/js/dataTables.responsive.min.js"></script>
    <script type="text/javascript" src="datatables/Responsive-2.0.2/js/responsive.bootstrap.min.js"></script>
    <script type="text/javascript" src="datatables/Scroller-1.4.1/js/dataTables.scroller.min.js"></script>
    <script type="text/javascript" src="datatables/Select-1.1.2/js/dataTables.select.min.js"></script>
    
    <!-- Custom Theme JavaScript -->
    <script src="js/sb-admin-2.js"></script>

    <!-- vSummary Datatables Load Javascript -->
    <script>
    
        $(document).ready(function() {

            // Setup - add a text input to each footer cell
            $('#dt-vsummary-<?php echo $view; ?> tfoot th').each( function () {
                var title = $(this).text();
                $(this).html( '<input type="text" placeholder="Search '+title+'" />' );
            });
            
            // Load Datatables
            var table = $('#dt-vsummary-<?php echo $view; ?>').DataTable({
                //dom: 'Blrtip',
                dom: "<'row'<'col-sm-6'l><'col-sm-6 text-right'B>><'row'<'col-sm-12'tr>><'row'<'col-sm-5'i><'col-sm-7'p>>",
                scrollY: '60vh',
                responsive: true,
                scrollX: true,
                stateSave: true,
                stateSaveParams: function (settings, data) {
                    // Loop through all columns and delete the search
                    for (var i = 0;i < data.columns.length; i++){
                        delete data.columns[i].search;
                    }
                },
                paging: true,
                pageLength: 15,
                lengthMenu: [[15, 25, 50, -1], [15, 25, 50, "All"]],
                scrollCollapse: true,
                processing: true,
                serverSide: true,
                ajax: "api/mysql_<?php echo $view; ?>.php",
                select: true,
                buttons: [
                  'copy', 
                  'csv',
                  { extend: 'colvis', className: 'colvis', text: 'Columns' }
                ]
            });
            
            // Apply the footer search
            table.columns().every( function () {
                var that = this;
                $( 'input', this.footer() ).on( 'keyup change', function () {
                    if ( that.search() !== this.value ) {
                        that
                            .search( this.value )
                            .draw();
                    }
                });
            });

        });

    </script>

    <!-- Sidebar Toggle Javascript -->
    <script type="text/javascript">
        $(document).ready(function() {
            $("#menu-toggle").click(function(e) {
                e.preventDefault();

                $("#wrapper").toggleClass("toggled");

                $('#wrapper.toggled').find("#sidebar-wrapper").find(".collapse").collapse('hide');

                // Redraw the datatables vsummary table
                $('#dt-vsummary-<?php echo $view; ?>').DataTable().draw();

            });
        });
    </script>


</body>

</html>