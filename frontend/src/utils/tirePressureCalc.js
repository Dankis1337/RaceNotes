/**
 * Offline tire pressure calculator — mirrors backend SRAM formula.
 * @param {Object} req - { rider_weight, bike_weight, tire_width, tire_type, surface, conditions }
 * @returns {{ front_pressure: number, rear_pressure: number, unit: string, recommendations: string[] }}
 */
export function calculateTirePressureOffline(req) {
  const totalWeight = req.rider_weight + req.bike_weight

  // Simplified SRAM formula
  let basePressure = totalWeight / (req.tire_width * 1.5)

  // Tire type adjustments
  const tireType = (req.tire_type || '').toLowerCase()
  if (tireType === 'tubeless') basePressure *= 0.93
  else if (tireType === 'tubular') basePressure *= 0.96

  // Surface adjustments
  const surface = (req.surface || '').toLowerCase()
  if (surface === 'gravel') basePressure *= 0.85
  else if (surface === 'mixed') basePressure *= 0.90
  else if (surface === 'cobblestone') basePressure *= 0.88

  // Conditions adjustments
  const conditions = (req.conditions || '').toLowerCase()
  if (conditions === 'wet') basePressure *= 0.95
  else if (conditions === 'mud') basePressure *= 0.88
  else if (conditions === 'snow') basePressure *= 0.85

  // Front/rear split: ~45/55 weight distribution
  let frontPressure = Math.round(basePressure * 0.95 * 100) / 100
  let rearPressure = Math.round(basePressure * 1.05 * 100) / 100

  // Clamp to reasonable range
  frontPressure = Math.max(1.5, Math.min(frontPressure, 9.0))
  rearPressure = Math.max(1.5, Math.min(rearPressure, 9.0))

  const recommendations = []

  if (tireType === 'tubeless') {
    recommendations.push('Tubeless setup allows lower pressures for better grip and comfort.')
  }
  if (conditions === 'wet') {
    recommendations.push('Wet conditions: consider reducing pressure by 0.1-0.2 bar for extra grip.')
  }
  if (surface === 'gravel') {
    recommendations.push('For gravel: lower pressure improves traction. Watch for pinch flats if using tubes.')
  }
  if (req.tire_width >= 35) {
    recommendations.push('Wide tires perform best at lower pressures — prioritize comfort and grip.')
  }
  if (req.tire_width <= 25) {
    recommendations.push('Narrow tires: be careful not to go too low to avoid pinch flats.')
  }

  const fLow = Math.max(1.5, frontPressure - 0.2).toFixed(1)
  const fHigh = (frontPressure + 0.2).toFixed(1)
  const rLow = Math.max(1.5, rearPressure - 0.2).toFixed(1)
  const rHigh = (rearPressure + 0.2).toFixed(1)
  recommendations.push(`Recommended range: front ${fLow}-${fHigh} bar, rear ${rLow}-${rHigh} bar.`)

  return {
    front_pressure: frontPressure,
    rear_pressure: rearPressure,
    unit: 'bar',
    recommendations,
  }
}
