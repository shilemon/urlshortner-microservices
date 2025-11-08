# URL Shortener - Microservice Architecture Demo

A microservice-based URL shortener demonstrating proper service separation with three independent services: Go for high-performance redirects, Python for analytics and dashboard, and Node.js for URL metadata enrichment.

## Architecture

This project demonstrates a realistic microservice architecture where different services handle their specific responsibilities:

### Services

**Go Service (Port 8000)**

- **Purpose**: Fast URL redirection and creation
- **Database**: `go.db` (SQLite)
- **Responsibilities**:
  - Generate and store short codes
  - Handle URL redirects with minimal latency
  - Send click events to Python service asynchronously
- **Technology**: Go with Gin framework

**Python Service (Port 5000)**

- **Purpose**: Analytics, data aggregation, and user interface
- **Database**: `python.db` (SQLite)
- **Responsibilities**:
  - Provide web dashboard for URL creation
  - Orchestrate URL creation (call Go) and metadata fetching (call Node.js)
  - Collect and aggregate click events
  - Display analytics and statistics with metadata
  - Generate visualizations
- **Technology**: Python with Flask

**Node.js Service (Port 3000)**

- **Purpose**: URL metadata enrichment
- **Database**: `node.db` (SQLite)
- **Responsibilities**:
  - Fetch page titles, descriptions, and favicons from URLs
  - Parse HTML content with Cheerio
  - Store and serve metadata via REST API
- **Technology**: Node.js with Express, Axios, Cheerio

### Microservice Communication

```
User → Python Dashboard 
         ↓
         ├→ Go Service (Port 8000) → Create Short URL → go.db
         └→ Node.js Service (Port 3000) → Fetch Metadata → node.db
         ↓
    Display URL + Metadata in UI

User clicks Short URL → Go Service (Port 8000)
                         ↓
                    Redirect User
                         ↓
               Async Event → Python Service (Port 5000) → python.db
```

- **Python → Go**: HTTP POST to create URLs
- **Python → Node.js**: HTTP POST to fetch metadata
- **Go → Python**: Async HTTP POST for click events (fire-and-forget)
- **No direct database sharing**: Each service owns its data

## Features

- ✅ Create short URLs through web dashboard
- ✅ Fast redirects handled by Go
- ✅ **URL metadata enrichment via Node.js (titles, descriptions, favicons)**
- ✅ Real-time analytics dashboard
- ✅ Click tracking and history
- ✅ Visual charts for click patterns
- ✅ Top URLs by popularity with page info
- ✅ Recent activity monitoring
- ✅ Auto-refreshing dashboard (every 5 seconds)
- ✅ **Visual indicators showing Node.js service status**

## Prerequisites

- **Go**: Version 1.21 or higher
- **Python**: Version 3.14 (or 3.8+)
- **Node.js**: Version 16 or higher (with npm)
- **SQLite**: Built-in with Go, Python, and Node.js

## Installation & Setup

### 1. Clone or navigate to the project

```bash
cd /home/xaadu/codes/urlshortner
```

### 2. Setup Go Service

```bash
cd go-service

# Download dependencies
go mod download

# Run the service
go run main.go
```

The Go service will start on `http://localhost:8000`

### 3. Setup Python Service

Open a new terminal:

```bash
cd /home/xaadu/codes/urlshortner/python-service

# Create virtual environment (following user preference)
python3.14 -m venv venv

# Activate virtual environment
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Run the service
python app.py
```

The Python service will start on `http://localhost:5000`

### 4. Setup Node.js Service

Open a new terminal:

```bash
cd /home/xaadu/codes/urlshortner/node-service

# Install dependencies
npm install

# Run the service
node server.js
```

The Node.js service will start on `http://localhost:3000`

## Usage

### Access the Dashboard

Open your browser and navigate to:

```
http://localhost:5000
```

### Create a Short URL

1. Enter a long URL in the input field
2. Click "Shorten"
3. Copy the generated short URL

### Test the Redirect

Visit the short URL in your browser:

```
http://localhost:8000/{short_code}
```

You'll be redirected to the original URL, and the click will be tracked in the analytics.

### View Analytics

The dashboard automatically shows:

- Total URLs created
- Total clicks
- **Page metadata (titles, favicons) fetched by Node.js**
- Clicks over time (24-hour chart)
- Top URLs by popularity with page info
- All created URLs with metadata status indicators
- Recent click activity

The dashboard refreshes every 5 seconds automatically.

**Visual Indicators:**
- ✅ Green badge "✓ Node.js" = Metadata successfully fetched
- ❌ Red badge "✗" = Metadata fetch failed
- Favicon icons displayed next to page titles

## API Endpoints

### Go Service (Port 8000)

**Create Short URL**

```bash
POST /api/shorten
Content-Type: application/json

{
  "long_url": "https://example.com/very/long/url"
}

Response:
{
  "short_code": "abc123",
  "short_url": "http://localhost:8000/abc123",
  "long_url": "https://example.com/very/long/url"
}
```

**Redirect**

```bash
GET /{short_code}
# Redirects to the long URL and sends event to Python service
```

### Python Service (Port 5000)

**Dashboard**

```bash
GET /
# Returns the web dashboard
```

**Create URL (from UI)**

```bash
POST /create
Content-Type: application/x-www-form-urlencoded

long_url=https://example.com
```

**Receive Click Event**

