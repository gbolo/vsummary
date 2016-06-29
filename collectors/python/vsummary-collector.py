#!/usr/bin/env python

#
#  Copyright (c) 2016 Frank Felhoffer, George Bolo
#
#  Permission is hereby granted, free of charge, to any person obtaining a 
#  copy of this software and associated documentation files (the "Software"),
#  to deal in the Software without restriction, including without limitation
#  the rights to use, copy, modify, merge, publish, distribute, sublicense, 
#  and/or sell copies of the Software, and to permit persons to whom the 
#  Software is furnished to do so, subject to the following conditions:
#
#  The above copyright notice and this permission notice shall be included 
#  in all copies or substantial portions of the Software.
#
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS 
#  OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
#  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL 
#  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
#  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING 
#  FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER 
#  DEALINGS IN THE SOFTWARE.
#

"""
Python program to dump information from the vCenter's Database
"""

from __future__ import print_function

from pyVim.connect import SmartConnect, Disconnect
from pyVmomi import vim

import argparse
import atexit
import getpass
import ssl
import json
import urllib2

from tools import cli
from tools import pchelper


def GetArgs():

   parser = argparse.ArgumentParser(description='')
   parser.add_argument('-s', '--host', required=True, action='store', help='')
   parser.add_argument('-o', '--port', type=int, default=443, action='store', help='')
   parser.add_argument('-u', '--user', required=True, action='store', help='')
   parser.add_argument('-p', '--password', required=False, action='store', help='')
   parser.add_argument('-a', '--api', required=True, action='store', help='')
   parser.add_argument('-d', '--dryrun', required=False, action='store_true', help='')
   parser.add_argument('-v', '--verbose', required=False, action='store_true', help='')
   args = parser.parse_args()
   return args


def send_vsummary_data(data, url):

    ret = {}

    if dryrun:
        ret['code'] = 0
        ret['reason'] = "DRYRUN!"
        return ret

    #
    #  Generating the JSON post data
    #

    json_post_data = json.dumps(data)


    #
    #  The POST request itself
    #

    try:
        req = urllib2.Request(url)
        req.add_header('Content-Type', 'application/json')

        response = urllib2.urlopen(req, json_post_data)

        res_code = response.getcode()

        if (res_code == 200):
            ret['code'] = res_code
            ret['reason'] = "OK!"
        else:
            ret['code'] = res_code
            ret['reason'] = "ERROR (" + res_code + ")"

        return ret

    except:
        ret['code'] = -1
        ret['reason'] = "FATAL!"
        return ret


