BEGIN {
	print "Percentile (%)", "Latency (us)"

	lat2us["us"] = 1
	lat2us["ms"] = 1000
	lat2us["s"] = 10000000
}

/^[[:space:]]+[0-9]+\.?[0-9]*\%/ {
	match ($1, /([0-9]+)\%/, perc)
	match ($2, /([0-9]+\.[0-9]*)([a-z]+)/, lat)
	lat_val = lat[1]
	lat_unit = lat[2]
	print perc[1],  lat_val*lat2us[lat_unit]
}
