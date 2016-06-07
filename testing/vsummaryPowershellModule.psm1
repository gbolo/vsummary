# VSUMMARY POWERSHELL MODULE

Function Connect-vcenter ($vc_fqdn){

    if ($global:DefaultVIServers.Count -gt 0) {
        Disconnect-VIServer -Server * -Force -Confirm:$false -WarningAction SilentlyContinue -ErrorAction SilentlyContinue | Out-Null
    }

    $creds = get-credential
    $c = Connect-VIServer $vc_fqdn -Credential $creds

    if ($c){
        $vc_uuid = $c.InstanceUuid
        Write-Host "Connected to vCenter: $vc_fqdn"
    } else{
        Write-Host "Could NOT Connect to vCenter: $vc_fqdn"
        Exit;
    }

}

Function Import-VM-Folders ( [object]$folderArray ){

    $datacenter = Read-Host "DESTINATION DataCenter Name"
    $folderArray | % {
     $startFolder = Get-Datacenter -Name $datacenter | Get-Folder -Name 'vm' -NoRecursion
        $path = $_
     
        $location = $startFolder
        echo $location
        $path.Split('/') | Select -skip 1 | %{
            # decode the folder name to reveal true folder name
            $folder=[System.Web.HttpUtility]::UrlDecode($_)
            Try {
                echo "GET: $folder LOC: $location"
                $location = Get-Folder -Name $folder -Location $location -NoRecursion -ErrorAction Stop
            }
            Catch{
                echo "NEW: $folder LOC: $location"
                $location = New-Folder -Name $folder -Location $location
            }
        } 
        echo "======="
    }

}


# source: http://www.lucd.info/2012/05/18/folder-by-path/
Function Get-FolderByPath {
  <# .SYNOPSIS Retrieve folders by giving a path .DESCRIPTION The function will retrieve a folder by it's path. The path can contain any type of leave (folder or datacenter). .NOTES Author: Luc Dekens .PARAMETER Path The path to the folder. This is a required parameter. .PARAMETER Path The path to the folder. This is a required parameter. .PARAMETER Separator The character that is used to separate the leaves in the path. The default is '/' .EXAMPLE PS> Get-FolderByPath -Path "Folder1/Datacenter/Folder2"
.EXAMPLE
  PS> Get-FolderByPath -Path "Folder1>Folder2" -Separator '>'
#>
 
  param(
  [CmdletBinding()]
  [parameter(Mandatory = $true)]
  [System.String[]]${Path},
  [char]${Separator} = '/'
  )
 
  process{
    if((Get-PowerCLIConfiguration).DefaultVIServerMode -eq "Multiple"){
      $vcs = $defaultVIServers
    }
    else{
      $vcs = $defaultVIServers[0]
    }
 
    foreach($vc in $vcs){
      foreach($strPath in $Path){
        $root = Get-Folder -Name Datacenters -Server $vc
        $strPath.Split($Separator) | %{
          $foldername = [System.Web.HttpUtility]::UrlDecode($_)
          $root = Get-Inventory -Name $foldername -Location $root -Server $vc -NoRecursion
          if((Get-Inventory -Location $root -NoRecursion | Select -ExpandProperty Name) -contains "vm"){
            $root = Get-Inventory -Name "vm" -Location $root -Server $vc -NoRecursion
          }
        }
        $root | where {$_ -is [VMware.VimAutomation.ViCore.Impl.V1.Inventory.FolderImpl]}|%{
          Get-Folder -Name $foldername -Location $root.Parent -NoRecursion -Server $vc
        }
      }
    }
  }
}