def vm_inventory(si, vc_uuid, api_url):

    print("VM Inventory:")

    vm_properties = ["name",
                     "config.files.vmPathName",
                     "config.hardware.numCPU",
                     "config.hardware.memoryMB",
                     "config.hardware.device",
                     "config.guestId",
                     "config.version",
                     "config.uuid", 
                     "config.instanceUuid",
                     "config.changeVersion",
                     "config.template",
                     "config.guestFullName",
                     "guest.toolsVersion",
                     "guest.toolsRunningStatus",
                     "guest.hostName",
                     "guest.ipAddress",
                     "guest.guestState",
                     "parent",
                     "parentVApp",
                     "resourcePool",
                     "summary.quickStats.overallCpuUsage",
                     "summary.quickStats.hostMemoryUsage",
                     "summary.quickStats.guestMemoryUsage",
                     "summary.quickStats.uptimeSeconds",
                     "runtime.powerState",
                     "runtime.host"]

    #
    #
    #

    content = si.RetrieveContent()

    view = pchelper.get_container_view(si, obj_type=[vim.VirtualMachine])
    vm_data = pchelper.collect_properties(si, view_ref=view,
                                          obj_type=vim.VirtualMachine,
                                          path_set=vm_properties,
                                          include_mors=True)

    #
    #  Creating variables matching the variables of the PowerCLI script
    #

    vm_data_compat = []
    vnic_data_compat = []
    vdisk_data_compat = []

    for vm in vm_data:

        vm_compat = {}

        vm_compat['vcenter_id'] = vc_uuid
        vm_compat['objecttype'] = "VM"

        vm_compat['name'] = vm['name'] if "name" in vm else None
        vm_compat['vmx_path'] = vm['config.files.vmPathName'] if "config.files.vmPathName" in vm else None
        vm_compat['vcpu'] = vm['config.hardware.numCPU'] if "config.hardware.numCPU" in vm else None
        vm_compat['memory_mb'] = vm['config.hardware.memoryMB'] if "config.hardware.memoryMB" in vm else None
        vm_compat['config_guest_os'] = vm['config.guestId'] if "config.guestId" in vm else None
        vm_compat['config_version'] = vm['config.version'] if "config.version" in vm else None
        vm_compat['smbios_uuid'] = vm['config.uuid'] if "config.uuid" in vm else None
        vm_compat['instance_uuid'] = vm['config.instanceUuid'] if "config.instanceUuid" in vm else None
        vm_compat['config_change_version'] = vm['config.changeVersion'] if "config.changeVersion" in vm else None
        vm_compat['template'] = vm['config.template'] if "config.template" in vm else None
        vm_compat['guest_tools_version'] = vm['guest.toolsVersion'] if "guest.toolsVersion" in vm else None
        vm_compat['guest_tools_running'] = vm['guest.toolsRunningStatus'] if "guest.toolsRunningStatus" in vm else None
        vm_compat['guest_hostname'] = vm['guest.hostName'] if "guest.hostName" in vm else None
        vm_compat['guest_ip'] = vm['guest.ipAddress'] if "guest.ipAddress" in vm else None
        vm_compat['config_guest_os'] = vm['config.guestId'] if "config.guestId" in vm else None
        vm_compat['folder_moref'] = vm['parent']._moId if "parent" in vm else None
        vm_compat['vapp_moref'] = vm['parentVApp']._moId if "parentVApp" in vm else None
        vm_compat['resourcepool_moref'] = vm['resourcePool']._moId if "resourcePool" in vm else None
        vm_compat['stat_cpu_usage'] = vm['summary.quickStats.overallCpuUsage'] if "summary.quickStats.overallCpuUsage" in vm else None
        vm_compat['stat_host_memory_usage'] = vm['summary.quickStats.hostMemoryUsage'] if "summary.quickStats.hostMemoryUsage" in vm else None
        vm_compat['stat_guest_memory_usage'] = vm['summary.quickStats.guestMemoryUsage'] if "summary.quickStats.guestMemoryUsage" in vm else None
        vm_compat['stat_uptime_sec'] = vm['summary.quickStats.uptimeSeconds'] if "summary.quickStats.uptimeSeconds" in vm else None
        vm_compat['esxi_moref'] = vm['runtime.host']._moId if "runtime.host" in vm else None
        vm_compat['moref'] = vm['obj']._moId if "obj" in vm else None

        if "runtime.powerState" in vm:
            vm_compat['power_state'] = 1 if vm["runtime.powerState"] == "poweredOn" else 0
        else:
            vm_compat['power_state'] = None

        vm_data_compat.append(vm_compat)


        #
        #  Processing the vnic information
        #

        if "config.hardware.device" in vm:
            for dev in vm["config.hardware.device"]:
                if isinstance(dev, vim.vm.device.VirtualEthernetCard):
                    dev_backing = dev.backing
                    if hasattr(dev_backing, 'port'):

                        portGroupKey = dev.backing.port.portgroupKey
                        dvsUuid = dev.backing.port.switchUuid
                        switch_type = "VmwareDistributedVirtualSwitch"

                        try:
                            dvs = content.dvSwitchManager.QueryDvsByUuid(dvsUuid)
                        except:
                            portGroup = "** Error: DVS not found **"
                            vlanId = "NA"
                            vSwitch = "NA"
                            portgroup_moref = "NA"
                        else:
                            pgObj = dvs.LookupDvPortGroup(portGroupKey)
                            portGroup = pgObj.config.name
                            vlanId = str(pgObj.config.defaultPortConfig.vlan.vlanId)
                            vSwitch = str(dvs.name)
                            portgroup_moref = pgObj._moId
                    else:
                        portGroup = dev.backing.network.name
                        switch_type = "HostVirtualSwitch"
                        pg_key = vm_compat['esxi_moref'] + "_" + portGroup
                        if pg_key in host_portgroups:
                            vSwitch, vlanId = host_portgroups[pg_key].split(':')
                        else:
                            vSwitch = None
                            vlanId = None
                        portgroup_moref = None

                    if portGroup is None:
                        portGroup = "NA"

                    #
                    #  Generating PowerCLI Compatible Output
                    #

                    vnic_compat = {}

                    vnic_compat["vm_moref"] = vm_compat['moref']
                    vnic_compat["esxi_moref"] = vm_compat['esxi_moref']
                    vnic_compat["vcenter_id"] = vm_compat['vcenter_id']
                    vnic_compat["objecttype"] = "VNIC"
                    vnic_compat["name"] = dev.deviceInfo.label
                    vnic_compat["mac"] = dev.macAddress
                    vnic_compat["connected"] = dev.connectable.connected
                    vnic_compat["status"] = dev.connectable.status
                    vnic_compat["portgroup_name"] = portGroup
                    vnic_compat["portgroup_moref"] = portgroup_moref
                    vnic_compat["vswitch_name"] = vSwitch
                    vnic_compat["vswitch_type"] = switch_type

                    if isinstance(dev, vim.vm.device.VirtualE1000):
                        vnic_compat["type"] = "VirtualE1000"
                    elif isinstance(dev, vim.vm.device.VirtualE1000e):
                        vnic_compat["type"] = "VirtualE1000e"
                    elif isinstance(dev, vim.vm.device.VirtualVmxnet3):
                        vnic_compat["type"] = "VirtualVmxnet3"
                    elif isinstance(dev, vim.vm.device.VirtualPCNet32):
                        vnic_compat["type"] = "VirtualPCNet32"
                    else:
                        vnic_compat["type"] = "N/A"
                        # vnic_compat["type"] = str(type(dev))

                    vnic_data_compat.append(vnic_compat)

                # Equal to Get-VSVirtualDisk function
                if isinstance(dev, vim.vm.device.VirtualDisk):

                    vdisk_compat = {}

                    vdisk_compat['vcenter_id'] = vc_uuid
                    vdisk_compat['objecttype'] = "VDISK"
                    vdisk_compat['name'] = dev.deviceInfo.label
                    vdisk_compat['capacity_bytes'] = dev.capacityInBytes
                    vdisk_compat['capacity_kb'] = dev.capacityInKB
                    vdisk_compat['path'] = dev.backing.fileName
                    vdisk_compat['thin_provisioned'] = dev.backing.thinProvisioned
                    vdisk_compat['datastore_moref'] = dev.backing.datastore._moId
                    vdisk_compat['uuid'] = dev.backing.uuid
                    vdisk_compat['disk_object_id'] = dev.diskObjectId
                    vdisk_compat['vm_moref'] = vm_compat['moref']
                    vdisk_compat['esxi_moref'] = vm_compat['esxi_moref']

                    vdisk_data_compat.append(vdisk_compat)

    print("  + Found {} VMs, {} vNICs, {} vDisks.".format(len(vm_data), len(vnic_data_compat), len(vdisk_data_compat)))

    #
    #  Sending over data
    #

    vm_ret = send_vsummary_data(vm_data_compat, api_url)
    vnic_ret = send_vsummary_data(vnic_data_compat, api_url)
    vdisk_ret = send_vsummary_data(vdisk_data_compat, api_url)

    print("  + Sending VMs: {}, vNICs: {}, vDisks: {}.".format(vm_ret['reason'], vnic_ret['reason'],
                                                               vdisk_ret['reason']))

    if verbose:
        print(json.dumps(vm_data_compat, indent=4, sort_keys=True))
        print(json.dumps(vnic_data_compat, indent=4, sort_keys=True))
        print(json.dumps(vdisk_data_compat, indent=4, sort_keys=True))

    
