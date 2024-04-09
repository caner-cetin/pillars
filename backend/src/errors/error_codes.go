package errors

const (
	InvalidRequestBody                  = "PIL-0001"
	FailedToInsertEarthquake            = "PIL-0002"
	FailedToGetEarthquakes              = "PIL-0003"
	FailedToDecodeData                  = "PIL-0004"
	FailedToDeleteDataFromMongo         = "PIL-0005"
	FailedToConnectToKafka              = "PIL-0006"
	FailedToReadFromKafka               = "PIL-0007"
	FailedToStreamToWebsocket           = "PIL-0008"
	FailedToGetEarthquakeCountFromMongo = "PIL-0009"
)
