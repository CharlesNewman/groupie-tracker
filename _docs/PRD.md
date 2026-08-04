# Product Requirements Document (PRD): Groupie Tracker

## 1. Overview

The **Groupie Tracker** is a Go web application that consumes the Groupie Trackers API and presents artist, band and concert information through a responsive website.

The API is divided into four connected datasets:

- Artists
- Locations
- Dates
- Relations

The product combines these datasets into artist cards, Quick View panels, detailed artist pages, concert lists and an interactive map.

## 2. Objectives

- Retrieve and use all four sections of the Groupie Trackers API.
- Display artist information in a clear and user-friendly interface.
- Connect concert locations with their matching dates.
- Implement a working client-server event.
- Use correct HTTP methods and status codes.
- Handle invalid input and server errors without crashing.
- Use only Go standard-library packages in the backend.
- Provide automated tests for important application behavior.

## 3. Target Users

- Users who want to browse artists and bands.
- Users who want to view artist members and historical information.
- Users who want to explore concert dates and locations.
- Zone01 auditors reviewing HTTP, API and client-server functionality.

## 4. Key Features

### 4.1 API Data Retrieval (`api.go`)

- Fetches artist data from `/api/artists`.
- Fetches relation data from `/api/relation`.
- Fetches location data from `/api/locations`.
- Fetches date data from `/api/dates`.
- Checks HTTP response status before decoding JSON.
- Returns errors when requests or decoding fail.
- Closes every API response body with `defer`.

### 4.2 Server Startup and Routing (`server.go`)

- Loads all required API data before starting the server.
- Rejects empty API responses.
- Uses `http.NewServeMux` for route registration.
- Uses port `8080` by default.
- Supports a custom port through the `PORT` environment variable.
- Serves static CSS, JavaScript and image files.

### 4.3 Artist List and Pagination (`Handler`)

- Displays artist cards on the home page.
- Shows 12 artists per page.
- Calculates previous and next pages dynamically.
- Rejects invalid page values with `400 Bad Request`.
- Rejects pages outside the available range with `404 Not Found`.

### 4.4 Search and Suggestions

- Filters artists by name without case sensitivity.
- Rejects empty search requests.
- Displays a clear no-results message when no artist matches.
- Returns live search suggestions as JSON.
- Allows a user to open an artist details page from a suggestion.

### 4.5 Quick View Client–Server Event

- Adds a Quick View button to every artist card.
- Uses JavaScript `fetch()` to request `/artist?id=<id>`.
- Returns artist and relation data as JSON.
- Uses location data to calculate the number of concert locations.
- Uses date data to calculate the total number of concerts.
- Opens a side panel without refreshing the page.
- Displays artist members, a concert statistic and an expandable location list.

### 4.6 Artist Details Page

- Displays artist image, name, creation year and first album date.
- Displays the complete member list.
- Displays concert locations with matching dates.
- Formats API location strings for normal reading.
- Provides buttons that focus the matching location on the map.

### 4.7 Interactive Concert Map

- Uses Leaflet in the browser.
- Uses OpenStreetMap map tiles.
- Uses Nominatim to convert location names into coordinates.
- Places a marker for each concert location.
- Fits the map to the available markers.
- Scrolls and zooms to a location when its button is selected.

### 4.8 Theme and Responsive Design

- Supports light and dark themes.
- Stores the selected theme in browser `localStorage`.
- Uses responsive artist grids for desktop, tablet and mobile.
- Converts the Quick View panel into a full-screen panel on smaller screens.
- Includes an animated background with reduced-motion support.

### 4.9 Error Handling

- Uses a styled HTML error page for full-page routes.
- Uses plain HTTP errors for JSON routes consumed by JavaScript.
- Supports:
  - `400 Bad Request`
  - `404 Not Found`
  - `405 Method Not Allowed`
  - `500 Internal Server Error`
- Prevents unknown paths from being treated as the home page.
- Returns instead of continuing after an error.

### 4.10 Automated Tests