def respool_inventory(si, vc_uuid, api_url):

    # TODO: vApp Support might be added

    print("Resource Pool Inventory")

    respool_properties = ["name",
                          "owner",
                          "parent",
                          "runtime.overallStatus",
                          "summary.configuredMemoryMB",
                          "summary.config.cpuAllocation.reservation",
                          "summary.config.cpuAllocation.limit",
                          "summary.config.memoryAllocation.reservation",
                          "summary.config.memoryAllocation.limit",
                          "summary.config.entity"]

    view = pchelper.get_container_view(si, obj_type=[vim.ResourcePool])
    respool_data = pchelper.collect_properties(si, view_ref=view,
                                               obj_type=vim.ResourcePool,
                                               path_set=respool_properties,
                                               include_mors=True)

    #
    #  Creating variables matching the variables of the PowerCLI script
    #

    respool_data_compat = []

    for respool in respool_data:

        respool_compat = {}

        respool_compat['vcenter_id'] = vc_uuid
        respool_compat['objecttype'] = "RES"
        respool_compat['type'] = "ResourcePool"

        respool_compat['name'] = respool['name'] if "name" in respool else None
        respool_compat['moref'] = respool["obj"]._moId if "obj" in respool else None
        respool_compat['status'] = respool['runtime.overallStatus'] if "runtime.overallStatus"  in respool else None
        respool_compat['vapp_state'] = None
        respool_compat['parent_moref'] = respool["parent"]._moId if "parent" in respool else None
        respool_compat['cluster_moref'] = respool["owner"]._moId if "owner" in respool else None
        respool_compat['configured_memory_mb'] = respool['summary.configuredMemoryMB'] if "summary.configuredMemoryMB" in respool else None
        respool_compat['cpu_reservation'] = respool['summary.config.cpuAllocation.reservation'] if "summary.config.cpuAllocation.reservation" in respool else None
        respool_compat['cpu_limit'] = respool['summary.config.cpuAllocation.limit'] if "summary.config.cpuAllocation.limit" in respool else None
        respool_compat['mem_reservation'] = respool['summary.config.memoryAllocation.reservation'] if "summary.config.memoryAllocation.reservation" in respool else None
        respool_compat['mem_limit'] = respool['summary.config.memoryAllocation.limit'] if "summary.config.memoryAllocation.limit" in respool else None

        respool_data_compat.append(respool_compat)

    print("  + Found {} Resource Pools.".format(len(respool_data_compat)))

    # Sending
    rpool_ret = send_vsummary_data(respool_data_compat, api_url)

    print("  + Sending Resource Pools: {}".format(rpool_ret['reason']))

    if verbose:
        print(json.dumps(respool_data_compat, indent=4, sort_keys=True))


