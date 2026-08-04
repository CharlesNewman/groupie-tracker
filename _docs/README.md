<h1 align="center">Groupie Tracker</h1>

<p align="center">
  <img src="icon.png" alt="Groupie Tracker logo" width="220">
</p>
<p align="center">
    <a href="https://go.dev/dl/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go" alt="Go Version" /></a>
    <img src="https://img.shields.io/badge/build-passing-brightgreen?style=for-the-badge" alt="Build Status" />
    <img src="https://img.shields.io/badge/Campus-zone01.gr-blue?style=for-the-badge&logo=gitea" alt="Campus" />
    <img src="https://img.shields.io/badge/Cohort-2.3-ed1c24?style=for-the-badge&labelColor=grey&logo=medusa" alt="Cohort" />
</p>

---

## 📖 Overview

**Groupie Tracker** is a Go web application that receives data from the Groupie Trackers API and displays information about artists and bands.

The application uses the API's four data sections:

- **Artists** — names, images, members, creation dates and first album dates.
- **Locations** — concert locations for each artist.
- **Dates** — concert dates for each artist.
- **Relations** — connects concert locations with their matching dates.

The website includes artist cards, search, live suggestions, pagination, a Quick View side panel, detailed artist pages, concert information, an interactive map, dark mode and custom error pages.

---

## 🚀 Usage

Run the project from the root directory:

```bash
go run ./cmd
```

The server starts at:

```text
http://localhost:8080
```

A different port can be selected with the `PORT` environment variable:

```bash
PORT=3000 go run ./cmd
```

---

## ✨ Features

- Artist cards with images and names.
- Pagination with 12 artists per page.
- Search by artist name.
- Live search suggestions using asynchronous JavaScript requests.
- Quick View side panel without reloading the page.
- Artist members, creation year and first album date.
- Dynamic total concert and location counts.
- Expandable concert-location list.
- Full artist details page.
- Concert locations connected to their dates.
- Interactive Leaflet and OpenStreetMap concert map.
- Light and dark modes stored in `localStorage`.
- Responsive layout for desktop, tablet and mobile.
- Custom `400`, `404`, `405` and `500` error handling.
- Automated tests for models, handlers, routes and server behavior.

---

## 🌐 Routes

| Route | Description |
|---|---|
| `/` | Displays the paginated artist list. |
| `/search?query=name` | Displays artists matching the search text. |
| `/suggestions?query=name` | Returns matching artists as JSON for live suggestions. |
| `/artist?id=1` | Returns artist Quick View data as JSON. |
| `/artist-details?id=1` | Displays the full artist details page. |
| `/static/` | Serves CSS, JavaScript and image files. |

---

## 🔄 Client–Server Event

The main client–server event is the **Quick View** button:

1. The user clicks a Quick View button.
2. JavaScript reads the artist ID from the button.
3. `fetch()` sends a GET request to `/artist?id=<id>`.
4. The Go handler finds the matching artist, relation, location and date data.
5. The server returns JSON.
6. JavaScript updates and opens the side panel without reloading the page.

Live search suggestions use the same client–server pattern through the `/suggestions` route.

---

## 📁 Structure

```text
groupie-tracker/
├── cmd/
│   ├── main.go                 # Application entry point
│   └── main_test.go            # Entry-point tests
├── internal/
│   ├── api.go                  # API requests and JSON decoding
│   ├── handlers.go             # HTTP handlers and error handling
│   ├── models.go               # Application data structures
│   ├── server.go               # Route setup and server startup
│   ├── handlers_test.go        # Handler status tests
│   ├── models_test.go          # Model tests
│   └── server_test.go          # Route and integration-style tests
├── templates/
│   ├── index.html              # Home, search results and side panel
│   ├── artist-details.html     # Full artist and concert page
│   └── error.html              # Custom error page
├── static/
│   ├── style.css               # Light, dark and responsive styling
│   ├── sidepanel.js            # Quick View client-server event
│   ├── search.js               # Live search suggestions
│   ├── map.js                  # Concert map and geocoding
│   ├── darkmode.js             # Theme switching
│   ├── background.js           # Interactive background
│   └── logo.png                # Project logo
├── _docs/
│   └── learning-guide.md
└── go.mod
```

---

## 🧠 Application Flow

1. **Startup:** `StartServer()` requests artists, relations, locations and dates from the API.
2. **Validation:** The server checks request errors, response status codes and empty API responses.
3. **Route setup:** `SetupRoutes()` registers the page, JSON and static-file handlers.
4. **Rendering:** Go templates generate the home, details and error pages.
5. **Quick View:** JavaScript requests selected artist data as JSON and fills the side panel.
6. **Search:** The server filters artist names; JavaScript displays live suggestions.
7. **Concert display:** Relation data connects each location with its concert dates.
8. **Map:** The browser geocodes concert locations and places markers on a Leaflet map.

---

## 🛠 Technical Stack

- **Backend:** Go 1.22.2
- **Backend packages:** Go standard library only
- **Frontend:** HTML, CSS and JavaScript
- **Templates:** `html/template`
- **HTTP server:** `net/http`
- **JSON:** `encoding/json`
- **Map:** Leaflet with OpenStreetMap tiles
- **Geocoding:** Nominatim OpenStreetMap search service

> The Go backend uses only standard-library packages. Leaflet is loaded in the browser for the map interface.

---

## ✅ Testing

Run all tests:

```bash
go test ./...
```

Run tests with detailed output:

```bash
go test -v ./...
```

The test suite covers:

- Successful routes.
- Invalid routes.
- Invalid IDs and page numbers.
- Wrong HTTP methods.
- Search and suggestion responses.
- Quick View JSON data.
- Location and concert counts.
- Artist details rendering.
- Model data structures.
- Main server success and error behavior.

---

## Build & Run

```bash
# Build the application
go build -o groupie-tracker ./cmd

# Run the compiled application
./groupie-tracker
```

---

## 📝 Author

**Charles Newman** | cnewman | Cohort 2.3 | Zone01 Athens

---

## 📄 License

Completed as part of the Zone01 Athens curriculum — Cohort 2.3.