# source: http://www.lucd.info/2016/06/03/vsphere-object-path
function Get-VIObjectByPath{
<#
.SYNOPSIS
  Retrieve a vSphere object by it's path.
.DESCRIPTION
  This function will retrieve a vSphere object from it's path.
  The path can be absolute or relative.
  When a relative path is provided, the StartNode needs
  to be provided
.NOTES
  Author:  Luc Dekens
.PARAMETER StartNode
  The vSphere Server (vCenter or ESXi) from which to retrieve
  the objects.
  The default is $Global:DefaultVIServer
.PARAMETER Path
  A string with the absolute or relative path.
  The path shall not contain any hidden folders.
.EXAMPLE
  PS> Get-VIObjectByPath -Path '/Datacenter/Folder/VM1'
.EXAMPLE
  PS> Get-InventoryPlus -StartNode $node -Path $path
#>
 
  param(
    [VMware.Vim.ManagedEntity]$StartNode = (
      Get-View -Id (Get-View -Id ServiceInstance).Content.RootFolder
    ),
    [String]$Path
  )
 
  function Get-NodeChild{
    param(
      [VMware.Vim.ManagedEntity]$Node
    )
   
    $hidden = 'vm','host','network','datastore','Resources'
    switch($Node){
      {$_ -is [VMware.Vim.Folder]}{
        if($Node.ChildEntity){
          Get-View -Id $Node.ChildEntity
        }
      }
      {$_ -is [VMware.Vim.Datacenter]}{
        $all = @()
        $all += Get-View -Id $Node.VmFolder
        $all += Get-View -Id $Node.HostFolder
        $all += Get-View -Id $Node.DatastoreFolder
        $all += Get-View -Id $Node.NetworkFolder
        $all | %{
          if($hidden -contains $_.Name){
            Get-NodeChild -Node $_
          }
          else{
            $_
          }
        }
      }
      {$_ -is [VMware.Vim.ClusterComputeResource]}{
        $all = @()
        $all += Get-View -Id $Node.Host
        $all += Get-View -Id $Node.ResourcePool 
        $all = $all | %{
          if($hidden -contains $_.Name){
            Get-NodeChild -Node $_
          }
          else{
            $_
          }
        }
        $all
      }
      {$_ -is [VMware.Vim.ResourcePool]}{
        $all = @()
        if($Node.ResourcePool){
          $all += Get-View -Id $Node.ResourcePool
        }
        if($Node.vm){
          $all += Get-View -Id $Node.vm
        }
        $all
      }
      {$_ -is [VMware.Vim.DistributedVirtualSwitch]}{
        Get-View -Id $Node.Portgroup
      }
    }
  }
 
  $found = $true
 
  # Loop through Path
  $node = $StartNode
  foreach($qualifier in $Path.TrimStart('/').Split('/',[StringSplitOptions]::RemoveEmptyEntries)){
    $nodeMatch = @($node) | %{
      Get-NodeChild -Node $_ | where{$_.Name -eq $qualifier}
    }
    if(!$nodeMatch){
      $found = $false
      $node = $null
      break
    }
    $node = $nodeMatch
  }
 
  New-Object PSObject -Property @{
    Path = $Path
    Found = $found
    Node = $node
  }
}

Function Export-Cluster-Vapps ( [string]$cluster_name ) {

    Write-Host "!! WARNING !! - STOPING ALL vApps IN CLUSTER: $cluster_name"
    Write-Host "  Make sure that there are no VMs in these vapps or they will be shutdown"
    $confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "
    if ($confirmation -eq 'y') {
        foreach($vapp in Get-Cluster $cluster_name | Get-VApp){
            $vapp | Stop-VApp -force
            $vapp | Export-VApp -destination ".\vapps"
        }
    }
}

Function Import-Cluster-Vapps ( [string]$cluster_name, [string]$folder) {

    Write-Host "IMPORTING vApps IN CLUSTER: $cluster_name"
    $confirmation = Read-Host "  Are you Sure You Want To Proceed?: (y/n) "
    if ($confirmation -eq 'y') {
        foreach( $vapp in Get-ChildItem $folder -recurse | Where {$_.extension -eq ".ovf"} ){
            $ovf_path = $vapp.FullName
            Import-vApp -Source "$ovf_path" -VMHost (get-vmhost "esxi1.lab1.midgar.local") -Location (get-cluster $cluster_name) -force
        }
    }
}


Function Import-Cluster-ResoucePools ( [object]$resourcePoolArray ) {

  $resourcePoolArray | % {
    $pool_name = $_.Name
    $pool_path = [System.Web.HttpUtility]::UrlDecode($_.ParentPath)
    $search = Get-VIObjectByPath -Path $pool_path

    # create resourcepool only if found
    if ($search.Found){
      $moref = $search.Node.Moref.ToString()
      $location = Get-ResourcePool -Id $moref
      Write-Host "CREATING RESOURCEPOOL: $pool_name"
      # create this pool
      New-ResourcePool -Location $location -Name $pool_name
    } else{
      Write-Host "ERROR - PARENT RESOURCEPOOL NOT FOUND: $pool_path"
    }

  }

}


