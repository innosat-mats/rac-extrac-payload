package innosat

import "time"

// Specification describes what version the current implementation follows
const Specification string = "IS-OSE-ICD-0005:1"

// Epoch is the CUC Epoch of the platform, GPS Time
var Epoch time.Time = time.Date(1980, time.January, 6, 0, 0, 1, 0, time.UTC)
