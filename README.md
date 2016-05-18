![Alt text](https://raw.githubusercontent.com/gbolo/vSummary/master/src/img/vsummary_logo.png "vSummary Logo")

vSummary is an open source  tool for collecting and displaying a summary of your vSphere Environment(s).

For a **LIVE DEMO**, please click this link: 
[http://vsummary.linuxctl.com/index.php?view=vm](http://vsummary.linuxctl.com/index.php?view=vm) 

### Screenshots
![Alt text](https://raw.githubusercontent.com/gbolo/vSummary/master/screenshots/screenshot_1.png "Screenshot 1")

### Architecture
vSummary is essentially a web application with both a frontend and backend. The baceknd accepts HTTP POST data in json format which it then normalizes and inserts/updates into mysql tables. The frontend is where it disaplys this data for users to see. Here is a basic architectural diagram to visualize this:

![Alt text](https://raw.githubusercontent.com/gbolo/vSummary/master/screenshots/vsummary_arch.png "Architecture")

### Requirements

The following requirements for vSummary have been identified so far:
* WEB SERVER (nginx,apache,...)
* PHP 5.3+ (datatables php lib)
* MYSQL 5.0+ (support create views)
* POWERSHELL 3.0+ (convert-json, http-request)
* POWERCLI 5.5+ (check api calls)
* vCenter 5.5+ (check api calls)

### Installation

1. Use an existing, or prepare a new web server that is able to execute php code. Ensure php can handle: pdo-mysql, json.
2. Use an existing, or prepare a new mysql database that is at least version 5.0.
3. Create a new database called vsummary with the following schema: [mysql_schema](https://github.com/gbolo/vSummary/blob/master/sql/vsummary_mysql_schema.sql)
4. Create a new mysql user and grant permissions to this database.
5. Deploy vsummary source code ([src](https://github.com/gbolo/vSummary/tree/master/src) folder) to the web root of the web server.
6. Modify the file [mysql_config.php](https://github.com/gbolo/vSummary/blob/master/src/api/lib/mysql_config.php) with the correct database information.
7. Prepare a Windows environment which has powershell version 3+ and vSphere PowerCLI 5.5+ installed
8. Allow execution of powershell files that are not signed: `Set-ExecutionPolicy -Scope "CurrentUser" -ExecutionPolicy "unrestricted"`
9. Download the powershell script [vsummary_collect.ps1](https://github.com/gbolo/vSummary/blob/master/powershell/vsummary_collect.ps1) and modify the vcenter server list and api endpoint located near the end of the script.
10. Execute the powershell script, then once complete visit your webserver address to see the results.
11. Create an automated job to run this script X amount of times per day.

### Docker

For a quicker deployment, a docker image is available (which does steps 1 to 6) with preinstalled nginx, php-fpm, mariadb, and vsummary source code. To run it please execute these commands:

**Start container and bind it to port 80 on local machine**
```
docker run --name vsummary -p 80:80 -d gbolo/vsummary
```
~~**(optional) Load sample data into the database for testing**~~
~~if you would like to load sample data into vsummary for testing, you may execute a php script inside the conatiner to do so:~~

**sample data generator script is broken for now until it's updated to support the newer data schema that has changed recently**
```
docker exec -it vsummary php /data/gen_sample_data.php

POSTING RANDOM SAMPLE DATA FOR VSUMMARY API: http://localhost/api/update.php
---
[vcenter] SUCCESS! RESPONSE: 200 TIME: 0.006185s
[esxi] SUCCESS! RESPONSE: 200 TIME: 0.00689s
[dvs] SUCCESS! RESPONSE: 200 TIME: 0.005439s
[datastore] SUCCESS! RESPONSE: 200 TIME: 0.004999s
[vm] SUCCESS! RESPONSE: 200 TIME: 0.010547s
[portgroup] SUCCESS! RESPONSE: 200 TIME: 0.007113s
[pnic] SUCCESS! RESPONSE: 200 TIME: 0.004621s
[vnic] SUCCESS! RESPONSE: 200 TIME: 0.008292s
[vdisk] SUCCESS! RESPONSE: 200 TIME: 0.010178s
```

### Development

This tool is under much development. **ANY CONTRIBUTIONS WILL BE GREATLY APPRECIATED**


License
----

MIT


**Free Software, Hell Yeah!**
