#!/bin/bash

# Portfolio Backend VM Setup Script for Ubuntu 22.04
# Run this script on your Azure VM to prepare for deployment

set -e

echo "=== Portfolio Backend VM Setup ==="

# Update system packages
echo "Updating system packages..."
sudo apt update && sudo apt upgrade -y

# Install required packages
echo "Installing required packages..."
sudo apt install -y curl wget git unzip nginx mysql-client

# Install Docker (optional - for containerized deployment)
echo "Installing Docker..."
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Create application directory structure
echo "Creating application directories..."
mkdir -p /home/$USER/apps/portfolio-backend/{bin,migrations,logs}

# Set up nginx reverse proxy (optional)
echo "Setting up nginx reverse proxy..."
sudo tee /etc/nginx/sites-available/portfolio-backend > /dev/null <<EOF
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
    }

    # Health check endpoint
    location /v1/health {
        proxy_pass http://localhost:8080/v1/health;
        access_log off;
    }
}
EOF

# Enable nginx site
sudo ln -sf /etc/nginx/sites-available/portfolio-backend /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl enable nginx
sudo systemctl restart nginx

# Configure firewall
echo "Configuring firewall..."
sudo ufw allow OpenSSH
sudo ufw allow 'Nginx Full'
sudo ufw allow 8080/tcp
sudo ufw --force enable

# Create log rotation configuration
echo "Setting up log rotation..."
sudo tee /etc/logrotate.d/portfolio-backend > /dev/null <<EOF
/home/$USER/apps/portfolio-backend/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    copytruncate
    su $USER $USER
}
EOF

# Set up MySQL connection test script
echo "Creating MySQL connection test script..."
tee /home/$USER/test-mysql.sh > /dev/null <<'EOF'
#!/bin/bash
# Test MySQL connection
source /home/$USER/apps/portfolio-backend/.env 2>/dev/null || {
    echo "Environment file not found. Run deployment first."
    exit 1
}

mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD -e "SELECT 1;" $DB_NAME
if [ $? -eq 0 ]; then
    echo "✅ MySQL connection successful"
else
    echo "❌ MySQL connection failed"
    exit 1
fi
EOF

chmod +x /home/$USER/test-mysql.sh

echo ""
echo "=== VM Setup Complete! ==="
echo ""
echo "Next steps:"
echo "1. Set up your MySQL database (Azure Database for MySQL)"
echo "2. Configure GitHub repository secrets"
echo "3. Push your code to trigger deployment"
echo ""
echo "Test MySQL connection after deployment:"
echo "  ./test-mysql.sh"
echo ""
echo "View application logs:"
echo "  tail -f /home/$USER/apps/portfolio-backend/logs/app.out.log"
echo "  sudo journalctl -u portfolio-backend -f"
echo ""
echo "Manage service:"
echo "  sudo systemctl status portfolio-backend"
echo "  sudo systemctl restart portfolio-backend"
echo ""