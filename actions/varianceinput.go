package actions

import "github.com/gobuffalo/buffalo"

//SecLocalVariance loads all local variance
func SecLocalVariance(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	currentQuarter := MakeQuarterFormalDate(request.Quarter, request.Year)
	variance := LoadSecLocalVariance(MakeLastQuarterFormalDate(currentQuarter), currentQuarter)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Local variance loaded", "data": map[string]interface{}{"variance": variance}}))
}

//NpraLocalVariance loads all local variance
func NpraLocalVariance(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	currentQuarter := MakeQuarterFormalDate(request.Quarter, request.Year)
	variance := LoadNpraLocalVariance(MakeLastQuarterFormalDate(currentQuarter), currentQuarter)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Local variance loaded", "data": map[string]interface{}{"variance": variance}}))
}

//SecForeignVariance loads all foreign variance
func SecForeignVariance(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	currentQuarter := MakeQuarterFormalDate(request.Quarter, request.Year)
	variance := LoadSecForeignVariance(MakeLastQuarterFormalDate(currentQuarter), currentQuarter)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Foreign variance loaded", "data": map[string]interface{}{"variance": variance}}))
}

//NpraForeignVariance loads all foreign variance
func NpraForeignVariance(c buffalo.Context) error {
	request := &QuarterAndYear{}
	if err := c.Bind(request); err != nil {
		return c.Render(200, r.JSON(map[string]interface{}{"error": true, "message": err.Error()}))
	}
	currentQuarter := MakeQuarterFormalDate(request.Quarter, request.Year)
	variance := LoadNpraForeignVariance(MakeLastQuarterFormalDate(currentQuarter), currentQuarter)
	return c.Render(200, r.JSON(map[string]interface{}{"error": false, "message": "Foreign variance loaded", "data": map[string]interface{}{"variance": variance}}))
}
