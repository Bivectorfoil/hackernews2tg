#!/bin/bash

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# 获取项目根目录
PROJECT_ROOT="$( cd "$SCRIPT_DIR/../.." && pwd )"

# 检查必要文件是否存在
if [ ! -f "$PROJECT_ROOT/dist/newsboy" ]; then
    echo "ERROR: No executable file found ($PROJECT_ROOT/dist/newsboy)"
    echo "Please build the project first"
    exit 1
fi

if [ ! -f "$PROJECT_ROOT/.env" ]; then
    echo "ERROR: No configuration file found ($PROJECT_ROOT/.env)"
    echo "Please create the configuration file first"
    exit 1
fi

# 创建安装目录
mkdir -p /opt/newsboy/bin
mkdir -p /opt/newsboy/config

# 复制二进制文件和配置
cp "$PROJECT_ROOT/dist/newsboy" /opt/newsboy/bin/
cp "$PROJECT_ROOT/.env" /opt/newsboy/config/
chmod +x /opt/newsboy/bin/newsboy

# 获取当前用户
CURRENT_USER=$(whoami)

# 设置权限
chown -R "$CURRENT_USER:$CURRENT_USER" /opt/newsboy

# 复制服务文件
cp "$PROJECT_ROOT/deploy/systemd/newsboy.service" /etc/systemd/system/

# 重新加载systemd
systemctl daemon-reload

# 启动服务
systemctl start newsboy

# 设置开机自启
systemctl enable newsboy

# 检查服务状态
systemctl status newsboy

echo "Installation completed!"
echo "Use the following commands to manage the service:"
echo "Start: systemctl start newsboy"
echo "Stop: systemctl stop newsboy"
echo "Restart: systemctl restart newsboy"
echo "View logs: journalctl -u newsboy"
