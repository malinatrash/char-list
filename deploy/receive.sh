#!/bin/sh
# Принимает бинарник по stdin (ssh ... < char-list), подменяет и перезапускает сервис.
set -e
tmp=/opt/char-list/char-list.new
cat > "$tmp"
chmod 755 "$tmp"
mv "$tmp" /opt/char-list/char-list
sudo /bin/systemctl restart char-list
echo "deployed $(date -u +%FT%TZ)"
