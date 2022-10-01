package main

import (
	"os"
	"strconv"
)

func GenerateGeo(f *os.File) {
	f.WriteString(
`[    
    {
        "hotelId": "1",
        "lat": 37.7867,
        "lon": -122.4112
    },
    {
        "hotelId": "2",
        "lat": 37.7854,
        "lon": -122.4005
    },
    {
        "hotelId": "3",
        "lat": 37.7834,
        "lon": -122.4071
    },
    {
        "hotelId": "4",
        "lat": 37.7936,
        "lon": -122.3930
    },
    {
        "hotelId": "5",
        "lat": 37.7831,
        "lon": -122.4181
    },
    {
        "hotelId": "6",
        "lat": 37.7863,
        "lon": -122.4015
    }`)
		
	// add up to 80 hotels
	for i := 7; i <= 80; i++ {
		hotel_id := strconv.Itoa(i)
		lat := strconv.FormatFloat(37.7835 + float64(i)/500.0*3, 'f', 4, 64)
		lon := strconv.FormatFloat(-122.41 + float64(i)/500.0*4, 'f', 4, 64)

		f.WriteString(",")
		f.WriteString(
`    
    {
        "hotelId": "` + hotel_id + `",
        "lat": ` + lat + `,
        "lon": ` + lon + `
    }`)
	}
	f.WriteString("\n]\n")
}