def host_inventory(si, vc_uuid, api_url):

    print("Host Inventory")

    host_properties = ["name",
                       "parent",
                       "summary.maxEVCModeKey",
                       "summary.currentEVCModeKey",
                       "summary.overallStatus",
                       "summary.runtime.powerState",
                       "summary.runtime.inMaintenanceMode",
                       "summary.hardware.vendor",
                       "summary.hardware.model",
                       "summary.hardware.uuid",
                       "summary.hardware.memorySize",
                       "summary.hardware.cpuModel",
                       "summary.hardware.cpuMhz",
                       "summary.hardware.numCpuPkgs",
                       "summary.hardware.numCpuCores",
                       "summary.hardware.numCpuThreads",
                       "summary.hardware.numNics",
                       "summary.hardware.numHBAs",
                       "summary.config.product.version",
                       "summary.config.product.build",
                       "summary.quickStats.overallCpuUsage",
                       "summary.quickStats.overallMemoryUsage",
                       "summary.quickStats.uptime",
                       "config.network.pnic",
                       "config.network.vswitch",
                       "config.network.portgroup"]

    view = pchelper.get_container_view(si, obj_type=[vim.HostSystem])

    host_data = pchelper.collect_properties(si, view_ref=view, obj_type=vim.HostSystem,
                                            path_set=host_properties, include_mors=True)

    host_data_compat = []
    pnic_data_compat = []
    vswitch_data_compat = []
    portgroup_data_compat = []

    for host in host_data:

        host_compat = {}

        host_compat['objecttype'] = "ESXI"
        host_compat['vcenter_id'] = vc_uuid

        host_compat['name'] = host['name'] if "name" in host else None
        host_compat['moref'] = host['obj']._moId if "obj" in host else None
        host_compat['max_evc'] = host['summary.maxEVCModeKey'] if "summary.maxEVCModeKey" in host else None
        host_compat['current_evc'] = host['summary.currentEVCModeKey'] if "summary.currentEVCModeKey" in host else None
        host_compat['status'] = host['summary.overallStatus'] if "summary.overallStatus" in host else None
        host_compat['in_maintenance_mode'] = host['summary.runtime.inMaintenanceMode'] if "summary.runtime.inMaintenanceMode" in host else None
        host_compat['vendor'] = host['summary.hardware.vendor'] if "summary.hardware.vendor" in host else None
        host_compat['model'] = host['summary.hardware.model'] if "summary.hardware.model" in host else None
        host_compat['uuid'] = host['summary.hardware.uuid'] if "summary.hardware.uuid" in host else None
        host_compat['memory_bytes'] = host['summary.hardware.memorySize'] if "summary.hardware.memorySize" in host else None
        host_compat['cpu_model'] = host['summary.hardware.cpuModel'] if "summary.hardware.cpuModel" in host else None
        host_compat['cpu_mhz'] = host['summary.hardware.cpuMhz'] if "summary.hardware.cpuMhz" in host else None
        host_compat['cpu_sockets'] = host['summary.hardware.numCpuPkgs'] if "summary.hardware.numCpuPkgs" in host else None
        host_compat['cpu_cores'] = host['summary.hardware.numCpuCores'] if "summary.hardware.numCpuCores" in host else None
        host_compat['cpu_threads'] = host['summary.hardware.numCpuThreads'] if "summary.hardware.numCpuThreads" in host else None
        host_compat['nics'] = host['summary.hardware.numNics'] if "summary.hardware.numNics" in host else None
        host_compat['hbas'] = host['summary.hardware.numHBAs'] if "summary.hardware.numHBAs" in host else None
        host_compat['version'] = host['summary.config.product.version'] if "summary.config.product.version" in host else None
        host_compat['build'] = host['summary.config.product.build'] if "summary.config.product.build" in host else None
        host_compat['stat_cpu_usage'] = host['summary.quickStats.overallCpuUsage'] if "summary.quickStats.overallCpuUsage" in host else None
        host_compat['stat_memory_usage'] = host['summary.quickStats.overallMemoryUsage'] if "summary.quickStats.overallMemoryUsage" in host else None
        host_compat['stat_uptime_sec'] = host['summary.quickStats.uptime'] if "summary.quickStats.uptime" in host else None
        host_compat['cluster_moref'] = host['parent']._moId if "parent" in host else None

        if "summary.runtime.powerState" in host:
            host_compat['power_state'] = 1 if host['summary.runtime.powerState'] == "poweredOn" else 0
        else:
            host_compat['power_state'] = None

        host_data_compat.append(host_compat)

        #
        #  Get-VSPhysicalNic function
        #

        if "config.network.pnic" in host:

            for pnic in host['config.network.pnic']:

                if isinstance(pnic, vim.host.PhysicalNic):

                    pnic_compat = {}

                    pnic_compat['vcenter_id'] = vc_uuid
                    pnic_compat['objecttype'] = "PNIC"

                    pnic_compat['esxi_moref'] = host_compat['moref']
                    pnic_compat['name'] = pnic.device
                    pnic_compat['mac'] = pnic.mac
                    pnic_compat['driver'] = pnic.driver
                    pnic_compat['link_speed'] = "Down"
                    if hasattr(pnic.linkSpeed, 'speedMb'):
                        pnic_compat['link_speed'] = pnic.linkSpeed.speedMb

                    pnic_data_compat.append(pnic_compat)

        #
        #  Get-VSStandardVswitch function
        #

        if "config.network.vswitch" in host:

            for vswitch in host['config.network.vswitch']:

                if isinstance(vswitch, vim.host.VirtualSwitch):

                    vswitch_compat = {}

                    vswitch_compat['vcenter_id'] = vc_uuid
                    vswitch_compat['objecttype'] = "SVS"

                    vswitch_compat['name'] = vswitch.name
                    vswitch_compat['ports'] = vswitch.spec.numPorts
                    vswitch_compat['max_mtu'] = vswitch.mtu
                    vswitch_compat['esxi_moref'] = host_compat['moref']

                    vswitch_data_compat.append(vswitch_compat)

        #
        #  Get-VSStandardPortGroup function
        #

        if "config.network.portgroup" in host:

            for pg in host['config.network.portgroup']:

                if isinstance(pg, vim.host.PortGroup):

                    pg_compat = {}

                    pg_compat['vcenter_id'] = vc_uuid
                    pg_compat['objecttype'] = "SVSPG"

                    pg_compat['name'] = pg.spec.name
                    pg_compat['vswitch_name'] = pg.spec.vswitchName
                    pg_compat['vlan'] = pg.spec.vlanId
                    pg_compat['esxi_moref'] = host_compat['moref']

                    portgroup_data_compat.append(pg_compat)

                    # Generating Port Group data for lookups
                    pg_key = host_compat['moref'] + "_" + pg_compat['name']
                    host_portgroups[pg_key] = pg_compat['vswitch_name'] + ":" + str(pg_compat['vlan'])

    print("  + Found {} Hosts, {} pNICs, {} vSwitces, {} Port Groups.".format(len(host_data_compat),
                                                                              len(pnic_data_compat),
                                                                              len(vswitch_data_compat),
                                                                              len(portgroup_data_compat)))

    # Sending
    host_ret = send_vsummary_data(host_data_compat, api_url)
    pnic_ret = send_vsummary_data(pnic_data_compat, api_url)
    vsw_ret = send_vsummary_data(vswitch_data_compat, api_url)
    pg_ret = send_vsummary_data(portgroup_data_compat, api_url)

    print("  + Sending Hosts: {}, pNICs: {}, vSwitches: {}, PortGroups: {}".format(host_ret['reason'],
                                                                                   pnic_ret['reason'],
                                                                                   vsw_ret['reason'],
                                                                                   pg_ret['reason']))

    if verbose:
        print(json.dumps(host_data_compat, indent=4, sort_keys=True))
        print(json.dumps(pnic_data_compat, indent=4, sort_keys=True))
        print(json.dumps(vswitch_data_compat, indent=4, sort_keys=True))
        print(json.dumps(portgroup_data_compat, indent=4, sort_keys=True))


