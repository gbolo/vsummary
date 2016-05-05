
<!DOCTYPE html>
<html lang="en">
    <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->

    <link rel="icon" href="favicon.ico">

    <title>vSummary</title>

    <!-- Bootstrap core CSS -->
    <link href="bootstrap/css/bootstrap.min.css" rel="stylesheet">
      
    <!-- Custom styles for this template -->
    <link href="css/dashboard.css" rel="stylesheet">


    <!-- Bootstrap -->
    <script type="text/javascript" src="datatables/jQuery-2.2.0/jquery-2.2.0.min.js"></script>
    <script src="bootstrap/js/bootstrap.min.js"></script>

    <!-- DataTables -->

  <link rel="stylesheet" type="text/css" href="datatables/DataTables-1.10.11/css/dataTables.bootstrap.min.css"/>
  <link rel="stylesheet" type="text/css" href="datatables/Buttons-1.1.2/css/buttons.bootstrap.min.css"/>
  <link rel="stylesheet" type="text/css" href="datatables/FixedColumns-3.2.1/css/fixedColumns.bootstrap.min.css"/>
  <link rel="stylesheet" type="text/css" href="datatables/FixedHeader-3.1.1/css/fixedHeader.bootstrap.min.css"/>
  <link rel="stylesheet" type="text/css" href="datatables/Responsive-2.0.2/css/responsive.bootstrap.min.css"/>
  <link rel="stylesheet" type="text/css" href="datatables/Scroller-1.4.1/css/scroller.bootstrap.min.css"/>
  <link rel="stylesheet" type="text/css" href="datatables/Select-1.1.2/css/select.bootstrap.min.css"/>
   
  
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


    <script type="text/javascript" language="javascript" class="init">
  


var editor; // use a global for the submit and return data rendering in the examples

$(document).ready(function() {

  // Setup - add a text input to each footer cell
  $('#example tfoot th').each( function () {
      var title = $(this).text();
      $(this).html( '<input type="text" placeholder="Search '+title+'" />' );
  } );


    // Activate an inline edit on click of a table cell
    /*
    $('#example').on( 'click', 'tbody td:not(:first-child)', function (e) {
        editor.inline( this );
    } );
    */
    
  var table = $('#example').DataTable( {
    //dom: 'Blrtip',
    dom: "<'row'<'col-sm-6'l><'col-sm-6 text-right'B>><'row'<'col-sm-12'tr>><'row'<'col-sm-5'i><'col-sm-7'p>>",
    scrollY: '60vh',
    scrollX: true,
    stateSave: true,
    paging: true,
    pageLength: 15,
    lengthMenu: [[15, 25, 50, -1], [15, 25, 50, "All"]],
    scrollCollapse: true,
    ajax: {
      url: "api/mysql_vdisk.php",
      type: "POST"
    },
    serverSide: true,
    columns: [
      { data: "vdisk.name", className: "dt-body" },
      { data: "vdisk.capacity_bytes", className: "dt-body" },
      { data: "vdisk.path", className: "dt-body" },
      { data: "vdisk.thin_provisioned", className: "dt-body" },
      { data: "vdisk.datastore_id", className: "dt-body" },
      { data: "vdisk.uuid", className: "dt-body" },
      { data: "vdisk.disk_object_id", className: "dt-body" },
      { data: "vdisk.vm_id", className: "dt-body" },
      { data: "vdisk.esxi_id", className: "dt-body" },
      { data: "vdisk.vcenter_id", className: "dt-body" }
    ],
    select: true,
    buttons: [
      'copy', 
      'excel', 
      { 
        extend: 'pdfHtml5', 
        className: 'pdf', 
        orientation: 'landscape',
        exportOptions: {
          modifier: {
            page: 'current'
          }
        },
        text: 'PDF' 
      },
      { 
        extend: 'colvis', 
        className: 'colvis', 
        text: 'Column Visability' 
      }
    ]
  } );


    // Apply the search
    table.columns().every( function () {
        var that = this;
 
        $( 'input', this.footer() ).on( 'keyup change', function () {
            if ( that.search() !== this.value ) {
                that
                    .search( this.value )
                    .draw();
            }
        } );
    } );


} );



  </script>

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>

    <?php require_once('include/navbar.php'); ?>	


    <div class="container-fluid">
      <div class="row">
        <div class="col-lg-1 sidebar">
          
          <?php require_once('include/sidebar.php'); ?>

        </div>
        <div class="col-lg-11 col-lg-offset-1 col-md-11 col-md-offset-1 col-sm-11 col-sm-offset-1 main">


          <h2 class="sub-header">Virtual Disk Summary</h2>

            <table id="example" class="table table-striped table-bordered" cellspacing="0" width="100%">
                <thead>
                    <tr>
                      <th>Name</th>
                      <th>Capacity</th>
                      <th>Path</th>
                      <th>ThinProvisioned</th>
                      <th>datastore_id</th>
                      <th>uuid</th>
                      <th>disk_object_id</th>
                      <th>vm_id</th>
                      <th>esxi_id</th>
                      <th>vcenter_id</th>
                    </tr>
                </thead>
         
                <tfoot>
                    <tr>
                      <th>Name</th>
                      <th>Capacity</th>
                      <th>Path</th>
                      <th>ThinProvisioned</th>
                      <th>datastore_id</th>
                      <th>uuid</th>
                      <th>disk_object_id</th>
                      <th>vm_id</th>
                      <th>esxi_id</th>
                      <th>vcenter_id</th>
                    </tr>
                </tfoot>
            </table>



          </div>
        </div>
      </div>
    </div>



    <!-- Bootstrap core JavaScript
    ================================================== -->
    <!-- Placed at the end of the document so the pages load faster -->
    
    








  </body>
</html>