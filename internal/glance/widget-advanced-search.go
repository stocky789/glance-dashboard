package glance

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"
)

type advancedSearchWidget struct {
	widgetBase
	Placeholder string `yaml:"placeholder"`
}

type searchResult struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	Score    float64 `json:"score"`
}

func (w *advancedSearchWidget) initialize() error {
	if w.Placeholder == "" {
		w.Placeholder = "Search dashboard..."
	}
	w.Title = "Advanced Search"
	w.cacheDuration = 1 * time.Hour
	w.cacheType = cacheTypeInfinite
	w.ContentAvailable = true
	return nil
}

func (w *advancedSearchWidget) update(ctx context.Context) {
	// Search widget doesn't need active updates
}

func (w *advancedSearchWidget) Render() template.HTML {
	tmpl := template.Must(template.New("search").Parse(`
		<div class="advanced-search-widget">
			<div class="search-container">
				<input 
					type="text" 
					class="search-input" 
					placeholder="{{.Placeholder}}"
					id="search-{{.ID}}"
					autocomplete="off"
				/>
				<div class="search-filters">
					<label><input type="checkbox" class="filter-type" value="widget"/> Widgets</label>
					<label><input type="checkbox" class="filter-type" value="page"/> Pages</label>
					<label><input type="checkbox" class="filter-type" value="bookmark"/> Bookmarks</label>
				</div>
			</div>
			<div class="search-results" id="results-{{.ID}}"></div>
			<style>
				.advanced-search-widget {
					padding: 1rem;
				}
				.search-container {
					display: flex;
					flex-direction: column;
					gap: 0.5rem;
				}
				.search-input {
					width: 100%;
					padding: 0.75rem;
					font-size: 1rem;
					border: 1px solid rgba(255,255,255,0.2);
					border-radius: 6px;
					background: rgba(0,0,0,0.3);
					color: inherit;
				}
				.search-input:focus {
					outline: none;
					border-color: var(--primary-color, #00ff00);
					box-shadow: 0 0 10px rgba(0,255,0,0.3);
				}
				.search-filters {
					display: flex;
					gap: 1rem;
					font-size: 0.85rem;
				}
				.search-filters label {
					display: flex;
					align-items: center;
					gap: 0.4rem;
					cursor: pointer;
				}
				.search-filters input[type="checkbox"] {
					cursor: pointer;
				}
				.search-results {
					margin-top: 1rem;
					max-height: 300px;
					overflow-y: auto;
				}
				.search-result-item {
					padding: 0.5rem;
					margin: 0.25rem 0;
					background: rgba(0,0,0,0.2);
					border-radius: 4px;
					cursor: pointer;
					transition: background 0.2s;
				}
				.search-result-item:hover {
					background: rgba(0,0,0,0.4);
				}
				.result-title {
					font-weight: bold;
					font-size: 0.9rem;
				}
				.result-type {
					font-size: 0.75rem;
					opacity: 0.7;
					text-transform: uppercase;
				}
				.result-score {
					float: right;
					color: var(--primary-color, #00ff00);
					font-size: 0.8rem;
				}
			</style>
			<script>
				(function() {
					const searchInput = document.getElementById('search-{{.ID}}');
					const resultsDiv = document.getElementById('results-{{.ID}}');
					
					searchInput.addEventListener('input', async (e) => {
						const query = e.target.value;
						if (query.length < 2) {
							resultsDiv.innerHTML = '';
							return;
						}
						
						try {
							const response = await fetch('/api/v1/search?q=' + encodeURIComponent(query));
							const results = await response.json();
							
						resultsDiv.innerHTML = results.map(r => '<div class="search-result-item"><div class="result-title">' + r.title + '<span class="result-score">' + (r.score * 100).toFixed(0) + '%</span></div><div class="result-type">' + r.type + '</div></div>').join('');
						} catch (err) {
							console.error('Search error:', err);
						}
					});
					
					searchInput.addEventListener('keydown', (e) => {
						if (e.key === 'Escape') {
							searchInput.value = '';
							resultsDiv.innerHTML = '';
						}
					});
				})();
			</script>
		</div>
	`))

	w.templateBuffer.Reset()
	tmpl.Execute(&w.templateBuffer, w)
	return template.HTML(w.templateBuffer.String())
}

func (w *advancedSearchWidget) GetType() string {
	return "advanced-search"
}

func (w *advancedSearchWidget) handleRequest(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(map[string]string{"status": "ok"})
}

func (w *advancedSearchWidget) setHideHeader(b bool) {
	w.HideHeader = b
}
