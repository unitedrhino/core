package commonlogic

import (
	"errors"
	"testing"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/def"
)

func TestSelectWeatherPositionPrefersProjectPositionOverRequestGPS(t *testing.T) {
	got, projectID, err := selectWeatherPosition(
		&sys.Point{Longitude: 116.39, Latitude: 39.90},
		0,
		1001,
		&sys.ProjectInfo{Position: &sys.Point{Longitude: 121.47, Latitude: 31.23}},
		nil,
	)
	if err != nil {
		t.Fatalf("selectWeatherPosition error: %v", err)
	}
	if projectID != 1001 {
		t.Fatalf("projectID = %d, want 1001", projectID)
	}
	if got.GetLongitude() != 121.47 || got.GetLatitude() != 31.23 {
		t.Fatalf("position = (%v,%v), want project position", got.GetLongitude(), got.GetLatitude())
	}
}

func TestSelectWeatherPositionFallsBackToGPSWhenImplicitProjectFails(t *testing.T) {
	wantErr := errors.New("project cache unavailable")
	got, projectID, err := selectWeatherPosition(
		&sys.Point{Longitude: 116.39, Latitude: 39.90},
		0,
		1001,
		nil,
		wantErr,
	)
	if err != nil {
		t.Fatalf("selectWeatherPosition error: %v", err)
	}
	if projectID != 0 {
		t.Fatalf("projectID = %d, want fallback GPS without project", projectID)
	}
	if got.GetLongitude() != 116.39 || got.GetLatitude() != 39.90 {
		t.Fatalf("position = (%v,%v), want request GPS", got.GetLongitude(), got.GetLatitude())
	}
}

func TestWeatherProjectIDIgnoresVirtualProject(t *testing.T) {
	projectID, explicit := weatherProjectID(0, def.NotClassified)

	if projectID != 0 {
		t.Fatalf("projectID = %d, want no project for virtual context", projectID)
	}
	if explicit {
		t.Fatal("explicit = true, want false for virtual context")
	}
}

func TestWeatherGeoCacheLocationUsesCityBucket(t *testing.T) {
	key, lookup := weatherGeoCacheKeyAndLookup(116.39123, 39.90456)

	if key != "sys:common:geo:39.9:116.4" {
		t.Fatalf("key = %q, want city bucket key", key)
	}
	if lookup != "116.4,39.9" {
		t.Fatalf("lookup = %q, want city bucket lookup", lookup)
	}
}
