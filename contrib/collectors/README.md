## vSummary - COLLECTORS
These collector scripts are used to pull data from the vcenter API, create a JSON object, then POST this data to the vSummary API for proccessing. Only one collector is needed, however you can run multiple collectors since they should produce the exact same data. The collectors leverage official vmware SDKs (powercli, python). These collectors should also be configured to run periodically (Linux cron / Widows Scheduler). The recommended interval between runs is at *least* once daily. 

### Architecture
![Alt text](https://raw.githubusercontent.com/gbolo/vSummary/master/screenshots/vsummary_arch.png "Architecture")
