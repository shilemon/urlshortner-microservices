# URL Shortener - Microservice Architecture Demo

A microservice-based URL shortener demonstrating proper service separation, with Go handling high-performance redirects and Python managing analytics and the dashboard UI.

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
  - Collect and aggregate click events
  - Display analytics and statistics
  - Generate visualizations
- **Technology**: Python with Flask

### Microservice Communication

```
User → Python Dashboard → Go Service → Database (go.db)
                              ↓
                        Redirect User
                              ↓
                    Async Event → Python Service → Database (python.db)
```

- **Python → Go**: HTTP POST to create URLs
- **Go → Python**: Async HTTP POST for click events (fire-and-forget)
- **No direct database sharing**: Each service owns its data

## Features

- ✅ Create short URLs through web dashboard
- ✅ Fast redirects handled by Go
- ✅ Real-time analytics dashboard
- ✅ Click tracking and history
- ✅ Visual charts for click patterns
- ✅ Top URLs by popularity
- ✅ Recent activity monitoring
- ✅ Auto-refreshing dashboard (every 5 seconds)

## Prerequisites

- **Go**: Version 1.21 or higher
- **Python**: Version 3.14 (or 3.8+)
- **SQLite**: Built-in with both Go and Python

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
- Clicks over time (24-hour chart)
- Top URLs by popularity
- All created URLs
- Recent click activity

The dashboard refreshes every 5 seconds automatically.

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
- top_urls
- recent_clicks
- clicks_over_time
- all_urls
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
    last_clicked DATETIME
);
```

## Microservice Design Principles Demonstrated

1. **Service Independence**: Each service has its own database and can run independently
2. **Single Responsibility**: Go focuses on performance (redirects), Python focuses on analytics
3. **API Communication**: Services communicate via REST APIs, not direct database access
4. **Asynchronous Operations**: Click events are sent asynchronously to avoid slowing redirects
5. **Data Ownership**: Each service owns and manages its own data
6. **Scalability**: Services can be scaled independently based on load

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

1. Create a URL through the Python dashboard
2. Check Go service logs - you should see the URL creation
3. Click the short URL
4. Check Go service logs - you should see the redirect and event sending
5. Check Python service logs - you should see the event received
6. Refresh the dashboard - you should see updated analytics

## Project Structure

```
xaadu/codes/urlshortner/
├── README.md
├── go-service/
│   ├── main.go           # Go application
│   ├── go.mod            # Go dependencies
│   └── go.db             # SQLite database (created at runtime)
└── python-service/
    ├── app.py            # Flask application
    ├── requirements.txt   # Python dependencies
    ├── python.db         # SQLite database (created at runtime)
    └── templates/
        └── dashboard.html # Web dashboard UI
```

## Technologies Used

- **Go 1.24+**: High-performance backend
  - Gin web framework
  - SQLite3 driver
- **Python 3.14**: Analytics and UI
  - Flask web framework
  - Requests library
  - SQLite3 (built-in)
- **SQLite**: Lightweight database for both services
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