- Tests models and response structures.
- Tests successful home, search, suggestion and details routes.
- Tests Quick View JSON output.
- Tests location and concert counts.
- Tests invalid IDs and page values.
- Tests unknown routes and unsupported HTTP methods.
- Tests server wrapper success and failure behavior.

## 5. Functional Requirements

### 5.1 Artist Data

The system must display:

- Artist name
- Image
- Members
- Creation date
- First album release date

### 5.2 Concert Data

The system must display:

- Concert locations
- Concert dates
- Location-to-date relationships
- Total location count
- Total concert count

### 5.3 Client–Server Communication

The system must include at least one browser action that communicates with the Go server.

The required implementation is the Quick View action:

1. User clicks Quick View.
2. Browser requests artist data from the Go server.
3. Server validates the ID and returns JSON.
4. Browser updates the side panel asynchronously.

### 5.4 HTTP Behavior

- Page and data routes use GET requests.
- Unsupported methods return `405`.
- Invalid user input returns `400`.
- Missing pages or artists return `404`.
- Template or internal failures return `500`.

## 6. Technical Stack

- **Language:** Go 1.22.2
- **Backend:** Go standard library
- **HTTP:** `net/http`
- **JSON:** `encoding/json`
- **Templates:** `html/template`
- **Frontend:** HTML, CSS and vanilla JavaScript
- **Map UI:** Leaflet
- **Map data:** OpenStreetMap
- **Geocoding:** Nominatim

## 7. Architecture

### 7.1 `cmd`

- Contains the application entry point.
- Calls the internal server startup function.
- Prints startup failures.

### 7.2 `internal/api.go`

- Handles external API requests.
- Decodes JSON into application models.

### 7.3 `internal/models.go`

- Defines artists, relations, locations and dates.
- Defines template pagination data.
- Defines the combined Quick View JSON response.

### 7.4 `internal/handlers.go`

- Handles pages, search and JSON endpoints.
- Performs artist lookups.
- Calculates location and concert totals.
- Renders templates and errors.

### 7.5 `internal/server.go`

- Loads API data.
- Registers routes.
- Starts the HTTP server.

### 7.6 `templates`

- Contains the home page, details page and custom error page.

### 7.7 `static`

- Contains styling, theme behavior, search, Quick View, map logic, background behavior and the logo.

## 8. Non-Functional Requirements

- **Reliability:** The application must handle invalid input without crashing.
- **Performance:** API data is fetched once during startup instead of on every page request.
- **Usability:** Users can browse, search and open artist information with clear controls.
- **Responsiveness:** The interface must work across desktop and mobile screen sizes.
- **Maintainability:** API logic, handlers, models, routes and frontend files are separated.
- **Testability:** Route setup accepts data as arguments, allowing handlers to be tested without calling the real API.
- **Accessibility:** Images use alternative text, controls use buttons or links, and reduced-motion preferences are respected.

## 9. Dependencies and Constraints

- The Go backend must use standard-library packages only.
- The initial server startup requires access to the Groupie Trackers API.
- The interactive map requires internet access to Leaflet, OpenStreetMap and Nominatim.
- If the API cannot be reached, the application returns a startup error instead of starting with missing data.

## 10. Success Criteria

The project is successful when:

- All four API datasets are retrieved and used.
- Required artist information is displayed correctly.
- Search, pagination, Quick View and details pages work.
- Concert locations and dates are presented correctly.
- The Quick View event completes client-server communication.
- Invalid requests return correct HTTP status codes.
- The server remains stable during normal and invalid requests.
- All automated tests pass with `go test ./...`.

## 11. Future Enhancements

- Add filters for creation year, member count and concert location.
- Cache geocoding results to reduce repeated external requests.
- Add request timeouts to API and geocoding calls.
- Add loading and visible error states to Quick View and map requests.
- Add keyboard navigation for live search suggestions.
- Deploy the application with a public domain.
- Add optional concurrent API loading with goroutines and channels.
