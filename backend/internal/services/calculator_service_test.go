package services

import (
	"testing"
)

func TestCalculateTirePressure_RoadClincher(t *testing.T) {
	req := CalculatorRequest{
		RiderWeight: 70,
		BikeWeight:  8,
		TireWidth:   28,
		TireType:    "clincher",
		Surface:     "road",
		Conditions:  "dry",
	}

	res := CalculateTirePressure(req)

	if res.Unit != "bar" {
		t.Errorf("expected unit bar, got %s", res.Unit)
	}
	if res.FrontPressure <= 0 || res.FrontPressure > 9 {
		t.Errorf("front pressure out of range: %f", res.FrontPressure)
	}
	if res.RearPressure <= 0 || res.RearPressure > 9 {
		t.Errorf("rear pressure out of range: %f", res.RearPressure)
	}
	if res.RearPressure <= res.FrontPressure {
		t.Errorf("rear pressure should be higher than front: front=%f rear=%f", res.FrontPressure, res.RearPressure)
	}
}

func TestCalculateTirePressure_TubelessLower(t *testing.T) {
	clincher := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 70, BikeWeight: 8, TireWidth: 28, TireType: "clincher", Surface: "road",
	})
	tubeless := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 70, BikeWeight: 8, TireWidth: 28, TireType: "tubeless", Surface: "road",
	})

	if tubeless.FrontPressure >= clincher.FrontPressure {
		t.Errorf("tubeless should have lower pressure than clincher: tubeless=%f clincher=%f",
			tubeless.FrontPressure, clincher.FrontPressure)
	}
}

func TestCalculateTirePressure_GravelLower(t *testing.T) {
	road := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 70, BikeWeight: 8, TireWidth: 28, TireType: "clincher", Surface: "road",
	})
	gravel := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 70, BikeWeight: 8, TireWidth: 28, TireType: "clincher", Surface: "gravel",
	})

	if gravel.FrontPressure >= road.FrontPressure {
		t.Errorf("gravel should have lower pressure than road: gravel=%f road=%f",
			gravel.FrontPressure, road.FrontPressure)
	}
}

func TestCalculateTirePressure_WetLower(t *testing.T) {
	dry := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 70, BikeWeight: 8, TireWidth: 28, TireType: "clincher", Surface: "road", Conditions: "dry",
	})
	wet := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 70, BikeWeight: 8, TireWidth: 28, TireType: "clincher", Surface: "road", Conditions: "wet",
	})

	if wet.FrontPressure >= dry.FrontPressure {
		t.Errorf("wet should have lower pressure than dry: wet=%f dry=%f",
			wet.FrontPressure, dry.FrontPressure)
	}
}

func TestCalculateTirePressure_HeavierRiderHigher(t *testing.T) {
	light := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 60, BikeWeight: 8, TireWidth: 28, TireType: "clincher", Surface: "road",
	})
	heavy := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 90, BikeWeight: 8, TireWidth: 28, TireType: "clincher", Surface: "road",
	})

	if heavy.FrontPressure <= light.FrontPressure {
		t.Errorf("heavier rider should have higher pressure: heavy=%f light=%f",
			heavy.FrontPressure, light.FrontPressure)
	}
}

func TestCalculateTirePressure_WiderTireLower(t *testing.T) {
	narrow := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 70, BikeWeight: 8, TireWidth: 25, TireType: "clincher", Surface: "road",
	})
	wide := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 70, BikeWeight: 8, TireWidth: 40, TireType: "clincher", Surface: "road",
	})

	if wide.FrontPressure >= narrow.FrontPressure {
		t.Errorf("wider tire should have lower pressure: wide=%f narrow=%f",
			wide.FrontPressure, narrow.FrontPressure)
	}
}

func TestCalculateTirePressure_ClampMin(t *testing.T) {
	res := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 30, BikeWeight: 5, TireWidth: 60, TireType: "tubeless", Surface: "gravel", Conditions: "mud",
	})

	if res.FrontPressure < 1.5 {
		t.Errorf("front pressure should be clamped to minimum 1.5: %f", res.FrontPressure)
	}
	if res.RearPressure < 1.5 {
		t.Errorf("rear pressure should be clamped to minimum 1.5: %f", res.RearPressure)
	}
}

func TestCalculateTirePressure_Recommendations(t *testing.T) {
	res := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 70, BikeWeight: 8, TireWidth: 28, TireType: "tubeless", Surface: "gravel", Conditions: "wet",
	})

	if len(res.Recommendations) == 0 {
		t.Error("expected recommendations, got none")
	}

	found := false
	for _, r := range res.Recommendations {
		if r == "Tubeless setup allows lower pressures for better grip and comfort." {
			found = true
		}
	}
	if !found {
		t.Error("expected tubeless recommendation")
	}
}

func TestCalculateTirePressure_NarrowTireRecommendation(t *testing.T) {
	res := CalculateTirePressure(CalculatorRequest{
		RiderWeight: 70, BikeWeight: 8, TireWidth: 23, TireType: "clincher", Surface: "road",
	})

	found := false
	for _, r := range res.Recommendations {
		if r == "Narrow tires: be careful not to go too low to avoid pinch flats." {
			found = true
		}
	}
	if !found {
		t.Error("expected narrow tire recommendation for width 23mm")
	}
}

func TestCalculateTirePressure_Surfaces(t *testing.T) {
	surfaces := []string{"road", "gravel", "mixed", "cobblestone"}
	for _, s := range surfaces {
		res := CalculateTirePressure(CalculatorRequest{
			RiderWeight: 70, BikeWeight: 8, TireWidth: 28, TireType: "clincher", Surface: s,
		})
		if res.FrontPressure <= 0 {
			t.Errorf("surface %s: got zero/negative front pressure", s)
		}
	}
}
