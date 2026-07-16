import { useEffect, useRef } from 'react'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

function decodePolyline(encoded: string): [number, number][] {
  const coords: [number, number][] = []
  let index = 0
  let lat = 0
  let lng = 0
  while (index < encoded.length) {
    let result = 0
    let shift = 0
    let b: number
    do {
      b = encoded.charCodeAt(index++) - 63
      result |= (b & 0x1f) << shift
      shift += 5
    } while (b >= 0x20)
    const dlat = (result & 1) ? ~(result >> 1) : result >> 1
    lat += dlat
    result = 0
    shift = 0
    do {
      b = encoded.charCodeAt(index++) - 63
      result |= (b & 0x1f) << shift
      shift += 5
    } while (b >= 0x20)
    const dlng = (result & 1) ? ~(result >> 1) : result >> 1
    lng += dlng
    coords.push([lat / 1e5, lng / 1e5])
  }
  return coords
}

export function LeafletMap({ polyline, accent, height = 200 }: { polyline: string; accent: string; height?: number }) {
  const ref = useRef<HTMLDivElement>(null)
  const mapRef = useRef<L.Map | null>(null)

  useEffect(() => {
    if (!ref.current || !polyline) return
    const coords = decodePolyline(polyline)
    if (coords.length < 2) return

    if (mapRef.current) {
      mapRef.current.remove()
      mapRef.current = null
    }

    const map = L.map(ref.current, { zoomControl: false, attributionControl: true })
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '&copy; OpenStreetMap',
    }).addTo(map)
    const line = L.polyline(coords, { color: accent, weight: 4 }).addTo(map)
    map.fitBounds(line.getBounds(), { padding: [20, 20] })
    mapRef.current = map

    return () => {
      map.remove()
      mapRef.current = null
    }
  }, [polyline, accent])

  return <div ref={ref} style={{ height, borderRadius: 14, overflow: 'hidden', zIndex: 0 }} />
}
