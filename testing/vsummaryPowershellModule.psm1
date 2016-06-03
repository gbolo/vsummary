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