def datastore_inventory(si, vc_uuid, api_url):

    print("Datastore Inventory")

    datastore_properties = ["name",
                            "overallStatus",
                            "summary.capacity",
                            "summary.freeSpace",
                            "summary.type",
                            "summary.uncommitted"]

    view = pchelper.get_container_view(si, obj_type=[vim.Datastore])

    datastore_data = pchelper.collect_properties(si, view_ref=view, obj_type=vim.Datastore,
                                                 path_set=datastore_properties, include_mors=True)

    datastore_data_compat = []

    for ds in datastore_data:

        ds_compat = {}

        ds_compat['vcenter_id'] = vc_uuid
        ds_compat['objecttype'] = "DS"
        ds_compat['name'] = ds['name'] if "name" in ds else None
        ds_compat['moref'] = ds['obj']._moId
        ds_compat['status'] = ds['overallStatus'] if "overallStatus" in ds else None
        ds_compat['capacity_bytes'] = ds['summary.capacity'] if "summary.capacity" in ds else None
        ds_compat['free_bytes'] = ds['summary.freeSpace'] if "summary.freeSpace" in ds else None
        ds_compat['uncommitted_bytes'] = ds['summary.uncommitted'] if "summary.uncommitted" in ds else None
        ds_compat['type'] = ds['summary.type'] if "summary.type" in ds else None

        datastore_data_compat.append(ds_compat)

    print("  + Found {} Data Stores.".format(len(datastore_data_compat)))

    # Sending
    ds_ret = send_vsummary_data(datastore_data_compat, api_url)

    print("  + Sending Data Stores: {}".format(ds_ret['reason']))

    if verbose:
        print(json.dumps(datastore_data_compat, indent=4, sort_keys=True))


