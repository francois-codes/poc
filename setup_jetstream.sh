#!/bin/bash

echo "🚀 Setting up JetStream streams for Users API"
echo "=============================================="

# Check if nats CLI is available
if ! command -v nats &> /dev/null; then
    echo "❌ NATS CLI not found. Installing..."
    # Install NATS CLI
    curl -sf https://binaries.nats.dev/nats-io/natscli/nats@latest | sh
    sudo mv nats /usr/local/bin/
fi

echo "📡 Creating JetStream streams..."

# Create stream for users.update (Frontend -> Backend)
nats stream add USERS_UPDATE \
    --subjects="users.update" \
    --storage=memory \
    --retention=limits \
    --max-msgs=1000 \
    --max-age=24h \
    --replicas=1 \
    --discard=old

# Create stream for users.broadcast (Backend -> Frontend) 
nats stream add USERS_BROADCAST \
    --subjects="users.broadcast" \
    --storage=memory \
    --retention=limits \
    --max-msgs=1000 \
    --max-age=24h \
    --replicas=1 \
    --discard=old

echo "✅ JetStream streams created successfully!"
echo ""
echo "📋 Stream information:"
nats stream ls
echo ""
echo "🎯 You can now run the Users API application"