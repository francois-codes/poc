#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🧪 Testing Users API with NATS Integration${NC}"
echo "=============================================="

BASE_URL="http://localhost:8080/api"

# Test 1: Create a user
echo -e "\n${YELLOW}Test 1: Creating a user${NC}"
RESPONSE=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -H "X-User-ID: admin" \
  -d '{
    "email": "thomas@cognyx.io",
    "status": "active",
    "role": "admin"
  }')

echo "Response: $RESPONSE"

# Extract user ID from response (assuming JSON response with id field)
USER_ID=$(echo $RESPONSE | jq -r '.user_data.id')
echo "Created User ID: $USER_ID"

if [ "$USER_ID" != "null" ]; then
    echo -e "${GREEN}✅ User created successfully${NC}"
else
    echo -e "${RED}❌ Failed to create user${NC}"
    exit 1
fi

# Test 2: Get the user
echo -e "\n${YELLOW}Test 2: Getting the user${NC}"
curl -s $BASE_URL/users/$USER_ID | jq .

# Test 3: Update the user
echo -e "\n${YELLOW}Test 3: Updating the user${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT $BASE_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "X-User-ID: admin" \
  -d '{
    "email": "john.smith@example.com",
    "status": "active",
    "role": "admin"
  }')

echo "Update Response: $UPDATE_RESPONSE"

# Test 4: Get all versions of the user
echo -e "\n${YELLOW}Test 4: Getting all user versions${NC}"
curl -s $BASE_URL/users/$USER_ID/versions | jq .

# Test 5: Get specific version
echo -e "\n${YELLOW}Test 5: Getting version 1 of the user${NC}"
curl -s $BASE_URL/users/$USER_ID/versions/1 | jq .

echo -e "\n${GREEN}🎉 API Testing Complete!${NC}"
echo -e "${YELLOW}📡 Check NATS messages on users.broadcast channel${NC}"