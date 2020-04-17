package innosat

import "time"

// Specification describes what version the current implementation follows
var Specification string = "IS-OSE-ICD-0005"

// Epoch is the CUC Epoch of the platform, GPS Time
var Epoch = time.Date(1980, time.January, 6, 0, 0, 1, 0, time.UTC)
