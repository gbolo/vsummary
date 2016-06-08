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
Python progream to dump information from the vCenter's Database
"""

from __future__ import print_function

from pyVim.connect import SmartConnect, Disconnect
from pyVmomi import vim

import argparse
import atexit
import getpass
import ssl

from tools import cli
from tools import pchelper


def GetArgs():

   parser = argparse.ArgumentParser(description='')
   parser.add_argument('-s', '--host', required=True, action='store', help='')
   parser.add_argument('-o', '--port', type=int, default=443, action='store', help='')
   parser.add_argument('-u', '--user', required=True, action='store', help='')
   parser.add_argument('-p', '--password', required=False, action='store', help='')
   args = parser.parse_args()
   return args


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

    vm_properties = ["name",
                     "config.files.vmPathName",
                     "config.hardware.numCPU",
                     "config.hardware.memoryMB",
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



    return 0

# Start program
if __name__ == "__main__":
    main()
