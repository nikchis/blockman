# Blockman
The scope of this application is to provide a simple secure tool to control the daemon's crashes. If daemon stops to work, the Blockman re-executes it and then opens wallet for staking.

The walletpassphrase and path to core directory stores locally.
Storage data encryption: AES-256.
Key computes using Scrypt.
Data for key generates from your machine's hardware information.
If it's not enough, you have two options to make it more secure:
-you can edit salt constant and recompile sources
-you can use additional storage password by providing it through flag 'spwd' (e.g. /blockman {other flags} -spwd your_secret)

###### Flags
-**path** set path to core files (blocknetdxd, blocknetdx-cli)
-**wpwd** set walletpassphrase
-**spwd** set additional storage passphrase (if you use it on set up, you must use it on print, execute too)
-**print** print data from storage (use it before execute to make sure everything is ok)
-**estake** execute staking with daemon recovering

###### Example 1 (try it)
```
./blockman/blockman -path ./blocknetcore
./blockman/blockman -wpwd mywlltsecpwd12345
./blockman/blockman -print
./blockman/blockman -estake
```

To use it in background just add 'nohup' and '&'.
Btw, you can use flags simultaneously.
###### Example 2 (one line with nohup and spwd)
```
nohup ./blockman/blockman -path ./blocknetcore -wpwd mywlltsecpwd12345 -spwd myst0rgsec -estake &
```

To use it 24/7 even after reboot you can create a systemd service.
(available in Centos 7)
###### Example 3 (service)
Set your data:
```
./blockman/blockman -path ./blocknetcore -wpwd mywlltsecpwd12345 -print
```
Create file *blockman.service*:
```
[Unit]
Description=Blockman Daemon
After=syslog.target
After=network.target

[Service]
Type=simple
User=root
ExecStart=/root/blockman/blockman -estake

[Install]
WantedBy=multi-user.target
```
Copy it to directory: */etc/systemd/system*
Enable autostart: *systemctl enable blockman*
Start it: *service blockman start*
Get the status: *service blockman status*
