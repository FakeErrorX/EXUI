#!/bin/sh

# Start fail2ban
[ $EXUI_ENABLE_FAIL2BAN == "true" ] && fail2ban-client -x start

# Run ex-ui
exec /app/ex-ui
