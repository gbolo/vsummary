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
        print("Name:                    {0}".format(vm["name"]))
        print("VM Path Name:            {0}".format(vm["config.files.vmPathName"]))
        print("CPUs:                    {0}".format(vm["config.hardware.numCPU"]))
        print("MemoryMB:                {0}".format(vm["config.hardware.memoryMB"]))
        print("Guest ID:                {0}".format(vm["config.guestId"]))
        print("Container Version:       {0}".format(vm["config.version"]))
        print("BIOS UUID:               {0}".format(vm["config.uuid"]))
        print("Instance UUID:           {0}".format(vm["config.instanceUuid"]))
        print("Change Version:          {0}".format(vm["config.changeVersion"]))
        print("Template:                {0}".format(vm["config.template"]))
        print("Guest Full Name:         {0}".format(vm["config.guestFullName"]))
        print("Guest Tools Version:     {0}".format(vm["guest.toolsVersion"]))
        print("Guest Tools Running:     {0}".format(vm["guest.toolsRunningStatus"]))
        print(vm["guest.hostName"])
	print("Guest Hostname:          {0}".format(vm["guest.hostName"]))
        print("Guest IP Address:        {0}".format(vm["guest.ipAddress"]))
        print("Guest PowerState:        {0}".format(vm["guest.guestState"]))
        print("Guest Container Type:    {0}".format(vm["config.guestId"]))
        print("Parent:                  {0}".format(vm["parent"]))
        print("Parent vApp:             {0}".format(vm["parentVApp"]))
        print("Resource Pool:           {0}".format(vm["resourcePool"]))
        print("Quickstat CPU Usage:     {0}".format(vm["summary.quickStats.overallCpuUsage"]))
        print("Quickstat Host Memory:   {0}".format(vm["summary.quickStats.hostMemoryUsage"]))
        print("Quickstat Guest Memory:  {0}".format(vm["summary.quickStats.guestMemoryUsage"]))
        print("Quickstat Uptime (sec):  {0}".format(vm["summary.quickStats.uptimeSeconds"]))
        print("Power State:             {0}".format(vm["runtime.powerState"]))
        print("Host:                    {0}".format(vm["runtime.host"]))

    print("")
    print("Found {0} VirtualMachines.".format(len(vm_data)))



    return 0

# Start program
if __name__ == "__main__":
    main()
