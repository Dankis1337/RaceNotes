import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Toast from '../Toast.vue'

describe('Toast', () => {
  it('is hidden when show is false', () => {
    const wrapper = mount(Toast, {
      props: { message: 'Test', type: 'success', show: false },
    })
    expect(wrapper.find('.fixed').exists()).toBe(false)
  })

  it('displays message when show is true', () => {
    const wrapper = mount(Toast, {
      props: { message: 'Saved!', type: 'success', show: true },
    })
    expect(wrapper.text()).toContain('Saved!')
  })

  it('applies error styling for error type', () => {
    const wrapper = mount(Toast, {
      props: { message: 'Error!', type: 'error', show: true },
    })
    const toast = wrapper.find('.bg-red-500')
    expect(toast.exists()).toBe(true)
  })

  it('applies success styling for success type', () => {
    const wrapper = mount(Toast, {
      props: { message: 'OK', type: 'success', show: true },
    })
    const toast = wrapper.find('.bg-primary')
    expect(toast.exists()).toBe(true)
  })

  it('emits close when X button clicked', async () => {
    const wrapper = mount(Toast, {
      props: { message: 'Test', type: 'success', show: true },
    })
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
