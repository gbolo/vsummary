### Installation on CentOS 7

#### 1. Install PIP
```
cd /tmp
curl "https://bootstrap.pypa.io/get-pip.py" -o "get-pip.py"
python get-pip.py
```

Verification:
```
pip --help
pip -V
```

#### 2. Install the vmware SDK
```
pip install pyvmomi
```


### Example output of the script
```
Host Inventory
  + Found 119 Hosts, 550 pNICs, 119 vSwitces, 234 Port Groups.
  + Sending Hosts: OK!, pNICs: OK!, vSwitches: OK!, PortGroups: OK!
VM Inventory:
  + Found 2900 VMs, 5262 vNICs, 3357 vDisks.
  + Sending VMs: OK!, vNICs: OK!, vDisks: OK!.
Datastore Inventory
  + Found 10 Data Stores.
  + Sending Data Stores: OK!
Resource Pool Inventory
  + Found 56 Resource Pools.
  + Sending Resource Pools: OK!
Data Center Inventory
  + Found 1 Data Centers.
  + Sending Data Centers: OK!
Distributed Virtual Switch Inventory
  + Found 1 DVS.
  + Sending DVS: OK!
DVS Port Group Inventory
  + Found 194 DVS Port Groups.
  + Sending DVS Port Groups: OK!
Folder Inventory
  + Found 628 Folders.
  + Sending Folders: OK!
Cluster Inventory
  + Found 9 Clusters.
  + Sending Clusters: OK!

-----------
Time spent: 135.0 seconds.
```