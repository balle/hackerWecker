package hackerWecker

import "testing"

func TestGetWeather(t *testing.T) {
	w := GetWeather()

	if w.Description == "" {
		t.Errorf("Description is empty")
	}

	if w.Temp == 0 {
		t.Errorf("Temp is empty")
	}

	if w.TempMin == 0 {
		t.Errorf("TempMin is empty")
	}

	if w.TempMax == 0 {
		t.Errorf("TempMax is empty")
	}

	if w.Sunrise == 0 {
		t.Errorf("Sunrise is empty")
	}

	if w.Sunset == 0 {
		t.Errorf("Sunset is empty")
	}

	if w.TempUnit == "" {
		t.Errorf("TempUnit is empty")
	}
}

func TestReadWeather(t *testing.T) {
	ReadWeather()
}