```bash
POST /api/events
Content-Type: application/json

{
  "short_code": "abc123",
  "clicked_at": "2025-11-08T12:00:00Z"
}
```

**Get Statistics**

```bash
GET /api/stats

Returns JSON with:
- total_urls
- total_clicks
- top_urls (with metadata)
- recent_clicks
- clicks_over_time
- all_urls (with metadata)
```

### Node.js Service (Port 3000)

**Fetch Metadata**

```bash
POST /api/metadata
Content-Type: application/json

{
  "short_code": "abc123",
  "long_url": "https://example.com"
}

Response:
{
  "short_code": "abc123",
  "url": "https://example.com",
  "title": "Example Domain",
  "description": "Example domain for documentation",
  "favicon_url": "https://example.com/favicon.ico",
  "status": "success"
}
```

**Get Metadata**

```bash
GET /api/metadata/{short_code}
# Returns stored metadata for a short code
```

**Health Check**

```bash
GET /health
# Returns service health status
```

## Database Schema

### Go Service (go.db)

```sql
CREATE TABLE urls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_code TEXT UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Python Service (python.db)

```sql
CREATE TABLE click_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_code TEXT NOT NULL,
    clicked_at DATETIME NOT NULL
);

CREATE TABLE url_metadata (
    short_code TEXT PRIMARY KEY,
    long_url TEXT NOT NULL,
    total_clicks INTEGER DEFAULT 0,
    first_seen DATETIME NOT NULL,
    last_clicked DATETIME,
    title TEXT,
    description TEXT,
    favicon_url TEXT,
    metadata_status TEXT DEFAULT 'pending'
);
```

### Node.js Service (node.db)

```sql
CREATE TABLE metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_code TEXT UNIQUE NOT NULL,
    url TEXT NOT NULL,
    title TEXT,
    description TEXT,
    favicon_url TEXT,
    fetched_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Microservice Design Principles Demonstrated

1. **Service Independence**: Each service has its own database and can run independently
2. **Single Responsibility**: Go=Redirects, Python=Analytics/UI, Node.js=Metadata enrichment
3. **API Communication**: Services communicate via REST APIs, not direct database access
4. **Service Orchestration**: Python orchestrates calls to both Go and Node.js
5. **Asynchronous Operations**: Click events are sent asynchronously to avoid slowing redirects
6. **Graceful Degradation**: System works even if Node.js service is unavailable
7. **Data Ownership**: Each service owns and manages its own data
8. **Scalability**: Services can be scaled independently based on load

## Testing the System

### Test URL Creation and Redirection

```bash
# Create a short URL
curl -X POST http://localhost:5000/create \
  -d "long_url=https://github.com"

# Test redirect (will open in browser)
curl -L http://localhost:8000/{returned_short_code}

# Check analytics
curl http://localhost:5000/api/stats
```

### Verify Microservice Communication

1. Create a URL through the Python dashboard (e.g., https://github.com)
2. Check Go service logs - you should see the URL creation
3. Check Node.js service logs - you should see metadata fetching
4. Check Python service logs - you should see metadata stored
5. Look at the dashboard - you should see the page title and favicon
6. Click the short URL
7. Check Go service logs - you should see the redirect and event sending
8. Check Python service logs - you should see the click event received
9. Refresh the dashboard - you should see updated analytics with metadata

**Testing Node.js Service Separately:**
```bash
# Test metadata fetching directly
curl -X POST http://localhost:3000/api/metadata \
  -H "Content-Type: application/json" \
  -d '{"short_code":"test123","long_url":"https://github.com"}'

# Check health
curl http://localhost:3000/health
```

## Project Structure

```
/home/xaadu/codes/urlshortner/
├── README.md
├── go-service/
│   ├── main.go           # Go application (redirects & URL creation)
│   ├── go.mod            # Go dependencies
│   ├── go.sum            # Go dependency checksums
│   └── go.db             # SQLite database (created at runtime)
├── python-service/
│   ├── app.py            # Flask application (analytics & orchestration)
│   ├── requirements.txt   # Python dependencies
│   ├── python.db         # SQLite database (created at runtime)
│   └── templates/
│       └── dashboard.html # Web dashboard UI with metadata display
└── node-service/
    ├── server.js         # Express application (metadata fetching)
    ├── package.json      # Node.js dependencies
    └── node.db           # SQLite database (created at runtime)
```

## Technologies Used

- **Go 1.21+**: High-performance backend
  - Gin web framework
  - SQLite3 driver
- **Python 3.14**: Analytics and UI
  - Flask web framework
  - Requests library
  - SQLite3 (built-in)
- **Node.js 16+**: Metadata service
  - Express web framework
  - Axios (HTTP client)
  - Cheerio (HTML parsing)
  - SQLite3 driver
- **SQLite**: Lightweight database for all three services
- **Chart.js**: Data visualization
- **Modern CSS**: Responsive dashboard design

## Future Enhancements

- Add Redis for message queue between services
- Implement rate limiting
- Add user authentication
- Support custom short codes
- Add geographic tracking
- Implement URL expiration
- Add bulk URL creation
- Export analytics reports

## Author

[Zayed](https://zayedabdullah.com) | [Email](mailto:contact@zayedabdullah.com) | [GitHub (xaadu)](https://github.com/xaadu) | [LinkedIn](https://www.linkedin.com/in/abdullahzayed01/)

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## Support

If you find this project useful, please consider supporting me with a star or a follow.

## License

MIT License - Free to use for educational purposes
