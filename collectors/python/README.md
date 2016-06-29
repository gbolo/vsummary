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