def datacenter_inventory(si, vc_uuid, api_url):

    print("Data Center Inventory")

    datacenter_properties = ["name",
                             "hostFolder",
                             "vmFolder"]

    view = pchelper.get_container_view(si, obj_type=[vim.Datacenter])

    datacenter_data = pchelper.collect_properties(si, view_ref=view, obj_type=vim.Datacenter,
                                                  path_set=datacenter_properties, include_mors=True)

    datacenter_data_compat = []

    for dc in datacenter_data:

        dc_compat = {}

        dc_compat['vcenter_id'] = vc_uuid
        dc_compat['objecttype'] = "DC"

        dc_compat['name'] = dc['name'] if "name" in dc else None
        dc_compat['moref'] = dc['obj']._moId if "obj" in dc else None
        dc_compat['vm_folder_moref'] = dc['vmFolder']._moId if "vmFolder" in dc else None
        dc_compat['esxi_folder_moref'] = dc['hostFolder']._moId if "hostFolder" in dc else None

        datacenter_data_compat.append(dc_compat)

    print("  + Found {} Data Centers.".format(len(datacenter_data_compat)))

    # Sending
    dc_ret = send_vsummary_data(datacenter_data_compat, api_url)

    print("  + Sending Data Centers: {}".format(dc_ret['reason']))

    if verbose:
        print(json.dumps(datacenter_data_compat, indent=4, sort_keys=True))


def folder_inventory(si, vc_uuid, api_url):

    print("Folder Inventory")

    folder_properties = ["name",
                         "parent",
                         "childType"]

    view = pchelper.get_container_view(si, obj_type=[vim.Folder])

    folder_data = pchelper.collect_properties(si, view_ref=view, obj_type=vim.Folder,
                                              path_set=folder_properties, include_mors=True)

    folder_data_compat = []

    for folder in folder_data:

        folder_compat = {}

        folder_compat['vcenter_id'] = vc_uuid
        folder_compat['objecttype'] = "FOLDER"
        folder_compat['name'] = folder['name'] if "name" in folder else None
        folder_compat['moref'] = folder['obj']._moId if "obj" in folder else None
        folder_compat['parent_moref'] = folder['parent']._moId if "parent" in folder else None

        if "childType" in folder:
            type_str = " ".join(folder['childType'])
            folder_compat['type'] = type_str
        else:
            folder_compat['type'] = None

        folder_data_compat.append(folder_compat)

    print("  + Found {} Folders.".format(len(folder_data_compat)))

    # Sending
    folder_ret = send_vsummary_data(folder_data_compat, api_url)

    print("  + Sending Folders: {}".format(folder_ret['reason']))

    if verbose:
        print(json.dumps(folder_data_compat, indent=4, sort_keys=True))


