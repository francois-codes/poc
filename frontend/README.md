# 🚀 RxDB + NATS Frontend

A lightweight frontend demonstrating real-time user management with RxDB (IndexedDB) and NATS messaging integration.

## Features

✅ **RxDB Integration**: Local IndexedDB storage with reactive queries
✅ **NATS Communication**: Real-time messaging with your backend
✅ **User Management**: Create, update, and manage users
✅ **Real-time Sync**: Automatic updates from backend broadcasts
✅ **Offline-First**: Works offline with local IndexedDB storage
✅ **Reactive UI**: Auto-updates when data changes

## Architecture

```
Frontend (Browser)           Backend (Go)
┌─────────────────┐         ┌──────────────────┐
│  RxDB           │         │  Users API       │
│  (IndexedDB)    │◄────────┤  + NATS          │
│                 │  HTTP   │                  │
│  ┌─────────────┐│         │ ┌──────────────┐ │
│  │ User Data   ││         │ │ PostgreSQL   │ │
│  │ Local Cache ││         │ │ + Versioning │ │
│  └─────────────┘│         │ └──────────────┘ │
└─────────────────┘         └──────────────────┘
         │                           │
         └───────── NATS ────────────┘
            users.broadcast
```

## Quick Start

### 1. Install Dependencies
```bash
cd frontend
npm install
```

### 2. Start the Frontend Server
```bash
npm start
```

The frontend will be available at: http://localhost:3000

### 3. Make Sure Backend is Running
Ensure your Users API is running on http://localhost:8080:
```bash
cd ../internal
go run cmd/users_api/main.go
```

### 4. Open in Browser
1. Go to http://localhost:3000
2. Click "Create User" to add users to both local RxDB and backend
3. Watch the logs to see real-time NATS messages
4. Users are automatically synced between IndexedDB and backend

## How It Works

### RxDB + IndexedDB
- **Local Storage**: All user data is stored locally in IndexedDB via RxDB
- **Reactive Queries**: UI automatically updates when data changes
- **Offline Support**: Works without internet connection
- **Schema Validation**: Type-safe user data with JSON schema

### NATS Integration
- **HTTP Fallback**: Uses HTTP API calls to your backend (WebSocket NATS coming soon)
- **Real-time Updates**: Receives broadcasts from `users.broadcast` channel
- **Bidirectional Sync**: Local changes trigger backend updates

### User Flow
1. **Create User**: Button click → HTTP POST → Backend → NATS broadcast → All clients updated
2. **Update User**: Button click → HTTP PUT → Backend → NATS broadcast → All clients updated
3. **Local Storage**: All operations update IndexedDB immediately
4. **Real-time Sync**: Other connected clients see changes instantly

## Features Demo

### 👤 User Management
- **Create User**: Add new users with email, status, and role
- **Update User**: Modify existing users (picks random user to update)
- **Local Display**: See all users stored in local IndexedDB
- **Delete Local**: Remove users from local storage (backend remains)

### 📡 NATS Messaging
- **Connection Status**: Shows NATS connection state
- **Message Logs**: Real-time display of all NATS messages
- **Auto-decoding**: Automatically decodes base64 payloads
- **Error Handling**: Graceful fallback to HTTP-only mode

### 💾 RxDB Features
- **Reactive UI**: List updates automatically when data changes
- **Local Persistence**: Data survives browser refresh
- **Schema Validation**: Ensures data integrity
- **Query Performance**: Fast local queries with indexing

## API Integration

The frontend communicates with your backend API:

### Create User
```javascript
POST http://localhost:8080/api/users
{
  "email": "user@example.com",
  "status": "active", 
  "role": "user"
}
```

### Update User
```javascript
PUT http://localhost:8080/api/users/:id
{
  "email": "updated@example.com",
  "status": "active",
  "role": "admin"
}
```

### NATS Messages
Receives messages from `users.broadcast`:
```json
{
  "id": "uuid",
  "user_id": 123,
  "operation": "create|update",
  "version": 1,
  "user_data": { /* complete user object */ },
  "timestamp": "2025-07-31T15:54:59Z",
  "created_by": "frontend-user"
}
```

## Development

### Project Structure
```
frontend/
├── index.html          # Main HTML file with RxDB + NATS integration
├── server.js           # Express server for serving static files
├── package.json        # Node.js dependencies
└── README.md          # This file
```

### Technologies Used
- **RxDB**: Reactive database with IndexedDB storage
- **NATS.js**: NATS client for browser (WebSocket mode)
- **Express**: Simple static file server
- **Vanilla JS**: No frameworks, pure JavaScript
- **Modern CSS**: Clean, responsive design

### Extending the Frontend

**Add new user fields**:
1. Update the RxDB schema in `setupRxDB()`
2. Add form fields in the HTML
3. Update `createUser()` and `updateUser()` functions

**Add more NATS channels**:
1. Subscribe to additional subjects in `connectToNATS()`
2. Handle different message types in `handleIncomingMessage()`

**Add real-time features**:
1. WebSocket connection to NATS server
2. Live user activity indicators
3. Collaborative editing

## Next Steps

🔄 **WebSocket NATS Connection**: Direct NATS WebSocket for real-time messaging
📊 **User Analytics**: Track user interactions and sync patterns  
🔍 **Search & Filter**: Local search with RxDB queries
🎨 **Better UI**: React/Vue integration with RxDB
📱 **Mobile Support**: Progressive Web App features
🧪 **Testing**: Unit tests for RxDB operations

## Troubleshooting

**"Failed to create user"**:
- Ensure backend API is running on http://localhost:8080
- Check CORS settings in backend
- Verify PostgreSQL connection

**"RxDB initialization failed"**:
- Check browser console for detailed errors
- Ensure IndexedDB is supported (modern browsers)
- Clear browser storage if schema changed

**"No NATS messages"**:
- Currently using HTTP fallback mode
- WebSocket NATS integration coming soon
- Check backend NATS connection

Your lightweight RxDB + NATS frontend is ready! 🎉