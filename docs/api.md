# API Documentation

Endpoints

Analyze Web page

##  GET localhost:8080/url_analyze?url={website_url}

  Description: Analyzes the provided web page URL and returns insights.

  Parameters:

   url (query parameter): The URL to analyze.

   Response:

   {
   "html_version": "HTML5",
   "title": "Example Page",
   "headings": { "h1": 2, "h2": 3 },
   "internal_links": 5,
   "external_links": 3,
    "broken_links": 1,
    "has_login_form": true
     }

Metrics

 ## GET localhost:8080/metrics

   Description: Exposes Prometheus metrics for monitoring.



   