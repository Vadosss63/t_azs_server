#!/bin/bash
#Build Terminal 
#sudo apt install golang
install_dir="/opt/t_azs_server"
go build ./main.go

if [ -d ${install_dir} ]; then
sudo systemctl stop t_azs_server.service
sudo systemctl disable t_azs_server
sudo rm -rf ${install_dir}
fi

### create pay_server.service
printf "[Unit]

Description=t_azs_server
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=forking
ExecStart=%s/main
ExecReload=%s/main
WorkingDirectory=%s/
TimeoutSec=300
Restart=always

[Install]
WantedBy=multi-user.target" ${install_dir} ${install_dir} ${install_dir} > t_azs_server.service
### end create t_azs_server.service 

sudo mkdir ${install_dir} 
sudo cp main ${install_dir}/main
sudo cp -r public/ ${install_dir}/public/
sudo cp t_azs_server.service ${install_dir}/t_azs_server.service
sudo cp settings.json ${install_dir}/

sudo ln -s -f ${install_dir}/t_azs_server.service /etc/systemd/system/t_azs_server.service

sudo systemctl enable t_azs_server
sudo systemctl restart t_azs_server.service
#sudo reboot now
