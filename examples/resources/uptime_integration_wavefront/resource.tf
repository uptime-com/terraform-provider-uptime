resource "uptime_integration_wavefront" "example" {
  name          = "My Wavefront Integration"
  wavefront_url = "https://example.wavefront.com"
  api_token     = "your-wavefront-api-token"
}
