import { ObjectId } from 'bson'
export interface Earthquake {
    _id: ObjectId
    mag: number
    place: string
    time: Date
    location: EarthquakeLocation
}

export interface EarthquakeLocation {
    type: string
    coordinates: number[]
}

export interface LatLng { lat: number, lng: number }

export interface EarthquakeCount {
    count: number
}