import { describe, it, expect } from 'vitest'
import { useToast } from '../useToast'

describe('useToast', () => {
  it('starts with show=false', () => {
    const { show } = useToast()
    expect(show.value).toBe(false)
  })

  it('shows toast with message', () => {
    const { toast, message, type, show } = useToast()
    toast('Saved!', 'success')
    expect(show.value).toBe(true)
    expect(message.value).toBe('Saved!')
    expect(type.value).toBe('success')
  })

  it('shows error toast', () => {
    const { toast, type } = useToast()
    toast('Failed', 'error')
    expect(type.value).toBe('error')
  })

  it('close hides toast', () => {
    const { toast, close, show } = useToast()
    toast('Test')
    expect(show.value).toBe(true)
    close()
    expect(show.value).toBe(false)
  })
})
