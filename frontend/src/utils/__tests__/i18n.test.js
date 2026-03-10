import { describe, it, expect, beforeEach, vi } from 'vitest'

// Mock localStorage
const storage = {}
vi.stubGlobal('localStorage', {
  getItem: (key) => storage[key] || null,
  setItem: (key, val) => { storage[key] = val },
  removeItem: (key) => { delete storage[key] },
})

describe('i18n', () => {
  let t, setLanguage, getLanguage

  beforeEach(async () => {
    // Fresh import each time
    const mod = await import('../../utils/i18n.js')
    t = mod.t
    setLanguage = mod.setLanguage
    getLanguage = mod.getLanguage
  })

  it('returns translation for known key', () => {
    setLanguage('en')
    expect(t('nav_races')).toBe('Races')
  })

  it('returns Russian translation', () => {
    setLanguage('ru')
    expect(t('nav_races')).toBe('Гонки')
  })

  it('returns key for unknown translation', () => {
    expect(t('unknown_key_xyz')).toBe('unknown_key_xyz')
  })

  it('interpolates params', () => {
    setLanguage('en')
    const result = t('hi_rider', { name: 'Alex' })
    expect(result).toBe('Hi, Alex!')
  })

  it('getLanguage returns current language', () => {
    setLanguage('en')
    expect(getLanguage()).toBe('en')
    setLanguage('ru')
    expect(getLanguage()).toBe('ru')
  })
})
