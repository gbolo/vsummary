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

# Change this to "1" if you wanna see debug output
debug = 0


def GetArgs():

   parser = argparse.ArgumentParser(description='')
   parser.add_argument('-s', '--host', required=True, action='store', help='')
   parser.add_argument('-o', '--port', type=int, default=443, action='store', help='')
   parser.add_argument('-u', '--user', required=True, action='store', help='')
   parser.add_argument('-p', '--password', required=False, action='store', help='')
   parser.add_argument('-a', '--api', required=True, action='store', help='')
   args = parser.parse_args()
   return args


def vm_inventory(si, vc_uuid, api_url):

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

    root_folder = si.content.rootFolder
    view = pchelper.get_container_view(si, obj_type=[vim.VirtualMachine])
    vm_data = pchelper.collect_properties(si, view_ref=view,
                                          obj_type=vim.VirtualMachine,
                                          path_set=vm_properties,
                                          include_mors=True)

    #
    #  Creating variables matching the variables of the PowerCLI script
    #

    vm_data_compat = []

    for vm in vm_data:

        vm_compat = {}

        vm_compat['objecttype'] = "VM"

        if "name" in vm:
            vm_compat['name'] = vm["name"]
        else:
            vm_compat['name'] = "null"

        if "config.files.vmPathName" in vm:
            vm_compat['vmx_path'] = vm["config.files.vmPathName"]
        else:
            vm_compat['vmx_path'] = "null"

        if "config.hardware.numCPU" in vm:
            vm_compat['vcpu'] = vm["config.hardware.numCPU"]
        else:
            vm_compat['vcpu'] = "null"

        if "config.hardware.memoryMB" in vm:
            vm_compat['memory_mb'] = vm["config.hardware.memoryMB"]
        else:
            vm_compat['memory_mb'] = "null"

        if "config.guestId" in vm:
            vm_compat['config_guest_os'] = vm["config.guestId"]
        else:
            vm_compat['config_guest_os'] = "null"

        if "config.version" in vm:
            vm_compat['config_version'] = vm["config.version"]
        else:
            vm_compat['config_version'] = "null"

        if "config.uuid" in vm:
            vm_compat['smbios_uuid'] = vm["config.uuid"]
        else:
            vm_compat['smbios_uuid'] = "null"

        if "config.instanceUuid" in vm:
            vm_compat['instance_uuid'] = vm["config.instanceUuid"]
        else:
            vm_compat['instance_uuid'] = "null"

        if "config.changeVersion" in vm:
            vm_compat['config_change_version'] = vm["config.changeVersion"]
        else:
            vm_compat['config_change_version'] = "null"

        if "config.template" in vm:
            vm_compat['template'] = vm["config.template"]
        else:
            vm_compat['template'] = "null"

        if "guest.toolsVersion" in vm:
            vm_compat['guest_tools_version'] = vm["guest.toolsVersion"]
        else:
            vm_compat['guest_tools_version'] = "null"

        if "guest.toolsRunningStatus" in vm:
            vm_compat['guest_tools_running'] = vm["guest.toolsRunningStatus"]
        else:
            vm_compat['guest_tools_running'] = "null"

        if "guest.hostName" in vm:
            vm_compat['guest_hostname'] = vm["guest.hostName"]
        else:
            vm_compat['guest_hostname'] = "null"

        if "guest.ipAddress" in vm:
            vm_compat['guest_ip'] = vm["guest.ipAddress"]
        else:
            vm_compat['guest_ip'] = "null"

        if "config.guestId" in vm:
            vm_compat['config_guest_os'] = vm["config.guestId"]
        else:
            vm_compat['config_guest_os'] = "null"

        if "parent" in vm:
            _, ref = str(vm["parent"]).replace("'","").split(":")
            vm_compat['folder_moref'] = ref
        else:
            vm_compat['folder_moref'] = "null"

        if "parentVApp" in vm:
            _, ref = str(vm["parentVApp"]).replace("'","").split(":")
            vm_compat['vapp_moref'] = ref
        else:
            vm_compat['vapp_moref'] = "null"

        if "resourcePool" in vm:
            _, ref = str(vm["resourcePool"]).replace("'","").split(":")
            vm_compat['resourcepool_moref'] = ref
        else:
            vm_compat['resourcepool_moref'] = "null"

        if "summary.quickStats.overallCpuUsage" in vm:
            vm_compat['stat_cpu_usage'] = vm["summary.quickStats.overallCpuUsage"]
        else:
            vm_compat['stat_cpu_usage'] = "null"

        if "summary.quickStats.hostMemoryUsage" in vm:
            vm_compat['stat_host_memory_usage'] = vm["summary.quickStats.hostMemoryUsage"]
        else:
            vm_compat['stat_host_memory_usage'] = "null"

        if "summary.quickStats.guestMemoryUsage" in vm:
            vm_compat['stat_guest_memory_usage'] = vm["summary.quickStats.guestMemoryUsage"]
        else:
            vm_compat['stat_guest_memory_usage'] = "null"

        if "summary.quickStats.uptimeSeconds" in vm:
            vm_compat['stat_uptime_sec'] = vm["summary.quickStats.uptimeSeconds"]
        else:
            vm_compat['stat_uptime_sec'] = "null"

        if "runtime.powerState" in vm:
            power_state = 0
            if vm["runtime.powerState"] == "poweredOn":
                power_state = 1
            vm_compat['power_state'] = power_state
        else:
            vm_compat['power_state'] = "null"

        if "runtime.host" in vm:
            _, ref = str(vm["runtime.host"]).replace("'","").split(":")
            vm_compat['esxi_moref'] = ref
        else:
            vm_compat['esxi_moref'] = "null"

        if "obj" in vm:
            _, ref = str(vm["obj"]).replace("'","").split(":")
            vm_compat['moref'] = ref
        else:
            vm_compat['moref'] = "null"

        if vc_uuid:
            vm_compat['vcenter_id'] = vc_uuid
        else:
            vm_compat['vcenter_id'] = "null"

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

                        print(portGroupKey)
                        print(dvsUuid)


    #
    #  Generating the JSON post data
    #

    json_post_data = json.dumps(vm_data_compat)


    #
    #  The POST request itself
    #

    try:
        req = urllib2.Request(api_url)
        req.add_header('Content-Type', 'application/json')

        response = urllib2.urlopen(req, json_post_data)
        print (response.getcode())

    except:
        print ("HTTP Post Failed!")


    #
    #  DEBUG
    #

    if debug:
        for vm in vm_data:
            print("-" * 70)
            if "name" in vm:
                print("Name:                    {0}".format(vm["name"]))
            if "config.files.vmPathName" in vm:
                print("VM Path Name:            {0}".format(vm["config.files.vmPathName"]))
            if "config.hardware.numCPU" in vm:
                print("CPUs:                    {0}".format(vm["config.hardware.numCPU"]))
            if "config.hardware.memoryMB" in vm:
                print("MemoryMB:                {0}".format(vm["config.hardware.memoryMB"]))
            if "config.guestId" in vm:
                print("Guest ID:                {0}".format(vm["config.guestId"]))
            if "config.version" in vm:
                print("Container Version:       {0}".format(vm["config.version"]))
            if "config.uuid" in vm:
                print("BIOS UUID:               {0}".format(vm["config.uuid"]))
            if "config.instanceUuid" in vm:
                print("Instance UUID:           {0}".format(vm["config.instanceUuid"]))
            if "config.changeVersion" in vm:
                print("Change Version:          {0}".format(vm["config.changeVersion"]))
            if "config.template" in vm:
                print("Template:                {0}".format(vm["config.template"]))
            if "config.guestFullName" in vm:
                print("Guest Full Name:         {0}".format(vm["config.guestFullName"]))
            if "guest.toolsVersion" in vm:
                print("Guest Tools Version:     {0}".format(vm["guest.toolsVersion"]))
            if "guest.toolsRunningStatus" in vm:
                print("Guest Tools Running:     {0}".format(vm["guest.toolsRunningStatus"]))
            if "guest.hostName" in vm:
                print("Guest Hostname:          {0}".format(vm["guest.hostName"]))
            if "guest.ipAddress" in vm:
                print("Guest IP Address:        {0}".format(vm["guest.ipAddress"]))
            if "guest.guestState" in vm:
                print("Guest PowerState:        {0}".format(vm["guest.guestState"]))
            if "config.guestId" in vm:
                print("Guest Container Type:    {0}".format(vm["config.guestId"]))
            if "parent" in vm:
                print("Parent:                  {0}".format(vm["parent"]))
            if "parentVApp" in vm:
                print("Parent vApp:             {0}".format(vm["parentVApp"]))
            if "resourcePool" in vm:
                print("Resource Pool:           {0}".format(vm["resourcePool"]))
            if "summary.quickStats.overallCpuUsage" in vm:
                print("Quickstat CPU Usage:     {0}".format(vm["summary.quickStats.overallCpuUsage"]))
            if "summary.quickStats.hostMemoryUsage" in vm:
                print("Quickstat Host Memory:   {0}".format(vm["summary.quickStats.hostMemoryUsage"]))
            if "summary.quickStats.guestMemoryUsage" in vm:
                print("Quickstat Guest Memory:  {0}".format(vm["summary.quickStats.guestMemoryUsage"]))
            if "summary.quickStats.uptimeSeconds" in vm:
                print("Quickstat Uptime (sec):  {0}".format(vm["summary.quickStats.uptimeSeconds"]))
            if "runtime.powerState" in vm:
                print("Power State:             {0}".format(vm["runtime.powerState"]))
            if "runtime.host" in vm:
                print("Host:                    {0}".format(vm["runtime.host"]))

    print("")
    print("Found {0} VirtualMachines.".format(len(vm_data)))

    
