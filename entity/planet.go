package entity

import "time"

type Planet struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Climate         string `json:"climate"`
	Terrain         string `json:"terrain"`
	MovieApparition int    `json:"movie_apparition"`
}

type PlanetSwapi struct {
	Name           string    `json:"name"`
	RotationPeriod string    `json:"rotation_period"`
	OrbitalPeriod  string    `json:"orbital_period"`
	Diameter       string    `json:"diameter"`
	Climate        string    `json:"climate"`
	Gravity        string    `json:"gravity"`
	Terrain        string    `json:"terrain"`
	SurfaceWater   string    `json:"surface_water"`
	Population     string    `json:"population"`
	Residents      []string  `json:"residents"`
	Films          []string  `json:"films"`
	Created        time.Time `json:"created"`
	Edited         time.Time `json:"edited"`
	URL            string    `json:"url"`
}
