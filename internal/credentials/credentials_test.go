package credentials

import "testing"

func TestGooglePlacesEnvironmentPrecedence(t *testing.T) {
	t.Setenv("GOOGLE_PLACES_API_KEY", "places-env")
	t.Setenv("GOOGLE_MAPS_API_KEY", "maps-env")
	v, source, err := ResolveGooglePlaces()
	if err != nil || v != "places-env" || source != "environment" {
		t.Fatalf("Google Places environment precedence failed: %q %q %v", v, source, err)
	}
	t.Setenv("GOOGLE_PLACES_API_KEY", "")
	v, source, err = ResolveGooglePlaces()
	if err != nil || v != "maps-env" || source != "environment" {
		t.Fatalf("Google Maps compatibility environment fallback failed: %q %q %v", v, source, err)
	}
}
