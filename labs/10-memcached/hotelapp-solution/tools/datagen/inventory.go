package main

import (
	"os"
	"strconv"
)

func GenerateInventory(f *os.File) {
	f.WriteString(
`[    
    {
        "hotelId": "1",
        "code": "RACK",
        "inDate": "2015-04-09",
        "outDate": "2015-04-10",
        "roomType": {
            "bookableRate": 109.00,
            "code": "KNG",
            "description": "King sized bed",
            "totalRate": 109.00,
            "totalRateInclusive": 123.17
        }
    },
    {
        "hotelId": "2",
        "code": "RACK",
        "inDate": "2015-04-09",
        "outDate": "2015-04-10",
        "roomType": {
            "bookableRate": 139.00,
            "code": "QN",
            "description": "Queen sized bed",
            "totalRate": 139.00,
            "totalRateInclusive": 153.09
        }
    },
    {
        "hotelId": "3",
        "code": "RACK",
        "inDate": "2015-04-09",
        "outDate": "2015-04-10",
        "roomType": {
            "bookableRate": 109.00,
            "code": "KNG",
            "description": "King sized bed",
            "totalRate": 109.00,
            "totalRateInclusive": 123.17
        }
    }`)

	// add up to 80 hotels
	for i := 7; i <= 80; i++ {
		if i%3 == 0 {
			hotel_id := strconv.Itoa(i)
			end_date := "2015-04-"
			rate := 109.00
			rate_inc := 123.17
			if i%2 == 0 {
				end_date = end_date + "17"
			} else {
				end_date = end_date + "24"
			}

			if i%5 == 1 {
				rate = 120.00
				rate_inc = 140.00
			} else if i%5 == 2 {
				rate = 124.00
				rate_inc = 144.00
			} else if i%5 == 3 {
				rate = 132.00
				rate_inc = 158.00
			} else if i%5 == 4 {
				rate = 232.00
				rate_inc = 258.00
			}

            rate_str := strconv.FormatFloat(rate, 'f', 2, 32)
            rate_inc_str := strconv.FormatFloat(rate_inc, 'f', 2, 32)
    
            f.WriteString(",")
            f.WriteString(
`    
    {
        "hotelId": "` + hotel_id + `",
        "code": "RACK",
        "inDate": "2015-04-09",
        "outDate": "` + end_date + `",
        "roomType": {
            "bookableRate": 109.00,
            "code": "KNG",
            "description": "King sized bed",
            "totalRate": ` + rate_str + `,
            "totalRateInclusive": ` + rate_inc_str + `
        }
    }`)
        }    
	}
    f.WriteString("\n]\n")
}