def cluster_inventory(si, vc_uuid, api_url):

    print("Cluster Inventory")

    cluster_properties = ["name",
                          "parent",
                          "overallStatus",
                          "configuration.dasConfig",
                          "configuration.drsConfig",
                          "summary"]

    view = pchelper.get_container_view(si, obj_type=[vim.ClusterComputeResource])

    cluster_data = pchelper.collect_properties(si, view_ref=view, obj_type=vim.ClusterComputeResource,
                                               path_set=cluster_properties, include_mors=True)

    cluster_data_compat = []

    for cluster in cluster_data:

        cluster_compat = {}

        cluster_compat['vcenter_id'] = vc_uuid
        cluster_compat['objecttype'] = "CLUSTER"

        cluster_compat['name'] = cluster['name'] if "name" in cluster else None
        cluster_compat['moref'] = cluster['obj']._moId if "obj" in cluster else None
        cluster_compat['datacenter_moref'] = cluster['parent']._moId if "parent" in cluster else None
        cluster_compat['total_cpu_threads'] = cluster['summary'].numCpuThreads if "summary" in cluster else None
        cluster_compat['total_cpu_mhz'] = cluster['summary'].totalCpu if "summary" in cluster else None
        cluster_compat['total_memory_bytes'] = cluster['summary'].totalMemory if "summary" in cluster else None
        cluster_compat['total_vmotions'] = cluster['summary'].numVmotions if "summary" in cluster else None
        cluster_compat['num_hosts'] = cluster['summary'].numHosts if "summary" in cluster else None
        cluster_compat['current_balance'] = cluster['summary'].currentBalance if "summary" in cluster else None
        cluster_compat['target_balance'] = cluster['summary'].targetBalance if "summary" in cluster else None
        cluster_compat['drs_enabled'] = cluster['configuration.drsConfig'].enabled if "configuration.drsConfig" in cluster else None
        cluster_compat['drs_behaviour'] = cluster['configuration.drsConfig'].defaultVmBehavior if "configuration.drsConfig" in cluster else None
        cluster_compat['ha_enabled'] = cluster['configuration.dasConfig'].enabled if "configuration.dasConfig" in cluster else None
        cluster_compat['status'] = cluster['overallStatus'] if "overallStatus" in cluster else None

        cluster_data_compat.append(cluster_compat)

    print("  + Found {} Clusters.".format(len(cluster_data_compat)))

    # Sending
    cluster_ret = send_vsummary_data(cluster_data_compat, api_url)

    print("  + Sending Clusters: {}".format(cluster_ret['reason']))

    if verbose:
        print(json.dumps(cluster_data_compat, indent=4, sort_keys=True))


def dvs_inventory(si, vc_uuid, api_url):

    print("Distributed Virtual Switch Inventory")

    dvs_properties = ["name",
                      "summary.productInfo.version",
                      "config"]

    view = pchelper.get_container_view(si, obj_type=[vim.DistributedVirtualSwitch])

    dvs_data = pchelper.collect_properties(si, view_ref=view, obj_type=vim.DistributedVirtualSwitch,
                                               path_set=dvs_properties, include_mors=True)

    dvs_data_compat = []

    for dvs in dvs_data:

        dvs_compat = {}

        dvs_compat['vcenter_id'] = vc_uuid
        dvs_compat['objecttype'] = "DVS"

        dvs_compat['name'] = dvs['name'] if "name" in dvs else None
        dvs_compat['moref'] = dvs['obj']._moId if "obj" in dvs else None
        dvs_compat['version'] = dvs['summary.productInfo.version'] if "summary.productInfo.version" in dvs else None
        dvs_compat['max_mtu'] = dvs['config'].maxMtu if "config" in dvs else None
        dvs_compat['ports'] = dvs['config'].numPorts if "config" in dvs else None

        dvs_data_compat.append(dvs_compat)

    print("  + Found {} DVS.".format(len(dvs_data_compat)))

    # Sending
    dvs_ret = send_vsummary_data(dvs_data_compat, api_url)

    print("  + Sending DVS: {}".format(dvs_ret['reason']))

    if verbose:
        print(json.dumps(dvs_data_compat, indent=4, sort_keys=True))