def respool_inventory(si, vc_uuid, api_url):

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

    root_folder = si.content.rootFolder
    view = pchelper.get_container_view(si, obj_type=[vim.ResourcePool])
    respool_data = pchelper.collect_properties(si, view_ref=view,
                                               obj_type=vim.ResourcePool,
                                               path_set=respool_properties,
                                               include_mors=True)

    #
    #  Creating variables matching the variables of the PowerCLI script
    #

    respool_data_compat = []

    # print (respool_data)

    for respool in respool_data:

        respool_compat = {}

        respool_compat['objecttype'] = "RES"
        respool_compat['type'] = "ResourcePool"

        if "name" in respool:
            respool_compat['name'] = respool['name']
        else:
            respool_compat['name'] = "null"

        if "obj" in respool:
            _, ref = str(respool["obj"]).replace("'", "").split(":")
            respool_compat['moref'] = ref
        else:
            respool_compat['moref'] = "null"

        if "runtime.overallStatus"  in respool:
            respool_compat['status'] = respool['runtime.overallStatus']
        else:
            respool_compat['status'] = "null"

        if "parent" in respool:
            _, ref = str(respool["parent"]).replace("'", "").split(":")
            respool_compat['parent_moref'] = ref
        else:
            respool_compat['parent_moref'] = "null"

        if "owner" in respool:
            _, ref = str(respool["owner"]).replace("'", "").split(":")
            respool_compat['cluster_moref'] = ref
        else:
            respool_compat['cluster_moref'] = "null"

        if "summary.configuredMemoryMB" in respool:
            respool_compat['configured_memory_mb'] = respool['summary.configuredMemoryMB']
        else:
            respool_compat['configured_memory_mb'] = "null"

        if "summary.config.cpuAllocation.reservation" in respool:
            respool_compat['cpu_reservation'] = respool['summary.config.cpuAllocation.reservation']
        else:
            respool_compat['cpu_reservation'] = "null"

        if "summary.config.cpuAllocation.limit" in respool:
            respool_compat['cpu_limit'] = respool['summary.config.cpuAllocation.limit']
        else:
            respool_compat['cpu_limit'] = "null"

        if "summary.config.memoryAllocation.reservation" in respool:
            respool_compat['mem_reservation'] = respool['summary.config.memoryAllocation.reservation']
        else:
            respool_compat['mem_reservation'] = "null"

        if "summary.config.memoryAllocation.limit" in respool:
            respool_compat['mem_limit'] = respool['summary.config.memoryAllocation.limit']
        else:
            respool_compat['mem_limit'] = "null"

        if vc_uuid:
            respool_compat['vcenter_id'] = vc_uuid
        else:
            respool_compat['vcenter_id'] = "null"

        if "summary.config.entity.name" in respool:
            print (respool['summary.config.entity'])

        respool_data_compat.append(respool_compat)


    #
    #  Generating the JSON post data
    #

    json_post_data = json.dumps(respool_data_compat)

    #
    #  The POST request itself
    #

    try:
        req = urllib2.Request(api_url)
        req.add_header('Content-Type', 'application/json')

        response = urllib2.urlopen(req, json_post_data)
        print(response.getcode())

    except:
        print("HTTP Post Failed!")


def main():

    args = GetArgs()
    if args.password:
        password = args.password
    else:
        password = getpass.getpass(prompt='Password: ' % (args.host,args.user))

    context = ssl.SSLContext(ssl.PROTOCOL_TLSv1)
    context.verify_mode = ssl.CERT_NONE
    
    si = SmartConnect(host=args.host, user=args.user, pwd=password, port=int(args.port), sslContext=context)

    if not si:
        print("Could not connect ...")
        return -1

    atexit.register(Disconnect, si)


    # vm_inventory(si, "null", args.api)
    respool_inventory(si, "null", args.api)

    return 0

# Start program
if __name__ == "__main__":
    main()
