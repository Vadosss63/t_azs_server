#!/bin/bash
#Build Terminal 
#sudo apt install golang
install_dir="/home/t_azs/t_azs_server/"
go build ./main.go

### create pay_server.service
printf "[Unit]

Description=t_azs_server
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=forking
ExecStart=/home/t_azs/t_azs_server/main
ExecReload=/home/t_azs/t_azs_server/main
WorkingDirectory=/home/t_azs/t_azs_server
TimeoutSec=300
Restart=always

[Install]
WantedBy=multi-user.target" >> t_azs_server.service
### end create t_azs_server.service 

sudo ln -s -f ${install_dir}t_azs_server.service /etc/systemd/system/t_azs_server.service
# sudo ln -s -f /opt/t_azs_server/t_azs_server.service /etc/systemd/system/t_azs_server.service


sudo systemctl enable t_azs_server
sudo systemctl restart t_azs_server.service

#sudo reboot now