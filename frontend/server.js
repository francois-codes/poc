const express = require('express');
const path = require('path');
const cors = require('cors');

const app = express();
const PORT = 3000;

// Enable CORS for all routes
app.use(cors());

// Serve static files from current directory
app.use(express.static(__dirname));

// Serve the main HTML file
app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'index.html'));
});

app.listen(PORT, () => {
    console.log(`🚀 Frontend server running at http://localhost:${PORT}`);
    console.log(`📱 Open your browser and go to http://localhost:${PORT}`);
    console.log(`🔗 Make sure your backend API is running on http://localhost:8080`);
    console.log(`📡 Make sure NATS JetStream is running on localhost:4222`);
});