def dvs_portgroup_inventory(si, vc_uuid, api_url):

    print("DVS Port Group Inventory")

    dvspg_properties = ["name",
                        "config.defaultPortConfig",
                        "config.distributedVirtualSwitch"]

    view = pchelper.get_container_view(si, obj_type=[vim.DistributedVirtualPortgroup])

    dvspg_data = pchelper.collect_properties(si, view_ref=view, obj_type=vim.DistributedVirtualPortgroup,
                                             path_set=dvspg_properties, include_mors=True)

    dvspg_data_compat = []

    for dvspg in dvspg_data:

        vlan = dvspg['config.defaultPortConfig'].vlan

        if isinstance(vlan, vim.dvs.VmwareDistributedVirtualSwitch.VlanIdSpec):
            vlan_type = "VmwareDistributedVirtualSwitchVlanIdSpec"
            vlan_id = vlan.vlanId
            vlan_start = "na"
            vlan_end = "na"

        # The API needs to be fixed to be able to implement this type of Port Groups
        elif isinstance(vlan, vim.dvs.VmwareDistributedVirtualSwitch.TrunkVlanSpec):
            vlan_type = "VmwareDistributedVirtualSwitchTrunkVlanSpec"
            vlan_id = "na"
            vlan_start = "na"
            vlan_end = "na"

        else:
            vlan_type = "TypeNotImplemented"
            vlan_id = "na"
            vlan_start = "na"
            vlan_end = "na"

        dvspg_compat = {}

        dvspg_compat['vcenter_id'] = vc_uuid
        dvspg_compat['objecttype'] = "DVSPG"

        dvspg_compat['name'] = dvspg['name'] if "name" in dvspg else None
        dvspg_compat['moref'] = dvspg['obj']._moId if "obj" in dvspg else None
        dvspg_compat['vlan_type'] = vlan_type
        dvspg_compat['vlan'] = vlan_id
        dvspg_compat['vlan_start'] = vlan_start
        dvspg_compat['vlan_end'] = vlan_end
        dvspg_compat['dvs_moref'] = dvspg['config.distributedVirtualSwitch']._moId if "config.distributedVirtualSwitch" in dvspg else None

        dvspg_data_compat.append(dvspg_compat)

    print("  + Found {} DVS Port Groups.".format(len(dvspg_data_compat)))

    # Sending
    dvspg_ret = send_vsummary_data(dvspg_data_compat, api_url)

    print("  + Sending DVS Port Groups: {}".format(dvspg_ret['reason']))

    if verbose:
        print(json.dumps(dvspg_data_compat, indent=4, sort_keys=True))


def main():

    global dryrun, verbose, host_portgroups
    host_portgroups = {}

    args = GetArgs()
    if args.password:
        password = args.password
    else:
        password = getpass.getpass(prompt='Password for %s@%s: ' % (args.user,args.host))

    if args.dryrun: dryrun = True
    else: dryrun = False

    if args.verbose: verbose = True
    else: verbose = False

    context = ssl.SSLContext(ssl.PROTOCOL_TLSv1)
    context.verify_mode = ssl.CERT_NONE
    
    si = SmartConnect(host=args.host, user=args.user, pwd=password, port=int(args.port), sslContext=context)

    if not si:
        print("Could not connect ...")
        return -1

    atexit.register(Disconnect, si)

    # Figuring out the UUID of the vcenter server
    content = si.RetrieveContent()

    if content.about.instanceUuid:
        vc_uuid = content.about.instanceUuid
    else:
        vc_uuid = None

    # host_inventory have to be run before vm_inventory as collected information here is being used
    # to get information about host related network objects such as standard switces, etc.
    host_inventory(si, vc_uuid, args.api)
    vm_inventory(si, vc_uuid, args.api)
    datastore_inventory(si, vc_uuid, args.api)
    respool_inventory(si, vc_uuid, args.api)
    datacenter_inventory(si, vc_uuid, args.api)
    dvs_inventory(si, vc_uuid, args.api)
    dvs_portgroup_inventory(si, vc_uuid, args.api)

    # Folder check needs to be done after datacenter check (API requirement)
    folder_inventory(si, vc_uuid, args.api)
    cluster_inventory(si, vc_uuid, args.api)

    return 0

# Start program
if __name__ == "__main__":
    main()
