#!/bin/bash

echo "🔍 NATS Message Monitor"
echo "======================"

echo "📊 Stream Information:"
./internal/nats stream info USERS_BROADCAST

echo ""
echo "📡 Subscribing to users.broadcast messages..."
echo "   (Press Ctrl+C to stop)"
echo ""

# Subscribe to messages and display them
./internal/nats subscribe users.broadcast --count=10